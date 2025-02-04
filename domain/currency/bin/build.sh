#!/bin/bash
set -e

# move to the root dir of the package
rd=$(git rev-parse --show-toplevel)
cd $rd

protoc \
  --proto_path=. "domain/currency/currency.proto" \
  "--go_out=." --go_opt=paths=source_relative

protoc \
  --proto_path=. "domain/currency/currency-grpc.proto" \
  "--go_out=." --go_opt=paths=source_relative \
  --go-grpc_opt=require_unimplemented_servers=false \
  "--go-grpc_out=." --go-grpc_opt=paths=source_relative
