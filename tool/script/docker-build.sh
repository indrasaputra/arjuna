#!/usr/bin/env bash

set -eo pipefail

DOCKERFILE=$1/dockerfile/$1.dockerfile

if [[ ! $1 = "gateway" ]] && [[ ! $1 = "blueprint" ]]; then
	DOCKERFILE=service/${DOCKERFILE}
fi

docker build --no-cache -t indrasaputra/arjuna-$1:latest -f ${DOCKERFILE} .
