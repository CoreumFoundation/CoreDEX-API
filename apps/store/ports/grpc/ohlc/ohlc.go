package ohlc

import (
	"context"

	pb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
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

func (s *GrpcServer) Upsert(ctx context.Context, in *ohlcgrpc.OHLC) (*pb.Empty, error) {
	err := s.store.OHLC.Upsert(in)
	if err != nil {
		logger.Errorf("OHLC: Upsert failed for %+v with error %v", *in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) Get(ctx context.Context, in *ohlcgrpc.OHLCFilter) (*ohlcgrpc.OHLCs, error) {
	st, err := s.store.OHLC.Get(in)
	if err != nil {
		logger.Errorf("Get failed for %+v with error %v", *in, err)
		return nil, err
	}
	return st, nil
}

func (s *GrpcServer) BatchUpsert(ctx context.Context, in *ohlcgrpc.OHLCs) (*pb.Empty, error) {
	err := s.store.OHLC.BatchUpsert(in)
	if err != nil {
		logger.Errorf("BatchUpsert for %+v failed with error %v", *in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) GetOHLCsForPeriods(ctx context.Context, in *ohlcgrpc.PeriodsFilter) (*ohlcgrpc.OHLCs, error) {
	st, err := s.store.OHLC.GetOHLCsForPeriods(in)
	if err != nil {
		logger.Errorf("GetSymbolsForPeriods failed for %+v with error %v", *in, err)
		return nil, err
	}
	return st, nil
}
