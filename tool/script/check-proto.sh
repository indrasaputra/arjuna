#!/usr/bin/env bash

set -euo pipefail

if [ ! -z "`git status -s | grep -e '.proto'`" ]; then
  echo "Proto files are not beautifully formatted for these files:"
  git status -s | grep -e '.proto'
  echo "Run 'make format' or 'make pretty' before commit and push"
  exit 1
fi
