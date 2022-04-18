#!/usr/bin/env bash

set -euo pipefail

if [[ $1 = 'cover' ]]; then
    for dir in `find . -type d`; do
    if [[ -f ${dir}/go.mod ]]; then
        (cd ${dir} && 
            go clean -testcache &&
            go test -count=1 -failfast -v -race -coverprofile=coverage.out ./... &&
            go tool cover -html=coverage.out -o coverage.html &&
            go tool cover -func coverage.out)
    fi
    done
elif [[ $1 = 'test' ]]; then
    for dir in `find . -type d`; do
    if [[ -f ${dir}/go.mod ]]; then
        (cd ${dir} && 
            go clean -testcache &&
            go test -count=1 -failfast -v -race ./...)
    fi
    done
fi
