#!/usr/bin/env bash
set -euo pipefail

make format
make check.import
