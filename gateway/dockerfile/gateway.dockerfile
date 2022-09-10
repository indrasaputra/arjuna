FROM golang:1.17 AS builder
ARG SERVICE=gateway
ARG OUTPUT_DIR=deploy/output
WORKDIR /app
COPY . .
RUN if [ ! -f {SERVICE}/${OUTPUT_DIR}/${SERVICE} ] ; then make compile svc=gateway ; fi

FROM alpine:3.13
ARG SERVICE=gateway
ARG OUTPUT_DIR=deploy/output
WORKDIR /app
COPY --from=builder /app/${SERVICE}/${OUTPUT_DIR}/${SERVICE} .
RUN apk add --update curl && \
    rm -rf /var/cache/apk/*
RUN chmod +x /app/${SERVICE}
EXPOSE 8000
CMD ["./gateway"]
