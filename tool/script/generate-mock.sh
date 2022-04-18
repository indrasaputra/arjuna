#!/usr/bin/env bash

set -euo pipefail

for dir in `find . -type d`; do
    if [[ -f ${dir}/go.mod ]]; then
        for file in `find ${dir} -name '*.go' | grep -v proto | grep -v test/mock`; do
            if `grep -q 'interface {' ${file}`; then
                dest=${file//internal\//}
                dest=${dest#"$dir"}
                mockgen -source=${file} -destination=${dir}/test/mock/${dest}
            fi
        done
    fi
done
