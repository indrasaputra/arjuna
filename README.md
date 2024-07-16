# Arjuna

Arjuna is a place or playground for me to write my state-of-the-art code.
I use monorepo purely because I want to implement microservices without managing multiple repositories.

## Owner

[Indra Saputra](https://github.com/indrasaputra)

## API

### gRPC

The API can be seen in proto files (`*.proto`) in each service api directory. It is always located in `service/<service-name>/api/` directory.

### HTTP

The API can be seen via web using [Stoplight Elements](https://github.com/stoplightio/elements).

There are two ways to see the API documentation.

1. Run the whole project and visit [http://localhost:4000](http://localhost:4000).
2. Copy-paste [arjuna.swagger.json](/openapiv2/arjuna.swagger.json) to [https://editor.swagger.io](https://editor.swagger.io)

## How to Run

Read [How to Run](docs/HOW_TO_RUN.md).

## Observability

Read [Observability](docs/OBSERVABILITY.md).

## Code Map

Read [Code Map](docs/CODE_MAP.md)
