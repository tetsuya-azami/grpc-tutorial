#!/bin/bash
DIRNAME=$(cd "$(dirname ${BASH_SOURCE[0]})" && pwd)

if [ ! -d "${DIRNAME}/../../pkg/grpc/hello" ]; then
	mkdir -p ${DIRNAME}/../../pkg/grpc/hello
fi

if [ -f "${DIRNAME}/../../pkg/grpc/hello/*.go" ]; then
	rm ${DIRNAME}/../../pkg/grpc/hello/*.go
fi

protoc --proto_path=${DIRNAME} --go_out=${DIRNAME}/../../pkg/grpc/hello --go_opt=paths=source_relative \
    --go-grpc_out=${DIRNAME}/../../pkg/grpc/hello --go-grpc_opt=paths=source_relative \
	${DIRNAME}/hello.proto
