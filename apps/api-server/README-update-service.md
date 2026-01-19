# Update service

The update service provides a websocket which refreshes the data based on actual changes. The use of the websocket reduces the numbers of reads required to produce the same data to the connected clients, thus lowering database load and related costs.

## Response messages

The response message echos the request so that the caller can sync up the response with the request.
The content of the message is set with the serialized JSON for the given response type, which is the same as for the restful call for the same data.

## Interaction with the websocket

Once connected to the websocket, the following messages can be sent:

* `subscribe` - Subscribe to a topic. The topic can be a generic topic, or a topic with an entity ID. The topic can be a single topic, or a list of topics.
* `unsubscribe` - Unsubscribe from a topic. The topic can be a generic topic, or a topic with an entity ID. The topic can be a single topic, or a list of topics.
* `close` - Close the connection to the websocket. The server will close the connection in that scenario. This close method is preferred so that any subscriptions can be cleaned up instead of garbage collected.

## App internals

The app on a subscription will produce data for that subscription every 10 seconds by polling the database.
The polls are done for OHLC, trades, tickers and order book, and produce the same objects are the restful API, however now wrapped in the update service response object.

### Initial connection

On the initial connection to the websocket, the client gets a response `Connected`

The websocket is accessible on the following URL:

```js
wss://api.coreum.com/ws
```

### Connection duration & load management

The connections are terminated when there is no activity for 15 minutes to prevent background usage from users to overload the application. The websocket client should reconnect if the frontend is still alive after that period. This 15 minutes is longer than the 10 minutes as defined in the accompanying frontend application for background tabs, thus essentially the frontend will always be able to reconnect when switching tabs.
Background tabs might also terminate the websocket connections (timeout) by going to sleep. This can lead to data loss in the frontend, and on re-activation of tab it is advised to reload the data from scratch (restful call) and re-subscribe to the topics.

The downside of this reconnect policy is that for the personal trades subscription, on no trading activity, the timeout will essentially always occur when a tab is open for a long time. The advice to a frontend developer is to open the subscription when the user does an actual trade. In most reasonable cases (>>99%) the trade is executed within the 15-minute time interval.

### Subscribe to a message (or topic)

The topics are enumerated in the `domain/update`. The subscription messages are also defined in the model. The choice for a separate model was made based on the previously mentioned stability risk (if this risk is real, then we do not have to refactor all, but can just implement a result kafka topic containing the required messages as defined in this `domain/update`).

A subscription consists of sending a message on the websocket with the following structure:

```json
{
    "Action": 0,
    "Subscription": {
        "Network": 2,
        "Method": 5,
        "ID": "denom1_denom2"
    }
}
```

This example represents a subscription to a ticker for a specific trade pair on testnet.

**Action Enum Values:**
- `0` - SUBSCRIBE
- `1` - UNSUBSCRIBE
- `2` - CLOSE
- `3` - RESPONSE

**Network Enum Values:**
- `1` - MAINNET
- `2` - TESTNET
- `3` - DEVNET

**Method Enum Values:**
- `1` - TRADES_FOR_SYMBOL (ID format: `{denom1}_{denom2}`)
- `2` - TRADES_FOR_ACCOUNT (ID format: `{account}`)
- `3` - TRADES_FOR_ACCOUNT_AND_SYMBOL (ID format: `{account}_{denom1}_{denom2}`)
- `4` - OHLC (ID format: `{denom1}_{denom2}_{interval}`)
- `5` - TICKER (ID format: `{denom1}_{denom2}`)
- `6` - ORDERBOOK (ID format: `{denom1}_{denom2}`)
- `7` - ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT (ID format: `{account}_{denom1}_{denom2}`)
- `8` - WALLET (ID format: `{account}`)

***Important:*** The actual messages use enums (integers), not strings. The package for this is provided in `domain/update`.

***Important:*** For the different ID encoding/decoding, see the Method in the `domain/update` package.

### Unsubscribe from a message

The unsubscribing works the same as the subscribing, however use `Action: 1` (UNSUBSCRIBE) instead of `Action: 0` (SUBSCRIBE).

### Close the connection

The close is the following JSON:

```json
{
    "Action": 2
}
```

which will close the websocket on the server side. `Action: 2` corresponds to CLOSE.

## Messages sent to the client

The messages sent to the client are in structure the same as the messages sent to the server.
The differences are:

* `Action`: `3` (RESPONSE enum value)
* `ID`: Always present since now the BE has to tell which ID should be used for a refresh (or if there is more FE logic: Which ID the potential refresh is related to)

A sample response:

