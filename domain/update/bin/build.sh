#!/bin/bash
set -e

# move to the root dir of the package
rd=$(git rev-parse --show-toplevel)
cd $rd

protoc \
  --proto_path=. "domain/update/update.proto" \
  "--go_out=." --go_opt=paths=source_relative
