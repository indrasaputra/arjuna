package interceptor_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/interceptor"
	mock_interceptor "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/grpc/interceptor"
)

func TestIdempotencyUnaryServerInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("no idempotency key provided", func(t *testing.T) {
		store := mock_interceptor.NewMockIdempotencyStore(ctrl)

		interceptorFunc := interceptor.IdempotencyUnaryServerInterceptor(store)
		handler := func(_ context.Context, _ any) (any, error) {
			return &emptypb.Empty{}, nil
		}

		ctx := context.Background()
		resp, err := interceptorFunc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("idempotency key provided but no cached response", func(t *testing.T) {
		store := mock_interceptor.NewMockIdempotencyStore(ctrl)
		store.EXPECT().
			Get(gomock.Any(), "/test.Service/Method:test-key").
			Return(nil, nil)
		store.EXPECT().
			Set(gomock.Any(), "/test.Service/Method:test-key", gomock.Any(), interceptor.DefaultIdempotencyTTL).
			Return(nil)

		interceptorFunc := interceptor.IdempotencyUnaryServerInterceptor(store)
		handler := func(_ context.Context, _ any) (any, error) {
			return &wrapperspb.StringValue{Value: "success"}, nil
		}

		md := metadata.New(map[string]string{
			interceptor.IdempotencyKeyHeader: "test-key",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptorFunc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("idempotency key provided with cached response", func(t *testing.T) {
		testResp := &wrapperspb.StringValue{Value: "cached"}
		cachedData := createCachedResponse(t, testResp, nil)
		store := mock_interceptor.NewMockIdempotencyStore(ctrl)

		store.EXPECT().
			Get(gomock.Any(), "/test.Service/Method:test-key").
			Return(cachedData, nil)

		interceptorFunc := interceptor.IdempotencyUnaryServerInterceptor(store)
		handler := func(_ context.Context, _ any) (any, error) {
			t.Fatal("Handler should not be called when cached response exists")
			return nil, nil
		}

		md := metadata.New(map[string]string{
			interceptor.IdempotencyKeyHeader: "test-key",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptorFunc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		_, ok := resp.(*wrapperspb.StringValue)
		assert.True(t, ok, "Response should be a StringValue proto message")
	})

	t.Run("idempotency key provided with cached error", func(t *testing.T) {
		testErr := grpcstatus.Error(codes.InvalidArgument, "test error")
		cachedData := createCachedResponse(t, nil, testErr)
		store := mock_interceptor.NewMockIdempotencyStore(ctrl)
		store.EXPECT().
			Get(gomock.Any(), "/test.Service/Method:test-key").
			Return(cachedData, nil)

		interceptorFunc := interceptor.IdempotencyUnaryServerInterceptor(store)
		handler := func(_ context.Context, _ any) (any, error) {
			t.Fatal("Handler should not be called")
			return nil, nil
		}

		md := metadata.New(map[string]string{
			interceptor.IdempotencyKeyHeader: "test-key",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptorFunc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := grpcstatus.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Equal(t, "test error", st.Message())
	})

	t.Run("redis get error continues processing", func(t *testing.T) {
		store := mock_interceptor.NewMockIdempotencyStore(ctrl)
		store.EXPECT().
			Get(gomock.Any(), "/test.Service/Method:test-key").
			Return(nil, errors.New("redis error"))
		store.EXPECT().
			Set(gomock.Any(), "/test.Service/Method:test-key", gomock.Any(), interceptor.DefaultIdempotencyTTL).
			Return(nil)

		interceptorFunc := interceptor.IdempotencyUnaryServerInterceptor(store)
		handler := func(_ context.Context, _ any) (any, error) {
			return &emptypb.Empty{}, nil
		}

		md := metadata.New(map[string]string{
			interceptor.IdempotencyKeyHeader: "test-key",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptorFunc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("redis set error does not fail request", func(t *testing.T) {
		store := mock_interceptor.NewMockIdempotencyStore(ctrl)
		store.EXPECT().
			Get(gomock.Any(), "/test.Service/Method:test-key").
			Return(nil, nil)
		store.EXPECT().
			Set(gomock.Any(), "/test.Service/Method:test-key", gomock.Any(), interceptor.DefaultIdempotencyTTL).
			Return(errors.New("redis error"))

		interceptorFunc := interceptor.IdempotencyUnaryServerInterceptor(store)
		handler := func(_ context.Context, _ any) (any, error) {
			return &emptypb.Empty{}, nil
		}

		md := metadata.New(map[string]string{
			interceptor.IdempotencyKeyHeader: "test-key",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptorFunc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("handler error is cached", func(t *testing.T) {
		store := mock_interceptor.NewMockIdempotencyStore(ctrl)
		store.EXPECT().
			Get(gomock.Any(), "/test.Service/Method:test-key").
			Return(nil, nil)
		store.EXPECT().
			Set(gomock.Any(), "/test.Service/Method:test-key", gomock.Any(), interceptor.DefaultIdempotencyTTL).
			Return(nil)

		interceptorFunc := interceptor.IdempotencyUnaryServerInterceptor(store)
		handlerErr := grpcstatus.Error(codes.NotFound, "not found")
		handler := func(_ context.Context, _ any) (any, error) {
			return nil, handlerErr
		}

		md := metadata.New(map[string]string{
			interceptor.IdempotencyKeyHeader: "test-key",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptorFunc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func createCachedResponse(t *testing.T, resp any, respErr error) []byte {
	t.Helper()

	type cacheableError struct {
		StatusProto []byte `json:"status_proto,omitempty"`
	}

	type serializableResponse struct {
		Error    *cacheableError `json:"error,omitempty"`
		TypeURL  string          `json:"type_url,omitempty"`
		Response []byte          `json:"response,omitempty"`
	}

	sr := &serializableResponse{}

	if respErr != nil {
		// Marshal the error status proto
		st, ok := grpcstatus.FromError(respErr)
		if ok {
			statusProto := st.Proto()
			data, err := proto.Marshal(statusProto)
			if err != nil {
				t.Fatal(err)
			}
			sr.Error = &cacheableError{
				StatusProto: data,
			}
		}
	} else if resp != nil {
		protoMsg, ok := resp.(proto.Message)
		if !ok {
			t.Fatal("response must be a proto.Message")
		}

		// Use anypb to wrap the message
		anyMsg, err := anypb.New(protoMsg)
		if err != nil {
			t.Fatal(err)
		}
		data, err := proto.Marshal(anyMsg)
		if err != nil {
			t.Fatal(err)
		}
		sr.Response = data
		sr.TypeURL = anyMsg.TypeUrl
	}

	data, err := json.Marshal(sr)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
