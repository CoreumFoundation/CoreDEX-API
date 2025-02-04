package rates

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRates map[string]*ExchangeRate

type ExchangeRate struct {
	Base   string
	Target string

	Rate decimal.Decimal
	Time time.Time
}

func newExchangeRate(base, target string, rate decimal.Decimal, time time.Time) *ExchangeRate {
	return &ExchangeRate{
		Base:   normalizeAssetID(base),
		Target: normalizeAssetID(target),
		Rate:   rate,
		Time:   time,
	}
}

func normalizeAssetID(assetID string) string {
	return strings.ToUpper(assetID)
}
