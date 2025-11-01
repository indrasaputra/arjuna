FROM golang:1.25 AS builder
ARG SERVICE=gateway
ARG OUTPUT_DIR=deploy/output
ARG CMD=server
WORKDIR /app
COPY . .
RUN if [ ! -f ${SERVICE}/${OUTPUT_DIR}/${CMD}/${SERVICE} ] ; then make compile svc=gateway ; fi

FROM alpine:3.16
ARG SERVICE=gateway
ARG OUTPUT_DIR=deploy/output
ARG CMD=server
WORKDIR /app
COPY --from=builder /app/${SERVICE}/${OUTPUT_DIR}/${CMD}/${SERVICE} .
RUN apk add --update curl && \
    rm -rf /var/cache/apk/*
RUN chmod +x /app/${SERVICE}
EXPOSE 8000
CMD ["./gateway"]
