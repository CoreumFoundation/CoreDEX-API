package grpc

import (
	"google.golang.org/grpc"

	"github.com/CoreumFoundation/CoreDEX-API/apps/store/ports/grpc/currency"
	"github.com/CoreumFoundation/CoreDEX-API/apps/store/ports/grpc/ohlc"
	order "github.com/CoreumFoundation/CoreDEX-API/apps/store/ports/grpc/order"
	state "github.com/CoreumFoundation/CoreDEX-API/apps/store/ports/grpc/state"
	trade "github.com/CoreumFoundation/CoreDEX-API/apps/store/ports/grpc/trade"
	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	stategrpc "github.com/CoreumFoundation/CoreDEX-API/domain/state"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
)

type GrpcServer struct {
	Server *grpc.Server
	store  store.StoreBase
}

func NewGrpcServer() *GrpcServer {
	s := grpc.NewServer()
	g := &GrpcServer{
		Server: s,
	}
	storeClient := store.NewStore()
	stategrpc.RegisterStateServiceServer(s, state.NewGrpcServer(storeClient))
	ordergrpc.RegisterOrderServiceServer(s, order.NewGrpcServer(storeClient))
	tradegrpc.RegisterTradeServiceServer(s, trade.NewGrpcServer(storeClient))
	ohlcgrpc.RegisterOHLCServiceServer(s, ohlc.NewGrpcServer(storeClient))
	currencygrpc.RegisterCurrencyServiceServer(s, currency.NewGrpcServer(storeClient))
	return g
}
