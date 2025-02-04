#!/bin/bash

running=`ps -ef|grep exe/data-aggregator |grep -v grep`

if [ "$running" ]; then
  echo "data-aggregator is running"
  exit 0
else
  echo "data-aggregator is starting"
fi

export NETWORKS='{
    "Node": [
        {
            "Network": "devnet",
            "GRPCHost": "full-node.devnet-1.coreum.dev:9090",
            "RPCHost": "https://full-node.devnet-1.coreum.dev:26657"
        }
    ]
}'
export STATE_STORE=localhost:50051
export TRADE_STORE=localhost:50051
export OHLC_STORE=localhost:50051
export ORDER_STORE=localhost:50051
export CURRENCY_STORE=localhost:50051
export LOG_LEVEL=info

# Start stores:
cd ../store
./bin/start.sh &
sleep 5
cd ../data-aggregator

# Takes care of any not downloaded/updated dependencies
go mod tidy

go run . > log.txt 2>&1 
