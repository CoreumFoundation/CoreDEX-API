# grpc client

This package contains the repeatable parts of the GRPC client for use within the Coreum-dex-api.

The client has the following base capabilities:

* Localhost unauthenticated behavior by providing the requested environment variable as `localhost:port` (usually `localhost:50051`)
* Authenticated behaviour (google cloud IAM) and https support by providing the requested environment variable as `host`

## Starting the a client

This client package is to be wrapped in the `models/{modelname}/client` to prevent the user from having to pass environment variables. The parameter to initClient in this package is the name of the environment variable to use, so not the value in itself. This way the checks for correctness can stay in this package and the `models/{modelname}/client` is as minimal as possible.

## Start parameters

The endpoint is parsed in as parameter.

If running in cloud run provide the endpoint as the service name, without any port number.
Provide the cloud run project dependent URL append as:

- `GRPC_APPEND` - `-abc-def` (so include the `-`)
