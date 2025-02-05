package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"cosmossdk.io/log"
	db "github.com/cosmos/cosmos-db"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txconfig "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/CoreumFoundation/coreum/v5/app"
	"github.com/CoreumFoundation/coreum/v5/pkg/client"
	coreumconfig "github.com/CoreumFoundation/coreum/v5/pkg/config"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

const (
	// txBase64       = "CnkKdwodL2NvcmV1bS5kZXgudjEuTXNnQ2FuY2VsT3JkZXISVgouZGV2Y29yZTE4NzhwazgyemxuZGhsZGdseDI2cjYwNnFjZDg4NjU2Mm1hZDU5eRIkZDkyY2M0YTYtMjRmMC00MmU2LWJiMmYtNmE2OWJlZjFmMmNlEmgKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQNYi6xN/kx0lUImWpFHUxU07lKjpZ7Zvm/jbIJCX6TcPRIECgIIARIWChAKCHVkZXZjb3JlEgQ2ODc1EMCaDBpA/HmZkyMiSVph4Qc00soEymAXmz3OqCuz2u33e6KE2tJVK9NBQ60rVU+6vbUBFebOzAcisu+wGGQ/zV0DG7lLaA=="
	// senderMnemonic = "silk loop drastic novel taste project mind dragon shock outside stove patrol immense car collect winter melody pizza all deputy kid during style ribbon"
	chainID = "coreum-devnet-1"

	grpcHost = "full-node.devnet-1.coreum.dev:9090"
	rpcHost  = "https://full-node.devnet-1.coreum.dev:26657"
)

func main() {
	flag.Parse()
	if flag.NArg() < 25 {
		panic("input should be base64 encoded transaction followed by mnemonics")
	}
	args := flag.Args()
	txBase64 := args[0]
	senderMnemonic := strings.Join(args[1:], " ")

	network, err := coreumconfig.NetworkConfigByChainID(chainID)
	if err != nil {
		panic(err)
	}
	network.SetSDKConfig()
	app.ChosenNetwork = network
	modules := auth.AppModuleBasic{}

	transportCredentials := credentials.NewTLS(&tls.Config{})
	grpcClient, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(err)
	}

	rpcClient, err := sdkclient.NewClientFromNode(rpcHost)
	if err != nil {
		panic(err)
	}

	clientCtx := client.NewContext(client.DefaultContextConfig(), modules).
		WithChainID(chainID).
		WithGRPCClient(grpcClient).
		WithClient(rpcClient).
		WithBroadcastMode(flags.BroadcastSync).
		WithAwaitTx(true)
	interfaceRegistry := clientCtx.InterfaceRegistry()
	dextypes.RegisterInterfaces(interfaceRegistry)

	// Create a keyring and generate private key from mnemonic
	kr := keyring.NewInMemory(clientCtx.Codec())
	senderInfo, err := kr.NewAccount(
		"key-name",
		senderMnemonic,
		"",
		sdk.GetConfig().GetFullBIP44Path(),
		hd.Secp256k1,
	)
	if err != nil {
		panic(err)
	}

	// Get the address from the keyring
	senderAddress, err := senderInfo.GetAddress()
	if err != nil {
		panic(err)
	}

	tempApp := app.New(log.NewNopLogger(), db.NewMemDB(), nil, true, sims.NewAppOptionsWithFlagHome(tempDir()))
	encodingConfig := coreumconfig.EncodingConfig{
		InterfaceRegistry: tempApp.InterfaceRegistry(),
		Codec:             tempApp.AppCodec(),
		TxConfig:          tempApp.TxConfig(),
		Amino:             tempApp.LegacyAmino(),
	}

	enabledSignModes := make([]signing.SignMode, 0)
	enabledSignModes = append(enabledSignModes, tx.DefaultSignModes...)
	enabledSignModes = append(enabledSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
	txConfigOpts := tx.ConfigOptions{
		EnabledSignModes:           tx.DefaultSignModes,
		TextualCoinMetadataQueryFn: txconfig.NewGRPCCoinMetadataQueryFn(clientCtx),
	}
	txConfig, err := tx.NewTxConfigWithOptions(
		encodingConfig.Codec,
		txConfigOpts,
	)
	if err != nil {
		panic(err)
	}

	// Create a client context
	clientCtx = client.NewContext(client.DefaultContextConfig(), modules).
		WithChainID(chainID).
		WithKeyring(kr).
		WithTxConfig(txConfig).
		WithInterfaceRegistry(interfaceRegistry).
		WithGRPCClient(grpcClient).
		WithClient(rpcClient).
		WithBroadcastMode(flags.BroadcastSync).
		WithAwaitTx(true)

	// Prepare the transaction factory
	txf := client.Factory{}.
		WithKeybase(clientCtx.Keyring()).
		WithChainID(clientCtx.ChainID()).
		WithTxConfig(clientCtx.TxConfig()).
		WithSimulateAndExecute(true)

	txBytes, err := base64.StdEncoding.DecodeString(txBase64)
	if err != nil {
		panic(err)
	}

	txWrapper, err := clientCtx.TxConfig().TxDecoder()(txBytes)
	if err != nil {
		panic(err)
	}

	builder, err := clientCtx.TxConfig().WrapTxBuilder(txWrapper)
	if err != nil {
		panic(err)
	}

	req := &authtypes.QueryAccountRequest{
		Address: senderAddress.String(),
	}
	authQueryClient := authtypes.NewQueryClient(clientCtx)
	ctx := context.Background()
	res, err := authQueryClient.Account(ctx, req)
	if err != nil {
		panic(err)
	}

	var acc sdk.AccountI
	if err := clientCtx.InterfaceRegistry().UnpackAny(res.Account, &acc); err != nil {
		panic(err)
	}

	txf = txf.
		WithAccountNumber(acc.GetAccountNumber()).
		WithSequence(acc.GetSequence())

	//Sign the transaction
	err = client.Sign(context.Background(), txf, "key-name", builder, true)
	if err != nil {
		panic(err)
	}

	//Get the signed transaction bytes
	txBytes, err = clientCtx.TxConfig().TxEncoder()(builder.GetTx())
	if err != nil {
		panic(err)
	}

	encodedTx := base64.StdEncoding.EncodeToString(txBytes)

	fmt.Printf("%s %s", encodedTx, senderAddress.String())
}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		dir = app.DefaultNodeHome
	}
	defer os.RemoveAll(dir)

	return dir
}
