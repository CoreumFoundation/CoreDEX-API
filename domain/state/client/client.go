package client

import (
	"context"

	grpcdef "github.com/CoreumFoundation/CoreDEX-API/domain/state"
	grpcclient "github.com/CoreumFoundation/CoreDEX-API/utils/grpc-client"
)

const endpoint = "STATE_STORE"

var client grpcdef.StateServiceClient
var grpcClient *grpcclient.GRPCClient

/*
Initialize the client.
Depending on the parameter, the environment is determined to be either in cluster of local by:
localhost:port => local
localhost => No port is not local
*/
func initClient() {
	grpcClient = grpcclient.InitClient(endpoint)
	client = grpcdef.NewStateServiceClient(grpcClient.Conn)
}

func Client() grpcdef.StateServiceClient {
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
