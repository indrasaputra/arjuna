module github.com/indrasaputra/arjuna/service/user

go 1.17

replace github.com/indrasaputra/arjuna/proto v0.0.0 => ../../proto

require (
	github.com/indrasaputra/arjuna/proto v0.0.0
	github.com/stretchr/testify v1.7.1
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/mock v1.4.4
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.8.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
