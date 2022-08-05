#!/usr/bin/env bash

set -euo pipefail

version=1.18

for dir in `find . -type d`; do
  if [[ -f ${dir}/go.mod ]]; then
    (cd ${dir} && go mod edit -go=${version} && go mod tidy -go=1.16 && go mod tidy -go=${version})
  fi
done
