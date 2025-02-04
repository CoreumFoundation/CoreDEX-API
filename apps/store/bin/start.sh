#!/bin/sh

running=`ps -ef|grep exe/store |grep -v grep`

if [ "$running" ]; then
  echo "store is running"
  exit 0
else
  echo "store is starting"
fi

export MYSQL_CONFIG='{"Username": "testuser","Password": "password","Host": "localhost","Port": 3306,"Database": "friendly_dex"}'
export LOG_LEVEL="info"
export GRPC_PORT=":50051"

# Takes care of any not downloaded/updated dependencies
go mod tidy

go run .
