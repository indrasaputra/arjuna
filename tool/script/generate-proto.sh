#!/usr/bin/env bash

set -euo pipefail

# buf breaking --against '.git#branch=main'
buf lint
buf build
buf generate
