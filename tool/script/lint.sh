#!/usr/bin/env bash

set -eo pipefail

buf lint
protolint lint -fix .

config=$(pwd)/.golangci.yml

if [[ -n $1 ]]; then
    (cd $1 && golangci-lint run --config=${config} ./...)
else
    for dir in `find . -type d`; do
        if [[ -f ${dir}/go.mod ]]; then
            (cd ${dir} && golangci-lint run --config=${config} ./...)
        fi
    done
fi
