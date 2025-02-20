#!/usr/bin/env bash

set -euo pipefail

for dir in `find . -type d`; do
    if [[ -f ${dir}/go.mod ]]; then
        for file in `find ${dir} -name '*.go' | grep -v proto | grep -v test/mock | grep -v internal/repository/db`; do
            if `grep -q 'interface {' ${file}`; then
                dest=${file//internal\//}
                dest=${dest#"$dir"}

                rm -rf ${dir}/test/mock/${dest}
                mockgen -source=${file} -destination=${dir}/test/mock/${dest}
            fi
        done
    fi
done
