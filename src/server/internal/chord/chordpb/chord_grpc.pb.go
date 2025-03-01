// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.3
// source: chord.proto

package chordpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ChordService_Notify_FullMethodName = "/chord.ChordService/Notify"
	ChordService_Health_FullMethodName = "/chord.ChordService/Health"
)

// ChordServiceClient is the client API for ChordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChordServiceClient interface {
	Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error)
	Health(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HealthResponse, error)
}

type chordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChordServiceClient(cc grpc.ClientConnInterface) ChordServiceClient {
	return &chordServiceClient{cc}
}

func (c *chordServiceClient) Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NotifyResponse)
	err := c.cc.Invoke(ctx, ChordService_Notify_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordServiceClient) Health(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HealthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, ChordService_Health_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChordServiceServer is the server API for ChordService service.
// All implementations must embed UnimplementedChordServiceServer
// for forward compatibility.
type ChordServiceServer interface {
	Notify(context.Context, *NotifyRequest) (*NotifyResponse, error)
	Health(context.Context, *Empty) (*HealthResponse, error)
	mustEmbedUnimplementedChordServiceServer()
}

// UnimplementedChordServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChordServiceServer struct{}

func (UnimplementedChordServiceServer) Notify(context.Context, *NotifyRequest) (*NotifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
}
func (UnimplementedChordServiceServer) Health(context.Context, *Empty) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedChordServiceServer) mustEmbedUnimplementedChordServiceServer() {}
func (UnimplementedChordServiceServer) testEmbeddedByValue()                      {}

// UnsafeChordServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChordServiceServer will
// result in compilation errors.
type UnsafeChordServiceServer interface {
	mustEmbedUnimplementedChordServiceServer()
}

func RegisterChordServiceServer(s grpc.ServiceRegistrar, srv ChordServiceServer) {
	// If the following call pancis, it indicates UnimplementedChordServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChordService_ServiceDesc, srv)
}

func _ChordService_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServiceServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChordService_Notify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServiceServer).Notify(ctx, req.(*NotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChordService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChordService_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServiceServer).Health(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ChordService_ServiceDesc is the grpc.ServiceDesc for ChordService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChordService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chord.ChordService",
	HandlerType: (*ChordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Notify",
			Handler:    _ChordService_Notify_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _ChordService_Health_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chord.proto",
}
