# Coreum util

Coreum related wrappers.

Isolated from the other utils due to the kind of code that is required to interact with the coreum network, and anticipated versioning where we would not have to update the apps/other locations for every update in the Coreum api version.

## Config for the websocket

The config is for multiple networks so that only 1 listener instance is required to process all the incoming events on all the networks.

The network itself is "abstract", however the most common ones are:

- `mainnet`
- `testnet`
- `devnet`

The config is a JSON file which contains the following fields:

GRPC:

- `Network` - Network to listen on
- `Host` - grpc host

### Sample config

A config with the above information can now look as follows:

```json
{
    "GRPC": [
        {
            "Network": "devnet",
            "Host": "full-node.devnet-1.coreum.dev:9090"
        }
    ]
}
```
