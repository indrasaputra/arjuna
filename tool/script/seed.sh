#!/usr/bin/env bash

set -eo pipefail

CMD=server
APP_DIR=service/$1

if [ -d ${APP_DIR}/cmd/${CMD} ]; then
    echo "running seeder for ${APP_DIR}..."
    (cd ${APP_DIR} && \
        go run cmd/${CMD}/main.go seed)
fi
