#!/usr/bin/env bash

set -eo pipefail

# In CI, limit parallelism to reduce memory pressure from race detector
if [ "${CI}" = "true" ] || [ "${GITHUB_ACTIONS}" = "true" ]; then
    # -p=1: run packages sequentially, -parallel=1: run tests within package sequentially
    PARALLEL_FLAG="-p=1 -parallel=1"
    # Enable CGO for race detector and configure TSAN for minimal memory usage
    export CGO_ENABLED=1
    # Ultra-aggressive TSAN options to minimize memory usage
    export TSAN_OPTIONS="halt_on_error=1:history_size=1:io_sync=0:flush_memory_ms=0:die_after_fork=0"
    export GORACE="halt_on_error=1"
    # Aggressive garbage collection to free memory quickly
    export GOGC=20
    # Reduce Go's memory limit
    export GOMEMLIMIT=2GiB
else
    PARALLEL_FLAG=""
fi

run_tests() {
    go clean -testcache
    go test -count=1 -failfast -v -race ${PARALLEL_FLAG} ./...
}

run_tests_cover() {
    go clean -testcache
    go test -count=1 -failfast -v -race ${PARALLEL_FLAG} -coverprofile=coverage.out $(go list ./... | grep -v /test/)
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
