// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: backup.proto

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

const (
	BackupService_Backup_FullMethodName = "/BackupService/Backup"
)

// BackupServiceClient is the client API for BackupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BackupServiceClient interface {
	Backup(ctx context.Context, opts ...grpc.CallOption) (BackupService_BackupClient, error)
}

type backupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBackupServiceClient(cc grpc.ClientConnInterface) BackupServiceClient {
	return &backupServiceClient{cc}
}

func (c *backupServiceClient) Backup(ctx context.Context, opts ...grpc.CallOption) (BackupService_BackupClient, error) {
	stream, err := c.cc.NewStream(ctx, &BackupService_ServiceDesc.Streams[0], BackupService_Backup_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &backupServiceBackupClient{stream}
	return x, nil
}

type BackupService_BackupClient interface {
	Send(*BackupRequest) error
	CloseAndRecv() (*BackupResponse, error)
	grpc.ClientStream
}

type backupServiceBackupClient struct {
	grpc.ClientStream
}

func (x *backupServiceBackupClient) Send(m *BackupRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *backupServiceBackupClient) CloseAndRecv() (*BackupResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(BackupResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BackupServiceServer is the server API for BackupService service.
// All implementations must embed UnimplementedBackupServiceServer
// for forward compatibility
type BackupServiceServer interface {
	Backup(BackupService_BackupServer) error
	mustEmbedUnimplementedBackupServiceServer()
}

// UnimplementedBackupServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBackupServiceServer struct {
}

func (UnimplementedBackupServiceServer) Backup(BackupService_BackupServer) error {
	return status.Errorf(codes.Unimplemented, "method Backup not implemented")
}
func (UnimplementedBackupServiceServer) mustEmbedUnimplementedBackupServiceServer() {}

// UnsafeBackupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BackupServiceServer will
// result in compilation errors.
type UnsafeBackupServiceServer interface {
	mustEmbedUnimplementedBackupServiceServer()
}

func RegisterBackupServiceServer(s grpc.ServiceRegistrar, srv BackupServiceServer) {
	s.RegisterService(&BackupService_ServiceDesc, srv)
}

func _BackupService_Backup_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BackupServiceServer).Backup(&backupServiceBackupServer{stream})
}

type BackupService_BackupServer interface {
	SendAndClose(*BackupResponse) error
	Recv() (*BackupRequest, error)
	grpc.ServerStream
}

type backupServiceBackupServer struct {
	grpc.ServerStream
}

func (x *backupServiceBackupServer) SendAndClose(m *BackupResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *backupServiceBackupServer) Recv() (*BackupRequest, error) {
	m := new(BackupRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BackupService_ServiceDesc is the grpc.ServiceDesc for BackupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BackupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "BackupService",
	HandlerType: (*BackupServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Backup",
			Handler:       _BackupService_Backup_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "backup.proto",
}
