#!/bin/bash
set -e

# move to the root dir of the package
rd=$(git rev-parse --show-toplevel)
cd $rd

protoc \
  --proto_path=. "domain/order-properties/order-properties.proto" \
  "--go_out=." --go_opt=paths=source_relative

cp domain/order-properties/package.json .
cp domain/order-properties/tsconfig.json .

rm -rf node_modules
npm i

protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto \
--proto_path=. "domain/order-properties/order-properties.proto" \
--ts_proto_out=. \
--ts_proto_opt=esModuleInterop=true \
--ts_proto_opt=outputServices=grpc-js \
domain/order-properties/order-properties.proto

npm run build-ts
git add build/

git add *.ts
rm -rf node_modules
rm package-lock.json
# rm package.json
# rm tsconfig.json
