// Code generated by protoc-gen-go. DO NOT EDIT.
// source: queue.proto

package queue

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PushRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                [][]byte `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushRequest) Reset()         { *m = PushRequest{} }
func (m *PushRequest) String() string { return proto.CompactTextString(m) }
func (*PushRequest) ProtoMessage()    {}
func (*PushRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e4d7d76a734cd8, []int{0}
}

func (m *PushRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushRequest.Unmarshal(m, b)
}
func (m *PushRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushRequest.Marshal(b, m, deterministic)
}
func (m *PushRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushRequest.Merge(m, src)
}
func (m *PushRequest) XXX_Size() int {
	return xxx_messageInfo_PushRequest.Size(m)
}
func (m *PushRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PushRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PushRequest proto.InternalMessageInfo

func (m *PushRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *PushRequest) GetValue() [][]byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type PullResponse struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                [][]byte `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PullResponse) Reset()         { *m = PullResponse{} }
func (m *PullResponse) String() string { return proto.CompactTextString(m) }
func (*PullResponse) ProtoMessage()    {}
func (*PullResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e4d7d76a734cd8, []int{1}
}

func (m *PullResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PullResponse.Unmarshal(m, b)
}
func (m *PullResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PullResponse.Marshal(b, m, deterministic)
}
func (m *PullResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PullResponse.Merge(m, src)
}
func (m *PullResponse) XXX_Size() int {
	return xxx_messageInfo_PullResponse.Size(m)
}
func (m *PullResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PullResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PullResponse proto.InternalMessageInfo

func (m *PullResponse) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *PullResponse) GetValue() [][]byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type AcknowledgePullRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcknowledgePullRequest) Reset()         { *m = AcknowledgePullRequest{} }
func (m *AcknowledgePullRequest) String() string { return proto.CompactTextString(m) }
func (*AcknowledgePullRequest) ProtoMessage()    {}
func (*AcknowledgePullRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e4d7d76a734cd8, []int{2}
}

