package market

import (
	"context"

	"github.com/CoreumFoundation/CoreDEX-API/coreum"
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
	logger.Infof("starting market scanner for %s", app.reader.Network.String())
	// Start the market scanner
	go app.scanMarkets(ctx)
}

func (app *Application) scanMarkets(ctx context.Context) {
	for {
		// Get the active markets:
		tps, err := app.tradeClient.GetTradePairs(ctx, &tradegrpc.TradePairFilter{})
		if err != nil {
			logger.Errorf("Error fetching trade pairs: %v", err)
			continue
		}
		for _, tp := range tps.TradePairs {
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
			ptf, _ := resp.PriceTick.Rat().Float64()
			tp.PriceTick = lo.ToPtr(int64(ptf))
			qts := resp.QuantityStep.Int64()
			tp.QuantityStep = lo.ToPtr(qts)
			_, err = app.tradeClient.UpsertTradePair(ctx, tp)
			if err != nil {
				logger.Errorf("Error upserting trade pair: %v. Error %v", tp, err)
				continue
			}
		}
	}
}
