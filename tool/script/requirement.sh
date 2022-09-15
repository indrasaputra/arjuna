#!/usr/bin/env bash

set -eo pipefail

req_dir="requirements"
timestamp=`date "+%Y%m%d%H%M%S"`

mkdir -p ${req_dir}
touch "${req_dir}/${timestamp}_${1}.md"
