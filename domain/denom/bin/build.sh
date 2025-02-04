#!/bin/bash
set -e
 
protoc \
  --proto_path=. "denom.proto" \
  --proto_path=$(dirname $(dirname "$rd")) \
  "--go_out=." --go_opt=paths=source_relative