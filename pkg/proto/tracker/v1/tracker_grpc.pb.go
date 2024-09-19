// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: trackerV1/v1/trackerV1.proto

package trackerV1

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
	TrackerService_GetTracker_FullMethodName = "/trackerV1.v1.TrackerService/GetTracker"
)

// TrackerServiceClient is the client API for TrackerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrackerServiceClient interface {
	// GetTracker
	GetTracker(ctx context.Context, in *GetTrackerRequest, opts ...grpc.CallOption) (*GetTrackerResponse, error)
}

type trackerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTrackerServiceClient(cc grpc.ClientConnInterface) TrackerServiceClient {
	return &trackerServiceClient{cc}
}

func (c *trackerServiceClient) GetTracker(ctx context.Context, in *GetTrackerRequest, opts ...grpc.CallOption) (*GetTrackerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTrackerResponse)
	err := c.cc.Invoke(ctx, TrackerService_GetTracker_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrackerServiceServer is the server API for TrackerService service.
// All implementations must embed UnimplementedTrackerServiceServer
// for forward compatibility.
type TrackerServiceServer interface {
	// GetTracker
	GetTracker(context.Context, *GetTrackerRequest) (*GetTrackerResponse, error)
	mustEmbedUnimplementedTrackerServiceServer()
}

// UnimplementedTrackerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTrackerServiceServer struct{}

func (UnimplementedTrackerServiceServer) GetTracker(context.Context, *GetTrackerRequest) (*GetTrackerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTracker not implemented")
}
func (UnimplementedTrackerServiceServer) mustEmbedUnimplementedTrackerServiceServer() {}
func (UnimplementedTrackerServiceServer) testEmbeddedByValue()                        {}

// UnsafeTrackerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrackerServiceServer will
// result in compilation errors.
type UnsafeTrackerServiceServer interface {
	mustEmbedUnimplementedTrackerServiceServer()
}

func RegisterTrackerServiceServer(s grpc.ServiceRegistrar, srv TrackerServiceServer) {
	// If the following call pancis, it indicates UnimplementedTrackerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TrackerService_ServiceDesc, srv)
}

func _TrackerService_GetTracker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTrackerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackerServiceServer).GetTracker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TrackerService_GetTracker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackerServiceServer).GetTracker(ctx, req.(*GetTrackerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TrackerService_ServiceDesc is the grpc.ServiceDesc for TrackerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TrackerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "trackerV1.v1.TrackerService",
	HandlerType: (*TrackerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTracker",
			Handler:    _TrackerService_GetTracker_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "trackerV1/v1/trackerV1.proto",
}
