package coreum

import (
	"crypto/tls"
	"strings"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	coreumconfig "github.com/CoreumFoundation/coreum/v5/pkg/config"

	"github.com/CoreumFoundation/coreum/v5/pkg/client"
	"github.com/CoreumFoundation/coreum/v5/pkg/config/constant"
)

var clients map[metadata.Network]*client.Context

func NewNodeConnections() map[metadata.Network]*client.Context {
	if clients != nil {
		return clients
	}
	clients = make(map[metadata.Network]*client.Context)

	// Parse the ENV variable NETWORKS
	networks := ParseConfig()
	for _, node := range networks.Node {
		logger.Infof("Connecting to GRPC interface: %s", node.GRPCHost)
		transportCredentials := credentials.NewTLS(&tls.Config{})
		if strings.HasPrefix(node.GRPCHost, "127.0.0.1") || strings.HasPrefix(node.GRPCHost, "localhost") {
			transportCredentials = insecure.NewCredentials()
		}
		grpcClient, err := grpc.NewClient(node.GRPCHost, grpc.WithTransportCredentials(transportCredentials))
		logger.Infof("Connected to GRPC interface: %s", node.GRPCHost)
		if err != nil {
			logger.Fatalf("error connecting to coreum GRPC interface: %v", err)
		}

		logger.Infof("Connecting to RPC interface: %s", node.RPCHost)
		rpcClient, err := sdkclient.NewClientFromNode(node.RPCHost)
		logger.Infof("Connected to RPC interface: %s", node.RPCHost)
		if err != nil {
			logger.Fatalf("error connecting to coreum RPC interface: %v", err)
		}

		// ChainID is set to mainnet (default)
		chainID := constant.ChainIDMain
		// And switched to testnet or devnet if needed
		switch node.Network {
		case "testnet":
			chainID = constant.ChainIDTest
		case "devnet":
			chainID = constant.ChainIDDev
		}

		network, err := coreumconfig.NetworkConfigByChainID(chainID)
		if err != nil {
			panic(err)
		}
		network.SetSDKConfig()

		modules := auth.AppModuleBasic{}
		encodingConfig := coreumconfig.NewEncodingConfig(modules)

		clientCtx := client.NewContext(client.DefaultContextConfig(), modules).
			WithChainID(string(chainID)).
			WithGRPCClient(grpcClient).
			WithClient(rpcClient).
			WithKeyring(keyring.NewInMemory(encodingConfig.Codec)).
			WithBroadcastMode(flags.BroadcastSync).
			WithAwaitTx(true)
		clients[metadata.Network(metadata.Network_value[strings.ToUpper(node.Network)])] = &clientCtx
	}
	return clients
}
