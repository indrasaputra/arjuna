module github.com/indrasaputra/arjuna/gateway

go 1.22.4

replace (
	github.com/indrasaputra/arjuna/pkg/sdk v0.0.0 => ../pkg/sdk
	github.com/indrasaputra/arjuna/proto v0.0.0 => ../proto
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0
	github.com/prometheus/client_golang v1.19.1
)

require (
	github.com/joeshaw/envdecode v0.0.0-20200121155833-099f1fc765bd
	github.com/joho/godotenv v1.4.0
	github.com/pkg/errors v0.9.1
)

require google.golang.org/grpc v1.65.0

require github.com/indrasaputra/arjuna/proto v0.0.0

require google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6

require github.com/stretchr/testify v1.9.0

require github.com/indrasaputra/arjuna/pkg/sdk v0.0.0

require (
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.53.0
	go.opentelemetry.io/otel v1.28.0
)

require go.uber.org/zap v1.27.0

require github.com/grpc-ecosystem/go-grpc-middleware v1.4.0

require github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
