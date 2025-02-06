package app

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptosecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/data-aggregator/domain"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/domain/order"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	"github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	dextypes "github.com/CoreumFoundation/coreum/v5/x/dex/types"
)

const testDataRoot = "test"

func TestApp(t *testing.T) {
	tests := []struct {
		name        string
		txFilenames []string
		wantTrades  []*trade.Trade
		wantOrders  []*order.Order
		wantErr     bool
	}{
		{
			name:        "limit order",
			txFilenames: []string{"test01_1_limit_order_created.json", "test01_2_limit_order_matched_and_closed.json"},
			wantTrades: []*trade.Trade{
				{
					TXID:     stringPtr("DC38A09467E6E621AFC5C08EDC4126B8382532C1444D53F4A7D4243F5FF4B646"),
					Account:  "devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
					OrderID:  "id1",
					Sequence: 3,
					Amount:   &decimal.Decimal{Value: 100, Exp: -6},
					Price:    10,
					Denom1: &denom.Denom{
						Currency:  "tknb755",
						Issuer:    "devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
						Precision: int32Ptr(6),
						IsIBC:     false,
						Denom:     "tknb755-devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
					},
					Denom2: &denom.Denom{
						Currency:  "tknb8b9",
						Issuer:    "devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
						Precision: int32Ptr(4),
						IsIBC:     false,
						Denom:     "tknb8b9-devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 468},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 468,
					USD:         nil,
					Enriched:    true,
				},
				{
					TXID:     stringPtr("DC38A09467E6E621AFC5C08EDC4126B8382532C1444D53F4A7D4243F5FF4B646"),
					Account:  "devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
					OrderID:  "id1",
					Sequence: 4,
					Amount:   &decimal.Decimal{Value: 100, Exp: -4},
					Price:    0.001,
					Denom1: &denom.Denom{
						Currency:  "tknb755",
						Issuer:    "devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
						Precision: int32Ptr(6),
						IsIBC:     false,
						Denom:     "tknb755-devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
					},
					Denom2: &denom.Denom{
						Currency:  "tknb8b9",
						Issuer:    "devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
						Precision: int32Ptr(4),
						IsIBC:     false,
						Denom:     "tknb8b9-devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
					},
					Side:      orderproperties.Side_SIDE_BUY,
					BlockTime: &timestamppb.Timestamp{Seconds: 468},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 468,
					USD:         nil,
					Enriched:    true,
				},
			},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("C6FA4DA015D6E72BEEAD566C6F331B0FB05A54A030084DBA4BF622E8F3E4C370"),
					Account:  "devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 3,
					BaseDenom: &denom.Denom{
						Currency:  "tknb755",
						Issuer:    "devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
						Precision: int32Ptr(6),
						IsIBC:     false,
						Denom:     "tknb755-devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tknb8b9",
						Issuer:    "devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
						Precision: int32Ptr(4),
						IsIBC:     false,
						Denom:     "tknb8b9-devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
					},
					Price:             10,
					Quantity:          &decimal.Decimal{Value: 100, Exp: -6},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 464},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 464,
					Enriched:    true,
				},
				{
					TXID:     stringPtr("DC38A09467E6E621AFC5C08EDC4126B8382532C1444D53F4A7D4243F5FF4B646"),
					Account:  "devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 4,
					BaseDenom: &denom.Denom{
						Currency:  "tknb755",
						Issuer:    "devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
						Precision: int32Ptr(6),
						IsIBC:     false,
						Denom:     "tknb755-devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tknb8b9",
						Issuer:    "devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
						Precision: int32Ptr(4),
						IsIBC:     false,
						Denom:     "tknb8b9-devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
					},
					Price:             0.0011,
					Quantity:          &decimal.Decimal{Value: 300, Exp: -6},
					RemainingQuantity: &decimal.Decimal{Value: 200, Exp: -6},
					Side:              orderproperties.Side_SIDE_BUY,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 468},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_OPEN,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 468,
					Enriched:    true,
				},
			},
			wantErr: false,
		},
		{
			name:        "market order",
			txFilenames: []string{"test02_1_market_order_created.json", "test02_2_market_order_matched_and_closed.json"},
			wantTrades: []*trade.Trade{
				{
					TXID:     stringPtr("B2685B9854E1C4E272FBEA5C9371771026621B5B89C955CD333DBED80C163FF9"),
					Account:  "devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
					OrderID:  "id1",
					Sequence: 4,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn3c74",
						Issuer:    "devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn3c74-devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
					},
					Denom2: &denom.Denom{
						Currency:  "tknc325",
						Issuer:    "devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tknc325-devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 103},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 103,
					USD:         nil,
				},
				{
					TXID:     stringPtr("B2685B9854E1C4E272FBEA5C9371771026621B5B89C955CD333DBED80C163FF9"),
					Account:  "devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
					OrderID:  "id2",
					Sequence: 5,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn3c74",
						Issuer:    "devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn3c74-devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
					},
					Denom2: &denom.Denom{
						Currency:  "tknc325",
						Issuer:    "devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tknc325-devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
					},
					Side:      orderproperties.Side_SIDE_BUY,
					BlockTime: &timestamppb.Timestamp{Seconds: 103},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 103,
					USD:         nil,
				},
			},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("3C8145893B12D63BAA4F999A99D976AF65936AD2EB63A9ABC72FEBC0F77BA00F"),
					Account:  "devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 4,
					BaseDenom: &denom.Denom{
						Currency:  "tkn3c74",
						Issuer:    "devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn3c74-devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tknc325",
						Issuer:    "devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tknc325-devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 100},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 99},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 99,
				},
				{
					TXID:     stringPtr("B2685B9854E1C4E272FBEA5C9371771026621B5B89C955CD333DBED80C163FF9"),
					Account:  "devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
					Type:     order.OrderType_ORDER_TYPE_MARKET,
					OrderID:  "id2",
					Sequence: 5,
					BaseDenom: &denom.Denom{
						Currency:  "tkn3c74",
						Issuer:    "devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn3c74-devcore1x6q0lr72atjxk3q5ve3nssmw8s8wylsnmx99np",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tknc325",
						Issuer:    "devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tknc325-devcore1tqmjf768ts93aju03qmjewajv5zuc5zfyucad2",
					},
					Quantity:          &decimal.Decimal{Value: 300},
					RemainingQuantity: &decimal.Decimal{Value: 200},
					Side:              orderproperties.Side_SIDE_BUY,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_UNSPECIFIED,
					BlockTime:         &timestamppb.Timestamp{Seconds: 103},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_OPEN,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 103,
				},
			},
			wantErr: false,
		},
		{
			name:        "limit order with good till height",
			txFilenames: []string{"test03_1_limit_order_with_good_till_height_created.json", "test03_2_limit_order_with_good_till_height_expired.json"},
			wantTrades:  []*trade.Trade{},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("4DEA9626B04A04409FB428F938F9E6AA3AA90319FB246CB6B9A4C187A7E8DC1E"),
					Account:  "devcore18w4cxydrcqftsde997prsrmv32sw65mydeqx8r",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 3,
					BaseDenom: &denom.Denom{
						Currency:  "tkn2827",
						Issuer:    "devcore18w4cxydrcqftsde997prsrmv32sw65mydeqx8r",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn2827-devcore18w4cxydrcqftsde997prsrmv32sw65mydeqx8r",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "denom2",
						Issuer:    "",
						Precision: nil,
						IsIBC:     false,
						Denom:     "denom2",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 100},
					RemainingQuantity: &decimal.Decimal{Value: 100},
					Side:              orderproperties.Side_SIDE_SELL,
					GoodTil: &order.GoodTil{
						BlockHeight: 180,
					},
					TimeInForce: order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:   &timestamppb.Timestamp{Seconds: 81},
					OrderStatus: order.OrderStatus_ORDER_STATUS_EXPIRED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 81,
				},
			},
			wantErr: false,
		},
		{
			name:        "multiple limit orders in one tx",
			txFilenames: []string{"test04_1_multiple_limit_orders_in_one_tx_sell1.json", "test04_2_multiple_limit_orders_in_one_tx_buy.json"},
			wantTrades: []*trade.Trade{
				{
					TXID:     stringPtr("1AC967655CC5F5EF61965A3A924A7FFBAE80B076E887C8C4444D6D9DF3D3143C"),
					Account:  "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					OrderID:  "id1",
					Sequence: 13,
					Amount:   &decimal.Decimal{Value: 50},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn2238",
						Issuer:    "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn2238-devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					},
					Denom2: &denom.Denom{
						Currency:  "tkn47ee",
						Issuer:    "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn47ee-devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 879},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 879,
					USD:         nil,
				},
				{
					TXID:     stringPtr("1AC967655CC5F5EF61965A3A924A7FFBAE80B076E887C8C4444D6D9DF3D3143C"),
					Account:  "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					OrderID:  "id2",
					Sequence: 14,
					Amount:   &decimal.Decimal{Value: 50},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn2238",
						Issuer:    "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn2238-devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					},
					Denom2: &denom.Denom{
						Currency:  "tkn47ee",
						Issuer:    "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn47ee-devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 879},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 879,
					USD:         nil,
				},
				{
					TXID:     stringPtr("1AC967655CC5F5EF61965A3A924A7FFBAE80B076E887C8C4444D6D9DF3D3143C"),
					Account:  "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					OrderID:  "id1",
					Sequence: 15,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn2238",
						Issuer:    "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn2238-devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					},
					Denom2: &denom.Denom{
						Currency:  "tkn47ee",
						Issuer:    "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn47ee-devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					},
					Side:      orderproperties.Side_SIDE_BUY,
					BlockTime: &timestamppb.Timestamp{Seconds: 879},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 879,
					USD:         nil,
				},
			},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("99A454ECD0728EABBC6503AE82856A8ADE7878B1B6FD21C89CBCC080FBA9EF54"),
					Account:  "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 13,
					BaseDenom: &denom.Denom{
						Currency:  "tkn2238",
						Issuer:    "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn2238-devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkn47ee",
						Issuer:    "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn47ee-devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 50},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 875},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 875,
				},
				{
					TXID:     stringPtr("99A454ECD0728EABBC6503AE82856A8ADE7878B1B6FD21C89CBCC080FBA9EF54"),
					Account:  "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id2",
					Sequence: 14,
					BaseDenom: &denom.Denom{
						Currency:  "tkn2238",
						Issuer:    "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn2238-devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkn47ee",
						Issuer:    "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn47ee-devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 50},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 875},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 875,
				},
				{
					TXID:     stringPtr("1AC967655CC5F5EF61965A3A924A7FFBAE80B076E887C8C4444D6D9DF3D3143C"),
					Account:  "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 15,
					BaseDenom: &denom.Denom{
						Currency:  "tkn2238",
						Issuer:    "devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn2238-devcore1zy8gf0ua96ndqwndhtzam2h25u6nam8a56estr",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkn47ee",
						Issuer:    "devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn47ee-devcore1n9wv5gxtrn42mks5vs2wly9qx3lj8c6h9yrpcg",
					},
					Price:             0.11,
					Quantity:          &decimal.Decimal{Value: 300},
					RemainingQuantity: &decimal.Decimal{Value: 200},
					Side:              orderproperties.Side_SIDE_BUY,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 879},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_OPEN,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 879,
				},
			},
			wantErr: false,
		},
		{
			name:        "multiple small limit orders matched by a large one",
			txFilenames: []string{"test05_1_multiple_small_limit_orders_matched_by_a_large_one_sell1.json", "test05_2_multiple_small_limit_orders_matched_by_a_large_one_sell2.json", "test05_3_multiple_small_limit_orders_matched_by_a_large_one_buy.json"},
			wantTrades: []*trade.Trade{
				{
					TXID:     stringPtr("E0033241D32C6E193B9DEA27BCBD305455E38A343CB28A307AABD190929A9685"),
					Account:  "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					OrderID:  "id1",
					Sequence: 16,
					Amount:   &decimal.Decimal{Value: 50},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn128c",
						Issuer:    "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn128c-devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					},
					Denom2: &denom.Denom{
						Currency:  "tkned51",
						Issuer:    "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkned51-devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 1235},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 1235,
					USD:         nil,
				},
				{
					TXID:     stringPtr("E0033241D32C6E193B9DEA27BCBD305455E38A343CB28A307AABD190929A9685"),
					Account:  "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					OrderID:  "id2",
					Sequence: 17,
					Amount:   &decimal.Decimal{Value: 50},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn128c",
						Issuer:    "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn128c-devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					},
					Denom2: &denom.Denom{
						Currency:  "tkned51",
						Issuer:    "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkned51-devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 1235},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 1235,
					USD:         nil,
				},
				{
					TXID:     stringPtr("E0033241D32C6E193B9DEA27BCBD305455E38A343CB28A307AABD190929A9685"),
					Account:  "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					OrderID:  "id1",
					Sequence: 18,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "tkn128c",
						Issuer:    "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn128c-devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					},
					Denom2: &denom.Denom{
						Currency:  "tkned51",
						Issuer:    "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkned51-devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					},
					Side:      orderproperties.Side_SIDE_BUY,
					BlockTime: &timestamppb.Timestamp{Seconds: 1235},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 1235,
					USD:         nil,
				},
			},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("93E2C5A86776EF4F35E003805F8F59296C81B42774D1CD2B1706A25D78603E20"),
					Account:  "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 16,
					BaseDenom: &denom.Denom{
						Currency:  "tkn128c",
						Issuer:    "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn128c-devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkned51",
						Issuer:    "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkned51-devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 50},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 1227},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 1227,
				},
				{
					TXID:     stringPtr("A8DAC0355E1884F9686666F87BF83DB22D0845B502502A90F34E3D2B83B91C84"),
					Account:  "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id2",
					Sequence: 17,
					BaseDenom: &denom.Denom{
						Currency:  "tkn128c",
						Issuer:    "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn128c-devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkned51",
						Issuer:    "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkned51-devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 50},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 1231},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 1231,
				},
				{
					TXID:     stringPtr("E0033241D32C6E193B9DEA27BCBD305455E38A343CB28A307AABD190929A9685"),
					Account:  "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 18,
					BaseDenom: &denom.Denom{
						Currency:  "tkn128c",
						Issuer:    "devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn128c-devcore1y0d0pkjz4ctkv29r04kcetdyet4gfy5t20qtf8",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkned51",
						Issuer:    "devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkned51-devcore1gl8c0au7kuv4y6c63y37whkzzkekk932en2g3r",
					},
					Price:             0.11,
					Quantity:          &decimal.Decimal{Value: 300},
					RemainingQuantity: &decimal.Decimal{Value: 200},
					Side:              orderproperties.Side_SIDE_BUY,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 1235},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_OPEN,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 1235,
				},
			},
			wantErr: false,
		},
		{
			name:        "one large limit order matched by multiple small ones",
			txFilenames: []string{"test06_1_one_large_limit_order_matched_by_multiple_small_ones_sell.json", "test06_2_one_large_limit_order_matched_by_multiple_small_ones_buy1.json", "test06_3_one_large_limit_order_matched_by_multiple_small_ones_buy2.json"},
			wantTrades: []*trade.Trade{
				{
					TXID:     stringPtr("B6F1973021319790729FBA5789F7CBC527CE24E624A4CE5CE67FF1905547B28C"),
					Account:  "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					OrderID:  "id1",
					Sequence: 7,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    1,
					Denom1: &denom.Denom{
						Currency:  "tkn0cea",
						Issuer:    "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn0cea-devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					},
					Denom2: &denom.Denom{
						Currency:  "tkn727c",
						Issuer:    "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn727c-devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 421},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 421,
					USD:         nil,
				},
				{
					TXID:     stringPtr("1D707A360FB8025C808EB558D682EBBD5D65BC535664DC9F29659E2F0B1E6E39"),
					Account:  "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					OrderID:  "id1",
					Sequence: 8,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    1,
					Denom1: &denom.Denom{
						Currency:  "tkn0cea",
						Issuer:    "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn0cea-devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					},
					Denom2: &denom.Denom{
						Currency:  "tkn727c",
						Issuer:    "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn727c-devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					},
					Side:      orderproperties.Side_SIDE_BUY,
					BlockTime: &timestamppb.Timestamp{Seconds: 417},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 417,
					USD:         nil,
				},
				{
					TXID:     stringPtr("B6F1973021319790729FBA5789F7CBC527CE24E624A4CE5CE67FF1905547B28C"),
					Account:  "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					OrderID:  "id2",
					Sequence: 9,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    1,
					Denom1: &denom.Denom{
						Currency:  "tkn0cea",
						Issuer:    "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn0cea-devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					},
					Denom2: &denom.Denom{
						Currency:  "tkn727c",
						Issuer:    "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn727c-devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					},
					Side:      orderproperties.Side_SIDE_BUY,
					BlockTime: &timestamppb.Timestamp{Seconds: 421},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 421,
					USD:         nil,
				},
			},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("779B4B5EBA59332C60C55D0C12D0FFA50737B90CAB0E2A5BFA016011F3B914E9"),
					Account:  "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 7,
					BaseDenom: &denom.Denom{
						Currency:  "tkn0cea",
						Issuer:    "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn0cea-devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkn727c",
						Issuer:    "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn727c-devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					},
					Price:             1,
					Quantity:          &decimal.Decimal{Value: 200},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 413},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 413,
				},
				{
					TXID:     stringPtr("1D707A360FB8025C808EB558D682EBBD5D65BC535664DC9F29659E2F0B1E6E39"),
					Account:  "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 8,
					BaseDenom: &denom.Denom{
						Currency:  "tkn0cea",
						Issuer:    "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn0cea-devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkn727c",
						Issuer:    "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn727c-devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					},
					Price:             1,
					Quantity:          &decimal.Decimal{Value: 100},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_BUY,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 417},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 417,
				},
				{
					TXID:     stringPtr("B6F1973021319790729FBA5789F7CBC527CE24E624A4CE5CE67FF1905547B28C"),
					Account:  "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id2",
					Sequence: 9,
					BaseDenom: &denom.Denom{
						Currency:  "tkn0cea",
						Issuer:    "devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn0cea-devcore1808eam6q6x4g3qxuwg5z4v3c64q4rwp43w0mjt",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkn727c",
						Issuer:    "devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn727c-devcore1vzzxmt20kyx42ygansq4r4heuwa8v3k5qsathz",
					},
					Price:             1,
					Quantity:          &decimal.Decimal{Value: 100},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_BUY,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 421},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 421,
				},
			},
			wantErr: false,
		},
		{
			name:        "cancel order",
			txFilenames: []string{"test07_1_limit_order_created_to_be_cancelled.json", "test07_2_cancel_order.json"},
			wantTrades:  []*trade.Trade{},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("CA18606228B7DD520CE1C6536DC5448CBF848EB8DE0079D2E25CAB049DCA28A0"),
					Account:  "devcore1zf86ygwf2tkwevl46f4vgydvexvym6mrk5fzxh",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "id1",
					Sequence: 23,
					BaseDenom: &denom.Denom{
						Currency:  "tkn02f0",
						Issuer:    "devcore1jtln26lek30ktrddkhnkt2n606k5qd3jk8pydg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn02f0-devcore1jtln26lek30ktrddkhnkt2n606k5qd3jk8pydg",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "tkn4243",
						Issuer:    "devcore1jtln26lek30ktrddkhnkt2n606k5qd3jk8pydg",
						Precision: nil,
						IsIBC:     false,
						Denom:     "tkn4243-devcore1jtln26lek30ktrddkhnkt2n606k5qd3jk8pydg",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 100},
					RemainingQuantity: &decimal.Decimal{Value: 100},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 4012},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_CANCELED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 4012,
				},
			},
			wantErr: false,
		},
		{
			name:        "limit order matched by opposite order book with the same side",
			txFilenames: []string{"test08_1_limit_order_created.json", "test08_2_limit_order_matched_by_opposite_order_book_with_the_same_side.json"},
			wantTrades: []*trade.Trade{
				{
					TXID:     stringPtr("EBA0B7C47E5C2CB76B0EFEB5CA901C0089CA7D89AFDFB67505C83FA3D17B3A26"),
					Account:  "devcore1vv97hnfxtar7a87szstpdlzexwu009cylfdsry",
					OrderID:  "854173c0-2831-4e51-ab5c-b64b823498c3",
					Sequence: 6167,
					Amount:   &decimal.Decimal{Value: 100},
					Price:    0.1,
					Denom1: &denom.Denom{
						Currency:  "dextestdenom1",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					Denom2: &denom.Denom{
						Currency:  "dextestdenom0",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 3246325},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 3246325,
					USD:         nil,
				},
				{
					TXID:     stringPtr("EBA0B7C47E5C2CB76B0EFEB5CA901C0089CA7D89AFDFB67505C83FA3D17B3A26"),
					Account:  "devcore1yegtec4a0zd9ksepr6hf8au95940p5hvy06064",
					OrderID:  "b7ebc57e-1cde-4ebc-925d-bcfd3609a166",
					Sequence: 6168,
					Amount:   &decimal.Decimal{Value: 10},
					Price:    10,
					Denom1: &denom.Denom{
						Currency:  "dextestdenom0",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					Denom2: &denom.Denom{
						Currency:  "dextestdenom1",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					Side:      orderproperties.Side_SIDE_SELL,
					BlockTime: &timestamppb.Timestamp{Seconds: 3246325},
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 3246325,
					USD:         nil,
				},
			},
			wantOrders: []*order.Order{
				{
					TXID:     stringPtr("834437BE045502656B6801090D31ECBB1DD753851AE3C896564566BBADDA19AE"),
					Account:  "devcore1vv97hnfxtar7a87szstpdlzexwu009cylfdsry",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "854173c0-2831-4e51-ab5c-b64b823498c3",
					Sequence: 6167,
					BaseDenom: &denom.Denom{
						Currency:  "dextestdenom1",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "dextestdenom0",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					Price:             0.1,
					Quantity:          &decimal.Decimal{Value: 100},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 3244942},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 3244942,
				},
				{
					TXID:     stringPtr("EBA0B7C47E5C2CB76B0EFEB5CA901C0089CA7D89AFDFB67505C83FA3D17B3A26"),
					Account:  "devcore1yegtec4a0zd9ksepr6hf8au95940p5hvy06064",
					Type:     order.OrderType_ORDER_TYPE_LIMIT,
					OrderID:  "b7ebc57e-1cde-4ebc-925d-bcfd3609a166",
					Sequence: 6168,
					BaseDenom: &denom.Denom{
						Currency:  "dextestdenom0",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					QuoteDenom: &denom.Denom{
						Currency:  "dextestdenom1",
						Issuer:    "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
						Precision: nil,
						IsIBC:     false,
						Denom:     "dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
					},
					Price:             10,
					Quantity:          &decimal.Decimal{Value: 10},
					RemainingQuantity: &decimal.Decimal{Value: 0},
					Side:              orderproperties.Side_SIDE_SELL,
					TimeInForce:       order.TimeInForce_TIME_IN_FORCE_GTC,
					BlockTime:         &timestamppb.Timestamp{Seconds: 3246325},
					OrderStatus:       order.OrderStatus_ORDER_STATUS_FILLED,
					MetaData: &metadata.MetaData{
						Network:   metadata.Network_DEVNET,
						UpdatedAt: timestamppb.Now(),
						CreatedAt: timestamppb.Now(),
					},
					BlockHeight: 3246325,
				},
			},
			wantErr: false,
		},
	}

	ctx := context.Background()

	interfaceRegistry := ctypes.NewInterfaceRegistry()
	dextypes.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &cryptosecp256k1.PubKey{})

	for _, tt := range tests {
		network := metadata.Network_DEVNET
		orderService := order.NewMockOrderServiceClient()
		tradeService := trade.NewMockTradeServiceClient()
		currencyService := currency.NewMockCurrencyServiceClient()
		agg := NewApplicationWithClients(ctx, orderService, tradeService, currencyService)

		_, err := currencyService.Upsert(ctx, &currency.Currency{
			Denom: &denom.Denom{
				Currency:  "tknb755",
				Issuer:    "devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
				Precision: int32Ptr(6),
				Denom:     "tknb755-devcore1j9vqw8mgc8zdtavjv786z969570k8dsq83uj24",
				Name:      stringPtr("tknb755"),
			},
			MetaData: &metadata.MetaData{
				Network:   network,
				UpdatedAt: timestamppb.Now(),
				CreatedAt: timestamppb.Now(),
			},
		})
		require.NoError(t, err)

		_, err = currencyService.Upsert(ctx, &currency.Currency{
			Denom: &denom.Denom{
				Currency:  "tknb8b9",
				Issuer:    "devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
				Precision: int32Ptr(4),
				Denom:     "tknb8b9-devcore1ydxhq3ccfv70gkrlc2jygh40v4czktnf679tp2",
				Name:      stringPtr("tknb8b9"),
			},
			MetaData: &metadata.MetaData{
				Network:   network,
				UpdatedAt: timestamppb.Now(),
				CreatedAt: timestamppb.Now(),
			},
		})
		require.NoError(t, err)

		t.Run(tt.name, func(t *testing.T) {
			for i, txFilename := range tt.txFilenames {
				file, err := filepath.Abs(filepath.Join(testDataRoot, txFilename))
				require.NoError(t, err)

				b, err := os.ReadFile(file)
				require.NoError(t, err)
				tx := &dmn.Result{}
				require.NoError(t, json.Unmarshal(b, tx))

				block := &coreum.ScannedBlock{
					Transactions: []*txtypes.GetTxResponse{
						{
							TxResponse: &types.TxResponse{},
						},
					},
					BlockHeight: int64(i + 1),
					BlockTime:   time.Unix(int64(i+1), 0),
				}
				if tx.Data.Value.TxResult != nil {
					height, err := strconv.ParseInt(tx.Data.Value.TxResult.Height, 10, 64)
					require.NoError(t, err)
					block.BlockHeight = height
					block.BlockTime = time.Unix(height, 0)

					txBytes, err := base64.StdEncoding.DecodeString(tx.Data.Value.TxResult.Tx)
					require.NoError(t, err)

					txHash := hash(txBytes)

					protoCodec := codec.NewProtoCodec(interfaceRegistry)
					transaction := &txtypes.Tx{}
					err = protoCodec.Unmarshal(txBytes, transaction)
					require.NoError(t, err)

					gasUsed, err := strconv.ParseInt(tx.Data.Value.TxResult.Result.GasUsed, 10, 64)
					require.NoError(t, err)

					block.Transactions[0].Tx = transaction
					block.Transactions[0].TxResponse.Events = tx.Data.Value.TxResult.Result.Events
					block.Transactions[0].TxResponse.TxHash = txHash
					block.Transactions[0].TxResponse.GasUsed = gasUsed
				}
				if tx.Data.Value.ResultFinalizeBlock != nil {
					block.BlockEvents = tx.Data.Value.ResultFinalizeBlock.Events
				}
				agg.scannerCoordinator(nil, block, network)
			}

			orders, err := orderService.GetAll(ctx, nil)
			require.NoError(t, err)

			trades, err := tradeService.GetAll(ctx, nil)
			require.NoError(t, err)

			assertOrdersEquality(t, tt.wantOrders, orders.Orders)
			assertTradesEquality(t, tt.wantTrades, trades.Trades)
		})
	}
}

func stringPtr(input string) *string {
	return &input
}

func int32Ptr(input int32) *int32 {
	return &input
}

func hash(txRaw []byte) string {
	h := sha256.New()
	h.Write(txRaw)
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}

func assertTradesEquality(t *testing.T, expected, actual []*trade.Trade) {
	// Compare with approximate time.
	cmpOpt := []cmp.Option{
		cmpopts.EquateApproxTime(3 * time.Second),
		cmpopts.IgnoreFields(timestamppb.Timestamp{}, "Nanos"),
		cmpopts.IgnoreUnexported(trade.Trade{}, decimal.Decimal{}, denom.Denom{}, timestamppb.Timestamp{}, metadata.MetaData{}, order.GoodTil{}),
	}

	require.True(t, cmp.Equal(expected, actual, cmpOpt...), cmp.Diff(expected, actual, cmpOpt...))
}

func assertOrdersEquality(t *testing.T, expected, actual []*order.Order) {
	// Compare with approximate time.
	cmpOpt := []cmp.Option{
		cmpopts.EquateApproxTime(3 * time.Second),
		cmpopts.IgnoreFields(timestamppb.Timestamp{}, "Nanos"),
		cmpopts.IgnoreUnexported(order.Order{}, decimal.Decimal{}, denom.Denom{}, timestamppb.Timestamp{}, metadata.MetaData{}, order.GoodTil{}),
	}

	require.True(t, cmp.Equal(expected, actual, cmpOpt...), cmp.Diff(expected, actual, cmpOpt...))
}
