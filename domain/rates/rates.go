package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradeclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

const (
	maxAge   = 15 * 60 // 15 minutes cache max
	BaseCoin = "BASE_COIN"
	BaseUSDC = "BASE_USDC"
)

type Fetcher struct {
	cl            *client
	exchangeRates ExchangeRates
	issuer        string
	currency      string
	tradeStore    tradegrpc.TradeServiceClient
	ohlcStore     ohlcgrpc.OHLCServiceClient
	graph         *weightedGraph
	usdc          string
	usdcIssuer    string
	network       metadata.Network
}

type Fetchers map[metadata.Network]*Fetcher // [Network]*Fetcher

// {\"BaseCoin\":[{{\"Network\": \"devnet\",\"Coin\": \"usara-devcore1wkwy0xh89ksdgj9hr347dyd2dw7zesmtrue6kfzyml4vdtz6e5wsyjwwgp\"}]}
type Issuers struct {
	BaseCoin []struct {
		Network string
		Coin    string
	}
}

func NewFetcher(tradeStore tradegrpc.TradeServiceClient, ohlcStore ohlcgrpc.OHLCServiceClient) *Fetchers {
	pi := os.Getenv(BaseCoin)
	if pi == "" {
		logger.Fatalf("%s env is required", BaseCoin)
	}
	d := Issuers{}
	err := json.Unmarshal([]byte(pi), &d)
	if err != nil {
		// FIXME: This error message is wrong
		logger.Fatalf("%s has to be set in format {\"BaseCoin\":[{{\"Network\": \"devnet\",\"Coin\": \"currency-issuer\"}]}", BaseCoin)
	}

	// Get the USDC base:
	// USDC tends to be only configured for mainnet, that is just to be accepted in the rest of the code: if there is no USDC, the code will just have to fall back to SARA/CORE for resolution of USD rates
	usdc := os.Getenv(BaseUSDC)
	if usdc == "" {
		logger.Fatalf("%s env is required", BaseUSDC)
	}
	usdcs := Issuers{}
	err = json.Unmarshal([]byte(usdc), &usdcs)
	if err != nil {
		logger.Fatalf("BaseUSDC %s has to be set in format {\"BaseCoin\":[{{\"Network\": \"mainnet\",\"Coin\": \"currency-issuer\"}]}", BaseUSDC)
	}

	// Split the issuer and currency to create an array of fetchers:
	f := make(Fetchers)
	for _, v := range d.BaseCoin {
		s := strings.Split(v.Coin, "-")
		if len(s) != 2 && (v.Coin != "ucore" && v.Coin != "utestcore" && v.Coin != "udevcore") {
			logger.Fatalf("Coin structure invalid: %s has to be set in format {\"BaseCoin\":[{{\"Network\": \"devnet\",\"Coin\": \"currency-issuer\"}]}", BaseCoin)
		}
		cur := s[0]
		iss := ""
		if len(s) == 2 {
			iss = s[1]
		}
		// metadata.Network is the key for the fetcher, string as input:
		nw := metadata.Network(metadata.Network_value[strings.ToUpper(v.Network)])
		f[nw] = &Fetcher{cl: newClient(tradeStore),
			issuer:     iss,
			currency:   cur,
			tradeStore: tradeStore,
			ohlcStore:  ohlcStore,
			network:    nw,
		}
		// Find the USDC base if present:
		for _, u := range usdcs.BaseCoin {
			if u.Network == v.Network {
				s := strings.Split(u.Coin, "-")
				if len(s) != 2 {
					logger.Fatalf("USDC Coin structure invalid: %s has to be set in format {\"BaseCoin\":[{{\"Network\": \"devnet\",\"Coin\": \"currency-issuer\"}]}", BaseUSDC)
				}
				f[nw].usdc = s[0]
				f[nw].usdcIssuer = s[1]
				break
			}
		}
		f[nw].initGraph()
	}

	return &f
}

