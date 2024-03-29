// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: order.proto

package grpc

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

const (
	OrderService_ChangeOrderPrice_FullMethodName         = "/myapp.OrderService/ChangeOrderPrice"
	OrderService_ChangeMultipleOrderPrice_FullMethodName = "/myapp.OrderService/ChangeMultipleOrderPrice"
)

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderServiceClient interface {
	ChangeOrderPrice(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	ChangeMultipleOrderPrice(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (OrderService_ChangeMultipleOrderPriceClient, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) ChangeOrderPrice(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, OrderService_ChangeOrderPrice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) ChangeMultipleOrderPrice(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (OrderService_ChangeMultipleOrderPriceClient, error) {
	stream, err := c.cc.NewStream(ctx, &OrderService_ServiceDesc.Streams[0], OrderService_ChangeMultipleOrderPrice_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &orderServiceChangeMultipleOrderPriceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OrderService_ChangeMultipleOrderPriceClient interface {
	Recv() (*OrderResponse, error)
	grpc.ClientStream
}

type orderServiceChangeMultipleOrderPriceClient struct {
	grpc.ClientStream
}

func (x *orderServiceChangeMultipleOrderPriceClient) Recv() (*OrderResponse, error) {
	m := new(OrderResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OrderServiceServer is the server API for OrderService service.
// All implementations must embed UnimplementedOrderServiceServer
// for forward compatibility
type OrderServiceServer interface {
	ChangeOrderPrice(context.Context, *OrderRequest) (*OrderResponse, error)
	ChangeMultipleOrderPrice(*OrderRequest, OrderService_ChangeMultipleOrderPriceServer) error
	mustEmbedUnimplementedOrderServiceServer()
}

// UnimplementedOrderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServiceServer struct {
}

func (UnimplementedOrderServiceServer) ChangeOrderPrice(context.Context, *OrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeOrderPrice not implemented")
}
func (UnimplementedOrderServiceServer) ChangeMultipleOrderPrice(*OrderRequest, OrderService_ChangeMultipleOrderPriceServer) error {
	return status.Errorf(codes.Unimplemented, "method ChangeMultipleOrderPrice not implemented")
}
func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}

// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServiceServer will
// result in compilation errors.
type UnsafeOrderServiceServer interface {
	mustEmbedUnimplementedOrderServiceServer()
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_ChangeOrderPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).ChangeOrderPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_ChangeOrderPrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).ChangeOrderPrice(ctx, req.(*OrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_ChangeMultipleOrderPrice_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OrderRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OrderServiceServer).ChangeMultipleOrderPrice(m, &orderServiceChangeMultipleOrderPriceServer{stream})
}

type OrderService_ChangeMultipleOrderPriceServer interface {
	Send(*OrderResponse) error
	grpc.ServerStream
}

type orderServiceChangeMultipleOrderPriceServer struct {
	grpc.ServerStream
}

func (x *orderServiceChangeMultipleOrderPriceServer) Send(m *OrderResponse) error {
	return x.ServerStream.SendMsg(m)
}

// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "myapp.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ChangeOrderPrice",
			Handler:    _OrderService_ChangeOrderPrice_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ChangeMultipleOrderPrice",
			Handler:       _OrderService_ChangeMultipleOrderPrice_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "order.proto",
}
