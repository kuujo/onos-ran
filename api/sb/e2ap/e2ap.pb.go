// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/sb/e2ap/e2ap.proto

package e2ap

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	e2sm "github.com/onosproject/onos-ric/api/sb/e2sm"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type RicSubscriptionRequest struct {
	RanFunctionID          e2sm.RanFunctionID      `protobuf:"varint,1,opt,name=ranFunctionID,proto3,enum=interface.e2sm.RanFunctionID" json:"ranFunctionID,omitempty"`
	RicSubscriptionDetails *RicSubscriptionDetails `protobuf:"bytes,2,opt,name=ricSubscriptionDetails,proto3" json:"ricSubscriptionDetails,omitempty"`
}

func (m *RicSubscriptionRequest) Reset()         { *m = RicSubscriptionRequest{} }
func (m *RicSubscriptionRequest) String() string { return proto.CompactTextString(m) }
func (*RicSubscriptionRequest) ProtoMessage()    {}
func (*RicSubscriptionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_69a5acb4866023e9, []int{0}
}
func (m *RicSubscriptionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RicSubscriptionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RicSubscriptionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RicSubscriptionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RicSubscriptionRequest.Merge(m, src)
}
func (m *RicSubscriptionRequest) XXX_Size() int {
	return m.Size()
}
func (m *RicSubscriptionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RicSubscriptionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RicSubscriptionRequest proto.InternalMessageInfo

func (m *RicSubscriptionRequest) GetRanFunctionID() e2sm.RanFunctionID {
	if m != nil {
		return m.RanFunctionID
	}
	return e2sm.RanFunctionID_RAN_FUNCTION_ID_INVALID
}

func (m *RicSubscriptionRequest) GetRicSubscriptionDetails() *RicSubscriptionDetails {
	if m != nil {
		return m.RicSubscriptionDetails
	}
	return nil
}

type RicSubscriptionDetails struct {
	RicEventTriggerDefinition *e2sm.RicEventTriggerDefinition `protobuf:"bytes,1,opt,name=ricEventTriggerDefinition,proto3" json:"ricEventTriggerDefinition,omitempty"`
}

func (m *RicSubscriptionDetails) Reset()         { *m = RicSubscriptionDetails{} }
func (m *RicSubscriptionDetails) String() string { return proto.CompactTextString(m) }
func (*RicSubscriptionDetails) ProtoMessage()    {}
func (*RicSubscriptionDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_69a5acb4866023e9, []int{1}
}
func (m *RicSubscriptionDetails) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RicSubscriptionDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RicSubscriptionDetails.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RicSubscriptionDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RicSubscriptionDetails.Merge(m, src)
}
func (m *RicSubscriptionDetails) XXX_Size() int {
	return m.Size()
}
func (m *RicSubscriptionDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_RicSubscriptionDetails.DiscardUnknown(m)
}

var xxx_messageInfo_RicSubscriptionDetails proto.InternalMessageInfo

func (m *RicSubscriptionDetails) GetRicEventTriggerDefinition() *e2sm.RicEventTriggerDefinition {
	if m != nil {
		return m.RicEventTriggerDefinition
	}
	return nil
}

