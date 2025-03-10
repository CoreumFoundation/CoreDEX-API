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
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	networklib "github.com/CoreumFoundation/CoreDEX-API/domain/network"
	dmnsymbol "github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

type rawTxBody struct {
	TX string
}

type SubmitResponse struct {
	TXHash string
}

// GoodTil is a good til order settings.
type GoodTil struct {
	// good_til_block_height means that order remains active until a specific blockchain block height is reached.
	GoodTilBlockHeight uint64 `json:"goodTilBlockHeight,omitempty"`
	// good_til_block_time means that order remains active until a specific blockchain block time is reached.
	GoodTilBlockTime *time.Time `json:"goodTilBlockTime,omitempty"`
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
	GoodTil     *GoodTil             `json:"goodTil,omitempty"`
	TimeInForce dextypes.TimeInForce `json:"TimeInForce,omitempty"`
}

type OrderResponse struct {
	Sequence  uint64
	OrderData OrderData
}

type OrderData struct {
	dextypes.MsgPlaceOrder
	BaseDenom   string               `json:"baseDenom"`
	QuoteDenom  string               `json:"quoteDenom"`
	TimeInForce dextypes.TimeInForce `json:"timeInForce"`
	GoodTil     *GoodTil             `json:"goodTil,omitempty"`
}

type OrderCancelResponse struct {
	Sequence    uint64
	OrderCancel dextypes.MsgCancelOrder
}

func (s *httpServer) createOrder() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		var orderReq MsgPlaceOrderRequest
		err := json.NewDecoder(r.Body).Decode(&orderReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		price := decimal.NewFromInt(0)
		if len(orderReq.Price) > 0 {
			price, err = decimal.NewFromString(orderReq.Price)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
		}
		var coreumPrice *dextypes.Price = nil
		if !price.IsZero() {
			parsedCoreumPrice, err := coreum.ParsePrice(price.String())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return nil
			}
			coreumPrice = &parsedCoreumPrice
		}
		baseCurrency, err := s.app.Currency.GetCurrency(r.Context(), network, orderReq.BaseDenom)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}
		baseDenomPrecision := int64(0)
		if baseCurrency.Denom != nil && baseCurrency.Denom.Precision != nil {
			baseDenomPrecision = int64(*baseCurrency.Denom.Precision)
		}

		quantity, err := decimal.NewFromString(orderReq.Quantity)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		quantity = quantity.Mul(decimal.New(1, int32(baseDenomPrecision)))
		if quantity.Rat().Denom().Cmp(math.OneInt().BigInt()) != 0 {
			// entered quantity is outside the precision range
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		// Generate a UUID for the ID:
		orderReq.OrderID = uuid.New().String()
		msgPlaceOrder := dextypes.MsgPlaceOrder{
			Sender:      orderReq.Sender,
			Type:        orderReq.Type,
			ID:          orderReq.OrderID,
			BaseDenom:   orderReq.BaseDenom,
			QuoteDenom:  orderReq.QuoteDenom,
			Price:       coreumPrice,
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
		sequence, err := s.app.Order.AccountSequence(network, orderReq.Sender)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		o := OrderData{
			MsgPlaceOrder: msgPlaceOrder,
			BaseDenom:     orderReq.BaseDenom,
			QuoteDenom:    orderReq.QuoteDenom,
			TimeInForce:   orderReq.TimeInForce,
			GoodTil:       orderReq.GoodTil,
		}
		return json.NewEncoder(w).Encode(OrderResponse{
			Sequence:  sequence,
			OrderData: o,
		})
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
			return nil
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		sequence, err := s.app.Order.AccountSequence(network, orderReq.Sender)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		msgCancelOrder := dextypes.MsgCancelOrder{
			Sender: orderReq.Sender,
			ID:     orderReq.OrderID,
		}
		return json.NewEncoder(w).Encode(OrderCancelResponse{
			Sequence:    sequence,
			OrderCancel: msgCancelOrder,
		})
	}
}

func (s *httpServer) submitOrder() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		rawTx, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		network, err := networklib.Network(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		txData := &rawTxBody{}
		err = json.Unmarshal(rawTx, txData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		rawTx, err = base64.StdEncoding.DecodeString(txData.TX)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		res, err := s.app.Order.SubmitTx(network, rawTx)
		if err != nil {
			logger.Errorf("Error submitting tx: %v", err)
			return nil
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
			return nil
		}
		account := q.Get("account")
		var res *coreum.OrderBookOrders
		if account == "" {
			res, err = s.app.Order.OrderBookRelevantOrders(network, denoms.Denom1.Denom, denoms.Denom2.Denom, limit, true)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return nil
			}
		} else {
			res, err = s.app.Order.OrderBookRelevantOrdersForAccount(network, denoms.Denom1.Denom, denoms.Denom2.Denom, account)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return nil
			}
		}
		return json.NewEncoder(w).Encode(res)
	}
}
