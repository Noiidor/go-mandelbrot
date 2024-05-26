// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: img_mandelbrot.proto

package mandelbrotv1

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
	ImgMandelbrotService_GetMandelbrotImage_FullMethodName = "/mandelbrot.v1.ImgMandelbrotService/GetMandelbrotImage"
)

// ImgMandelbrotServiceClient is the client API for ImgMandelbrotService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ImgMandelbrotServiceClient interface {
	GetMandelbrotImage(ctx context.Context, in *GetMandelbrotImageRequest, opts ...grpc.CallOption) (*GetMandelbrotImageResponse, error)
}

type imgMandelbrotServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewImgMandelbrotServiceClient(cc grpc.ClientConnInterface) ImgMandelbrotServiceClient {
	return &imgMandelbrotServiceClient{cc}
}

func (c *imgMandelbrotServiceClient) GetMandelbrotImage(ctx context.Context, in *GetMandelbrotImageRequest, opts ...grpc.CallOption) (*GetMandelbrotImageResponse, error) {
	out := new(GetMandelbrotImageResponse)
	err := c.cc.Invoke(ctx, ImgMandelbrotService_GetMandelbrotImage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImgMandelbrotServiceServer is the server API for ImgMandelbrotService service.
// All implementations must embed UnimplementedImgMandelbrotServiceServer
// for forward compatibility
type ImgMandelbrotServiceServer interface {
	GetMandelbrotImage(context.Context, *GetMandelbrotImageRequest) (*GetMandelbrotImageResponse, error)
	mustEmbedUnimplementedImgMandelbrotServiceServer()
}

// UnimplementedImgMandelbrotServiceServer must be embedded to have forward compatible implementations.
type UnimplementedImgMandelbrotServiceServer struct {
}

func (UnimplementedImgMandelbrotServiceServer) GetMandelbrotImage(context.Context, *GetMandelbrotImageRequest) (*GetMandelbrotImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMandelbrotImage not implemented")
}
func (UnimplementedImgMandelbrotServiceServer) mustEmbedUnimplementedImgMandelbrotServiceServer() {}

// UnsafeImgMandelbrotServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ImgMandelbrotServiceServer will
// result in compilation errors.
type UnsafeImgMandelbrotServiceServer interface {
	mustEmbedUnimplementedImgMandelbrotServiceServer()
}

func RegisterImgMandelbrotServiceServer(s grpc.ServiceRegistrar, srv ImgMandelbrotServiceServer) {
	s.RegisterService(&ImgMandelbrotService_ServiceDesc, srv)
}

func _ImgMandelbrotService_GetMandelbrotImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMandelbrotImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImgMandelbrotServiceServer).GetMandelbrotImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ImgMandelbrotService_GetMandelbrotImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImgMandelbrotServiceServer).GetMandelbrotImage(ctx, req.(*GetMandelbrotImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ImgMandelbrotService_ServiceDesc is the grpc.ServiceDesc for ImgMandelbrotService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ImgMandelbrotService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mandelbrot.v1.ImgMandelbrotService",
	HandlerType: (*ImgMandelbrotServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMandelbrotImage",
			Handler:    _ImgMandelbrotService_GetMandelbrotImage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "img_mandelbrot.proto",
}