type RicIndication struct {
	Hdr *e2sm.RicIndicationHeader  `protobuf:"bytes,1,opt,name=hdr,proto3" json:"hdr,omitempty"`
	Msg *e2sm.RicIndicationMessage `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (m *RicIndication) Reset()         { *m = RicIndication{} }
func (m *RicIndication) String() string { return proto.CompactTextString(m) }
func (*RicIndication) ProtoMessage()    {}
func (*RicIndication) Descriptor() ([]byte, []int) {
	return fileDescriptor_69a5acb4866023e9, []int{2}
}
func (m *RicIndication) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RicIndication) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RicIndication.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RicIndication) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RicIndication.Merge(m, src)
}
func (m *RicIndication) XXX_Size() int {
	return m.Size()
}
func (m *RicIndication) XXX_DiscardUnknown() {
	xxx_messageInfo_RicIndication.DiscardUnknown(m)
}

var xxx_messageInfo_RicIndication proto.InternalMessageInfo

func (m *RicIndication) GetHdr() *e2sm.RicIndicationHeader {
	if m != nil {
		return m.Hdr
	}
	return nil
}

func (m *RicIndication) GetMsg() *e2sm.RicIndicationMessage {
	if m != nil {
		return m.Msg
	}
	return nil
}

type RicControlRequest struct {
	Hdr *e2sm.RicControlHeader  `protobuf:"bytes,1,opt,name=hdr,proto3" json:"hdr,omitempty"`
	Msg *e2sm.RicControlMessage `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (m *RicControlRequest) Reset()         { *m = RicControlRequest{} }
func (m *RicControlRequest) String() string { return proto.CompactTextString(m) }
func (*RicControlRequest) ProtoMessage()    {}
func (*RicControlRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_69a5acb4866023e9, []int{3}
}
func (m *RicControlRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RicControlRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RicControlRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RicControlRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RicControlRequest.Merge(m, src)
}
func (m *RicControlRequest) XXX_Size() int {
	return m.Size()
}
func (m *RicControlRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RicControlRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RicControlRequest proto.InternalMessageInfo

func (m *RicControlRequest) GetHdr() *e2sm.RicControlHeader {
	if m != nil {
		return m.Hdr
	}
	return nil
}

func (m *RicControlRequest) GetMsg() *e2sm.RicControlMessage {
	if m != nil {
		return m.Msg
	}
	return nil
}

type RicControlResponse struct {
	Hdr *e2sm.RicControlHeader  `protobuf:"bytes,1,opt,name=hdr,proto3" json:"hdr,omitempty"`
	Msg *e2sm.RicControlOutcome `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (m *RicControlResponse) Reset()         { *m = RicControlResponse{} }
func (m *RicControlResponse) String() string { return proto.CompactTextString(m) }
func (*RicControlResponse) ProtoMessage()    {}
func (*RicControlResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_69a5acb4866023e9, []int{4}
}
func (m *RicControlResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RicControlResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RicControlResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RicControlResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RicControlResponse.Merge(m, src)
}
func (m *RicControlResponse) XXX_Size() int {
	return m.Size()
}
func (m *RicControlResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RicControlResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RicControlResponse proto.InternalMessageInfo

func (m *RicControlResponse) GetHdr() *e2sm.RicControlHeader {
	if m != nil {
		return m.Hdr
	}
	return nil
}

func (m *RicControlResponse) GetMsg() *e2sm.RicControlOutcome {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterType((*RicSubscriptionRequest)(nil), "interface.e2ap.RicSubscriptionRequest")
	proto.RegisterType((*RicSubscriptionDetails)(nil), "interface.e2ap.RicSubscriptionDetails")
	proto.RegisterType((*RicIndication)(nil), "interface.e2ap.RicIndication")
	proto.RegisterType((*RicControlRequest)(nil), "interface.e2ap.RicControlRequest")
	proto.RegisterType((*RicControlResponse)(nil), "interface.e2ap.RicControlResponse")
}

func init() { proto.RegisterFile("api/sb/e2ap/e2ap.proto", fileDescriptor_69a5acb4866023e9) }

var fileDescriptor_69a5acb4866023e9 = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0x3f, 0x8f, 0xd3, 0x30,
	0x14, 0x8f, 0xef, 0x10, 0x48, 0x3e, 0xdd, 0x21, 0x3c, 0x54, 0xa5, 0x52, 0xa3, 0x36, 0x20, 0x54,
	0x06, 0x12, 0x94, 0x8a, 0xee, 0xd0, 0x16, 0xd1, 0x01, 0x01, 0x86, 0x85, 0xa5, 0x92, 0xe3, 0xba,
	0xa9, 0x51, 0x63, 0x07, 0xdb, 0x61, 0x02, 0x24, 0xbe, 0x01, 0x03, 0x5f, 0x87, 0x9d, 0xb1, 0x23,
	0x23, 0x6a, 0xbf, 0x08, 0xca, 0x1f, 0x68, 0x42, 0x93, 0x02, 0xc3, 0x2d, 0x4f, 0xd1, 0xcb, 0xef,
	0xdf, 0x7b, 0x79, 0x81, 0x2d, 0x12, 0x73, 0x4f, 0x07, 0x1e, 0xf3, 0x49, 0x9c, 0x15, 0x37, 0x56,
	0xd2, 0x48, 0x74, 0xc1, 0x85, 0x61, 0x6a, 0x49, 0x28, 0x73, 0xd3, 0x6e, 0x67, 0x14, 0x72, 0xb3,
	0x4a, 0x02, 0x97, 0xca, 0xc8, 0x93, 0x42, 0xea, 0x58, 0xc9, 0x37, 0x8c, 0x9a, 0xec, 0xf9, 0x9e,
	0xe2, 0xd4, 0xfb, 0xad, 0xa3, 0xa3, 0xac, 0xe4, 0x3a, 0xce, 0x57, 0x00, 0x5b, 0x98, 0xd3, 0x97,
	0x49, 0xa0, 0xa9, 0xe2, 0xb1, 0xe1, 0x52, 0x60, 0xf6, 0x36, 0x61, 0xda, 0xa0, 0x31, 0x3c, 0x57,
	0x44, 0x3c, 0x4e, 0x04, 0x4d, 0xbb, 0xb3, 0x49, 0x1b, 0xf4, 0xc0, 0xe0, 0xc2, 0xef, 0xba, 0x65,
	0x6b, 0x1d, 0xb9, 0xb8, 0x0c, 0xc2, 0x55, 0x0e, 0x9a, 0xc3, 0x96, 0xaa, 0xca, 0x4f, 0x98, 0x21,
	0x7c, 0xad, 0xdb, 0x27, 0x3d, 0x30, 0x38, 0xf3, 0xef, 0xb8, 0xd5, 0x41, 0x5c, 0x5c, 0x8b, 0xc6,
	0x0d, 0x2a, 0xce, 0xa7, 0xc3, 0xfc, 0xc5, 0x2b, 0x14, 0xc2, 0x9b, 0x8a, 0xd3, 0xe9, 0x3b, 0x26,
	0xcc, 0x2b, 0xc5, 0xc3, 0x90, 0xa9, 0x09, 0x5b, 0x72, 0xc1, 0x53, 0x4c, 0x36, 0xcb, 0x99, 0x7f,
	0xf7, 0x60, 0x96, 0x26, 0x02, 0x6e, 0xd6, 0x72, 0x3e, 0xc2, 0x73, 0xcc, 0xe9, 0x4c, 0x2c, 0x38,
	0x25, 0x69, 0x03, 0x3d, 0x80, 0xa7, 0xab, 0x85, 0x2a, 0x3c, 0x6e, 0xd5, 0x78, 0xec, 0xb1, 0x4f,
	0x18, 0x59, 0x30, 0x85, 0x53, 0x3c, 0x1a, 0xc1, 0xd3, 0x48, 0x87, 0xc5, 0x62, 0x6e, 0x1f, 0xa5,
	0x3d, 0x65, 0x5a, 0x93, 0x90, 0xe1, 0x94, 0xe0, 0xbc, 0x87, 0x37, 0x30, 0xa7, 0x63, 0x29, 0x8c,
	0x92, 0xeb, 0x5f, 0x5f, 0xcf, 0x2f, 0x67, 0xe8, 0xd5, 0x88, 0x15, 0xf8, 0x72, 0x80, 0x61, 0x39,
	0x40, 0xbf, 0x99, 0x53, 0x71, 0xff, 0x00, 0x51, 0xd9, 0x5d, 0xc7, 0x52, 0x68, 0x76, 0xa9, 0xf6,
	0xcf, 0x12, 0x43, 0x65, 0x94, 0xdb, 0xfb, 0x5f, 0x4e, 0xe0, 0x95, 0xa9, 0xff, 0xf0, 0x39, 0x7a,
	0x01, 0xaf, 0xa5, 0x90, 0x15, 0x11, 0xa8, 0x5f, 0x73, 0x54, 0xd5, 0xf5, 0x74, 0xba, 0x35, 0x90,
	0xfd, 0x7a, 0x1d, 0x6b, 0x00, 0xee, 0x03, 0x34, 0x87, 0xd7, 0xff, 0xb8, 0x2d, 0xf4, 0xb7, 0x7b,
	0xfd, 0x2f, 0xfd, 0xd7, 0x10, 0xee, 0x93, 0xfd, 0x4b, 0x6a, 0xe7, 0x18, 0x24, 0xdf, 0x7c, 0x2e,
	0xfd, 0xa8, 0xfd, 0x6d, 0x6b, 0x83, 0xcd, 0xd6, 0x06, 0x3f, 0xb6, 0x36, 0xf8, 0xbc, 0xb3, 0xad,
	0xcd, 0xce, 0xb6, 0xbe, 0xef, 0x6c, 0x2b, 0xb8, 0x9a, 0xfd, 0xf8, 0xc3, 0x9f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x3e, 0xf1, 0xa0, 0x5e, 0x5a, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// E2APClient is the client API for E2AP service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type E2APClient interface {
	// RicChan is a bi-directonal stream for all messaging between RIC and E2-Node
	RicChan(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicChanClient, error)
	// ----
	// Below two rpcs are to be removed
	RicSubscription(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicSubscriptionClient, error)
	RicControl(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicControlClient, error)
}

type e2APClient struct {
	cc *grpc.ClientConn
}

func NewE2APClient(cc *grpc.ClientConn) E2APClient {
	return &e2APClient{cc}
}

func (c *e2APClient) RicChan(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicChanClient, error) {
	stream, err := c.cc.NewStream(ctx, &_E2AP_serviceDesc.Streams[0], "/interface.e2ap.E2AP/RicChan", opts...)
	if err != nil {
		return nil, err
	}
	x := &e2APRicChanClient{stream}
	return x, nil
}

type E2AP_RicChanClient interface {
	Send(*RicControlRequest) error
	Recv() (*RicIndication, error)
	grpc.ClientStream
}

type e2APRicChanClient struct {
	grpc.ClientStream
}

func (x *e2APRicChanClient) Send(m *RicControlRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *e2APRicChanClient) Recv() (*RicIndication, error) {
	m := new(RicIndication)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *e2APClient) RicSubscription(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicSubscriptionClient, error) {
	stream, err := c.cc.NewStream(ctx, &_E2AP_serviceDesc.Streams[1], "/interface.e2ap.E2AP/RicSubscription", opts...)
	if err != nil {
		return nil, err
	}
	x := &e2APRicSubscriptionClient{stream}
	return x, nil
}

type E2AP_RicSubscriptionClient interface {
	Send(*RicSubscriptionRequest) error
	Recv() (*RicIndication, error)
	grpc.ClientStream
}

type e2APRicSubscriptionClient struct {
	grpc.ClientStream
}

func (x *e2APRicSubscriptionClient) Send(m *RicSubscriptionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *e2APRicSubscriptionClient) Recv() (*RicIndication, error) {
	m := new(RicIndication)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *e2APClient) RicControl(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicControlClient, error) {
	stream, err := c.cc.NewStream(ctx, &_E2AP_serviceDesc.Streams[2], "/interface.e2ap.E2AP/RicControl", opts...)
	if err != nil {
		return nil, err
	}
	x := &e2APRicControlClient{stream}
	return x, nil
}

type E2AP_RicControlClient interface {
	Send(*RicControlRequest) error
	Recv() (*RicControlResponse, error)
	grpc.ClientStream
}

type e2APRicControlClient struct {
	grpc.ClientStream
}

func (x *e2APRicControlClient) Send(m *RicControlRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *e2APRicControlClient) Recv() (*RicControlResponse, error) {
	m := new(RicControlResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// E2APServer is the server API for E2AP service.
type E2APServer interface {
	// RicChan is a bi-directonal stream for all messaging between RIC and E2-Node
	RicChan(E2AP_RicChanServer) error
	// ----
	// Below two rpcs are to be removed
	RicSubscription(E2AP_RicSubscriptionServer) error
	RicControl(E2AP_RicControlServer) error
}

// UnimplementedE2APServer can be embedded to have forward compatible implementations.
type UnimplementedE2APServer struct {
}

func (*UnimplementedE2APServer) RicChan(srv E2AP_RicChanServer) error {
	return status.Errorf(codes.Unimplemented, "method RicChan not implemented")
}
func (*UnimplementedE2APServer) RicSubscription(srv E2AP_RicSubscriptionServer) error {
	return status.Errorf(codes.Unimplemented, "method RicSubscription not implemented")
}
func (*UnimplementedE2APServer) RicControl(srv E2AP_RicControlServer) error {
	return status.Errorf(codes.Unimplemented, "method RicControl not implemented")
}

func RegisterE2APServer(s *grpc.Server, srv E2APServer) {
	s.RegisterService(&_E2AP_serviceDesc, srv)
}

func _E2AP_RicChan_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(E2APServer).RicChan(&e2APRicChanServer{stream})
}

type E2AP_RicChanServer interface {
	Send(*RicIndication) error
	Recv() (*RicControlRequest, error)
	grpc.ServerStream
}

type e2APRicChanServer struct {
	grpc.ServerStream
}

func (x *e2APRicChanServer) Send(m *RicIndication) error {
	return x.ServerStream.SendMsg(m)
}

func (x *e2APRicChanServer) Recv() (*RicControlRequest, error) {
	m := new(RicControlRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _E2AP_RicSubscription_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(E2APServer).RicSubscription(&e2APRicSubscriptionServer{stream})
}

type E2AP_RicSubscriptionServer interface {
	Send(*RicIndication) error
	Recv() (*RicSubscriptionRequest, error)
	grpc.ServerStream
}

type e2APRicSubscriptionServer struct {
	grpc.ServerStream
}

func (x *e2APRicSubscriptionServer) Send(m *RicIndication) error {
	return x.ServerStream.SendMsg(m)
}

func (x *e2APRicSubscriptionServer) Recv() (*RicSubscriptionRequest, error) {
	m := new(RicSubscriptionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _E2AP_RicControl_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(E2APServer).RicControl(&e2APRicControlServer{stream})
}

type E2AP_RicControlServer interface {
	Send(*RicControlResponse) error
	Recv() (*RicControlRequest, error)
	grpc.ServerStream
}

type e2APRicControlServer struct {
	grpc.ServerStream
}

func (x *e2APRicControlServer) Send(m *RicControlResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *e2APRicControlServer) Recv() (*RicControlRequest, error) {
	m := new(RicControlRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _E2AP_serviceDesc = grpc.ServiceDesc{
	ServiceName: "interface.e2ap.E2AP",
	HandlerType: (*E2APServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RicChan",
			Handler:       _E2AP_RicChan_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "RicSubscription",
			Handler:       _E2AP_RicSubscription_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "RicControl",
			Handler:       _E2AP_RicControl_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/sb/e2ap/e2ap.proto",
}

func (m *RicSubscriptionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RicSubscriptionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RicSubscriptionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RicSubscriptionDetails != nil {
		{
			size, err := m.RicSubscriptionDetails.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RanFunctionID != 0 {
		i = encodeVarintE2Ap(dAtA, i, uint64(m.RanFunctionID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RicSubscriptionDetails) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RicSubscriptionDetails) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RicSubscriptionDetails) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RicEventTriggerDefinition != nil {
		{
			size, err := m.RicEventTriggerDefinition.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RicIndication) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RicIndication) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RicIndication) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Msg != nil {
		{
			size, err := m.Msg.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Hdr != nil {
		{
			size, err := m.Hdr.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RicControlRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RicControlRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RicControlRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Msg != nil {
		{
			size, err := m.Msg.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Hdr != nil {
		{
			size, err := m.Hdr.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RicControlResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RicControlResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RicControlResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Msg != nil {
		{
			size, err := m.Msg.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Hdr != nil {
		{
			size, err := m.Hdr.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintE2Ap(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintE2Ap(dAtA []byte, offset int, v uint64) int {
	offset -= sovE2Ap(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RicSubscriptionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RanFunctionID != 0 {
		n += 1 + sovE2Ap(uint64(m.RanFunctionID))
	}
	if m.RicSubscriptionDetails != nil {
		l = m.RicSubscriptionDetails.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	return n
}

func (m *RicSubscriptionDetails) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RicEventTriggerDefinition != nil {
		l = m.RicEventTriggerDefinition.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	return n
}

func (m *RicIndication) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Hdr != nil {
		l = m.Hdr.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	if m.Msg != nil {
		l = m.Msg.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	return n
}

func (m *RicControlRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Hdr != nil {
		l = m.Hdr.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	if m.Msg != nil {
		l = m.Msg.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	return n
}

func (m *RicControlResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Hdr != nil {
		l = m.Hdr.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	if m.Msg != nil {
		l = m.Msg.Size()
		n += 1 + l + sovE2Ap(uint64(l))
	}
	return n
}

func sovE2Ap(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozE2Ap(x uint64) (n int) {
	return sovE2Ap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RicSubscriptionRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowE2Ap
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RicSubscriptionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RicSubscriptionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RanFunctionID", wireType)
			}
			m.RanFunctionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RanFunctionID |= e2sm.RanFunctionID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RicSubscriptionDetails", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RicSubscriptionDetails == nil {
				m.RicSubscriptionDetails = &RicSubscriptionDetails{}
			}
			if err := m.RicSubscriptionDetails.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipE2Ap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RicSubscriptionDetails) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowE2Ap
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RicSubscriptionDetails: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RicSubscriptionDetails: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RicEventTriggerDefinition", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RicEventTriggerDefinition == nil {
				m.RicEventTriggerDefinition = &e2sm.RicEventTriggerDefinition{}
			}
			if err := m.RicEventTriggerDefinition.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipE2Ap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RicIndication) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowE2Ap
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RicIndication: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RicIndication: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hdr", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Hdr == nil {
				m.Hdr = &e2sm.RicIndicationHeader{}
			}
			if err := m.Hdr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Msg == nil {
				m.Msg = &e2sm.RicIndicationMessage{}
			}
			if err := m.Msg.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipE2Ap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RicControlRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowE2Ap
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RicControlRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RicControlRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hdr", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Hdr == nil {
				m.Hdr = &e2sm.RicControlHeader{}
			}
			if err := m.Hdr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Msg == nil {
				m.Msg = &e2sm.RicControlMessage{}
			}
			if err := m.Msg.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipE2Ap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RicControlResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowE2Ap
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RicControlResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RicControlResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hdr", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Hdr == nil {
				m.Hdr = &e2sm.RicControlHeader{}
			}
			if err := m.Hdr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthE2Ap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthE2Ap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Msg == nil {
				m.Msg = &e2sm.RicControlOutcome{}
			}
			if err := m.Msg.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipE2Ap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthE2Ap
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipE2Ap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowE2Ap
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowE2Ap
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthE2Ap
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupE2Ap
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthE2Ap
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthE2Ap        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowE2Ap          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupE2Ap = fmt.Errorf("proto: unexpected end of group")
)
