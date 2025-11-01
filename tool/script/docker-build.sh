#!/usr/bin/env bash

set -eo pipefail

CMD_DIRS=(server worker relayer)
DOCKERFILE=$1/dockerfile/$1.dockerfile

set_app_dir() {
	if [[ $1 = "gateway" ]]; then
		APP_DIR=$1
	else
		APP_DIR=service/$1
	fi
}

set_dockerfile() {
	if [[ ! $1 = "gateway" ]] && [[ ! $1 = "apidoc" ]]; then
		DOCKERFILE=service/${DOCKERFILE}
	fi
}

build_apidoc() {
	echo "building apidoc server..."
	docker build --no-cache -t indrasaputra/arjuna-$1-server:latest -f ${DOCKERFILE} .
}

build_non_apidoc() {
	for cmd in "${CMD_DIRS[@]}"; do
		if [ -d ${APP_DIR}/cmd/${cmd} ]; then
			echo "building $1 ${cmd}..."
			docker build --no-cache --build-arg CMD=${cmd} -t indrasaputra/arjuna-$1-${cmd}:latest -f ${DOCKERFILE} .
		fi
	done
}

build_docker() {
	if [[ $1 == "apidoc" ]]; then
		build_apidoc $1
	else
		build_non_apidoc $1
	fi
}

set_app_dir $1
set_dockerfile $1
build_docker $1
