package grpc

import (
	"context"

	pb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store"
	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
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

func (s *GrpcServer) Upsert(ctx context.Context, in *ordergrpc.Order) (*pb.Empty, error) {
	err := s.store.Order.Upsert(in)
	if err != nil {
		logger.Errorf("Order: Upsert failed for %+v with error %v", *in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) Get(ctx context.Context, in *ordergrpc.ID) (*ordergrpc.Order, error) {
	st, err := s.store.Order.Get(in)
	if err != nil {
		logger.Errorf("Get failed for %+v with error %v", *in, err)
		return nil, err
	}
	return st, nil
}

func (s *GrpcServer) GetAll(ctx context.Context, in *ordergrpc.Filter) (*ordergrpc.Orders, error) {
	st, err := s.store.Order.GetAll(in)
	if err != nil {
		logger.Errorf("GetAll with filter %+v failed with error %v", *in, err)
		return nil, err
	}
	return st, nil
}

func (s *GrpcServer) BatchUpsert(ctx context.Context, in *ordergrpc.Orders) (*pb.Empty, error) {
	err := s.store.Order.BatchUpsert(in)
	if err != nil {
		logger.Errorf("BatchUpsert for %+v failed with error %v", *in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}
