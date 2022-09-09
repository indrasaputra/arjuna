#!/bin/sh

set -euo pipefail

echo "wait for postgres @" $POSTGRES_HOST:$POSTGRES_PORT
./wait-for -t 20 $POSTGRES_HOST:$POSTGRES_PORT

./$1
