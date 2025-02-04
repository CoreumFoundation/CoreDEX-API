package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	dmnsymbol "github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

// GoodTil is a good til order settings.
type GoodTil struct {
	// good_til_block_height means that order remains active until a specific blockchain block height is reached.
	GoodTilBlockHeight uint64
	// good_til_block_time means that order remains active until a specific blockchain block time is reached.
	GoodTilBlockTime *time.Time
}

type MsgPlaceOrderRequest struct {
	Sender      string
	Type        dextypes.OrderType
	OrderID     string `json:"ID,omitempty"`
	BaseDenom   string
	QuoteDenom  string
	Price       string
	Quantity    string
	Side        dextypes.Side
	GoodTil     *GoodTil             `json:"GoodTil,omitempty"`
	TimeInForce dextypes.TimeInForce `json:"TimeInForce,omitempty"`
}

type Order struct {
	TXBytes string
}

func (s *httpServer) createOrder() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		var orderReq MsgPlaceOrderRequest
		err := json.NewDecoder(r.Body).Decode(&orderReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		price, err := decimal.NewFromString(orderReq.Price)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		baseCurrency, err := s.app.Currency.GetCurrency(r.Context(), network, orderReq.BaseDenom)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		quoteCurrency, err := s.app.Currency.GetCurrency(r.Context(), network, orderReq.QuoteDenom)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		baseDenomPrecision := int64(0)
		if baseCurrency.Denom.Precision != nil {
			baseDenomPrecision = int64(*baseCurrency.Denom.Precision)
		}
		quoteDenomPrecision := int64(0)
		if quoteCurrency.Denom.Precision != nil {
			quoteDenomPrecision = int64(*quoteCurrency.Denom.Precision)
		}

		priceExponent := quoteDenomPrecision - baseDenomPrecision
		if orderReq.Side == dextypes.SIDE_BUY {
			priceExponent = baseDenomPrecision - quoteDenomPrecision
		}
		price = price.Mul(decimal.New(1, int32(priceExponent)))

		coreumPrice, err := coreum.ParsePrice(price.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		quantity, err := decimal.NewFromString(orderReq.Quantity)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		quantity = quantity.Mul(decimal.New(1, int32(baseDenomPrecision)))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		if quantity.Rat().Denom().Cmp(math.OneInt().BigInt()) != 0 {
			// entered quantity is outside the precision range
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		// Generate a UUID for the ID:
		orderReq.OrderID = uuid.New().String()
		msgPlaceOrder := dextypes.MsgPlaceOrder{
			Sender:      orderReq.Sender,
			Type:        orderReq.Type,
			ID:          orderReq.OrderID,
			BaseDenom:   orderReq.BaseDenom,
			QuoteDenom:  orderReq.QuoteDenom,
			Price:       &coreumPrice,
			Quantity:    math.NewIntFromBigInt(quantity.Rat().Num()),
			Side:        orderReq.Side,
			TimeInForce: orderReq.TimeInForce,
		}
		if orderReq.GoodTil != nil {
			msgPlaceOrder.GoodTil = &dextypes.GoodTil{
				GoodTilBlockHeight: orderReq.GoodTil.GoodTilBlockHeight,
				GoodTilBlockTime:   orderReq.GoodTil.GoodTilBlockTime,
			}
		}

		addr, err := types.AccAddressFromBech32(orderReq.Sender)
		if err != nil {
			return err
		}
		txBytes, err := s.app.Order.EncodeTx(network, addr, &msgPlaceOrder)
		if err != nil {
			logger.Errorf("Error encoding tx: %v", err)
			return err
		}
		// We do not want to handle bytes in the return, and we like to have a structured response:
		// base64 encode the txBytes:
		res := Order{
			TXBytes: base64.StdEncoding.EncodeToString(txBytes),
		}
		return json.NewEncoder(w).Encode(res)
	}
}

func (s *httpServer) cancelOrder() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		var orderReq struct {
			Sender  string
			OrderID string
		}
		err := json.NewDecoder(r.Body).Decode(&orderReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		msgCancelOrder := dextypes.MsgCancelOrder{
			Sender: orderReq.Sender,
			ID:     orderReq.OrderID,
		}
		addr, err := types.AccAddressFromBech32(orderReq.Sender)
		if err != nil {
			return err
		}
		txBytes, err := s.app.Order.EncodeTx(network, addr, &msgCancelOrder)
		if err != nil {
			logger.Infof("Error encoding cancel tx: %v", err)
			return err
		}
		// We do not want to handle bytes in the return, and we like to have a structured response:
		// base64 encode the txBytes:
		res := Order{
			TXBytes: base64.StdEncoding.EncodeToString(txBytes),
		}
		return json.NewEncoder(w).Encode(res)
	}
}

type rawTxBody struct {
	TX string
}

type SubmitResponse struct {
	TXHash string
}

func (s *httpServer) submitOrder() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		rawTx, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		txData := &rawTxBody{}
		err = json.Unmarshal(rawTx, txData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		rawTx, err = base64.StdEncoding.DecodeString(txData.TX)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		res, err := s.app.Order.SubmitTx(network, rawTx)
		if err != nil {
			logger.Errorf("Error submitting tx: %v", err)
			return err
		}
		submitResponse := SubmitResponse{
			TXHash: res.TxHash,
		}
		return json.NewEncoder(w).Encode(submitResponse)
	}
}

func (s *httpServer) getOrders() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		q := r.URL.Query()
		symbol := q.Get("symbol")
		if len(symbol) == 0 || !strings.Contains(symbol, "_") {
			w.WriteHeader(http.StatusBadRequest)
			return fmt.Errorf("symbol %q is not provided in the correct format", symbol)
		}
		denoms, err := dmnsymbol.NewSymbol(symbol)
		if err != nil {
			return fmt.Errorf("symbol %q is not provided in the correct format: %v", symbol, err)
		}
		limitStr := q.Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 50
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		account := q.Get("account")
		var res *coreum.OrderBookOrders
		if account == "" {
			res, err = s.app.Order.OrderBookRelevantOrders(network, denoms.Denom1.Denom, denoms.Denom2.Denom, limit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return err
			}
		} else {
			res, err = s.app.Order.OrderBookRelevantOrdersForAccount(network, denoms.Denom1.Denom, denoms.Denom2.Denom, account, limit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return err
			}
		}
		return json.NewEncoder(w).Encode(res)
	}
}
