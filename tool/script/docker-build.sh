#!/usr/bin/env bash

set -eo pipefail

CMD_DIRS=(server worker)
DOCKERFILE=$1/dockerfile/$1.dockerfile

set_app_dir() {
	if [[ $1 = "gateway" ]]; then
		APP_DIR=$1
	else
		APP_DIR=service/$1
	fi
}

set_dockerfile() {
	if [[ ! $1 = "gateway" ]] && [[ ! $1 = "blueprint" ]]; then
		DOCKERFILE=service/${DOCKERFILE}
	fi
}

build_blueprint() {
	echo "building blueprint server..."
	docker build --no-cache -t indrasaputra/arjuna-$1-server:latest -f ${DOCKERFILE} .
}

build_non_blueprint() {
	for cmd in "${CMD_DIRS[@]}"; do
		if [ -d ${APP_DIR}/cmd/${cmd} ]; then
			echo "building $1 ${cmd}..."
			docker build --no-cache --build-arg CMD=${cmd} -t indrasaputra/arjuna-$1-${cmd}:latest -f ${DOCKERFILE} .
		fi
	done
}

build_docker() {
	if [[ $1 == "blueprint" ]]; then
		build_blueprint $1
	else
		build_non_blueprint $1
	fi
}

set_app_dir $1
set_dockerfile $1
build_docker $1
