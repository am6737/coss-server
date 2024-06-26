// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: api/grpc/v1/group_relation.proto

package v1

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

const (
	GroupRelationService_GetGroupUserIDs_FullMethodName                    = "/relation_v1.GroupRelationService/GetGroupUserIDs"
	GroupRelationService_GetUserManageGroupID_FullMethodName               = "/relation_v1.GroupRelationService/GetUserManageGroupID"
	GroupRelationService_GetGroupRelation_FullMethodName                   = "/relation_v1.GroupRelationService/GetGroupRelation"
	GroupRelationService_GetBatchGroupRelation_FullMethodName              = "/relation_v1.GroupRelationService/GetBatchGroupRelation"
	GroupRelationService_DeleteGroupRelationByGroupId_FullMethodName       = "/relation_v1.GroupRelationService/DeleteGroupRelationByGroupId"
	GroupRelationService_DeleteGroupRelationByGroupIdRevert_FullMethodName = "/relation_v1.GroupRelationService/DeleteGroupRelationByGroupIdRevert"
	GroupRelationService_CreateGroupAndInviteUsers_FullMethodName          = "/relation_v1.GroupRelationService/CreateGroupAndInviteUsers"
	GroupRelationService_CreateGroupAndInviteUsersRevert_FullMethodName    = "/relation_v1.GroupRelationService/CreateGroupAndInviteUsersRevert"
)

