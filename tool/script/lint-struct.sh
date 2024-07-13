#!/usr/bin/env bash

set -eo pipefail

if [[ -n $1 ]]; then
    echo "running in $1"
    (cd $1 && fieldalignment -fix ./...)
else
    for dir in `find . -type d`; do
        if [[ -f ${dir}/go.mod ]]; then
            echo "running in ${dir}"
            (cd ${dir} && fieldalignment -fix ./...)
        fi
    done
fi
