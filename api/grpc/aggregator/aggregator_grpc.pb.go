// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: aggregator/aggregator.proto

package aggregator

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AggregatorClient is the client API for Aggregator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AggregatorClient interface {
	InitOperator(ctx context.Context, in *InitOperatorRequest, opts ...grpc.CallOption) (*InitOperatorResponse, error)
	CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponse, error)
	ProcessSignedTaskResponse(ctx context.Context, in *SignedTaskRespRequest, opts ...grpc.CallOption) (*SignedTaskRespResponse, error)
}

type aggregatorClient struct {
	cc grpc.ClientConnInterface
}

func NewAggregatorClient(cc grpc.ClientConnInterface) AggregatorClient {
	return &aggregatorClient{cc}
}

func (c *aggregatorClient) InitOperator(ctx context.Context, in *InitOperatorRequest, opts ...grpc.CallOption) (*InitOperatorResponse, error) {
	out := new(InitOperatorResponse)
	err := c.cc.Invoke(ctx, "/aggregator.Aggregator/InitOperator", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aggregatorClient) CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponse, error) {
	out := new(CreateTaskResponse)
	err := c.cc.Invoke(ctx, "/aggregator.Aggregator/CreateTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aggregatorClient) ProcessSignedTaskResponse(ctx context.Context, in *SignedTaskRespRequest, opts ...grpc.CallOption) (*SignedTaskRespResponse, error) {
	out := new(SignedTaskRespResponse)
	err := c.cc.Invoke(ctx, "/aggregator.Aggregator/ProcessSignedTaskResponse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AggregatorServer is the server API for Aggregator service.
// All implementations must embed UnimplementedAggregatorServer
// for forward compatibility
type AggregatorServer interface {
	InitOperator(context.Context, *InitOperatorRequest) (*InitOperatorResponse, error)
	CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponse, error)
	ProcessSignedTaskResponse(context.Context, *SignedTaskRespRequest) (*SignedTaskRespResponse, error)
	mustEmbedUnimplementedAggregatorServer()
}

// UnimplementedAggregatorServer must be embedded to have forward compatible implementations.
type UnimplementedAggregatorServer struct {
}

func (UnimplementedAggregatorServer) InitOperator(context.Context, *InitOperatorRequest) (*InitOperatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitOperator not implemented")
}
func (UnimplementedAggregatorServer) CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}
func (UnimplementedAggregatorServer) ProcessSignedTaskResponse(context.Context, *SignedTaskRespRequest) (*SignedTaskRespResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessSignedTaskResponse not implemented")
}
func (UnimplementedAggregatorServer) mustEmbedUnimplementedAggregatorServer() {}

// UnsafeAggregatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AggregatorServer will
// result in compilation errors.
type UnsafeAggregatorServer interface {
	mustEmbedUnimplementedAggregatorServer()
}

func RegisterAggregatorServer(s grpc.ServiceRegistrar, srv AggregatorServer) {
	s.RegisterService(&Aggregator_ServiceDesc, srv)
}

func _Aggregator_InitOperator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitOperatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AggregatorServer).InitOperator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aggregator.Aggregator/InitOperator",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AggregatorServer).InitOperator(ctx, req.(*InitOperatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Aggregator_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AggregatorServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aggregator.Aggregator/CreateTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AggregatorServer).CreateTask(ctx, req.(*CreateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Aggregator_ProcessSignedTaskResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignedTaskRespRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AggregatorServer).ProcessSignedTaskResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aggregator.Aggregator/ProcessSignedTaskResponse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AggregatorServer).ProcessSignedTaskResponse(ctx, req.(*SignedTaskRespRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Aggregator_ServiceDesc is the grpc.ServiceDesc for Aggregator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Aggregator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aggregator.Aggregator",
	HandlerType: (*AggregatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitOperator",
			Handler:    _Aggregator_InitOperator_Handler,
		},
		{
			MethodName: "CreateTask",
			Handler:    _Aggregator_CreateTask_Handler,
		},
		{
			MethodName: "ProcessSignedTaskResponse",
			Handler:    _Aggregator_ProcessSignedTaskResponse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aggregator/aggregator.proto",
}
