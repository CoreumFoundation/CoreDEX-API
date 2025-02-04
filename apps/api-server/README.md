# Coreum DEX friendly API server

The API server provides a RESTFUL API for the Coreum DEX.
The goal of the API server is to flatten the learning curve for the Coreum DEX (and with that of the Cosmos SDK) and to provide a simple and easy to use API for developers to interact with the Coreum DEX.
The only part where there is some Cosmos knowledge required is where the developer has to interact with a wallet application to sign the transaction.

## Requirement for running the API

The API depends on the store to be running. The data in the store is populated by the `data-aggregator`, which is also present in this repository (see `/apps/data-aggregator`).

## Contents of the API

The following is provided:

* OHLC data
* Trade history & currencies
* Order book
* Update service using a websocket for real-time updates
* Endpoint for transaction code generation
* Endpoint for order submission

In the following paragraphs, the endpoints are described in more detail.

### Routes

The API server provides the following routes:

* GET /ohlc : Returns OHLC data
* GET /trades : Returns trade data (filterable)
* GET /tickers?symbols=base64encoded : Returns the latest price for one or more trading pairs
* GET /currencies : Returns the currencies
* GET /ws : Websocket for real-time updates
* POST /order/create : Create an order
* POST /order/cancel : Cancel an order
* POST /order/submit : Submit an order
* GET /order/orderbook : Returns the order book
* GET /wallet/assets : Returns the assets of a wallet

#### GET /ohlc

Params:

- `symbol` *required* - symbol for which OHLC should be returned. NOTE: `symbol` should be urlsafe encoded.
- `period` *required* - one of `["1m","3m","5m","15m","30m","1h","3h","6h","12h","1d","3d","1w"]`
- `from` *required* - unix timestamp of OHLC start
- `to` *required* - unix timestamp of OHLC end

**NOTE: API has limitation of 2000 points per request**

It means that you can request 2000 of 1m points (2000/60/24) ~1.3 days. Or 2000 of 1h points ~83.3 days etc.

Returns:

```json5
[
  [
    1611069600,         // timestamp
    "5.54016620498615", // open
    "5.54016620498615", // high
    "5.54016620498615", // low
    "5.54016620498615", // close
    "1.9855"            // volume
  ],
  [
    1611069660,
    "5.54016620498615",
    "5.54016620498615",
    "5.54016620498615",
    "5.54016620498615",
    "0"
  ],
]
```

#### Example

```bash
curl "https://api.coreum.com/ohlc?symbol=dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs&period=1m&from=1611007200&to=1611070980" \
--header "Network: mainnet"
```

#### Get /trades

Params:

- `symbol` *required* - symbol for which trades should be returned. NOTE: `symbol` should be urlsafe encoded.
- `from` *required* - unix timestamp of trades start
- `to` *required* - unix timestamp of trades end
- `account` *optional* - account address for which trades should be returned

Returns:

```json5
{
  "Trades": [
   {
        "Account": "devcore1fpdgztw4aepgy8vezs9hx27yqua4fpewygdspc",
        "OrderID": "8b341e25-482e-487f-b9e2-9467d98c16ac",
        "Sequence": 27388,
        "Amount": {
            "Value": 2080,
            "Exp": -6
        },
        "Price": 35.015385,
        "Denom1": {
            "Currency": "dextestdenom8",
            "Issuer": "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
            "Precision": 6,
            "Denom": "dextestdenom8-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs"
        },
        "Denom2": {
            "Currency": "dextestdenom3",
            "Issuer": "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
            "Precision": 6,
            "Denom": "dextestdenom3-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs"
        },
        "Side": 1,
        "BlockTime": {
            "seconds": 1736358800,
            "nanos": 634506142
        },
        "MetaData": {
            "Network": 3,
            "UpdatedAt": {
                "seconds": 1737747689,
                "nanos": 498590000
            },
            "CreatedAt": {
                "seconds": 1737747689,
                "nanos": 498029000
            }
        },
        "TXID": "29E2362BE19BE53B5A38CFAAB4B777484F5956972C656A4378D7620A6E8F4A36",
        "BlockHeight": 6714462,
        "HumanReadablePrice": "35.015385",
        "SymbolAmount": "0.002080"
    },
    //...
  ],
  "Offset": 0
}
```

