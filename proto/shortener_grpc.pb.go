// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package url_shortener

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

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	GetOriginalURL(ctx context.Context, in *ShortURL, opts ...grpc.CallOption) (*OriginalURL, error)
	GetShortURL(ctx context.Context, in *OriginalURL, opts ...grpc.CallOption) (*ShortURL, error)
	GetBatchShortURL(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*BatchResponse, error)
	GetAllUserURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllUrlsResponse, error)
	DeleteURLs(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetStatistic(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatisticResposne, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) GetOriginalURL(ctx context.Context, in *ShortURL, opts ...grpc.CallOption) (*OriginalURL, error) {
	out := new(OriginalURL)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetOriginalURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetShortURL(ctx context.Context, in *OriginalURL, opts ...grpc.CallOption) (*ShortURL, error) {
	out := new(ShortURL)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetShortURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetBatchShortURL(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*BatchResponse, error) {
	out := new(BatchResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetBatchShortURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetAllUserURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AllUrlsResponse, error) {
	out := new(AllUrlsResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetAllUserURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) DeleteURLs(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/DeleteURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetStatistic(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatisticResposne, error) {
	out := new(StatisticResposne)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetStatistic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	GetOriginalURL(context.Context, *ShortURL) (*OriginalURL, error)
	GetShortURL(context.Context, *OriginalURL) (*ShortURL, error)
	GetBatchShortURL(context.Context, *BatchRequest) (*BatchResponse, error)
	GetAllUserURL(context.Context, *emptypb.Empty) (*AllUrlsResponse, error)
	DeleteURLs(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	GetStatistic(context.Context, *emptypb.Empty) (*StatisticResposne, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) GetOriginalURL(context.Context, *ShortURL) (*OriginalURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedShortenerServer) GetShortURL(context.Context, *OriginalURL) (*ShortURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortURL not implemented")
}
func (UnimplementedShortenerServer) GetBatchShortURL(context.Context, *BatchRequest) (*BatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBatchShortURL not implemented")
}
func (UnimplementedShortenerServer) GetAllUserURL(context.Context, *emptypb.Empty) (*AllUrlsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUserURL not implemented")
}
func (UnimplementedShortenerServer) DeleteURLs(context.Context, *DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURLs not implemented")
}
func (UnimplementedShortenerServer) GetStatistic(context.Context, *emptypb.Empty) (*StatisticResposne, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatistic not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortURL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetOriginalURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetOriginalURL(ctx, req.(*ShortURL))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OriginalURL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetShortURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetShortURL(ctx, req.(*OriginalURL))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetBatchShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetBatchShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetBatchShortURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetBatchShortURL(ctx, req.(*BatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetAllUserURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetAllUserURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetAllUserURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetAllUserURL(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_DeleteURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/DeleteURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteURLs(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetStatistic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetStatistic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetStatistic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetStatistic(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOriginalURL",
			Handler:    _Shortener_GetOriginalURL_Handler,
		},
		{
			MethodName: "GetShortURL",
			Handler:    _Shortener_GetShortURL_Handler,
		},
		{
			MethodName: "GetBatchShortURL",
			Handler:    _Shortener_GetBatchShortURL_Handler,
		},
		{
			MethodName: "GetAllUserURL",
			Handler:    _Shortener_GetAllUserURL_Handler,
		},
		{
			MethodName: "DeleteURLs",
			Handler:    _Shortener_DeleteURLs_Handler,
		},
		{
			MethodName: "GetStatistic",
			Handler:    _Shortener_GetStatistic_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}