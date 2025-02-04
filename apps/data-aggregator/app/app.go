package app

import (
	"context"
	"time"

	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/app/state"
	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain"
	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain/dex"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencyclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	denomproto "github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderclient "github.com/CoreumFoundation/CoreDEX-API/domain/order/client"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradeclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptosecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/app/ohlc"

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
		// We run scanCurrencies once here in blocking fashion to make sure
		// we have all the currencies in the database before starting to scan blocks
		if err := l.scanCurrencies(ctx, reader); err != nil {
			logger.Errorf("scanning currencies of %s failed: %v", reader.Network.String(), err)
		}
		go l.startCurrenciesScan(ctx, reader)
		go l.startBlocksScan(ctx, reader)
	}
}

func (l *Application) startBlocksScan(ctx context.Context, reader *coreum.Reader) {
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

// Rescan the currencies every 30 minutes
func (l *Application) startCurrenciesScan(ctx context.Context, reader *coreum.Reader) {
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := l.scanCurrencies(ctx, reader); err != nil {
				logger.Errorf("scanning currencies of %s failed: %v", reader.Network.String(), err)
			}
		}
	}
}

func (l *Application) scanCurrencies(ctx context.Context, reader *coreum.Reader) (err error) {
	currencyClient := currencyclient.Client()

	tokenRegistryEntries, err := dmn.GetTokenRegistryEntries(ctx, reader.Network)
	if err != nil {
		logger.Errorf("could not get token registry entries : %v", err)
	}

	var metadataList []banktypes.Metadata
	var paginationKey []byte = nil
	for {
		metadataList, paginationKey, err = reader.QueryDenomsMetadata(ctx, paginationKey)
		if err != nil {
			return err
		}
		for _, meta := range metadataList {
			meta := meta
			parsedDenom, err := denomproto.NewDenom(meta.Base)
			if err != nil {
				logger.Errorf("could not parse denom %s : %v", meta.Base, err)
				continue
			}
			parsedDenom.Name = &meta.Name
			parsedDenom.Description = &meta.Description
			for _, denomUnit := range meta.DenomUnits {
				denomUnit := denomUnit
				if denomUnit.Denom == meta.Symbol {
					precision := int32(denomUnit.Exponent)
					parsedDenom.Precision = &precision
				}
			}
			c := &currency.Currency{
				Denom:          parsedDenom,
				SendCommission: nil,
				BurnRate:       nil,
				InitialAmount:  nil,
				Chain:          "",
				OriginChain:    "",
				ChainSupply:    "",
				Description:    meta.Description,
				SkipDisplay:    false,
				MetaData: &metadata.MetaData{
					Network:   reader.Network,
					UpdatedAt: timestamppb.Now(),
					CreatedAt: timestamppb.Now(),
				},
			}
			cur, err := currencyClient.Get(ctx, &currency.ID{
				Network: reader.Network,
				Denom:   meta.Base,
			})
			if err != nil || cur.Denom == nil {
				logger.Warnf("could not find denom %s in database : %v", meta.Base, err)
			} else {
				c = cur
			}
			// This occurs on certain denoms, debug line to see which ones
			if c.Denom == nil {
				logger.Warnf("denom is nil for %s", meta.Base)
				continue
			}
			c.MetaData.UpdatedAt = timestamppb.Now()
			_, err = currencyClient.Upsert(ctx, c)
			if err != nil {
				logger.Errorf("could not upsert denom %s : %v", meta.Base, err)
				continue
			}
		}
		if paginationKey == nil {
			break
		}
	}

	var denomList types.Coins
	paginationKey = nil
	for {
		denomList, paginationKey, err = reader.QueryDenoms(ctx, paginationKey)
		if err != nil {
			return err
		}

		for _, currentDenom := range denomList {
			currentDenom := currentDenom
			d, err := denomproto.NewDenom(currentDenom.Denom)
			if err != nil {
				logger.Errorf("could not parse denom %s : %v", currentDenom.Denom, err)
				continue
			}
			c := &currency.Currency{
				Denom: d,
				MetaData: &metadata.MetaData{
					Network:   reader.Network,
					UpdatedAt: timestamppb.Now(),
					CreatedAt: timestamppb.Now(),
				},
			}
			cur, err := currencyClient.Get(ctx, &currency.ID{
				Network: reader.Network,
				Denom:   currentDenom.Denom,
			})
			if err != nil {
				logger.Infof("could not get denom %s from database, initializing new currency: %v", currentDenom.Denom, err)
			} else {
				c = cur
			}
			if c.MetaData == nil {
				c.MetaData = &metadata.MetaData{
					Network:   reader.Network,
					UpdatedAt: timestamppb.Now(),
					CreatedAt: timestamppb.Now(),
				}
			}
			c.ChainSupply = currentDenom.Amount.String()
			c.MetaData.UpdatedAt = timestamppb.Now()
			if token, ok := tokenRegistryEntries[currentDenom.Denom]; ok {
				tokenName := token.TokenName
				c.Denom.Name = &tokenName
				tokenPrecision := int32(token.Decimals)
				c.Denom.Precision = &tokenPrecision
				tokenIcon := token.LogoURIs.Png
				c.Denom.Icon = &tokenIcon
				tokenDescription := token.Description
				c.Denom.Description = &tokenDescription
			}
			if c.Denom == nil {
				logger.Warnf("denom is nil for %s, unable to persist", currentDenom.Denom)
				continue
			}
			_, err = currencyClient.Upsert(ctx, c)
			if err != nil {
				logger.Errorf("could not upsert denom %s : %v", currentDenom.Denom, err)
				continue
			}
		}
		if paginationKey == nil {
			break
		}
	}
	return nil
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
