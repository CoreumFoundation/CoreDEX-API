package grpc

import (
	"context"

	pb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/store/store"
	stategrpc "github.com/CoreumFoundation/CoreDEX-API/domain/state"
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

func (s *GrpcServer) Upsert(ctx context.Context, in *stategrpc.State) (*pb.Empty, error) {
	err := s.store.State.Upsert(ctx, in)
	if err != nil {
		logger.Errorf("State: Upsert failed for %+v with error %v", in, err)
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *GrpcServer) Get(ctx context.Context, in *stategrpc.StateQuery) (*stategrpc.State, error) {
	st, err := s.store.State.Get(ctx, in)
	if err != nil {
		logger.Errorf("Get failed for %s with error %v", in, err)
		return nil, err
	}
	return st, nil
}
