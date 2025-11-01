#!/usr/bin/env bash

set -eo pipefail

CMD=server
APP_DIR=service/$1

if [ -d ${APP_DIR}/cmd/${CMD} ]; then
    echo "running atlas migration for ${APP_DIR}..."
    (cd ${APP_DIR} && \
        atlas migrate diff --dir file://db/migrations --dev-url docker://postgres/16 --to file://schema.sql && \
        sqlfluff fix -d postgres -e LT05,AM04 && \
        atlas migrate hash --dir file://db/migrations)
fi