func (m *AcknowledgePullRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcknowledgePullRequest.Unmarshal(m, b)
}
func (m *AcknowledgePullRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcknowledgePullRequest.Marshal(b, m, deterministic)
}
func (m *AcknowledgePullRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcknowledgePullRequest.Merge(m, src)
}
func (m *AcknowledgePullRequest) XXX_Size() int {
	return xxx_messageInfo_AcknowledgePullRequest.Size(m)
}
func (m *AcknowledgePullRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AcknowledgePullRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AcknowledgePullRequest proto.InternalMessageInfo

func (m *AcknowledgePullRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func init() {
	proto.RegisterType((*PushRequest)(nil), "queue.PushRequest")
	proto.RegisterType((*PullResponse)(nil), "queue.PullResponse")
	proto.RegisterType((*AcknowledgePullRequest)(nil), "queue.AcknowledgePullRequest")
}

func init() {
	proto.RegisterFile("queue.proto", fileDescriptor_96e4d7d76a734cd8)
}

var fileDescriptor_96e4d7d76a734cd8 = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x8f, 0xdf, 0x4b, 0xc3, 0x30,
	0x10, 0xc7, 0x99, 0xb3, 0x82, 0xd9, 0x40, 0x39, 0x65, 0x8c, 0x8a, 0x30, 0xf6, 0x34, 0x15, 0x12,
	0xd8, 0xd0, 0x77, 0x05, 0xc1, 0x17, 0x61, 0x2b, 0x3e, 0xf9, 0x22, 0xb6, 0x9e, 0xd9, 0xe8, 0xad,
	0xe9, 0x9a, 0x9c, 0xd2, 0x7f, 0xcd, 0xbf, 0x4e, 0xd2, 0xa8, 0x0c, 0xb1, 0xe0, 0x5b, 0x2e, 0x77,
	0xdf, 0x1f, 0x1f, 0xd1, 0xdb, 0x30, 0x32, 0xca, 0xb2, 0x32, 0xce, 0x40, 0xd4, 0x0c, 0xf1, 0x89,
	0x36, 0x46, 0x13, 0xaa, 0xe6, 0x33, 0xe5, 0x57, 0x85, 0xeb, 0xd2, 0xd5, 0xe1, 0x66, 0x7c, 0x29,
	0x7a, 0x73, 0xb6, 0xcb, 0x04, 0x37, 0x8c, 0xd6, 0xc1, 0xa1, 0xe8, 0xe6, 0x58, 0x0f, 0x3b, 0xa3,
	0xce, 0x64, 0x3f, 0xf1, 0x4f, 0x38, 0x16, 0xd1, 0xdb, 0x33, 0x31, 0x0e, 0x77, 0x46, 0xdd, 0x49,
	0x3f, 0x09, 0xc3, 0xf8, 0x4a, 0xf4, 0xe7, 0x4c, 0x94, 0xa0, 0x2d, 0x4d, 0x61, 0xf1, 0xdf, 0xba,
	0x73, 0x31, 0xb8, 0xce, 0xf2, 0xc2, 0xbc, 0x13, 0xbe, 0x68, 0x0c, 0x16, 0x2d, 0xc9, 0xd3, 0x8f,
	0x8e, 0x88, 0x16, 0x9e, 0x00, 0xa6, 0x62, 0xd7, 0x97, 0x04, 0x90, 0x01, 0x6f, 0xab, 0x71, 0x3c,
	0x90, 0x01, 0x4f, 0x7e, 0xe3, 0xc9, 0x5b, 0x8f, 0x07, 0x33, 0xaf, 0x21, 0x82, 0x96, 0x7d, 0x7c,
	0xf4, 0xe3, 0xb5, 0x85, 0x71, 0x27, 0x0e, 0x7e, 0xd5, 0x83, 0xd3, 0xaf, 0xbb, 0xbf, 0x6b, 0xb7,
	0xc5, 0xdf, 0x5c, 0x3c, 0x9e, 0xe9, 0x95, 0x5b, 0x72, 0x2a, 0x33, 0xb3, 0x56, 0x79, 0x6d, 0x2b,
	0x54, 0x0f, 0x5c, 0x39, 0xc2, 0xfb, 0x85, 0xca, 0x68, 0x85, 0x85, 0x7b, 0xd2, 0x46, 0x35, 0xde,
	0xe9, 0x5e, 0x23, 0x9e, 0x7d, 0x06, 0x00, 0x00, 0xff, 0xff, 0xd5, 0xeb, 0xd0, 0xf6, 0xbe, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// QueueClient is the client API for Queue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueueClient interface {
	Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Pull(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PullResponse, error)
	AcknowledgePull(ctx context.Context, in *AcknowledgePullRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type queueClient struct {
	cc grpc.ClientConnInterface
}

func NewQueueClient(cc grpc.ClientConnInterface) QueueClient {
	return &queueClient{cc}
}

func (c *queueClient) Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/queue.Queue/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) Pull(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PullResponse, error) {
	out := new(PullResponse)
	err := c.cc.Invoke(ctx, "/queue.Queue/Pull", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) AcknowledgePull(ctx context.Context, in *AcknowledgePullRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/queue.Queue/AcknowledgePull", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueueServer is the server API for Queue service.
type QueueServer interface {
	Push(context.Context, *PushRequest) (*empty.Empty, error)
	Pull(context.Context, *empty.Empty) (*PullResponse, error)
	AcknowledgePull(context.Context, *AcknowledgePullRequest) (*empty.Empty, error)
}

// UnimplementedQueueServer can be embedded to have forward compatible implementations.
type UnimplementedQueueServer struct {
}

func (*UnimplementedQueueServer) Push(ctx context.Context, req *PushRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (*UnimplementedQueueServer) Pull(ctx context.Context, req *empty.Empty) (*PullResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pull not implemented")
}
func (*UnimplementedQueueServer) AcknowledgePull(ctx context.Context, req *AcknowledgePullRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcknowledgePull not implemented")
}

func RegisterQueueServer(s *grpc.Server, srv QueueServer) {
	s.RegisterService(&_Queue_serviceDesc, srv)
}

func _Queue_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue.Queue/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).Push(ctx, req.(*PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_Pull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).Pull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue.Queue/Pull",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).Pull(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_AcknowledgePull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcknowledgePullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).AcknowledgePull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue.Queue/AcknowledgePull",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).AcknowledgePull(ctx, req.(*AcknowledgePullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Queue_serviceDesc = grpc.ServiceDesc{
	ServiceName: "queue.Queue",
	HandlerType: (*QueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _Queue_Push_Handler,
		},
		{
			MethodName: "Pull",
			Handler:    _Queue_Pull_Handler,
		},
		{
			MethodName: "AcknowledgePull",
			Handler:    _Queue_AcknowledgePull_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "queue.proto",
}
