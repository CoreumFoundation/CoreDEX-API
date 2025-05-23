package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	gomath "math"
	"math/big"
	"math/bits"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"github.com/CoreumFoundation/coreum/v5/pkg/client"
	coreumconfig "github.com/CoreumFoundation/coreum/v5/pkg/config"
	"github.com/CoreumFoundation/coreum/v5/pkg/config/constant"
	cbig "github.com/CoreumFoundation/coreum/v5/pkg/math/big"
	assetfttypes "github.com/CoreumFoundation/coreum/v5/x/asset/ft/types"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type AccountWallet struct {
	Address  string
	Mnemonic string
}

type AppConfig struct {
	Network                   string
	Fund                      string
	GRPCHost                  string
	Issuer                    AccountWallet
	AccountsWallet            []AccountWallet
	AssetFTDefaultDenomsCount int
}

const (
	// maxWordLen defines the maximum word length supported by Int and Uint types.
	maxSDKIntWordLen = math.MaxBitLen / bits.UintSize
)

type currency struct {
	name     string
	currency string
}

type App struct {
	cfg                                          AppConfig
	issuer                                       types.AccAddress
	accounts                                     []types.AccAddress
	denoms                                       []string
	sides                                        []dextypes.Side
	clientCtx                                    client.Context
	txFactory                                    tx.Factory
	iteration                                    int // Track iteration count
	previousPrice, baseVolatility, trendStrength float64
}

/*
array with 2 elements of currency-issuer
*/
var (
	currencyArray = []currency{
		{"NOR", "nor"}, {"ALB", "alb"}}
)

func NewApp(
	ctx context.Context,
	cfg AppConfig,
) (App, error) {
	logger.Infof("initializing app")

	app := App{cfg: cfg}
	transportCredentials := credentials.NewTLS(&tls.Config{})
	if strings.HasPrefix(cfg.GRPCHost, "127.0.0.1") || strings.HasPrefix(cfg.GRPCHost, "localhost") {
		transportCredentials = insecure.NewCredentials()
	}

	modules := auth.AppModuleBasic{}
	encodingConfig := coreumconfig.NewEncodingConfig(modules)

	pc, ok := encodingConfig.Codec.(codec.GRPCCodecProvider)
	if !ok {
		logger.Fatalf("failed to cast codec to codec.GRPCCodecProvider")
	}

	grpcClient, err := grpc.NewClient(
		cfg.GRPCHost,
		grpc.WithDefaultCallOptions(grpc.ForceCodec(pc.GRPCCodec())),
		grpc.WithTransportCredentials(transportCredentials),
	)
	logger.Infof("Connected to GRPC interface %s", cfg.GRPCHost)
	if err != nil {
		return App{}, fmt.Errorf("error connecting to coreum GRPC interface: %v", err)
	}
	// ChainID is set to devnet (default)
	chainID := constant.ChainIDDev
	if cfg.Network != "" {
		if cfg.Network == "testnet" {
			chainID = constant.ChainIDTest
		}
	}

	network, err := coreumconfig.NetworkConfigByChainID(chainID)
	if err != nil {
		return App{}, err
	}
	network.SetSDKConfig()

	clientCtx := client.NewContext(client.DefaultContextConfig(), modules).
		WithChainID(string(chainID)).
		WithGRPCClient(grpcClient).
		WithKeyring(keyring.NewInMemory(encodingConfig.Codec)).
		WithBroadcastMode(flags.BroadcastSync).
		WithAwaitTx(true)

	txFactory := client.Factory{}.
		WithKeybase(clientCtx.Keyring()).
		WithChainID(clientCtx.ChainID()).
		WithTxConfig(clientCtx.TxConfig()).
		WithSimulateAndExecute(true).
		WithGasAdjustment(1.5)

	bankClient := banktypes.NewQueryClient(clientCtx)
	issuer := addToKeyring(clientCtx, cfg.Issuer)

	accounts := lo.Map(cfg.AccountsWallet, func(item AccountWallet, index int) types.AccAddress {
		return addToKeyring(clientCtx, item)
	})
	// Make certain the all accounts have funds (could be new accounts)
	app.addFunds(issuer.String())
	app.addFunds(accounts[0].String())
	app.addFunds(accounts[1].String())

	denoms := make([]string, 0)
	denoms = append(denoms, lo.RepeatBy(cfg.AssetFTDefaultDenomsCount, func(i int) string {
		denom := denom(i, issuer.String())
		supply, err := bankClient.SupplyOf(ctx, &banktypes.QuerySupplyOfRequest{Denom: denom})
		if err != nil {
			panic(err)
		}
		if supply.Amount.IsZero() {

			issueMsg := &assetfttypes.MsgIssue{
				Issuer:        issuer.String(),
				Symbol:        currencyArray[i].name,
				Subunit:       currencyArray[i].currency,
				Precision:     6,
				Description:   currencyArray[i].name,
				InitialAmount: math.NewInt(1000000000000),
				Features: []assetfttypes.Feature{
					assetfttypes.Feature_minting,
				},
			}

			_, err = client.BroadcastTx(
				ctx,
				clientCtx.WithFromAddress(issuer),
				txFactory,
				issueMsg,
			)
			if err != nil {
				logger.Infof("Error issuing %s: %s", denom, err)
			}
		}

		return denom
	})...)
	sides := []dextypes.Side{
		dextypes.SIDE_SELL,
		dextypes.SIDE_BUY,
	}
	logger.Infof("app initialized")

	app.clientCtx = clientCtx
	app.txFactory = txFactory
	app.issuer = issuer
	app.accounts = accounts
	app.denoms = denoms
	app.sides = sides
	app.baseVolatility = 0.0004  // Makes prices oscillate ±4% around the trend
	app.trendStrength = 0.000035 // Upward trend to push prices higher over time
	app.previousPrice = 75.0     // Also the initial price for the first order of the simulation
	return app, nil
}

