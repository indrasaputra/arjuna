module github.com/indrasaputra/arjuna/service/transaction

go 1.22.4

replace (
	github.com/indrasaputra/arjuna/pkg/sdk v0.0.0 => ../../pkg/sdk
	github.com/indrasaputra/arjuna/proto v0.0.0 => ../../proto
	github.com/indrasaputra/arjuna/service/auth v0.0.0 => ../../service/auth
)

require (
	github.com/indrasaputra/arjuna/pkg/sdk v0.0.0
	github.com/indrasaputra/arjuna/proto v0.0.0
	github.com/segmentio/ksuid v1.0.4
	github.com/shopspring/decimal v1.4.0
	github.com/stretchr/testify v1.9.0
	go.uber.org/mock v0.4.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d
	google.golang.org/grpc v1.65.0
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
