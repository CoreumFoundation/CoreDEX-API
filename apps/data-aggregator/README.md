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

## Notes on application start

The application does a scan of the currencies to make certain all currencies are present in the database. This is done on start of the application, while certain go routines are also already running. Depending on the speed of the scan, it can look as if the application is hanging. The log will show BlockHeight as logged is not increasing, channel capacity left is 0. This is not a problem: The application will continue processing once the currencies have been scanned. At this moment this behaviour is mainly visible on testnet, which has over 4000 currencies to process on start of the application.