#### tickers

The tickers endpoint returns the latest price for all trading pairs.

Params:

* symbols *mandatory* - a base64 encoded list of symbols for which the ticker should be returned.

Maximum 20 symbols. Watch out for overflow of the URL in certain browsers: The symbol strings are quite long, so most likely the queries should be limited to 10 symbols or even less

The content to be base64 encode is a JSON array: `["USD-..._BTC-...]`.

Example call:

```bash
curl -H "Network: mainnet" \
-X "GET" "https://api.coreum.com/tickers?symbols=WyJVU0QrckQ5VzdVTHZlYXZ6OHFCR00xUjVqTWdLMlFLc0VEUFFWaS9YUlAiLCAiWFJQL1VTRCtyRDlXN1VMdmVhdno4cUJHTTFSNWpNZ0syUUtzRURQUVZpIl0="
```

Returns: 

```json
{
    "Tickers": {
        "dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs": {
            "OpenTime": 1723084712,
            "CloseTime": 1723171112,
            "OpenPrice": 0.17726004,
            "HighPrice": 3135.5825912454534,
            "LowPrice": 0.1730303609235772,
            "LastPrice": 0.1783099502981715,
            "FirstPrice": 0.1730303609235772,
            "Volume": 170639.98819271981,
            "InvertedVolume": 30319.373563999994,
            "Inverted": false
        }
    },
    "USDTickers": {
        "dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs": {
            "OpenTime": 1723084712,
            "CloseTime": 1723171112,
            "OpenPrice": 1.24082028,
            "HighPrice": 21949.078138718174,
            "LowPrice": 1.2112125264650404,
            "LastPrice": 1.2481696520872005,
            "FirstPrice": 1.2112125264650404,
            "Volume": 24377.141170388546,
            "InvertedVolume": 212235.61494799994,
            "Inverted": false
        }
    }
}
```

#### Currencies

Retrieve the currencies.
Returns a list of currencies.

```json5
{
    "Currencies": [
        {
            "Denom": {
                "Currency": "ATOM",
                "Issuer": "cosmoshub-4",
                "Precision": 6,
                "IsIBC": false,
                "Denom": "uatom",
                "Name": "Cosmos",
                "Description": "Cosmos Hub native token",
                "Icon": "https://cosmos.network/logo.png"
            },
            "SendCommission": {
                "Value": 1,
                "Exp": -2
            },
            "BurnRate": {
                "Value": 1,
                "Exp": -3
            },
            "InitialAmount": {
                "Value": 1000000,
                "Exp": 0
            },
            "Chain": "cosmoshub-4",
            "OriginChain": "cosmoshub-4",
            "ChainSupply": "100000000",
            "Description": "Cosmos Hub native token",
            "SkipDisplay": false,
            "MetaData": {
                "Network": 1,
                "UpdatedAt": "2023-10-01T12:00:00Z",
                "CreatedAt": "2020-01-01T00:00:00Z"
            }
        },
        //... 
    ],
    "Offset": 0
}
```

Example call:

```bash
curl -H "Network: mainnet" \
-X "GET" "https://api.coreum.com/currencies"
```

#### /order/create

The order create uses a POST request to create an order. The function returns a to be signed transaction. The function does not persist the order in the order book (persistence is done only after submitting the order to the blockchain).

The call takes a POST:

```json5
{
  "Sender": "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
  "Type": ORDER_TYPE_LIMIT, // Enum value
  "OrderID": "string", // Optional, can be an empty string
  "BaseDenom": "dextestdenom5-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
  "QuoteDenom": " dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
  // price is value of one unit of the base_denom expressed in terms of the quote_denom.
  "Price": "0.25",
  "Quantity": "1000",
  "Side": SIDE_BUY, // Enum value
  "GoodTil": {
    "GoodTilBlockHeight": 1000,
    "GoodTilBlockTime": "2023-10-01T12:00:00Z"
  },
  "TimeInForce": TIME_IN_FORCE_GTC // Enum value
}
```

