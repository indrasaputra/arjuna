package interceptor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	// IdempotencyKeyHeader is the header key for idempotency key.
	IdempotencyKeyHeader = "x-idempotency-key"
	// DefaultIdempotencyTTL is the default TTL for idempotency key in Redis.
	DefaultIdempotencyTTL = 1 * time.Hour
)

// IdempotencyStore defines the interface for storing and retrieving idempotency responses.
type IdempotencyStore interface {
	// Get retrieves a response from the store by key.
	Get(ctx context.Context, key string) ([]byte, error)
	// Set stores a response in the store with the given key and TTL.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// IdempotencyUnaryServerInterceptor creates a unary server interceptor for idempotency check.
// It extracts the idempotency key from the request metadata, checks if it exists in the store,
// and returns the cached response if found. Otherwise, it proceeds with the request and stores the response.
func IdempotencyUnaryServerInterceptor(store IdempotencyStore) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		idpKey, err := extractIdempotencyKeyHeader(ctx)
		if err != nil {
			return handler(ctx, req)
		}

		key := fmt.Sprintf("%s:%s", info.FullMethod, idpKey)

		cachedResp, err := store.Get(ctx, key)
		// if error is nil, use normal flow
		// only when cache returns non error and has cached response, use cached response
		if err == nil && cachedResp != nil {
			resp, errResp := unmarshalCachedResponse(cachedResp)
			if errResp == nil && resp.Error != nil {
				return nil, resp.Error
			}
			if errResp == nil && resp.Response != nil {
				return resp.Response, nil
			}
		}

		resp, err := handler(ctx, req)

		if key != "" {
			if cacheErr := cacheResponse(ctx, store, key, resp, err); cacheErr != nil {
				slog.ErrorContext(ctx, "[Idempotency] failed to cache response", "key", key, "error", cacheErr)
			}
		}
		return resp, err
	}
}

func extractIdempotencyKeyHeader(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", grpcstatus.Error(codes.InvalidArgument, "missing metadata")
	}

	keys := md.Get(IdempotencyKeyHeader)
	if len(keys) == 0 {
		return "", nil
	}
	if keys[0] == "" {
		return "", nil
	}
	return keys[0], nil
}

// cachedResponse represents a cached response structure.
type cachedResponse struct {
	Response any
	Error    error
}

// cacheableError represents an error that can be serialized and cached.
type cacheableError struct {
	StatusProto []byte `json:"status_proto,omitempty"` // Marshaled google.rpc.status
}

// serializableResponse represents a serializable cached response.
type serializableResponse struct {
	Error    *cacheableError `json:"error,omitempty"`
	TypeURL  string          `json:"type_url,omitempty"`
	Response []byte          `json:"response,omitempty"`
}

// cacheResponse caches the response in the store.
func cacheResponse(ctx context.Context, store IdempotencyStore, key string, resp any, err error) error {
	sr := &serializableResponse{}

	if err != nil {
		st, ok := grpcstatus.FromError(err)
		if ok {
			statusProto := st.Proto()
			data, errMarshal := proto.Marshal(statusProto)
			if errMarshal == nil {
				sr.Error = &cacheableError{
					StatusProto: data,
				}
			}
		} else {
			st := grpcstatus.New(codes.Internal, err.Error())
			statusProto := st.Proto()
			data, _ := proto.Marshal(statusProto)
			sr.Error = &cacheableError{
				StatusProto: data,
			}
		}
	} else if resp != nil {
		protoMsg, ok := resp.(proto.Message)
		if !ok {
			return grpcstatus.Error(codes.Internal, "response is not a proto.Message")
		}

		anyMsg, errProto := anypb.New(protoMsg)
		if errProto != nil {
			return err
		}

		data, errMarshal := proto.Marshal(anyMsg)
		if errMarshal != nil {
			return err
		}

		sr.Response = data
		sr.TypeURL = anyMsg.TypeUrl
	}

	data, err := json.Marshal(sr)
	if err != nil {
		return err
	}

	return store.Set(ctx, key, data, DefaultIdempotencyTTL)
}

// unmarshalCachedResponse unmarshals a cached response.
func unmarshalCachedResponse(data []byte) (*cachedResponse, error) {
	var sr serializableResponse
	if err := json.Unmarshal(data, &sr); err != nil {
		return nil, err
	}

	cr := &cachedResponse{}

	if sr.Error != nil {
		if len(sr.Error.StatusProto) > 0 {
			var statusProto status.Status
			if err := proto.Unmarshal(sr.Error.StatusProto, &statusProto); err == nil {
				st := grpcstatus.FromProto(&statusProto)
				cr.Error = st.Err()
			} else {
				cr.Error = grpcstatus.Error(codes.Internal, "failed to unmarshal cached error")
			}
		} else {
			cr.Error = grpcstatus.Error(codes.Internal, "invalid cached error")
		}
	} else if sr.Response != nil {
		var anyMsg anypb.Any
		if err := proto.Unmarshal(sr.Response, &anyMsg); err != nil {
			return nil, err
		}

		msg, err := anyMsg.UnmarshalNew()
		if err != nil {
			return nil, err
		}

		cr.Response = msg
	}

	return cr, nil
}
