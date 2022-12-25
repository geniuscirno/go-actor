#!/usr/bin/env bash

set -e

if ! [[ "$0" =~ scripts/genproto.sh ]]; then
    echo "must be run from repository root"
    exit 255
fi

protoc -I . --go_out=paths=source_relative:. \
 actor/*.proto

protoc -I . --go_out=paths=source_relative:. \
 --go-grpc_out=paths=source_relative:. \
 remote/*.proto
