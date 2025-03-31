package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
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
	// Translate the from to the period start from:
	symbol := query.Get("symbol")
	account := query.Get("account")
	if account == "" && symbol == "" {
		return nil, handler.NewAPIError(422, "external_id.invalid")
	}
	tf := &tradegrpc.Filter{}
	from := query.Get("from")
	if from != "" {
		fr, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			return nil, handler.NewAPIError(422, "from.invalid")
		}
		tf.From = timestamppb.New(time.Unix(fr, 0))
	}
	to := query.Get("to")
	if to != "" {
		t, err := strconv.ParseInt(to, 10, 64)
		if err != nil {
			return nil, handler.NewAPIError(422, "to.invalid")
		}
		tf.To = timestamppb.New(time.Unix(t, 0))
	}
	if tf.From != nil && tf.From.AsTime().After(tf.To.AsTime()) {
		return nil, handler.NewAPIError(422, "from.after.to")
	}
	// Limit interval to 24hrs max to prevent overflows (in all reasonable scenarios)
	if tf.From != nil && tf.To != nil && tf.To.AsTime().Sub(tf.From.AsTime()) > 24*time.Hour {
		return nil, handler.NewAPIError(422, "interval.too.long")
	}
	side := query.Get("side")
	if side != "" {
		// Parse side into a valid trade side:
		sideInt, err := strconv.Atoi(side)
		if err != nil {
			return nil, handler.NewAPIError(422, "side.invalid")
		}
		// Parse into orderproperties.Side:
		tf.Side = lo.ToPtr(orderproperties.Side(sideInt))
		if *tf.Side == orderproperties.Side_SIDE_UNSPECIFIED {
			return nil, handler.NewAPIError(422, "side.invalid")
		}
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
