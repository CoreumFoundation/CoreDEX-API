package coreum

import (
	"crypto/tls"
	"strings"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"github.com/CoreumFoundation/coreum/v5/pkg/client"
	coreumconfig "github.com/CoreumFoundation/coreum/v5/pkg/config"
	"github.com/CoreumFoundation/coreum/v5/pkg/config/constant"
)

var (
	clients  map[metadata.Network]*client.Context
	chainID  = constant.ChainIDMain
	networks = ParseConfig()
)

func NewNodeConnections() map[metadata.Network]*client.Context {
	if clients != nil {
		return clients
	}
	clients = make(map[metadata.Network]*client.Context)

	// ChainID is set to mainnet (hardcoded here since we do not use the chain in such a way that we need to change it
	// plus the cosmos sdk has a limitation in being able to support more than 1 chain at the time, and we want to connect
	// several chains here (less resources, plus there are use cases where we might want to connect to several chains)
	network, err := coreumconfig.NetworkConfigByChainID(chainID)
	if err != nil {
		logger.Fatalf("error getting network config: %v", err)
	}
	network.SetSDKConfig()

	for _, node := range networks.Node {
		NodeConnection(node.Network)
	}
	return clients
}

// Replaces or creates a node connection for a given node and chainID
func NodeConnection(network string) *client.Context {
	network = strings.ToUpper(network)
	node := &Node{}
	// Locate the Node in the networks config
	for _, n := range networks.Node {
		logger.Infof("network: %s", n.Network)
		if n.Network == network {
			node = n
			break
		}
	}
	if node.Network == "" {
		logger.Fatalf("node not found in networks config: %s", network)
	}

	transportCredentials := credentials.NewTLS(&tls.Config{})
	if strings.HasPrefix(node.GRPCHost, "127.0.0.1") || strings.HasPrefix(node.GRPCHost, "localhost") {
		transportCredentials = insecure.NewCredentials()
	}

	modules := auth.AppModuleBasic{}
	encodingConfig := coreumconfig.NewEncodingConfig(modules)

	pc, ok := encodingConfig.Codec.(codec.GRPCCodecProvider)
	if !ok {
		logger.Fatalf("failed to cast codec to codec.GRPCCodecProvider")
	}

	grpcClient, err := grpc.NewClient(
		node.GRPCHost,
		grpc.WithDefaultCallOptions(grpc.ForceCodec(pc.GRPCCodec())),
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		logger.Fatalf("error connecting to coreum GRPC interface: %v", err)
	}
	logger.Infof("Connected to GRPC interface: %s", node.GRPCHost)

	rpcClient, err := sdkclient.NewClientFromNode(node.RPCHost)
	if err != nil {
		logger.Fatalf("error connecting to coreum RPC interface: %v", err)
	}
	logger.Infof("Connected to RPC interface: %s", node.RPCHost)

	clientCtx := client.NewContext(client.DefaultContextConfig(), modules).
		WithChainID(string(chainID)).
		WithGRPCClient(grpcClient).
		WithClient(rpcClient).
		WithKeyring(keyring.NewInMemory(encodingConfig.Codec)).
		WithBroadcastMode(flags.BroadcastSync).
		WithAwaitTx(true)
	clients[metadata.Network(metadata.Network_value[strings.ToUpper(node.Network)])] = &clientCtx
	return &clientCtx
}
