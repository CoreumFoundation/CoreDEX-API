#!/bin/sh


running=`ps -ef|grep exe/api-server |grep -v grep`

if [ "$running" ]; then
  echo "api-server is running"
  exit 0
else
  echo "api-server is starting"
fi

cd ../store
./bin/start.sh &
cd ../api-server 

export NETWORKS="{\"Node\": [{\"Network\": \"devnet\",\"GRPCHost\": \"full-node.devnet-1.coreum.dev:9090\",\"RPCHost\": \"https://full-node.devnet-1.coreum.dev:26657\"}]}"
export STATE_STORE="localhost:50051"
export TRADE_STORE="localhost:50051"
export OHLC_STORE="localhost:50051"
export ORDER_STORE="localhost:50051"
export CURRENCY_STORE="localhost:50051"
export LOG_LEVEL="info"
export HTTP_CONFIG="{\"port\": \":8080\",\"cors\": {\"allowedOrigins\":[\"http://localhost:3000\",\"http://localhost:3001\"]},\"timeouts\": {\"read\": \"10s\",\"write\": \"10s\",\"idle\": \"10s\",\"shutdown\": \"10s\"}}"
export BASE_COIN="{\"BaseCoin\":[{\"Network\": \"mainnet\",\"Coin\": \"ucore\"},{\"Network\": \"testnet\",\"Coin\": \"utestcore\"},{\"Network\": \"devnet\",\"Coin\": \"udevcore\"}]}"
export BASE_USDC="{\"BaseCoin\":[{\"Network\": \"mainnet\",\"Coin\": \"uusdc-E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D\"}]}"

# Takes care of any not downloaded/updated dependencies
go mod tidy

go run .