func (f *Fetcher) ParseExchangeRate(ctx context.Context, denom1 *denom.Denom, quoteAsset string) (float64, error) {
	// Resolve the currency to use for retrieval:
	rateUSD, err := f.ParseExchangeRateVolume(ctx, denom1, quoteAsset, 1, 0)
	if err != nil {
		return 0.0, err
	}
	return rateUSD, nil
}

// ParseExchangeRate uses ParseExchangeRateVolume to calculate the exchange rate for a given currency and issuer, however multiplies the output by the transaction size
// (which ParseTradeExchangeRate does not do)
func (f *Fetcher) ParseExchangeRateVolume(ctx context.Context, denom1 *denom.Denom, quoteAsset string, val int64, exp int32) (float64, error) {
	r, err := f.ParseTradeExchangeRate(ctx, denom1, quoteAsset, val, exp)
	if err != nil {
		return 0.0, err
	}
	// transaction size is the decimal representing val and exp:
	d := decimal.New(val, exp)
	return r * d.InexactFloat64(), nil
}

func (f *Fetcher) ParseTradeExchangeRate(ctx context.Context, denom1 *denom.Denom, quoteAsset string, val int64, exp int32) (float64, error) {
	// The provided currency has to resolve to:
	/*
		USDC => Comes from env var BASE_USDC
		CORE => Is ucore by definition

		To be able to find a USD rate in the end
		Preferred order for resolution is
		USDC
		CORE
	*/
	// Currency path should always resolve to a value. If no value is possible as path, we can not convert the value and return 0
	currencyPath := getCurrencyPath(f.graph, denom1, key(&denom.Denom{Currency: f.usdc, Issuer: f.usdcIssuer}))
	if currencyPath == nil {
		logger.Warnf("No currency path for %s-%s to %s-%s", denom1.Currency, denom1.Issuer, f.usdc, f.usdcIssuer)
		return 0.0, nil
	}
	/* When there is a path, we can calculate the rate

	the last node should be usdc
	Every pair is a set of two nodes representing a tradepair (with an OHLC).
	Pools are alphabetically ordered by the currencies and issuer, so when converting from A -> B -> C that needs to be taken into account

	To resolve the whole process:
	Iterate over the path and calculate the rate for each pair
	*/
	var fRes float64 = 1.0
	for i := 0; i < len(currencyPath)-1; i++ { // len of array -1 since we need the next currencyPath to resolve the last pool in a potential list of pools for the calculation
		// Get the translate for the pair:
		denoms := strings.Split(currencyPath[i], "_")
		denom1, err := denom.NewDenom(denoms[0])
		if err != nil {
			return 0.0, err
		}
		denom2, err := denom.NewDenom(denoms[1])
		if err != nil {
			return 0.0, err
		}

		p, err := f.ohlcStore.Get(tradeclient.AuthCtx(ctx), &ohlcgrpc.OHLCFilter{
			Symbol:   currencyPath[i],
			Backfill: true,
			Period: &ohlcgrpc.Period{
				PeriodType: ohlcgrpc.PeriodType_PERIOD_TYPE_HOUR,
				Duration:   1,
			},
			Network: f.network,
			From:    timestamppb.Now(),
		})
		if err != nil {
			return 0.0, err
		}
		if len(p.OHLCs) != 1 {
			return 0.0, fmt.Errorf("no or incorrect OHLC for %s to %s", denom1.ToString(), denom2.ToString())
		}
		// Calculate the rate for the pair:
		logger.Infof("Looking to convert %s_%s", denom1.ToString(), denom2.ToString())
		fRes = fRes * p.OHLCs[0].Close
	}
	logger.Infof("Final conversion rate for %s to %s-%s is %f", denom1.ToString(), f.usdc, f.usdcIssuer, fRes)
	return fRes, nil
}

func Key(base, target string) string {
	return base + "-" + target
}