func addToKeyring(clientCtx client.Context, item AccountWallet) types.AccAddress {
	keyInfo, err := clientCtx.Keyring().NewAccount(
		uuid.New().String(),
		item.Mnemonic,
		"",
		hd.CreateHDPath(constant.CoinType, 0, 0).String(),
		hd.Secp256k1,
	)
	if err != nil {
		panic(err)
	}

	address, err := keyInfo.GetAddress()
	if err != nil {
		panic(err)
	}

	if address.String() != item.Address {
		panic(fmt.Errorf("generated address %q is not equal to expected address %q", address.String(), item.Address))
	}

	return address
}

func (fa *App) CreateOrder(
	ctx context.Context,
	accounts []types.AccAddress,
) error {
	startTime := time.Now()

	// One side always sells, the other side always buys
	msgIssueSell, msgIssueBuy, msgPlaceSellOrder, msgPlaceBuyOrder, err := fa.genOrder(accounts)
	if err != nil {
		return err
	}

	_, err = client.BroadcastTx(
		ctx,
		fa.clientCtx.WithFromAddress(fa.issuer),
		fa.txFactory,
		msgIssueSell,
	)
	if err != nil {
		logger.Errorf("Error broadcasting TX: account=%s, denom=%s, error=%s", msgIssueSell.Sender, msgIssueSell.Coin.Denom, err.Error())
		return err
	}

	_, err = client.BroadcastTx(
		ctx,
		fa.clientCtx.WithFromAddress(fa.issuer),
		fa.txFactory,
		msgIssueBuy,
	)
	if err != nil {
		logger.Errorf("Error broadcasting TX: account=%s, denom=%s, error=%s", msgIssueBuy.Sender, msgIssueBuy.Coin.Denom, err.Error())
		return err
	}
	res, err := client.BroadcastTx(
		ctx,
		fa.clientCtx.WithFromAddress(accounts[0]),
		fa.txFactory,
		msgPlaceSellOrder,
	)
	if err != nil {
		if !strings.Contains(err.Error(), "it's prohibited to save more than 100 orders per denom") {
			logger.Errorf("Unknown error (SELL): error=%v", err)
			return err
		}
		logger.Warnf("Error: it is prohibited to save more than 100 orders per denom: account=%s, denom=%s", msgPlaceSellOrder.Sender, msgPlaceSellOrder.BaseDenom)
	}

	logger.Infof("Info: new order SELL: TX hash=%v, Block Height=%d, Gas Used=%d, price=%s", res.TxHash, res.Height, res.GasUsed, msgPlaceSellOrder.Price.String())

	res, err = client.BroadcastTx(
		ctx,
		fa.clientCtx.WithFromAddress(accounts[1]),
		fa.txFactory,
		msgPlaceBuyOrder,
	)
	if err != nil {
		if !strings.Contains(err.Error(), "it's prohibited to save more than 100 orders per denom") {
			logger.Errorf("Unknown error (BUY): error=%v", err)
			return err
		}
		logger.Errorf("Error: it's prohibited to save more than 100 orders per denom: account=%s, denom=%s", msgPlaceBuyOrder.Sender, msgPlaceBuyOrder.BaseDenom)
	}
	logger.Infof("Info: new order BUY: TX hash=%v, Block Height=%d, Gas Used=%d, price=%s, broadcasting took: %s", res.TxHash, res.Height, res.GasUsed, msgPlaceBuyOrder.Price.String(), time.Since(startTime).String())
	return nil
}

func (fa *App) GetAccounts() []types.AccAddress {
	return fa.accounts
}

