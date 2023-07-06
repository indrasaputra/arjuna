#!/usr/bin/env bash

set -eo pipefail

CMD_DIRS=(server worker relayer)
OUTPUT_DIR=deploy/output
APP_DIR=

if [[ $1 = "gateway" ]]; then
	APP_DIR=$1
else
    APP_DIR=service/$1
fi

for cmd in "${CMD_DIRS[@]}"; do
    if [ -d ${APP_DIR}/cmd/${cmd} ]; then
        echo "compiling ${APP_DIR}/cmd/${cmd}..."
        (cd ${APP_DIR} && \
            GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
            go build -a -installsuffix cgo -ldflags '-extldflags "-static"' \
            -o ${OUTPUT_DIR}/${cmd}/$1 \
            cmd/${cmd}/main.go)
    fi
done
