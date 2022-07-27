// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: alliance.proto

package protobuf

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

// AllianceStorageClient is the client API for AllianceStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AllianceStorageClient interface {
	Cmd(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type allianceStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewAllianceStorageClient(cc grpc.ClientConnInterface) AllianceStorageClient {
	return &allianceStorageClient{cc}
}

func (c *allianceStorageClient) Cmd(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/pb.AllianceStorage/Cmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AllianceStorageServer is the server API for AllianceStorage service.
// All implementations should embed UnimplementedAllianceStorageServer
// for forward compatibility
type AllianceStorageServer interface {
	Cmd(context.Context, *Request) (*Response, error)
}

// UnimplementedAllianceStorageServer should be embedded to have forward compatible implementations.
type UnimplementedAllianceStorageServer struct {
}

func (UnimplementedAllianceStorageServer) Cmd(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cmd not implemented")
}

// UnsafeAllianceStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AllianceStorageServer will
// result in compilation errors.
type UnsafeAllianceStorageServer interface {
	mustEmbedUnimplementedAllianceStorageServer()
}

func RegisterAllianceStorageServer(s grpc.ServiceRegistrar, srv AllianceStorageServer) {
	s.RegisterService(&AllianceStorage_ServiceDesc, srv)
}

func _AllianceStorage_Cmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AllianceStorageServer).Cmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AllianceStorage/Cmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AllianceStorageServer).Cmd(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// AllianceStorage_ServiceDesc is the grpc.ServiceDesc for AllianceStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AllianceStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AllianceStorage",
	HandlerType: (*AllianceStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Cmd",
			Handler:    _AllianceStorage_Cmd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "alliance.proto",
}