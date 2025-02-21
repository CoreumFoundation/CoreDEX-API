# Store

This store implementation uses MySQL to persist the data. The store is a grpc server such that essentially any other database can be used.
The tables in the store are independent of each other: There are no merge operations by design. By avoiding merge operations, it becomes easier to scale the store by backing it, for example, with several databases of the same or different architecture.

The server sizing can be a challenge since the overall application can produce a high load, partially dependent on the number of end users, and on trading activity on the blockchain.

The store contains multiple GRPC services, all exposed on the same port:

- `State` - Used to store and retrieve the state of the application - See main README.md for usage
- `OrderData` - Used to store and retrieve orders
- `OrderDataHistory` - Used to store and retrieve order history
- `Trade` - Used to store and retrieve trades (executed orders either whole or partial)
- `TradePairs` - Used to store and retrieve trade pairs (can be used to populating a drop-down with active markets)
- `OHLC` - Used to store and retrieve OHLC data (Open High Low Close = OHLC)
- `Currency` - Used to retrieve denom/currency information

## Setting up the store

The setup of the store requires a MySQL database with the following:

- A database, name of your choice
- A user with read/write access to that database

All tables, indexes, and foreign keys will be created by the store itself on first use of such an entity.

## Start parameters

- `MYSQL_CONFIG` - See utils/mysqlstore for connection description
- `LOG_LEVEL` (optional) - Set the log level to one of the following values: `debug`, `info`, `warn`, `error`
- `GRPC_PORT` - Format `:{port number}`. Port the grpc server will listen on