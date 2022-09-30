module github.com/indrasaputra/arjuna/gateway

go 1.18

replace github.com/indrasaputra/arjuna/proto v0.0.0 => ../proto

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.8.0
	github.com/prometheus/client_golang v1.12.1
)

require (
	github.com/joeshaw/envdecode v0.0.0-20200121155833-099f1fc765bd
	github.com/joho/godotenv v1.4.0
	github.com/pkg/errors v0.9.1
)

require google.golang.org/grpc v1.45.0

require github.com/indrasaputra/arjuna/proto v0.0.0

require google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106

require github.com/stretchr/testify v1.7.0

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
