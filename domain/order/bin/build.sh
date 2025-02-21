#!/bin/bash
set -e

# move to the root dir of the package
rd=$(git rev-parse --show-toplevel)
cd $rd

protoc \
  --proto_path=. "domain/order/order.proto" \
  "--go_out=." --go_opt=paths=source_relative
  
protoc \
  --proto_path=. "domain/order/order-grpc.proto" \
  "--go_out=." --go_opt=paths=source_relative \
  --go-grpc_opt=require_unimplemented_servers=false \
  "--go-grpc_out=." --go-grpc_opt=paths=source_relative

cp domain/order/package.json .
cp domain/order/tsconfig.json .

rm -rf node_modules
npm i

protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto \
--proto_path=. "domain/order/order.proto" \
--ts_proto_out=. \
--ts_proto_opt=esModuleInterop=true \
--ts_proto_opt=outputServices=grpc-js \
domain/order/order.proto

npm run build-ts
git add build/

git add *.ts
rm -rf node_modules
rm package-lock.json
rm package.json
rm tsconfig.json