The response is a transaction that needs to be signed and submitted to the blockchain.

Example call:

```bash
curl -H "Network: devnet" \
-X "POST" "https://api.coreum.com/order/create" \
-d '{
    "Sender": "devcore1878pk82zlndhldglx26r606qcd886562mad59y",
    "Type": 1,
    "BaseDenom": "dextestdenom5-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
    "QuoteDenom": "dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
    "Price": "0.25",
    "Quantity": "1000",
    "Side": 1,
    "GoodTil": {
        "GoodTilBlockTime": "2025-12-30T12:00:00Z"
    },
    "TimeInForce": 1
    }'
```

Example response:

```json
{
  "TXBytes":"CvwBCvkBChwvY29yZXVtLmRleC52MS5Nc2dQbGFjZU9yZGVyEtgBCi5kZXZjb3JlMXAwZWR6eXpwYXpwdDY4dmRyankyMGM0Mmx2d3NqcHZmemFoeWdzEAEaBnN0cmluZyI8ZGV4dGVzdGRlbm9tNS1kZXZjb3JlMXAwZWR6eXpwYXpwdDY4dmRyankyMGM0Mmx2d3NqcHZmemFoeWdzKjxkZXh0ZXN0ZGVub205LWRldmNvcmUxcDBlZHp5enBhenB0Njh2ZHJqeTIwYzQybHZ3c2pwdmZ6YWh5Z3MyBTI1ZS0yOgoxMDAwMDAwMDAwQAFKCwjoBxIGCMCv9MYGEgISAA=="
}
```

In which the TXBytes is the base64 encoded transaction that needs to be signed and submitted to the blockchain.

#### /order/cancel

The order cancel uses a POST request to cancel an order. The function returns a transaction hash.

Example call:

```bash
curl -H "Network: devnet" \
-X "POST" "https://api.coreum.com/order/cancel" \
-d '{
    "Sender": "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
    "OrderID": "8b341e25-482e-487f-b9e2-9467d98c16ac"
}'
```

#### /order/submit

The order submit uses a POST request to submit an order. The function returns a transaction hash.

Example call:

```bash
curl -H "Network: devnet" \
-X "POST" "https://api.coreum.com/order/submit" \
-d '{
  "TX": "CqcCCqQCChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMCCvwBCvkBChwvY29yZXVtLmRleC52MS5Nc2dQbGFjZU9yZGVyEtgBCi5kZXZjb3JlMXAwZWR6eXpwYXpwdDY4dmRyankyMGM0Mmx2d3NqcHZmemFoeWdzEAEaBnN0cmluZyI8ZGV4dGVzdGRlbm9tNS1kZXZjb3JlMXAwZWR6eXpwYXpwdDY4dmRyankyMGM0Mmx2d3NqcHZmemFoeWdzKjxkZXh0ZXN0ZGVub205LWRldmNvcmUxcDBlZHp5enBhenB0Njh2ZHJqeTIwYzQybHZ3c2pwdmZ6YWh5Z3MyBTI1ZS0yOgoxMDAwMDAwMDAwQAFKCwjoBxIGCMCv9MYGEgISABJSCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDWIusTf5MdJVCJlqRR1MVNO5So6We2b5v42yCQl+k3D0SBAoCCAESABpAa+STFz77MqaUIU2bbS5uKtXeOyGW/m2qLcf96WEAJPJCSnO/3em/QZtGKnF1eg4ylCawYOH/8gpDLrxI1ybsiw=="
}
```

Where the TX is the signed TX in base64 encoding.

Response: 

```json
{
  "TXHash": "0x1234567890"
}
```

#### /order/orderbook

Params:

- `symbol` *required* - symbol for which order book should be returned. NOTE: `symbol` should be urlsafe encoded.
- `account` *optional* - account address for which the orderbook should be returned

Symbol is defined as `denom1_denom2` where `denom1` and `denom2` are the denom strings of the two assets in the trading pair.

The order book returns the first 50 buys, and the first 50 sells around the spread.
In this buy is defined as **I want to buy this from you at this price** and sell is defined as **You want to sell this to you at this price**

Sample return:

```json5
{
    "Buy": [
        {
            "Price": "3140",
            "HumanReadablePrice": "3140",
            "Amount": "2272",
            "SymbolAmount": "0.002272",
            "Sequence": 41567,
            "OrderID": "8b341e25-482e-487f-b9e2-9467d98c16ac"
        }
    ],
    "Sell": [
        {
            "Price": "3360",
            "HumanReadablePrice": "3360",
            "Amount": "2071",
            "SymbolAmount": "0.002071",
            "Sequence": 41760,
            "OrderID": "8b341e25-482e-487f-b9e2-9467d98c16ac"
        }
    ]
} 
```

#### /wallet/assets

Returns the assets for a certain account.

Params:
- `address` *required* - the address of the account for which the assets should be returned.

Sample return:

```json5
[
    {
        "Denom": "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
        "Amount": "1000000000000",
        "SymbolAmount": "1000.0000"
    },
    {
        "Denom": "dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
        "Amount": "1000000000000",
        "SymbolAmount": "1000.0000"
    },
    //...
]
```

Example call:

```bash
curl -H "Network: devnet" \
-X "GET" "https://api.coreum.com/wallet/assets?address=devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs"
```

#### Update service for real-time updates

The update service uses a websocket for a subscription system in which the user subscribes to certain information.

See the [README-update-service.md](README-update-service.md) for extensive information.

## Application start parameters

The API server can be started with the following parameters:

- `NETWORKS` - The networks handled by the blockchain scanner. 
- `STATE_STORE` - Location where the store is running, typically at grpc port 50051, syntax for non-google cloudrun usage: {host}:50051
- `TRADE_STORE` - Location where the store is running, typically at grpc port 50051, syntax for non-google cloudrun usage: {host}:50051
- `OHLC_STORE` - Location where the store is running, typically at grpc port 50051, syntax for non-google cloudrun usage: {host}:50051
- `ORDER_STORE` - Location where the store is running, typically at grpc port 50051, syntax for non-google cloudrun usage: {host}:50051
- `CURRENCY_STORE` - Location where the store is running, typically at grpc port 50051, syntax for non-google cloudrun usage: {host}:50051
- `LOG_LEVEL` - info. Other values: debug, error, warn
- `HTTP_CONFIG` - HTTP configuration with CORS settings
- `BASE_COIN` - Native/system coin configuration
- `BASE_USDC` - USDC reference for resolving equivalent USD values

### Networks

The networks are configured in a json array which depending on usage might require escaping. The network configuration is as follows:

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

For production purposes it is advised to use a private node, and the API servfer functions are proxy for that scenario.

### HTTP configuration

The HTTP configuration is configured in a json array which depending on usage might require escaping. The HTTP configuration is as follows:

```json
{
  "HTTP": {
    "port": ":8080",
    "cors": {
      "allowedOrigins": [
        "http://localhost:3000",
        "http://localhost:3001"
      ]
    },
    "timeouts": {
      "read": "10s",
      "write": "10s",
      "idle": "10s",
      "shutdown": "10s"
    }
  }
}
```

The most often changed values are the `cors.AllowOrigins` and the `port`.

### Base coin

The base coin is configured in a json array which depending on usage might require escaping. The base coin configuration is as follows:

```json
{
  "BaseCoin": [
    {
      "Network": "mainnet",
      "Coin": "ucore"
    },
    {
      "Network": "testnet",
      "Coin": "utestcore"
    },
    {
      "Network": "devnet",
      "Coin": "udevcore"
    }
  ]
}
```

### Base USDC

Base USDC is used by a path resolving algorithm (non-weighted Dijkstra alghorithm) to find the path between two assets and be able to resolve that to the representative USD value using a USDC as base.
Often devnet and testnet will not have an IBC USDC currency available, resulting in a mainnet configuration only. On devnet/testnet the user will most likely not see any USD values in the data.
The base USDC is configured in a json array which depending on usage might require escaping. The base USDC configuration is as follows:

```json
{
  "BaseCoin": [
    {
      "Network": "mainnet",
      "Coin": "uusdc-E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D"
    }
  ]
}
```