func (fa *App) genOrder(accounts []types.AccAddress) (*assetfttypes.MsgMint,
	*assetfttypes.MsgMint, *dextypes.MsgPlaceOrder, *dextypes.MsgPlaceOrder, error) {
	baseDenom, quoteDenom := fa.denoms[0], fa.denoms[1]

	price := fa.getNextPrice(fa.previousPrice)
	logger.Infof("Info: Previous price: %.6f, new price: %s", fa.previousPrice, price.String())
	quantity := 10 * int64(gomath.Pow(10, 6))
	coinsToMint := types.NewCoin(baseDenom, math.NewInt(quantity))

	sellOrder := &dextypes.MsgPlaceOrder{
		Sender:     accounts[0].String(),
		Type:       dextypes.ORDER_TYPE_LIMIT,
		ID:         uuid.New().String(),
		BaseDenom:  baseDenom,
		QuoteDenom: quoteDenom,
		Price:      &price,
		Quantity:   math.NewInt(quantity),
		Side:       dextypes.SIDE_SELL,
		GoodTil: &dextypes.GoodTil{
			GoodTilBlockHeight: 0,
			GoodTilBlockTime:   lo.ToPtr(time.Now().Add(time.Hour)),
		},
		TimeInForce: dextypes.TIME_IN_FORCE_GTC,
	}
	mintSell := &assetfttypes.MsgMint{
		Sender:    fa.issuer.String(),
		Coin:      coinsToMint,
		Recipient: accounts[0].String(),
	}

	// The opposing order:
	amount, err := mulCeil(math.NewInt(quantity), price)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	coinsToMint = types.NewCoin(quoteDenom, amount)

	buyOrder := &dextypes.MsgPlaceOrder{
		Sender:     accounts[1].String(),
		Type:       dextypes.ORDER_TYPE_LIMIT,
		ID:         uuid.New().String(),
		BaseDenom:  baseDenom,
		QuoteDenom: quoteDenom,
		Price:      &price,
		Quantity:   math.NewInt(quantity),
		Side:       dextypes.SIDE_BUY,
		GoodTil: &dextypes.GoodTil{
			GoodTilBlockHeight: 0,
			GoodTilBlockTime:   lo.ToPtr(time.Now().Add(time.Hour)),
		},
		TimeInForce: dextypes.TIME_IN_FORCE_GTC,
	}
	mintBuy := &assetfttypes.MsgMint{
		Sender:    fa.issuer.String(),
		Coin:      coinsToMint,
		Recipient: accounts[1].String(),
	}
	// Increment iteration for next order
	fa.iteration++
	// Every 1000 iterations fund the accounts (way sufficient to keep running)
	if fa.iteration%1000 == 0 {
		fa.addFunds(accounts[0].String())
		fa.addFunds(accounts[1].String())
	}
	return mintSell,
		mintBuy,
		sellOrder,
		buyOrder, nil
}

func (fa *App) getNextPrice(price float64) dextypes.Price {
	// Introduce random volatility
	currentVolatility := fa.baseVolatility * (0.8 + rand.Float64()*0.4) // Random volatility between 0.8x and 1.2x base
	direction := 1 - 2*rand.Float64()                                   // Generates either -1 or 1
	priceChange := price * (direction * currentVolatility)
	trend := price * fa.trendStrength
	price += priceChange + trend
	if price < 0 {
		price = 0 // Ensure prices don't go negative
	}
	fa.previousPrice = price
	return buildNumExpPrice(price)
}

func buildNumExpPrice(
	num float64,
) dextypes.Price {
	numPart := fmt.Sprintf("%.4f", num)

	fl, _ := strconv.ParseFloat(numPart, 64)
	var priceStr string
	// Convert float to exponent based price string
	parts := strings.Split(numPart, ".")
	exp2 := 0
	if len(parts) == 2 {
		exp2 = len(parts[1])
		priceStr = fmt.Sprintf("%d", uint64(fl*gomath.Pow(10, float64(exp2))))
	}
	// PriceStr is not allowed to end in a 0: We have to strip the 0s and add them to the exponent
	for priceStr[len(priceStr)-1] == '0' {
		priceStr = priceStr[:len(priceStr)-1]
		exp2--
	}
	if exp2 != 0 {
		priceStr = priceStr + fmt.Sprintf("e-%d", exp2)
	}
	return dextypes.MustNewPriceFromString(priceStr)
}

func mulCeil(quantity math.Int, price dextypes.Price) (math.Int, error) {
	balance, remainder := cbig.IntMulRatWithRemainder(quantity.BigInt(), price.Rat())
	if !cbig.IntEqZero(remainder) {
		balance = cbig.IntAdd(balance, big.NewInt(1))
	}
	if isBigIntOverflowsSDKInt(balance) {
		return math.Int{}, errors.New("invalid order quantity and price, out of supported math.Int range")
	}

	return math.NewIntFromBigInt(balance), nil
}

// isBigIntOverflowsSDKInt checks if the big int overflows the sdkmath.Int.
// copy form sdkmath.Int.
func isBigIntOverflowsSDKInt(i *big.Int) bool {
	if len(i.Bits()) > maxSDKIntWordLen {
		return i.BitLen() > math.MaxBitLen
	}
	return false
}

func denom(denom int, issuer string) string {
	return fmt.Sprintf("%s-%s", currencyArray[denom].currency, issuer)
}
