#!/usr/bin/env bash

set -euo pipefail

for file in `find . -name '*.go'`; do
  # Defensive, just in case.
  if [[ -f ${file} ]]; then
    awk '/^import \($/,/^\)$/{if($0=="")next}{print}' ${file} > /tmp/file
    mv /tmp/file ${file}
  fi
done

for dir in `find . -type d`; do
  if [[ -f ${dir}/go.mod ]]; then
    (cd ${dir} &&
      goimports -w -local github.com/indrasaputra/arjuna $(go list -f '{{.Dir}}' -tags integration ./...) &&
      gofmt -s -w .)
  fi
done

buf format -w
