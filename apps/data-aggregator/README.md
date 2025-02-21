# Data aggregator

The data aggregator scans blocks near real time as they are produced by the blockchain.

The scan results are per transaction processed by different associated handlers.

## Start parameters

The data aggregator can be started with the following parameters:

- `NETWORKS` : Network config JSON similar to `"{\"Node\": [{\"Network\": \"devnet\",\"GRPCHost\": \"full-node.devnet-1.coreum.dev:9090\",\"RPCHost\": \"https://full-node.devnet-1.coreum.dev:26657\"}]}"`,`
- `STATE_STORE` - Store connection host:port format
- `TRADE_STORE` - Store connection host:port format
- `OHLC_STORE` - Store connection host:port format
- `ORDER_STORE` - Store connection host:port format
- `CURRENCY_STORE` - Store connection host:port format
- `LOG_LEVEL` - Optional

### NETWORKS

The NETWORKS env is a JSON that can hold one or more Coreum nodes, 1 per network.
In most user facing applications the users will only use `mainnet` however for development and testing purposes you can add `devnet` and `testnet` nodes.

Sample:

```json
{
    "Node": [
        {
            "Network": "devnet",
            "GRPCHost": "full-node.devnet-1.coreum.dev:9090",
            "RPCHost": "https://full-node.devnet-1.coreum.dev:26657"
        }
    ]
}
```

The use of this style of configuration allows operations to be connected to either a public node (not recommended for production) or a private node (recommended for production).
