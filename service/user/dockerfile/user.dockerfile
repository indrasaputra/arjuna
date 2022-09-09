FROM golang:1.17 AS builder
WORKDIR /app
COPY . .
RUN make compile svc=user
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.5 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
RUN WAIT_FOR_VERSION=v2.1.2 && \
    wget -qO/bin/wait-for https://github.com/eficode/wait-for/releases/download/${WAIT_FOR_VERSION}/wait-for && \
    chmod +x /bin/wait-for

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /bin/grpc_health_probe ./grpc_health_probe
COPY --from=builder /bin/wait-for ./wait-for
COPY --from=builder /app/service/user/deploy/output/user .
COPY --from=builder /app/tool/script/start.sh ./start.sh
RUN chmod +x /app/start.sh /app/wait-for /app/user
EXPOSE 8001
CMD ["./start.sh", "user"]
