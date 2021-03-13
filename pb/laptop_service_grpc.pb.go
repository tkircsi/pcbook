// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// LaptopServiceClient is the client API for LaptopService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LaptopServiceClient interface {
	CreateLaptop(ctx context.Context, in *CreateLaptopRequest, opts ...grpc.CallOption) (*CreateLaptopResponse, error)
	SearchLaptop(ctx context.Context, in *SearchLaptopRequest, opts ...grpc.CallOption) (LaptopService_SearchLaptopClient, error)
	UploadImage(ctx context.Context, opts ...grpc.CallOption) (LaptopService_UploadImageClient, error)
	RateLaptop(ctx context.Context, opts ...grpc.CallOption) (LaptopService_RateLaptopClient, error)
}

type laptopServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLaptopServiceClient(cc grpc.ClientConnInterface) LaptopServiceClient {
	return &laptopServiceClient{cc}
}

func (c *laptopServiceClient) CreateLaptop(ctx context.Context, in *CreateLaptopRequest, opts ...grpc.CallOption) (*CreateLaptopResponse, error) {
	out := new(CreateLaptopResponse)
	err := c.cc.Invoke(ctx, "/pcbook.pb.LaptopService/CreateLaptop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *laptopServiceClient) SearchLaptop(ctx context.Context, in *SearchLaptopRequest, opts ...grpc.CallOption) (LaptopService_SearchLaptopClient, error) {
	stream, err := c.cc.NewStream(ctx, &LaptopService_ServiceDesc.Streams[0], "/pcbook.pb.LaptopService/SearchLaptop", opts...)
	if err != nil {
		return nil, err
	}
	x := &laptopServiceSearchLaptopClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LaptopService_SearchLaptopClient interface {
	Recv() (*SearchLaptopResponse, error)
	grpc.ClientStream
}

type laptopServiceSearchLaptopClient struct {
	grpc.ClientStream
}

func (x *laptopServiceSearchLaptopClient) Recv() (*SearchLaptopResponse, error) {
	m := new(SearchLaptopResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *laptopServiceClient) UploadImage(ctx context.Context, opts ...grpc.CallOption) (LaptopService_UploadImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &LaptopService_ServiceDesc.Streams[1], "/pcbook.pb.LaptopService/UploadImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &laptopServiceUploadImageClient{stream}
	return x, nil
}

type LaptopService_UploadImageClient interface {
	Send(*UploadImageRequest) error
	CloseAndRecv() (*UploadImageResponse, error)
	grpc.ClientStream
}

type laptopServiceUploadImageClient struct {
	grpc.ClientStream
}

func (x *laptopServiceUploadImageClient) Send(m *UploadImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *laptopServiceUploadImageClient) CloseAndRecv() (*UploadImageResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadImageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *laptopServiceClient) RateLaptop(ctx context.Context, opts ...grpc.CallOption) (LaptopService_RateLaptopClient, error) {
	stream, err := c.cc.NewStream(ctx, &LaptopService_ServiceDesc.Streams[2], "/pcbook.pb.LaptopService/RateLaptop", opts...)
	if err != nil {
		return nil, err
	}
	x := &laptopServiceRateLaptopClient{stream}
	return x, nil
}

type LaptopService_RateLaptopClient interface {
	Send(*RateLaptopRequest) error
	Recv() (*RateLaptopResponse, error)
	grpc.ClientStream
}

type laptopServiceRateLaptopClient struct {
	grpc.ClientStream
}

func (x *laptopServiceRateLaptopClient) Send(m *RateLaptopRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *laptopServiceRateLaptopClient) Recv() (*RateLaptopResponse, error) {
	m := new(RateLaptopResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LaptopServiceServer is the server API for LaptopService service.
// All implementations must embed UnimplementedLaptopServiceServer
// for forward compatibility
type LaptopServiceServer interface {
	CreateLaptop(context.Context, *CreateLaptopRequest) (*CreateLaptopResponse, error)
	SearchLaptop(*SearchLaptopRequest, LaptopService_SearchLaptopServer) error
	UploadImage(LaptopService_UploadImageServer) error
	RateLaptop(LaptopService_RateLaptopServer) error
	mustEmbedUnimplementedLaptopServiceServer()
}

// UnimplementedLaptopServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLaptopServiceServer struct {
}

func (UnimplementedLaptopServiceServer) CreateLaptop(context.Context, *CreateLaptopRequest) (*CreateLaptopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLaptop not implemented")
}
func (UnimplementedLaptopServiceServer) SearchLaptop(*SearchLaptopRequest, LaptopService_SearchLaptopServer) error {
	return status.Errorf(codes.Unimplemented, "method SearchLaptop not implemented")
}
func (UnimplementedLaptopServiceServer) UploadImage(LaptopService_UploadImageServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadImage not implemented")
}
func (UnimplementedLaptopServiceServer) RateLaptop(LaptopService_RateLaptopServer) error {
	return status.Errorf(codes.Unimplemented, "method RateLaptop not implemented")
}
func (UnimplementedLaptopServiceServer) mustEmbedUnimplementedLaptopServiceServer() {}

// UnsafeLaptopServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LaptopServiceServer will
// result in compilation errors.
type UnsafeLaptopServiceServer interface {
	mustEmbedUnimplementedLaptopServiceServer()
}

func RegisterLaptopServiceServer(s grpc.ServiceRegistrar, srv LaptopServiceServer) {
	s.RegisterService(&LaptopService_ServiceDesc, srv)
}

func _LaptopService_CreateLaptop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLaptopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LaptopServiceServer).CreateLaptop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pcbook.pb.LaptopService/CreateLaptop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LaptopServiceServer).CreateLaptop(ctx, req.(*CreateLaptopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LaptopService_SearchLaptop_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchLaptopRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LaptopServiceServer).SearchLaptop(m, &laptopServiceSearchLaptopServer{stream})
}

type LaptopService_SearchLaptopServer interface {
	Send(*SearchLaptopResponse) error
	grpc.ServerStream
}

type laptopServiceSearchLaptopServer struct {
	grpc.ServerStream
}

func (x *laptopServiceSearchLaptopServer) Send(m *SearchLaptopResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _LaptopService_UploadImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LaptopServiceServer).UploadImage(&laptopServiceUploadImageServer{stream})
}

type LaptopService_UploadImageServer interface {
	SendAndClose(*UploadImageResponse) error
	Recv() (*UploadImageRequest, error)
	grpc.ServerStream
}

type laptopServiceUploadImageServer struct {
	grpc.ServerStream
}

func (x *laptopServiceUploadImageServer) SendAndClose(m *UploadImageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *laptopServiceUploadImageServer) Recv() (*UploadImageRequest, error) {
	m := new(UploadImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _LaptopService_RateLaptop_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LaptopServiceServer).RateLaptop(&laptopServiceRateLaptopServer{stream})
}

type LaptopService_RateLaptopServer interface {
	Send(*RateLaptopResponse) error
	Recv() (*RateLaptopRequest, error)
	grpc.ServerStream
}

type laptopServiceRateLaptopServer struct {
	grpc.ServerStream
}

func (x *laptopServiceRateLaptopServer) Send(m *RateLaptopResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *laptopServiceRateLaptopServer) Recv() (*RateLaptopRequest, error) {
	m := new(RateLaptopRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LaptopService_ServiceDesc is the grpc.ServiceDesc for LaptopService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LaptopService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pcbook.pb.LaptopService",
	HandlerType: (*LaptopServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLaptop",
			Handler:    _LaptopService_CreateLaptop_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SearchLaptop",
			Handler:       _LaptopService_SearchLaptop_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadImage",
			Handler:       _LaptopService_UploadImage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RateLaptop",
			Handler:       _LaptopService_RateLaptop_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "laptop_service.proto",
}
