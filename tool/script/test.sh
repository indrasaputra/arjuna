#!/usr/bin/env bash

set -eo pipefail

run_tests() {
    go clean -testcache
    go test -count=1 -failfast -v -race ./...
}

run_tests_cover() {
    go clean -testcache
    go test -count=1 -failfast -v -race -coverprofile=coverage.out $(go list ./... | grep -v /test/)
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func coverage.out
}

if [[ $1 = 'cover' ]]; then
    if [ $2 ]; then
        (cd ./$2 && run_tests_cover) # from https://github.com/indrasaputra/arjuna/pull/55#pullrequestreview-2328123308
    else
        find . -type d | while IFS= read -r dir; do # from https://github.com/indrasaputra/arjuna/pull/55#pullrequestreview-2328123308
            if [[ -f ${dir}/go.mod ]]; then
                (cd ${dir} && run_tests_cover) # from https://github.com/indrasaputra/arjuna/pull/55#pullrequestreview-2328123308
            fi
        done
    fi
elif [[ $1 = 'unit' ]]; then
    if [ $2 ]; then
        (cd ./$2 && run_tests) # from https://github.com/indrasaputra/arjuna/pull/55#pullrequestreview-2328123308
    else
        find . -type d | while IFS= read -r dir; do # from https://github.com/indrasaputra/arjuna/pull/55#pullrequestreview-2328123308
            if [[ -f ${dir}/go.mod ]]; then
                (cd ${dir} && run_tests) # from https://github.com/indrasaputra/arjuna/pull/55#pullrequestreview-2328123308
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
