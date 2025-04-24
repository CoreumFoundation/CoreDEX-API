package app

import (
	"context"

	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptosecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	currencyapp "github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/app/currency"
	marketapp "github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/app/market"
	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/app/ohlc"
	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/app/state"
	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain"
	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain/dex"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencyclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradeclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type Application struct {
	state          *state.State
	registry       *dmn.Registry
	tradeChan      chan *tradegrpc.Trade
	orderClient    order.OrderServiceClient
	tradeClient    tradegrpc.TradeServiceClient
	currencyClient currency.CurrencyServiceClient
}

type App interface {
	Handle(ctx context.Context, tx *dmn.Result, network string, height int64, eventIndex int) error
}

func NewApplication(ctx context.Context) *Application {
	orderClient := orderclient.Client()
	tradeClient := tradeclient.Client()
	currencyClient := currencyclient.Client()
	return NewApplicationWithClients(ctx, orderClient, tradeClient, currencyClient)
}

func NewApplicationWithClients(
	ctx context.Context,
	orderClient order.OrderServiceClient,
	tradeClient tradegrpc.TradeServiceClient,
	currencyClient currency.CurrencyServiceClient,
) *Application {
	interfaceRegistry := ctypes.NewInterfaceRegistry()
	dextypes.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &cryptosecp256k1.PubKey{})

	registry := dmn.NewRegistry(interfaceRegistry)
	dex.NewMsgPlaceOrderHandler(interfaceRegistry, registry)
	dex.NewMsgCancelOrderHandler(interfaceRegistry, registry)

	return &Application{
		state:          state.NewApplication(ctx),
		registry:       registry,
		tradeChan:      make(chan *tradegrpc.Trade, 1000),
		orderClient:    orderClient,
		tradeClient:    tradeClient,
		currencyClient: currencyClient,
	}
}

func (l *Application) StartScanners(ctx context.Context) {
	readers := coreum.InitReaders()
	// Get the state for the readers:
	for _, reader := range readers {
		reader.BlockHeight = l.state.GetState(ctx, reader.Network)
		logger.Infof("Start: Last scanned height for network %s is %d", reader.Network, reader.BlockHeight)
	}
	// Start the readers
	readers.Start()
	// Add a channel listener for the readers
	for _, reader := range readers {
		logger.Infof("Start: Started aggregator for network %s", reader.Network)
		currencyApp := currencyapp.NewApplication(ctx, reader)
		marketApp := marketapp.NewApplication(reader, l.tradeClient)
		go currencyApp.Start(ctx)
		go marketApp.Start(ctx)
		go l.startBlocksScan(ctx, reader)
	}
}

func (l *Application) startBlocksScan(ctx context.Context, reader *coreum.Reader) {
	logger.Infof("Start: Started scanner for network %s", reader.Network)
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-reader.ProcessBlockChannel:
			l.scannerCoordinator(ctx, block, reader.Network)
			l.state.SetState(reader.Network, reader.BlockHeight)
		}
	}
}

// Processing the actual content of the messages using the registry.HandleAction method
func (l *Application) scannerCoordinator(ctx context.Context, block *coreum.ScannedBlock, network metadata.Network) {
	for _, transaction := range block.Transactions {
		if transaction.Tx == nil {
			continue
		}
		for _, msg := range transaction.Tx.Body.Messages {
			meta := dmn.Metadata{
				Network:     network,
				BlockHeight: block.BlockHeight,
				BlockTime:   block.BlockTime,
				TxHash:      transaction.TxResponse.TxHash,
				GasUsed:     transaction.TxResponse.GasUsed,
			}
			message := l.registry.ParseMsg(msg.TypeUrl, msg.Value, meta)
			l.registry.ParseActions(ctx, l.orderClient, l.tradeClient, l.currencyClient, message, transaction.TxResponse.Events, meta, l.tradeChan)
		}
	}

	meta := dmn.Metadata{
		Network:         network,
		BlockHeight:     block.BlockHeight,
		BlockTime:       block.BlockTime,
		IsEndBlockEvent: true,
	}

	if len(block.BlockEvents) > 0 {
		// Process the block Events
		l.registry.HandleBlockEvent(ctx, l.orderClient, l.tradeClient, l.currencyClient, block.BlockEvents, meta, l.tradeChan)
	}
}

func (l *Application) StartOHLCProcessor(ctx context.Context) {
	ohlc.NewApplication(ctx, l.tradeChan)
}
