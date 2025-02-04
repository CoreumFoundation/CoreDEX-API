// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: domain/state/state-grpc.proto

package state

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StateServiceClient is the client API for StateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StateServiceClient interface {
	Upsert(ctx context.Context, in *State, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Get(ctx context.Context, in *StateQuery, opts ...grpc.CallOption) (*State, error)
}

type stateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStateServiceClient(cc grpc.ClientConnInterface) StateServiceClient {
	return &stateServiceClient{cc}
}

func (c *stateServiceClient) Upsert(ctx context.Context, in *State, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/state.StateService/Upsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateServiceClient) Get(ctx context.Context, in *StateQuery, opts ...grpc.CallOption) (*State, error) {
	out := new(State)
	err := c.cc.Invoke(ctx, "/state.StateService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StateServiceServer is the server API for StateService service.
// All implementations should embed UnimplementedStateServiceServer
// for forward compatibility
type StateServiceServer interface {
	Upsert(context.Context, *State) (*emptypb.Empty, error)
	Get(context.Context, *StateQuery) (*State, error)
}

// UnimplementedStateServiceServer should be embedded to have forward compatible implementations.
type UnimplementedStateServiceServer struct {
}

func (UnimplementedStateServiceServer) Upsert(context.Context, *State) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upsert not implemented")
}
func (UnimplementedStateServiceServer) Get(context.Context, *StateQuery) (*State, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

// UnsafeStateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StateServiceServer will
// result in compilation errors.
type UnsafeStateServiceServer interface {
	mustEmbedUnimplementedStateServiceServer()
}

func RegisterStateServiceServer(s grpc.ServiceRegistrar, srv StateServiceServer) {
	s.RegisterService(&StateService_ServiceDesc, srv)
}

func _StateService_Upsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(State)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServiceServer).Upsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/state.StateService/Upsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServiceServer).Upsert(ctx, req.(*State))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StateQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/state.StateService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServiceServer).Get(ctx, req.(*StateQuery))
	}
	return interceptor(ctx, in, info, handler)
}

// StateService_ServiceDesc is the grpc.ServiceDesc for StateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "state.StateService",
	HandlerType: (*StateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upsert",
			Handler:    _StateService_Upsert_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _StateService_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "domain/state/state-grpc.proto",
}
