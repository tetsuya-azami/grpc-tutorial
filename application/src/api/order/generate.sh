#!/bin/bash
DIRNAME=$(cd "$(dirname ${BASH_SOURCE[0]})" && pwd)

if [ ! -d "${DIRNAME}/../../pkg/grpc/order" ]; then
	mkdir -p ${DIRNAME}/../../pkg/grpc/order
fi

if [ -f "${DIRNAME}/../../pkg/grpc/order/*.go" ]; then
	rm ${DIRNAME}/../../pkg/grpc/order/*.go
fi

protoc --proto_path=${DIRNAME} --go_out=${DIRNAME}/../../pkg/grpc/order --go_opt=paths=source_relative \
    --go-grpc_out=${DIRNAME}/../../pkg/grpc/order --go-grpc_opt=paths=source_relative \
	${DIRNAME}/order.proto
