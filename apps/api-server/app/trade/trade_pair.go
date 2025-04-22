package trade

import (
	"context"
	"fmt"

	metadata "github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	dmnsymbol "github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
	dmntrade "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
)

func (app *Application) GetMarket(ctx context.Context, symbol *dmnsymbol.Symbol, network metadata.Network) (*dmntrade.TradePair, error) {
	tps, err := app.tradeClient.GetTradePairs(ctx, &dmntrade.TradePairFilter{
		Denom1: symbol.Denom1,
		Denom2: symbol.Denom2,
	})
	if err != nil {
		return nil, err
	}
	if len(tps.TradePairs) == 0 || len(tps.TradePairs) > 1 {
		return nil, fmt.Errorf("trade pair not found")
	}
	return tps.TradePairs[0], nil
}
