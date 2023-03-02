module github.com/indrasaputra/arjuna/gateway

go 1.18

replace (
	github.com/indrasaputra/arjuna/pkg/sdk v0.0.0 => ../pkg/sdk
	github.com/indrasaputra/arjuna/proto v0.0.0 => ../proto
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.8.0
	github.com/prometheus/client_golang v1.14.0
)

require (
	github.com/joeshaw/envdecode v0.0.0-20200121155833-099f1fc765bd
	github.com/joho/godotenv v1.4.0
	github.com/pkg/errors v0.9.1
)

require google.golang.org/grpc v1.52.3

require github.com/indrasaputra/arjuna/proto v0.0.0

require google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6

require github.com/stretchr/testify v1.8.1

require github.com/indrasaputra/arjuna/pkg/sdk v0.0.0

require (
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.38.0
	go.opentelemetry.io/otel v1.13.0
)

require go.uber.org/zap v1.13.0

require github.com/grpc-ecosystem/go-grpc-middleware v1.3.0

require github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.13.0 // indirect
	go.opentelemetry.io/otel/metric v0.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.13.0 // indirect
	go.opentelemetry.io/otel/trace v1.13.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
