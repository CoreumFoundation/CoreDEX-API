// Scans the currencies and updates the database at regular intervals.
package currency

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	currencyclient "github.com/CoreumFoundation/CoreDEX-API/domain/currency/client"
	denomproto "github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type Application struct {
	reader                         *coreum.Reader
	currencyClient                 currencygrpc.CurrencyServiceClient
	bankClient                     map[metadata.Network]banktypes.QueryClient
	currencyScannerPoolInitialized bool
	denomScannerPoolInitialized    bool
	denomScannerChan               chan banktypes.Metadata
	tokenRegistryEntries           map[string]dmn.Asset
	currencyDenomScannerChan       chan types.Coin
}

func NewApplication(ctx context.Context, reader *coreum.Reader) *Application {
	app := &Application{
		reader:                         reader,
		currencyClient:                 currencyclient.Client(),
		bankClient:                     make(map[metadata.Network]banktypes.QueryClient),
		tokenRegistryEntries:           make(map[string]dmn.Asset),
		currencyDenomScannerChan:       make(chan types.Coin, 1000),
		denomScannerChan:               make(chan banktypes.Metadata, 1000),
		currencyScannerPoolInitialized: false,
		denomScannerPoolInitialized:    false,
	}
	// We run scanCurrencies once here in blocking fashion to make sure
	// we have all the currencies in the database before starting to scan blocks
	if err := app.scanCurrencies(ctx); err != nil {
		logger.Errorf("scanning currencies of %s failed: %v", reader.Network.String(), err)
	}
	return app
}

/*
Rescan currencies every 30 minutes
*/
func (l *Application) Start(ctx context.Context) {
	logger.Infof("Starting currency scanner for %s", l.reader.Network.String())
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := l.scanCurrencies(ctx); err != nil {
				logger.Errorf("scanning currencies of %s failed: %v", l.reader.Network.String(), err)
			}
		}
	}
}

func (l *Application) scanCurrencies(ctx context.Context) (err error) {
	if l.bankClient[l.reader.Network] == nil {
		l.bankClient[l.reader.Network] = banktypes.NewQueryClient(l.reader.ClientContext)
	}

	err = l.denomScanner(ctx)
	if err != nil {
		return err
	}

	err = l.currencyDenomScanner(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (l *Application) denomScanner(ctx context.Context) (err error) {
	if !l.currencyScannerPoolInitialized {
		l.currencyScannerPoolInitialized = true
		for i := 0; i < 20; i++ {
			go l.denomScannerPool(ctx)
		}
	}
	if !l.denomScannerPoolInitialized {
		l.denomScannerPoolInitialized = true
		for i := 0; i < 20; i++ {
			go l.currencyDenomScannerPool(ctx)
		}
	}
	var metadataList []banktypes.Metadata
	var paginationKey []byte = nil
	for {
		metadataList, paginationKey, err = l.reader.QueryDenomsMetadata(ctx, l.bankClient[l.reader.Network], paginationKey)
		if err != nil {
			return err
		}
		for _, meta := range metadataList {
			l.denomScannerChan <- meta
		}
		if paginationKey == nil {
			break
		}
	}
	return nil
}

// Pool of denomInternalScanner workers
func (l *Application) denomScannerPool(ctx context.Context) {
	for {
		select {
		case meta := <-l.denomScannerChan:
			l.denomScannerInternal(ctx, meta)
		}
	}
}

func (l *Application) denomScannerInternal(ctx context.Context, meta banktypes.Metadata) {
	parsedDenom, err := denomproto.NewDenom(meta.Base)
	if err != nil {
		logger.Errorf("could not parse denom %s : %v", meta.Base, err)
		return
	}
	parsedDenom.Name = &meta.Display
	parsedDenom.Description = &meta.Description
	for _, denomUnit := range meta.DenomUnits {
		denomUnit := denomUnit
		if denomUnit.Denom == meta.Display {
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
			Network:   l.reader.Network,
			UpdatedAt: timestamppb.Now(),
			CreatedAt: timestamppb.Now(),
		},
	}
	cur, err := l.currencyClient.Get(ctx, &currency.ID{
		Network: l.reader.Network,
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
		return
	}
	c.MetaData.UpdatedAt = timestamppb.Now()
	_, err = l.currencyClient.Upsert(ctx, c)
	if err != nil {
		logger.Errorf("could not upsert denom %s : %v", meta.Base, err)
		return
	}
}

func (l *Application) currencyDenomScanner(ctx context.Context) error {
	var err error
	l.tokenRegistryEntries, err = dmn.GetTokenRegistryEntries(ctx, l.reader.Network)
	if err != nil {
		logger.Errorf("could not get token registry entries : %v", err)
	}

	var denomList types.Coins
	var paginationKey []byte = nil
	for {
		denomList, paginationKey, err = l.reader.QueryDenoms(ctx, l.bankClient[l.reader.Network], paginationKey)
		if err != nil {
			return err
		}

		for _, currentDenom := range denomList {
			l.currencyDenomScannerChan <- currentDenom
		}
		if paginationKey == nil {
			break
		}
	}
	return nil
}

func (l *Application) currencyDenomScannerPool(ctx context.Context) {
	for {
		select {
		case currentDenom := <-l.currencyDenomScannerChan:
			l.currencyDenomScannerInternal(ctx, currentDenom)
		}
	}
}

func (l *Application) currencyDenomScannerInternal(ctx context.Context, currentDenom types.Coin) {
	d, err := denomproto.NewDenom(currentDenom.Denom)
	if err != nil {
		logger.Errorf("could not parse denom %s : %v", currentDenom.Denom, err)
		return
	}
	c := &currency.Currency{
		Denom: d,
		MetaData: &metadata.MetaData{
			Network:   l.reader.Network,
			UpdatedAt: timestamppb.Now(),
			CreatedAt: timestamppb.Now(),
		},
	}
	cur, err := l.currencyClient.Get(ctx, &currency.ID{
		Network: l.reader.Network,
		Denom:   currentDenom.Denom,
	})
	if err != nil {
		logger.Infof("could not get denom %s from database, initializing new currency: %v", currentDenom.Denom, err)
	} else {
		c = cur
	}
	if c.MetaData == nil {
		c.MetaData = &metadata.MetaData{
			Network:   l.reader.Network,
			UpdatedAt: timestamppb.Now(),
			CreatedAt: timestamppb.Now(),
		}
	}
	c.ChainSupply = currentDenom.Amount.String()
	c.MetaData.UpdatedAt = timestamppb.Now()
	if token, ok := l.tokenRegistryEntries[currentDenom.Denom]; ok {
		tokenName := token.TokenName
		if c.Denom.Name == nil || *c.Denom.Name == "" {
			c.Denom.Name = &tokenName
		}
		tokenPrecision := int32(token.Decimals)
		if c.Denom.Precision == nil || *c.Denom.Precision == 0 {
			c.Denom.Precision = &tokenPrecision
		}
		tokenIcon := token.LogoURIs.Png
		if c.Denom.Icon == nil || *c.Denom.Icon == "" {
			c.Denom.Icon = &tokenIcon
		}
		tokenDescription := token.Description
		if c.Denom.Description == nil || *c.Denom.Description == "" {
			c.Denom.Description = &tokenDescription
		}
	}
	if c.Denom == nil {
		logger.Warnf("denom is nil for %s, unable to persist", currentDenom.Denom)
		return
	}
	_, err = l.currencyClient.Upsert(ctx, c)
	if err != nil {
		logger.Errorf("could not upsert denom %s : %v", currentDenom.Denom, err)
		return
	}
}