```json
{
    "Action": 3,
    "Subscription": {
        "Network": 1,
        "Method": 6,
        "ID": "core1abcdef",
        "Content": "{\"Buy\":[{\"Price\":\"50000/91\",\"Amount\":\"4\",\"Sequence\":1630},{\"Price\":\"50000/11\",\"Amount\":\"1\",\"Sequence\":4304},{\"Price\":\"45e-7\",\"Amount\":\"3944\",\"Sequence\":6017},{\"Price\":\"4000/9\",\"Amount\":\"7\",\"Sequence\":3386},{\"Price\":\"387e-5\",\"Amount\":\"4124\",\"Sequence\":6093},{\"Price\":\"382e-8\",\"Amount\":\"3090\",\"Sequence\":5627},{\"Price\":\"2500000/7\",\"Amount\":\"1\",\"Sequence\":2871},{\"Price\":\"20000000/73\",\"Amount\":\"1\",\"Sequence\":3337},{\"Price\":\"10000000/187\",\"Amount\":\"1\",\"Sequence\":3881},{\"Price\":\"1000000/9\",\"Amount\":\"1\",\"Sequence\":1179}],\"Sell\":[{\"Price\":\"1/1850000000\",\"Amount\":\"756650000000\",\"Sequence\":4968},{\"Price\":\"1/37500000000\",\"Amount\":\"58687500000000\",\"Sequence\":4234},{\"Price\":\"1/4060000000\",\"Amount\":\"7413560000000\",\"Sequence\":4170},{\"Price\":\"1/4820000000\",\"Amount\":\"15028760000000\",\"Sequence\":1432},{\"Price\":\"1/6140000000\",\"Amount\":\"12998380000000\",\"Sequence\":1173},{\"Price\":\"1/65100000\",\"Amount\":\"209231400000\",\"Sequence\":5362},{\"Price\":\"1/759000\",\"Amount\":\"1022373000\",\"Sequence\":4871},{\"Price\":\"1/76600000000\",\"Amount\":\"294756800000000\",\"Sequence\":5315},{\"Price\":\"1/77300000000\",\"Amount\":\"86034900000000\",\"Sequence\":4354},{\"Price\":\"1/88400000\",\"Amount\":\"424938800000\",\"Sequence\":3393}]}"
    }
}
```

### Subscription keys

There are several subscriptions possible, and all have their respective key format:

* `OHLC` (Method: 4): `{denom1}_{denom2}_{interval}`
* `TRADES_FOR_SYMBOL` (Method: 1): `{denom1}_{denom2}`
* `TRADES_FOR_ACCOUNT` (Method: 2): `{account}`
* `TRADES_FOR_ACCOUNT_AND_SYMBOL` (Method: 3): `{account}_{denom1}_{denom2}`
* `TICKER` (Method: 5): `{denom1}_{denom2}`
* `ORDERBOOK` (Method: 6): `{denom1}_{denom2}`
* `ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT` (Method: 7): `{account}_{denom1}_{denom2}`
* `WALLET` (Method: 8): `{account}`

Where:

* `denom1` and `denom2` are the denominations (e.g., `nor-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57` and `alb-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57`)
* `interval` is the period for the OHLC data (1m, 5m, 15m, 1h, 4h, 1d, 1w, 1M)
* `account` is the account address

#### OHLC

The OHLC call (Method: 4) is the most complex call of the set since it requires the period. Other timestamps and parameters are determined or defaulted by the backend (compared to the RESTFUL call where more retrieval manipulation is possible).

The OHLC subscription will (re)produce the data for the last interval and flows over into the next interval when that interval is reached. The interval is defined by the period.

The OHLC call contains the last 10 minutes. The data needs to be merged into the existing OHLC and replace the previous records if present.

The ID format is: `{denom1}_{denom2}_{interval}`

#### TRADES

There are few different retrievals possible for trades, depending on the required information. The filters are:

* `TRADES_FOR_ACCOUNT` (Method: 2): `{account}`
* `TRADES_FOR_SYMBOL` (Method: 1): `{denom1}_{denom2}`
* `TRADES_FOR_ACCOUNT_AND_SYMBOL` (Method: 3): `{account}_{denom1}_{denom2}`

The trades are produced in the same way as the restful API, however now wrapped in the update service response object.
Trades are read from the last 10 minutes, and need to be deduplicated by the receiver against their base set.

#### ORDERBOOK

There are 2 orderbook filters:

* `ORDERBOOK` (Method: 6): `{denom1}_{denom2}`
* `ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT` (Method: 7): `{account}_{denom1}_{denom2}`

The response of both is the same.

The orderbook is a complete replacement of the previous order book and represents the current state of the order book.
