package market

import (
	"context"
	"time"

	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	"github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
	"github.com/samber/lo"
)

type Application struct {
	reader      *coreum.Reader
	tradeClient tradegrpc.TradeServiceClient
}

func NewApplication(reader *coreum.Reader, tradeClient tradegrpc.TradeServiceClient) *Application {
	return &Application{
		reader:      reader,
		tradeClient: tradeClient,
	}
}

func (app *Application) Start(ctx context.Context) {
	// Start the market scanner
	// TODO: Hash code back in when issue with coreum package is resolved
	go app.scanMarkets(ctx, app.reader.Network)
	logger.Infof("Started market scanner for %s", app.reader.Network.String())
}

func (app *Application) scanMarkets(ctx context.Context, network metadata.Network) {
	for {
		// Get the active markets:
		retrieveRecords := true
		var offset int32 = 0
		tps := make([]*tradegrpc.TradePair, 0)
		for retrieveRecords {
			pairs, err := app.tradeClient.GetTradePairs(ctx, &tradegrpc.TradePairFilter{Network: network, Offset: &offset})
			if err != nil {
				logger.Errorf("Error fetching trade pairs: %v", err)
				retrieveRecords = false
				continue
			}
			retrieveRecords = false
			if pairs.Offset != nil && *pairs.Offset > 0 {
				offset = *pairs.Offset
				retrieveRecords = true
			}
			if len(pairs.TradePairs) > 0 {
				tps = append(tps, pairs.TradePairs...)
			}
		}
		if len(tps) == 0 {
			time.Sleep(30 * time.Minute)
			continue
		}
		for _, tp := range tps {
			dexClient := dextypes.NewQueryClient(app.reader.ClientContext)
			resp, err := dexClient.OrderBookParams(ctx, &dextypes.QueryOrderBookParamsRequest{
				BaseDenom:  tp.Denom1.Denom,
				QuoteDenom: tp.Denom2.Denom,
			})
			if err != nil {
				logger.Errorf("Error fetching order book params: %v", err)
				continue
			}
			// Process the order book params into the tradepair:
			f, _ := resp.PriceTick.Rat().Float64()
			ptf := decimal.FromFloat64(f)
			tp.PriceTick = ptf
			qts := resp.QuantityStep.BigInt()
			if qts == nil {
				logger.Errorf("Error fetching quantity step for trade pair %s-%s: %v", tp.Denom1.Denom, tp.Denom2.Denom, err)
				continue
			}
			qtsf, _ := qts.Float64()
			tp.QuantityStep = lo.ToPtr(int64(qtsf))
			_, err = app.tradeClient.UpsertTradePair(ctx, tp)
			if err != nil {
				logger.Errorf("Error upserting trade pair: %v. Error %v", tp, err)
				continue
			}
		}
		time.Sleep(30 * time.Minute)
	}
}
