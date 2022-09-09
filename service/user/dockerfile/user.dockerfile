FROM golang:1.17 AS builder
ARG SERVICE=user
ARG OUTPUT_DIR=deploy/output
WORKDIR /app
COPY . .
RUN if [ ! -f service/${SERVICE}/${OUTPUT_DIR}/${SERVICE} ] ; then make compile svc=${SERVICE} ; fi
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.5 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
RUN WAIT_FOR_VERSION=v2.1.2 && \
    wget -qO/bin/wait-for https://github.com/eficode/wait-for/releases/download/${WAIT_FOR_VERSION}/wait-for && \
    chmod +x /bin/wait-for

FROM alpine:3.13
ARG SERVICE=user
ARG OUTPUT_DIR=deploy/output
WORKDIR /app
COPY --from=builder /bin/grpc_health_probe ./grpc_health_probe
COPY --from=builder /bin/wait-for ./wait-for
COPY --from=builder /app/service/${SERVICE}/${OUTPUT_DIR}/${SERVICE} .
COPY --from=builder /app/tool/script/start.sh ./start.sh
RUN chmod +x /app/start.sh /app/wait-for /app/${SERVICE}
EXPOSE 8001
CMD ["./start.sh", ${SERVICE}]
