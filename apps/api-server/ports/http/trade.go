package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
)

type TradeOptionsFromParams struct {
	Symbol string `valid:"required~symbol.missing,symbol~symbol.invalid"`
	From   string `valid:"required~from.missing,unixtime~from.invalid"`
	To     string `valid:"required~to.missing,unixtime~to.invalid"`
}

func init() {
	govalidator.TagMap["symbol"] = func(str string) bool {
		return domain.ValidSymbol(str)
	}

	govalidator.TagMap["unixtime"] = func(str string) bool {
		return govalidator.IsUnixTime(str) &&
			govalidator.InRangeInt(str, 0, time.Now().Unix()+30) // Allow 30 seconds clock drift.
	}
}

func (s *httpServer) getTrades() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		network, err := networklib.Network(r)
		if err != nil {
			return handler.NewAPIError(401, "network.invalid")
		}
		opt, err := validateTradeParams(r.URL.Query())
		if err != nil {
			return handler.NewAPIError(422, err.Error())
		}
		opt.Network = network
		retvals, err := s.app.Trade.GetTrades(r.Context(), opt)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(retvals)
	}
}

// validateTradeParams validates the parameters for the trade endpoint.
// Returns a correct query, period and the original from timestamp.
// This originalFrom timestamp is used to prevent confusing the FE graph which does not seem to be able to handle anything different than what is outputs.
func validateTradeParams(query url.Values) (*tradegrpc.Filter, error) {
	beforeIDStr := query.Get("before_time")
	beforeID, err := strconv.ParseInt(beforeIDStr, 10, 64)
	if err != nil && beforeIDStr != "" {
		return nil, handler.NewAPIError(422, "before_id is not a valid integer")
	}
	// Translate the from to the period start from:
	afterIDStr := query.Get("after_time")
	afterID, err := strconv.ParseInt(afterIDStr, 10, 64)
	if err != nil && afterIDStr != "" {
		return nil, handler.NewAPIError(422, "external_id.invalid")
	}

	symbol := query.Get("symbol")
	account := query.Get("account")
	if account == "" && symbol == "" {
		return nil, handler.NewAPIError(422, "external_id.invalid")
	}
	tf := &tradegrpc.Filter{
		From: timestamppb.New(time.Unix(afterID, 0)),
		To:   timestamppb.New(time.Unix(beforeID, 0)),
	}
	if account != "" {
		tf.Account = &account
	}
	if symbol != "" {
		// Parse the symbol to see if it is valid:
		dmnSymbol, err := domain.NewSymbolFromString(symbol)
		if err != nil {
			return nil, handler.NewAPIError(422, "symbol.invalid")
		}
		b, err := denom.NewDenom(dmnSymbol.Base)
		if err != nil {
			return nil, handler.NewAPIError(422, "symbol.invalid")
		}
		q, err := denom.NewDenom(dmnSymbol.Quote)
		if err != nil {
			return nil, handler.NewAPIError(422, "symbol.invalid")
		}
		tf.Denom1 = b
		tf.Denom2 = q
	}
	return tf, nil
}
