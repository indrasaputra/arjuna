#!/usr/bin/env bash

set -eo pipefail

# ignore LT05 because it will be handled by sqlc, AM04 because it's not a problem
sqlfluff fix -d postgres -e LT05,AM04
buf lint
protolint lint -fix .

config=$(pwd)/.golangci.yml

if [[ -n $1 ]]; then
    echo "running in $1"
    (cd $1 && golangci-lint run --config=${config} ./...)
else
    for dir in `find . -type d`; do
        if [[ -f ${dir}/go.mod ]]; then
            echo "running in ${dir}"
            (cd ${dir} && golangci-lint run --config=${config} ./...)
        fi
    done
fi
