#!/usr/bin/env bash

set -euo pipefail

for dir in `find . -type d`; do
  if [[ -f ${dir}/go.mod ]]; then
    (cd ${dir} && go mod tidy)
  fi
done

go work sync
