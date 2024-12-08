// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.2
// source: v1/gostream.proto

package gostreamv1

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

// WatchersServiceClient is the client API for WatchersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WatchersServiceClient interface {
	// Watch connects to server streams
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (WatchersService_WatchClient, error)
	// Advertise returns the server stats useful for client-side loadbalancing
	Advertise(ctx context.Context, in *AdvertiseRequest, opts ...grpc.CallOption) (*AdvertiseResponse, error)
}

type watchersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWatchersServiceClient(cc grpc.ClientConnInterface) WatchersServiceClient {
	return &watchersServiceClient{cc}
}

func (c *watchersServiceClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (WatchersService_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &WatchersService_ServiceDesc.Streams[0], "/com.lordvidex.gostream.v1.WatchersService/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &watchersServiceWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WatchersService_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type watchersServiceWatchClient struct {
	grpc.ClientStream
}

func (x *watchersServiceWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *watchersServiceClient) Advertise(ctx context.Context, in *AdvertiseRequest, opts ...grpc.CallOption) (*AdvertiseResponse, error) {
	out := new(AdvertiseResponse)
	err := c.cc.Invoke(ctx, "/com.lordvidex.gostream.v1.WatchersService/Advertise", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WatchersServiceServer is the server API for WatchersService service.
// All implementations must embed UnimplementedWatchersServiceServer
// for forward compatibility
type WatchersServiceServer interface {
	// Watch connects to server streams
	Watch(*WatchRequest, WatchersService_WatchServer) error
	// Advertise returns the server stats useful for client-side loadbalancing
	Advertise(context.Context, *AdvertiseRequest) (*AdvertiseResponse, error)
	mustEmbedUnimplementedWatchersServiceServer()
}

// UnimplementedWatchersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWatchersServiceServer struct {
}

func (UnimplementedWatchersServiceServer) Watch(*WatchRequest, WatchersService_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedWatchersServiceServer) Advertise(context.Context, *AdvertiseRequest) (*AdvertiseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Advertise not implemented")
}
func (UnimplementedWatchersServiceServer) mustEmbedUnimplementedWatchersServiceServer() {}

// UnsafeWatchersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WatchersServiceServer will
// result in compilation errors.
type UnsafeWatchersServiceServer interface {
	mustEmbedUnimplementedWatchersServiceServer()
}

func RegisterWatchersServiceServer(s grpc.ServiceRegistrar, srv WatchersServiceServer) {
	s.RegisterService(&WatchersService_ServiceDesc, srv)
}

func _WatchersService_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WatchersServiceServer).Watch(m, &watchersServiceWatchServer{stream})
}

type WatchersService_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type watchersServiceWatchServer struct {
	grpc.ServerStream
}

func (x *watchersServiceWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _WatchersService_Advertise_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdvertiseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WatchersServiceServer).Advertise(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.lordvidex.gostream.v1.WatchersService/Advertise",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WatchersServiceServer).Advertise(ctx, req.(*AdvertiseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WatchersService_ServiceDesc is the grpc.ServiceDesc for WatchersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WatchersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.lordvidex.gostream.v1.WatchersService",
	HandlerType: (*WatchersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Advertise",
			Handler:    _WatchersService_Advertise_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _WatchersService_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "v1/gostream.proto",
}

// PetServiceClient is the client API for PetService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PetServiceClient interface {
	CreatePet(ctx context.Context, in *CreatePetRequest, opts ...grpc.CallOption) (*CreatePetResponse, error)
	UpdatePet(ctx context.Context, in *UpdatePetRequest, opts ...grpc.CallOption) (*UpdatePetResponse, error)
	ListPets(ctx context.Context, in *ListPetsRequest, opts ...grpc.CallOption) (*ListPetsResponse, error)
}

type petServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPetServiceClient(cc grpc.ClientConnInterface) PetServiceClient {
	return &petServiceClient{cc}
}

