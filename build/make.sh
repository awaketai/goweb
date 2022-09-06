#!/bin/bash

folder=${1}
serviceName=${2}
appFile=${serviceName}

# exit immediately if a command exists with a non-zero status
set -e
go mod tidy
# build
mkdir -p ./bin/
go build -o ./bin/${appFile} ../main
cp ./bin/${appFile} ..
