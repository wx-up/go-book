// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: inter/interactive.proto

// buf:lint:ignore PACKAGE_VERSION_SUFFIX

package inter

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
	InteractiveService_IncrReadCnt_FullMethodName = "/inter.InteractiveService/IncrReadCnt"
	InteractiveService_Get_FullMethodName         = "/inter.InteractiveService/Get"
)

// InteractiveServiceClient is the client API for InteractiveService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InteractiveServiceClient interface {
	IncrReadCnt(ctx context.Context, in *IncrReadCntRequest, opts ...grpc.CallOption) (*IncrReadCntResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type interactiveServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInteractiveServiceClient(cc grpc.ClientConnInterface) InteractiveServiceClient {
	return &interactiveServiceClient{cc}
}

func (c *interactiveServiceClient) IncrReadCnt(ctx context.Context, in *IncrReadCntRequest, opts ...grpc.CallOption) (*IncrReadCntResponse, error) {
	out := new(IncrReadCntResponse)
	err := c.cc.Invoke(ctx, InteractiveService_IncrReadCnt_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *interactiveServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, InteractiveService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InteractiveServiceServer is the server API for InteractiveService service.
// All implementations must embed UnimplementedInteractiveServiceServer
// for forward compatibility
type InteractiveServiceServer interface {
	IncrReadCnt(context.Context, *IncrReadCntRequest) (*IncrReadCntResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	mustEmbedUnimplementedInteractiveServiceServer()
}

// UnimplementedInteractiveServiceServer must be embedded to have forward compatible implementations.
type UnimplementedInteractiveServiceServer struct {
}

func (UnimplementedInteractiveServiceServer) IncrReadCnt(context.Context, *IncrReadCntRequest) (*IncrReadCntResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IncrReadCnt not implemented")
}
func (UnimplementedInteractiveServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedInteractiveServiceServer) mustEmbedUnimplementedInteractiveServiceServer() {}

// UnsafeInteractiveServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InteractiveServiceServer will
// result in compilation errors.
type UnsafeInteractiveServiceServer interface {
	mustEmbedUnimplementedInteractiveServiceServer()
}

func RegisterInteractiveServiceServer(s grpc.ServiceRegistrar, srv InteractiveServiceServer) {
	s.RegisterService(&InteractiveService_ServiceDesc, srv)
}

func _InteractiveService_IncrReadCnt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncrReadCntRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InteractiveServiceServer).IncrReadCnt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InteractiveService_IncrReadCnt_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InteractiveServiceServer).IncrReadCnt(ctx, req.(*IncrReadCntRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InteractiveService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InteractiveServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InteractiveService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InteractiveServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InteractiveService_ServiceDesc is the grpc.ServiceDesc for InteractiveService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InteractiveService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "inter.InteractiveService",
	HandlerType: (*InteractiveServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IncrReadCnt",
			Handler:    _InteractiveService_IncrReadCnt_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _InteractiveService_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "inter/interactive.proto",
}
