FROM golang:1.22 AS builder
ARG SERVICE=user
ARG OUTPUT_DIR=deploy/output
ARG CMD=server
WORKDIR /app
COPY . .
RUN if [ ! -f service/${SERVICE}/${OUTPUT_DIR}/${CMD}/${SERVICE} ] ; then make compile svc=${SERVICE} ; fi
RUN chmod +x /app/bin/grpc_health_probe-linux-amd64-v0.4.28 /app/service/${SERVICE}/${OUTPUT_DIR}/${CMD}/${SERVICE}

FROM alpine:3.16
ARG SERVICE=user
ARG OUTPUT_DIR=deploy/output
ARG CMD=server
WORKDIR /app
COPY --from=builder /app/bin/grpc_health_probe-linux-amd64-v0.4.28 ./grpc_health_probe
COPY --from=builder /app/service/${SERVICE}/${OUTPUT_DIR}/${CMD}/${SERVICE} .
EXPOSE 8001
