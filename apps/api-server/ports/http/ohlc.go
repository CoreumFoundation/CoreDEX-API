package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"google.golang.org/protobuf/types/known/timestamppb"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
)

type OHLCOptionsFromParams struct {
	Symbol string `valid:"required~symbol.missing,symbol~symbol.invalid"`
	Period string `valid:"required~period.missing,in(1m|3m|5m|15m|30m|1h|3h|6h|12h|1d|3d|1w)~period.invalid"`
	From   string `valid:"required~from.missing,unixtime~from.invalid"`
	To     string `valid:"required~to.missing,unixtime~to.invalid"`
}

func init() {
	govalidator.TagMap["symbol"] = func(str string) bool {
		return dmn.ValidSymbol(str)
	}

	govalidator.TagMap["unixtime"] = func(str string) bool {
		return govalidator.IsUnixTime(str) &&
			govalidator.InRangeInt(str, 0, time.Now().Unix()+30) // Allow 30 seconds leeway.
	}
}

func (s *httpServer) getOHLC() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		network, err := networklib.Network(r)
		if err != nil {
			return err
		}
		ohlcOpt, period, err := validateOHLCParams(r.URL.Query())
		if err != nil {
			return handler.NewAPIError(422, err.Error())
		}
		ohlcOpt.Network = network
		ohlcOpt.Period = period
		retvals, err := s.app.OHLC.Get(r.Context(), ohlcOpt)
		if err != nil {
			return err
		}
		// To get a better response pattern, round the from and to to the period time periods.
		return json.NewEncoder(w).Encode(retvals)
	}
}

// validateOHLCParams validates the parameters for the OHLC endpoint.
// Returns a correct query, period and the original from timestamp.
// This originalFrom timestamp is used to prevent confusing the FE graph which does not seem to be able to handle anything different than what is outputs.
// The symnol parameter looks like:
// dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs
// {currency1}-{issuer}_{currency2}-{issuer}
func validateOHLCParams(query url.Values) (*ohlcgrpc.OHLCFilter, *ohlcgrpc.Period, error) {
	period, err := dmn.HttpPeriodToPeriod(query.Get("period"))
	if err != nil {
		return nil, nil, dmn.ErrIncorrectRequestParm
	}
	from, err := strconv.ParseInt(query.Get("from"), 10, 64)
	if err != nil {
		return nil, nil, handler.NewAPIError(422, "from is not a valid integer")
	}
	from = time.Unix(from, 0).UnixNano()
	// Translate the from to the period start from:
	queryFrom := period.ToOHLCKeyTimestampFrom(from)
	to, err := strconv.ParseInt(query.Get("to"), 10, 64)
	if err != nil {
		return nil, nil, handler.NewAPIError(422, "to is not a valid integer")
	}
	to = time.Unix(to, 0).UnixNano()
	to = period.ToOHLCKeyTimestampTo(to)

	symbol := query.Get("symbol")
	return &ohlcgrpc.OHLCFilter{
		Symbol: symbol,
		Period: period,
		From:   timestamppb.New(time.Unix(0, queryFrom)),
		To:     timestamppb.New(time.Unix(0, to)),
	}, period, nil
}
