#!/usr/bin/env bash

set -euo pipefail

for dir in `find . -type d`; do
    if [[ -f ${dir}/go.mod && -f ${dir}/sqlc.yaml ]]; then
        (cd ${dir} && sqlc generate)
    fi
done
