package currency

import (
	"context"

	pb "google.golang.org/protobuf/types/known/emptypb"

	store "github.com/CoreumFoundation/CoreDEX-API/apps/store/store"
	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
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

func (s *GrpcServer) Upsert(ctx context.Context, in *currencygrpc.Currency) (*pb.Empty, error) {
	err := s.store.Currency.Upsert(in)
	if err != nil {
		logger.Errorf("Currency: Upsert failed for %+v with error %v", *in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) Get(ctx context.Context, in *currencygrpc.ID) (*currencygrpc.Currency, error) {
	st, err := s.store.Currency.Get(in)
	if err != nil {
		logger.Errorf("Currency: Get failed for %+v with error %v", *in, err)
		return nil, err
	}
	return st, nil
}

func (s *GrpcServer) BatchUpsert(ctx context.Context, in *currencygrpc.Currencies) (*pb.Empty, error) {
	err := s.store.Currency.BatchUpsert(in)
	if err != nil {
		logger.Errorf("Currency: BatchUpsert for %+v failed with error %v", *in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) GetAll(ctx context.Context, in *currencygrpc.Filter) (*currencygrpc.Currencies, error) {
	st, err := s.store.Currency.GetAll(in)
	if err != nil {
		logger.Errorf("Currency: GetAll with filter %+v failed with error %v", *in, err)
		return nil, err
	}
	return st, nil
}
