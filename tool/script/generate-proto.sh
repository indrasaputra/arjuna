#!/usr/bin/env bash

set -euo pipefail

# Download buf dependencies (googleapis, grpc-gateway, etc.)
buf dep update

# buf breaking --against '.git#branch=main'
buf lint
buf build
buf generate
