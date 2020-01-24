// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/sb/sb.proto

package sb

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func init() { proto.RegisterFile("api/sb/sb.proto", fileDescriptor_b10b7e1317dc95eb) }

var fileDescriptor_b10b7e1317dc95eb = []byte{
	// 88 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x2c, 0xc8, 0xd4,
	0x2f, 0x4e, 0xd2, 0x2f, 0x4e, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x4a, 0xcc,
	0xd3, 0x2b, 0x4e, 0x32, 0x62, 0xe1, 0x62, 0x72, 0x35, 0x72, 0x92, 0x38, 0xf1, 0x48, 0x8e, 0xf1,
	0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e,
	0xe1, 0xc6, 0x63, 0x39, 0x86, 0x24, 0x36, 0xb0, 0x72, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x88, 0xe6, 0xc0, 0x29, 0x41, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// E2Client is the client API for E2 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type E2Client interface {
}

type e2Client struct {
	cc *grpc.ClientConn
}

func NewE2Client(cc *grpc.ClientConn) E2Client {
	return &e2Client{cc}
}

// E2Server is the server API for E2 service.
type E2Server interface {
}

// UnimplementedE2Server can be embedded to have forward compatible implementations.
type UnimplementedE2Server struct {
}

func RegisterE2Server(s *grpc.Server, srv E2Server) {
	s.RegisterService(&_E2_serviceDesc, srv)
}

var _E2_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ran.sb.E2",
	HandlerType: (*E2Server)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "api/sb/sb.proto",
}