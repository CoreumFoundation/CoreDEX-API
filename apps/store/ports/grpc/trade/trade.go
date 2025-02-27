package grpc

import (
	"context"

	pb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type GrpcServer struct {
	store *store.StoreBase
}

func NewGrpcServer(store *store.StoreBase) *GrpcServer {
	return &GrpcServer{
		store: store,
	}
}

func (s *GrpcServer) Upsert(ctx context.Context, in *tradegrpc.Trade) (*pb.Empty, error) {
	err := s.store.Trade.Upsert(in)
	if err != nil {
		logger.Errorf("Trade: Upsert failed for %s with error %v", *in.TXID, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) Get(ctx context.Context, in *tradegrpc.ID) (*tradegrpc.Trade, error) {
	st, err := s.store.Trade.Get(in)
	if err != nil {
		logger.Errorf("Get failed for %+v with error %v", *in, err)
		return nil, err
	}
	return st, nil
}

func (s *GrpcServer) BatchUpsert(ctx context.Context, in *tradegrpc.Trades) (*pb.Empty, error) {
	err := s.store.Trade.BatchUpsert(in)
	if err != nil {
		logger.Errorf("BatchUpsert for %+v failed with error %v", *in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) GetAll(ctx context.Context, in *tradegrpc.Filter) (*tradegrpc.Trades, error) {
	st, err := s.store.Trade.GetAll(in)
	if err != nil {
		logger.Errorf("GetAll with filter %+v failed with error %v", *in, err)
		return nil, err
	}
	return st, nil
}

func (s *GrpcServer) GetTradePairs(ctx context.Context, filter *tradegrpc.TradePairFilter) (*tradegrpc.TradePairs, error) {
	tradePairs, err := s.store.Trade.GetTradePairs(filter)
	if err != nil {
		logger.Errorf("GetTradePairs failed for %+v with error %v", *filter, err)
		return nil, err
	}
	return tradePairs, nil
}