// GroupRelationServiceClient is the client API for GroupRelationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GroupRelationServiceClient interface {
	// 获取群聊成员ID列表
	GetGroupUserIDs(ctx context.Context, in *GroupIDRequest, opts ...grpc.CallOption) (*UserIdsResponse, error)
	// 获取用户管理的群聊ID列表
	GetUserManageGroupID(ctx context.Context, in *GetUserManageGroupIDRequest, opts ...grpc.CallOption) (*GetUserManageGroupIDResponse, error)
	// 获取用户与群聊关系信息
	GetGroupRelation(ctx context.Context, in *GetGroupRelationRequest, opts ...grpc.CallOption) (*GetGroupRelationResponse, error)
	// 批量获取用户与群聊关系信息
	GetBatchGroupRelation(ctx context.Context, in *GetBatchGroupRelationRequest, opts ...grpc.CallOption) (*GetBatchGroupRelationResponse, error)
	// 根据群聊ID删除群聊的所有关系
	DeleteGroupRelationByGroupId(ctx context.Context, in *GroupIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteGroupRelationByGroupId回滚操作
	DeleteGroupRelationByGroupIdRevert(ctx context.Context, in *GroupIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 创建群聊并邀请多个用户进群（包括自己）
	CreateGroupAndInviteUsers(ctx context.Context, in *CreateGroupAndInviteUsersRequest, opts ...grpc.CallOption) (*CreateGroupAndInviteUsersResponse, error)
	// CreateGroupAndInviteUsers 回滚操作
	CreateGroupAndInviteUsersRevert(ctx context.Context, in *CreateGroupAndInviteUsersRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type groupRelationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupRelationServiceClient(cc grpc.ClientConnInterface) GroupRelationServiceClient {
	return &groupRelationServiceClient{cc}
}

func (c *groupRelationServiceClient) GetGroupUserIDs(ctx context.Context, in *GroupIDRequest, opts ...grpc.CallOption) (*UserIdsResponse, error) {
	out := new(UserIdsResponse)
	err := c.cc.Invoke(ctx, GroupRelationService_GetGroupUserIDs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupRelationServiceClient) GetUserManageGroupID(ctx context.Context, in *GetUserManageGroupIDRequest, opts ...grpc.CallOption) (*GetUserManageGroupIDResponse, error) {
	out := new(GetUserManageGroupIDResponse)
	err := c.cc.Invoke(ctx, GroupRelationService_GetUserManageGroupID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupRelationServiceClient) GetGroupRelation(ctx context.Context, in *GetGroupRelationRequest, opts ...grpc.CallOption) (*GetGroupRelationResponse, error) {
	out := new(GetGroupRelationResponse)
	err := c.cc.Invoke(ctx, GroupRelationService_GetGroupRelation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupRelationServiceClient) GetBatchGroupRelation(ctx context.Context, in *GetBatchGroupRelationRequest, opts ...grpc.CallOption) (*GetBatchGroupRelationResponse, error) {
	out := new(GetBatchGroupRelationResponse)
	err := c.cc.Invoke(ctx, GroupRelationService_GetBatchGroupRelation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupRelationServiceClient) DeleteGroupRelationByGroupId(ctx context.Context, in *GroupIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, GroupRelationService_DeleteGroupRelationByGroupId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupRelationServiceClient) DeleteGroupRelationByGroupIdRevert(ctx context.Context, in *GroupIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, GroupRelationService_DeleteGroupRelationByGroupIdRevert_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupRelationServiceClient) CreateGroupAndInviteUsers(ctx context.Context, in *CreateGroupAndInviteUsersRequest, opts ...grpc.CallOption) (*CreateGroupAndInviteUsersResponse, error) {
	out := new(CreateGroupAndInviteUsersResponse)
	err := c.cc.Invoke(ctx, GroupRelationService_CreateGroupAndInviteUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupRelationServiceClient) CreateGroupAndInviteUsersRevert(ctx context.Context, in *CreateGroupAndInviteUsersRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, GroupRelationService_CreateGroupAndInviteUsersRevert_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupRelationServiceServer is the server API for GroupRelationService service.
// All implementations should embed UnimplementedGroupRelationServiceServer
// for forward compatibility
type GroupRelationServiceServer interface {
	// 获取群聊成员ID列表
	GetGroupUserIDs(context.Context, *GroupIDRequest) (*UserIdsResponse, error)
	// 获取用户管理的群聊ID列表
	GetUserManageGroupID(context.Context, *GetUserManageGroupIDRequest) (*GetUserManageGroupIDResponse, error)
	// 获取用户与群聊关系信息
	GetGroupRelation(context.Context, *GetGroupRelationRequest) (*GetGroupRelationResponse, error)
	// 批量获取用户与群聊关系信息
	GetBatchGroupRelation(context.Context, *GetBatchGroupRelationRequest) (*GetBatchGroupRelationResponse, error)
	// 根据群聊ID删除群聊的所有关系
	DeleteGroupRelationByGroupId(context.Context, *GroupIDRequest) (*emptypb.Empty, error)
	// DeleteGroupRelationByGroupId回滚操作
	DeleteGroupRelationByGroupIdRevert(context.Context, *GroupIDRequest) (*emptypb.Empty, error)
	// 创建群聊并邀请多个用户进群（包括自己）
	CreateGroupAndInviteUsers(context.Context, *CreateGroupAndInviteUsersRequest) (*CreateGroupAndInviteUsersResponse, error)
	// CreateGroupAndInviteUsers 回滚操作
	CreateGroupAndInviteUsersRevert(context.Context, *CreateGroupAndInviteUsersRequest) (*emptypb.Empty, error)
}

// UnimplementedGroupRelationServiceServer should be embedded to have forward compatible implementations.
type UnimplementedGroupRelationServiceServer struct {
}

func (UnimplementedGroupRelationServiceServer) GetGroupUserIDs(context.Context, *GroupIDRequest) (*UserIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupUserIDs not implemented")
}
func (UnimplementedGroupRelationServiceServer) GetUserManageGroupID(context.Context, *GetUserManageGroupIDRequest) (*GetUserManageGroupIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserManageGroupID not implemented")
}
func (UnimplementedGroupRelationServiceServer) GetGroupRelation(context.Context, *GetGroupRelationRequest) (*GetGroupRelationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupRelation not implemented")
}
func (UnimplementedGroupRelationServiceServer) GetBatchGroupRelation(context.Context, *GetBatchGroupRelationRequest) (*GetBatchGroupRelationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBatchGroupRelation not implemented")
}
func (UnimplementedGroupRelationServiceServer) DeleteGroupRelationByGroupId(context.Context, *GroupIDRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroupRelationByGroupId not implemented")
}
func (UnimplementedGroupRelationServiceServer) DeleteGroupRelationByGroupIdRevert(context.Context, *GroupIDRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroupRelationByGroupIdRevert not implemented")
}
func (UnimplementedGroupRelationServiceServer) CreateGroupAndInviteUsers(context.Context, *CreateGroupAndInviteUsersRequest) (*CreateGroupAndInviteUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroupAndInviteUsers not implemented")
}
func (UnimplementedGroupRelationServiceServer) CreateGroupAndInviteUsersRevert(context.Context, *CreateGroupAndInviteUsersRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroupAndInviteUsersRevert not implemented")
}

// UnsafeGroupRelationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GroupRelationServiceServer will
// result in compilation errors.
type UnsafeGroupRelationServiceServer interface {
	mustEmbedUnimplementedGroupRelationServiceServer()
}

func RegisterGroupRelationServiceServer(s grpc.ServiceRegistrar, srv GroupRelationServiceServer) {
	s.RegisterService(&GroupRelationService_ServiceDesc, srv)
}

func _GroupRelationService_GetGroupUserIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).GetGroupUserIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_GetGroupUserIDs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).GetGroupUserIDs(ctx, req.(*GroupIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupRelationService_GetUserManageGroupID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserManageGroupIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).GetUserManageGroupID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_GetUserManageGroupID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).GetUserManageGroupID(ctx, req.(*GetUserManageGroupIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupRelationService_GetGroupRelation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupRelationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).GetGroupRelation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_GetGroupRelation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).GetGroupRelation(ctx, req.(*GetGroupRelationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupRelationService_GetBatchGroupRelation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBatchGroupRelationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).GetBatchGroupRelation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_GetBatchGroupRelation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).GetBatchGroupRelation(ctx, req.(*GetBatchGroupRelationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupRelationService_DeleteGroupRelationByGroupId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).DeleteGroupRelationByGroupId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_DeleteGroupRelationByGroupId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).DeleteGroupRelationByGroupId(ctx, req.(*GroupIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupRelationService_DeleteGroupRelationByGroupIdRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).DeleteGroupRelationByGroupIdRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_DeleteGroupRelationByGroupIdRevert_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).DeleteGroupRelationByGroupIdRevert(ctx, req.(*GroupIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupRelationService_CreateGroupAndInviteUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupAndInviteUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).CreateGroupAndInviteUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_CreateGroupAndInviteUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).CreateGroupAndInviteUsers(ctx, req.(*CreateGroupAndInviteUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupRelationService_CreateGroupAndInviteUsersRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupAndInviteUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupRelationServiceServer).CreateGroupAndInviteUsersRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupRelationService_CreateGroupAndInviteUsersRevert_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupRelationServiceServer).CreateGroupAndInviteUsersRevert(ctx, req.(*CreateGroupAndInviteUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GroupRelationService_ServiceDesc is the grpc.ServiceDesc for GroupRelationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GroupRelationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relation_v1.GroupRelationService",
	HandlerType: (*GroupRelationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGroupUserIDs",
			Handler:    _GroupRelationService_GetGroupUserIDs_Handler,
		},
		{
			MethodName: "GetUserManageGroupID",
			Handler:    _GroupRelationService_GetUserManageGroupID_Handler,
		},
		{
			MethodName: "GetGroupRelation",
			Handler:    _GroupRelationService_GetGroupRelation_Handler,
		},
		{
			MethodName: "GetBatchGroupRelation",
			Handler:    _GroupRelationService_GetBatchGroupRelation_Handler,
		},
		{
			MethodName: "DeleteGroupRelationByGroupId",
			Handler:    _GroupRelationService_DeleteGroupRelationByGroupId_Handler,
		},
		{
			MethodName: "DeleteGroupRelationByGroupIdRevert",
			Handler:    _GroupRelationService_DeleteGroupRelationByGroupIdRevert_Handler,
		},
		{
			MethodName: "CreateGroupAndInviteUsers",
			Handler:    _GroupRelationService_CreateGroupAndInviteUsers_Handler,
		},
		{
			MethodName: "CreateGroupAndInviteUsersRevert",
			Handler:    _GroupRelationService_CreateGroupAndInviteUsersRevert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/grpc/v1/group_relation.proto",
}
