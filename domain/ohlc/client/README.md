# client

The package provides a wrapper to include the grpc service in your application.

## Start the service

The client is self initializing. For this it requires the following environment variable:

- `OHLC_STORE_ENDPOINT` - The host of the service. Host is in the format `host:port`

If the client is included and the env variable is not set, the client panics.
