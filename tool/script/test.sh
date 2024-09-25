#!/usr/bin/env bash

set -eo pipefail

if [[ $1 = 'cover' ]]; then
    if [ $2 ]; then
        (cd ./$2 && 
            go clean -testcache &&
            go test -count=1 -failfast -v -race -coverprofile=coverage.out $(go list ./... | grep -v /test/) &&
            go tool cover -html=coverage.out -o coverage.html &&
            go tool cover -func coverage.out)
    else
        for dir in `find . -type d`; do
            if [[ -f ${dir}/go.mod ]]; then
                (cd ${dir} && 
                    go clean -testcache &&
                    go test -count=1 -failfast -v -race -coverprofile=coverage.out $(go list ./... | grep -v /test/) &&
                    go tool cover -html=coverage.out -o coverage.html &&
                    go tool cover -func coverage.out)
            fi
        done
    fi
elif [[ $1 = 'unit' ]]; then
    if [ $2 ]; then
        (cd ./$2 && 
            go clean -testcache &&
            go test -count=1 -failfast -v -race ./...)
    else
        for dir in `find . -type d`; do
            if [[ -f ${dir}/go.mod ]]; then
                (cd ${dir} && 
                    go clean -testcache &&
                    go test -count=1 -failfast -v -race ./...)
            fi
        done
    fi
elif [[ $1 = 'e2e' ]]; then
    for dir in `find . -type d | grep service`; do
        if [[ -f ${dir}/go.mod ]] && [ -d ${dir}/test/integration ]; then
            (cd ${dir} && 
                go clean -testcache &&
                go test --tags integration -v ./test/integration)
        fi
    done
fi
