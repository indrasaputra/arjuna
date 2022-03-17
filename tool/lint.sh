#!/usr/bin/env bash

set -euo pipefail

buf lint
protolint lint -fix .

config=$(pwd)/.golangci.yml
for dir in `find . -type d`; do
  if [[ -f ${dir}/go.mod ]]; then
    (cd ${dir} && golangci-lint run --config ${config} ./...)
  fi
done
