#!/bin/bash
set -e

# move to the root dir of the package
rd=$(git rev-parse --show-toplevel)
cd $rd

protoc \
  --proto_path=. "domain/update/update.proto" \
  "--go_out=." --go_opt=paths=source_relative
  
# protoc \
#   --proto_path=. "domain/update/update-grpc.proto" \
#   "--go_out=." --go_opt=paths=source_relative \
#   --go-grpc_opt=require_unimplemented_servers=false \
#   "--go-grpc_out=." --go-grpc_opt=paths=source_relative

cp domain/update/package.json .
cp domain/update/tsconfig.json .

rm -rf node_modules
npm i

protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto \
--proto_path=. "domain/update/update.proto" \
--ts_proto_out=. \
--ts_proto_opt=esModuleInterop=true \
--ts_proto_opt=outputServices=grpc-js \
domain/update/update.proto

npm run build-ts
git add build/

# git add *.ts
# rm -rf node_modules
# rm package-lock.json
# rm package.json
# rm tsconfig.json
