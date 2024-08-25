#!/usr/bin/env bash

set -eo pipefail

CMD=server
OUTPUT_DIR=deploy/output
APP_DIR=

if [[ $1 = "gateway" ]]; then
	APP_DIR=$1
else
    APP_DIR=service/$1
fi

if [ -d ${APP_DIR}/cmd/${CMD} ]; then
    echo "compiling ${APP_DIR}/cmd/${CMD}..."
    (cd ${APP_DIR} && \
        GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
        go build -a -installsuffix cgo -ldflags '-extldflags "-static"' \
        -o ${OUTPUT_DIR}/${CMD}/$1 \
        cmd/${CMD}/main.go)
fi
