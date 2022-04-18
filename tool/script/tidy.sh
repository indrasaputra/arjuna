#!/usr/bin/env bash

set -euo pipefail

for dir in `find . -type d`; do
  if [[ -f ${dir}/go.mod ]]; then
    (cd ${dir} && go mod tidy -go=1.16 && go mod tidy -go=1.17)
  fi
done
