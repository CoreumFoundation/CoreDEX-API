/*
The config:
- Parses the config as provided to the app
- Can only parse the config parts relevant to this middleware
- Depends on providing the config as environment variables so that init() can run independent per component and no coordination is required
*/
package client

import (
	"context"

	grpcdef "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	grpcclient "github.com/CoreumFoundation/CoreDEX-API/utils/grpc-client"
)

const endpoint = "CURRENCY_STORE"

var (
	client     grpcdef.CurrencyServiceClient
	grpcClient *grpcclient.GRPCClient
)

func initClient() {
	grpcClient = grpcclient.InitClient(endpoint)
	cl := grpcdef.NewCurrencyServiceClient(grpcClient.Conn)
	client = cl
}

func Client() grpcdef.CurrencyServiceClient {
	if client == nil {
		initClient()
	}
	return client
}

func AuthCtx(ctx context.Context) context.Context {
	if grpcClient == nil {
		initClient()
	}
	return grpcClient.AuthCtx(ctx)
}
