package rates

import (
	"net/http"
	"net/url"
	"time"

	"github.com/shopspring/decimal"

	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
)

const defaultAPIRootURL = "https://api.coingecko.com"

// Coingecko client doesn't have public interface and is not supposed to be
// used outside of coingecko package for now.
type client struct {
	apiRootURL *url.URL
	httpClient *http.Client
	tradeStore tradegrpc.TradeServiceClient
}

func newClient(tradeStore tradegrpc.TradeServiceClient) *client {
	httpCl := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     50,
		},
	}

	defaultUrl, err := url.Parse(defaultAPIRootURL)
	if err != nil {
		panic(err)
	}

	return &client{
		apiRootURL: defaultUrl,
		httpClient: httpCl,
		tradeStore: tradeStore,
	}
}

type coinsMarketsRequestParams struct {
	VsCurrency string
	ID         string
}

type coinsMarketsResponseEntity struct {
	ID           string          `json:"id"`
	Symbol       string          `json:"symbol"`
	Name         string          `json:"name"`
	CurrentPrice decimal.Decimal `json:"current_price"`
	LastUpdated  time.Time       `json:"last_updated"`
	// NOTE:
	// Only required response fields are defined here.
	// If needed other can be added later in the future.
}