func (c *petServiceClient) CreatePet(ctx context.Context, in *CreatePetRequest, opts ...grpc.CallOption) (*CreatePetResponse, error) {
	out := new(CreatePetResponse)
	err := c.cc.Invoke(ctx, "/com.lordvidex.gostream.v1.PetService/CreatePet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petServiceClient) UpdatePet(ctx context.Context, in *UpdatePetRequest, opts ...grpc.CallOption) (*UpdatePetResponse, error) {
	out := new(UpdatePetResponse)
	err := c.cc.Invoke(ctx, "/com.lordvidex.gostream.v1.PetService/UpdatePet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petServiceClient) ListPets(ctx context.Context, in *ListPetsRequest, opts ...grpc.CallOption) (*ListPetsResponse, error) {
	out := new(ListPetsResponse)
	err := c.cc.Invoke(ctx, "/com.lordvidex.gostream.v1.PetService/ListPets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PetServiceServer is the server API for PetService service.
// All implementations must embed UnimplementedPetServiceServer
// for forward compatibility
type PetServiceServer interface {
	CreatePet(context.Context, *CreatePetRequest) (*CreatePetResponse, error)
	UpdatePet(context.Context, *UpdatePetRequest) (*UpdatePetResponse, error)
	ListPets(context.Context, *ListPetsRequest) (*ListPetsResponse, error)
	mustEmbedUnimplementedPetServiceServer()
}

// UnimplementedPetServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPetServiceServer struct {
}

func (UnimplementedPetServiceServer) CreatePet(context.Context, *CreatePetRequest) (*CreatePetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePet not implemented")
}
func (UnimplementedPetServiceServer) UpdatePet(context.Context, *UpdatePetRequest) (*UpdatePetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePet not implemented")
}
func (UnimplementedPetServiceServer) ListPets(context.Context, *ListPetsRequest) (*ListPetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPets not implemented")
}
func (UnimplementedPetServiceServer) mustEmbedUnimplementedPetServiceServer() {}

// UnsafePetServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PetServiceServer will
// result in compilation errors.
type UnsafePetServiceServer interface {
	mustEmbedUnimplementedPetServiceServer()
}

func RegisterPetServiceServer(s grpc.ServiceRegistrar, srv PetServiceServer) {
	s.RegisterService(&PetService_ServiceDesc, srv)
}

func _PetService_CreatePet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetServiceServer).CreatePet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.lordvidex.gostream.v1.PetService/CreatePet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetServiceServer).CreatePet(ctx, req.(*CreatePetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetService_UpdatePet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetServiceServer).UpdatePet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.lordvidex.gostream.v1.PetService/UpdatePet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetServiceServer).UpdatePet(ctx, req.(*UpdatePetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetService_ListPets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetServiceServer).ListPets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.lordvidex.gostream.v1.PetService/ListPets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetServiceServer).ListPets(ctx, req.(*ListPetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PetService_ServiceDesc is the grpc.ServiceDesc for PetService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PetService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.lordvidex.gostream.v1.PetService",
	HandlerType: (*PetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePet",
			Handler:    _PetService_CreatePet_Handler,
		},
		{
			MethodName: "UpdatePet",
			Handler:    _PetService_UpdatePet_Handler,
		},
		{
			MethodName: "ListPets",
			Handler:    _PetService_ListPets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/gostream.proto",
}

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserRequest, error)
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/com.lordvidex.gostream.v1.UserService/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserRequest, error) {
	out := new(UpdateUserRequest)
	err := c.cc.Invoke(ctx, "/com.lordvidex.gostream.v1.UserService/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, "/com.lordvidex.gostream.v1.UserService/ListUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserRequest, error)
	ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServiceServer) ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.lordvidex.gostream.v1.UserService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.lordvidex.gostream.v1.UserService/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.lordvidex.gostream.v1.UserService/ListUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ListUsers(ctx, req.(*ListUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.lordvidex.gostream.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "ListUsers",
			Handler:    _UserService_ListUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/gostream.proto",
}
