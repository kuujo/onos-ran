// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/sb/e2ap/e2ap.proto

package e2ap

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	sb "github.com/onosproject/onos-ric/api/sb"
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

type RicControlRequest struct {
	Hdr *e2sm.RicControlHeader  `protobuf:"bytes,1,opt,name=hdr,proto3" json:"hdr,omitempty"`
	Msg *e2sm.RicControlMessage `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (m *RicControlRequest) Reset()         { *m = RicControlRequest{} }
func (m *RicControlRequest) String() string { return proto.CompactTextString(m) }
func (*RicControlRequest) ProtoMessage()    {}
func (*RicControlRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_69a5acb4866023e9, []int{0}
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
	return fileDescriptor_69a5acb4866023e9, []int{1}
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
	proto.RegisterType((*RicControlRequest)(nil), "interface.e2ap.RicControlRequest")
	proto.RegisterType((*RicControlResponse)(nil), "interface.e2ap.RicControlResponse")
}

func init() { proto.RegisterFile("api/sb/e2ap/e2ap.proto", fileDescriptor_69a5acb4866023e9) }

var fileDescriptor_69a5acb4866023e9 = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x92, 0xbd, 0x4e, 0xfb, 0x30,
	0x14, 0xc5, 0xe3, 0xf6, 0xaf, 0xff, 0xe0, 0x0a, 0x24, 0x3c, 0xa0, 0x2a, 0x08, 0xab, 0x64, 0xea,
	0x52, 0xa7, 0x4a, 0x25, 0x24, 0x46, 0xa8, 0x90, 0x18, 0xa8, 0x8a, 0x02, 0x0c, 0x8c, 0x6e, 0x72,
	0x9b, 0x06, 0x35, 0xb1, 0xb1, 0xdd, 0x01, 0x09, 0x26, 0x5e, 0x80, 0xc7, 0x62, 0xec, 0xc8, 0x88,
	0xda, 0x17, 0x41, 0x49, 0xa0, 0x69, 0xf8, 0x12, 0x0b, 0xcb, 0x95, 0xe5, 0x7b, 0x7e, 0xe7, 0xf8,
	0xe3, 0xe2, 0x6d, 0x2e, 0x63, 0x57, 0x8f, 0x5c, 0xf0, 0xb8, 0xcc, 0x0b, 0x93, 0x4a, 0x18, 0x41,
	0x36, 0xe3, 0xd4, 0x80, 0x1a, 0xf3, 0x00, 0x58, 0xb6, 0x6b, 0x1f, 0x44, 0xb1, 0x99, 0xcc, 0x46,
	0x2c, 0x10, 0x89, 0x2b, 0x52, 0xa1, 0xa5, 0x12, 0xd7, 0x10, 0x98, 0x7c, 0xdd, 0x51, 0x71, 0xe0,
	0xae, 0x7c, 0x3a, 0x25, 0x99, 0x5b, 0xd9, 0xfb, 0xbf, 0x46, 0x75, 0x92, 0x97, 0x82, 0x73, 0xee,
	0xf0, 0x96, 0x1f, 0x07, 0x7d, 0x91, 0x1a, 0x25, 0xa6, 0x3e, 0xdc, 0xcc, 0x40, 0x1b, 0xe2, 0xe1,
	0xfa, 0x24, 0x54, 0x4d, 0xd4, 0x42, 0xed, 0x86, 0xd7, 0x62, 0xeb, 0xa7, 0xd4, 0x09, 0x2b, 0xf5,
	0x27, 0xc0, 0x43, 0x50, 0x7e, 0x26, 0x26, 0x3d, 0x5c, 0x4f, 0x74, 0xd4, 0xac, 0xe5, 0xcc, 0xde,
	0xf7, 0xcc, 0x00, 0xb4, 0xe6, 0x11, 0xf8, 0x99, 0xda, 0xb9, 0xc7, 0x64, 0x3d, 0x5d, 0x4b, 0x91,
	0x6a, 0xf8, 0xd3, 0xf8, 0xe1, 0xcc, 0x04, 0x22, 0x29, 0xe2, 0xbd, 0x87, 0x1a, 0xfe, 0x77, 0xec,
	0x1d, 0x9e, 0x91, 0x2b, 0x8c, 0x4b, 0x09, 0xa9, 0xe2, 0x5c, 0xb2, 0x4f, 0x2f, 0x64, 0x3b, 0x3f,
	0x49, 0x8a, 0x6b, 0x38, 0x56, 0x1b, 0x75, 0x11, 0x19, 0xe2, 0xc6, 0x39, 0xa4, 0xe1, 0xbb, 0xf7,
	0x6e, 0x05, 0x64, 0x1f, 0x18, 0x7b, 0xe7, 0xcb, 0xf6, 0xa5, 0x0c, 0xb9, 0x29, 0x0d, 0x37, 0x32,
	0xc3, 0x0b, 0x98, 0x42, 0x02, 0x46, 0xdd, 0x12, 0xbb, 0xca, 0x9c, 0x7a, 0x03, 0xe0, 0xba, 0x2f,
	0xd2, 0x71, 0x1c, 0xd9, 0xb4, 0xda, 0x5b, 0x41, 0x6f, 0xbf, 0xe0, 0x58, 0x5d, 0x74, 0xd4, 0x7c,
	0x5a, 0x50, 0x34, 0x5f, 0x50, 0xf4, 0xb2, 0xa0, 0xe8, 0x71, 0x49, 0xad, 0xf9, 0x92, 0x5a, 0xcf,
	0x4b, 0x6a, 0x8d, 0xfe, 0xe7, 0x33, 0xd2, 0x7b, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x2a, 0x63, 0x6f,
	0x4a, 0xc0, 0x02, 0x00, 0x00,
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
	//
	//ricControl E2AP-ELEMENTARY-PROCEDURE ::= {
	//INITIATING MESSAGE        RICcontrolRequest
	//SUCCESSFUL OUTCOME        RICcontrolAcknowledge
	//UNSUCCESSFUL OUTCOME    RICcontrolFailure
	//PROCEDURE CODE            id-RICcontrol
	//CRITICALITY                reject
	//}
	RicControl(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicControlClient, error)
	// TODO - convert rest of the services to new E2
	SendControl(ctx context.Context, opts ...grpc.CallOption) (E2AP_SendControlClient, error)
	SendTelemetry(ctx context.Context, in *sb.L2MeasConfig, opts ...grpc.CallOption) (E2AP_SendTelemetryClient, error)
}

type e2APClient struct {
	cc *grpc.ClientConn
}

func NewE2APClient(cc *grpc.ClientConn) E2APClient {
	return &e2APClient{cc}
}

func (c *e2APClient) RicControl(ctx context.Context, opts ...grpc.CallOption) (E2AP_RicControlClient, error) {
	stream, err := c.cc.NewStream(ctx, &_E2AP_serviceDesc.Streams[0], "/interface.e2ap.E2AP/RicControl", opts...)
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

func (c *e2APClient) SendControl(ctx context.Context, opts ...grpc.CallOption) (E2AP_SendControlClient, error) {
	stream, err := c.cc.NewStream(ctx, &_E2AP_serviceDesc.Streams[1], "/interface.e2ap.E2AP/SendControl", opts...)
	if err != nil {
		return nil, err
	}
	x := &e2APSendControlClient{stream}
	return x, nil
}

type E2AP_SendControlClient interface {
	Send(*sb.ControlResponse) error
	Recv() (*sb.ControlUpdate, error)
	grpc.ClientStream
}

type e2APSendControlClient struct {
	grpc.ClientStream
}

func (x *e2APSendControlClient) Send(m *sb.ControlResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *e2APSendControlClient) Recv() (*sb.ControlUpdate, error) {
	m := new(sb.ControlUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *e2APClient) SendTelemetry(ctx context.Context, in *sb.L2MeasConfig, opts ...grpc.CallOption) (E2AP_SendTelemetryClient, error) {
	stream, err := c.cc.NewStream(ctx, &_E2AP_serviceDesc.Streams[2], "/interface.e2ap.E2AP/SendTelemetry", opts...)
	if err != nil {
		return nil, err
	}
	x := &e2APSendTelemetryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type E2AP_SendTelemetryClient interface {
	Recv() (*sb.TelemetryMessage, error)
	grpc.ClientStream
}

type e2APSendTelemetryClient struct {
	grpc.ClientStream
}

func (x *e2APSendTelemetryClient) Recv() (*sb.TelemetryMessage, error) {
	m := new(sb.TelemetryMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// E2APServer is the server API for E2AP service.
type E2APServer interface {
	//
	//ricControl E2AP-ELEMENTARY-PROCEDURE ::= {
	//INITIATING MESSAGE        RICcontrolRequest
	//SUCCESSFUL OUTCOME        RICcontrolAcknowledge
	//UNSUCCESSFUL OUTCOME    RICcontrolFailure
	//PROCEDURE CODE            id-RICcontrol
	//CRITICALITY                reject
	//}
	RicControl(E2AP_RicControlServer) error
	// TODO - convert rest of the services to new E2
	SendControl(E2AP_SendControlServer) error
	SendTelemetry(*sb.L2MeasConfig, E2AP_SendTelemetryServer) error
}

// UnimplementedE2APServer can be embedded to have forward compatible implementations.
type UnimplementedE2APServer struct {
}

func (*UnimplementedE2APServer) RicControl(srv E2AP_RicControlServer) error {
	return status.Errorf(codes.Unimplemented, "method RicControl not implemented")
}
func (*UnimplementedE2APServer) SendControl(srv E2AP_SendControlServer) error {
	return status.Errorf(codes.Unimplemented, "method SendControl not implemented")
}
func (*UnimplementedE2APServer) SendTelemetry(req *sb.L2MeasConfig, srv E2AP_SendTelemetryServer) error {
	return status.Errorf(codes.Unimplemented, "method SendTelemetry not implemented")
}

func RegisterE2APServer(s *grpc.Server, srv E2APServer) {
	s.RegisterService(&_E2AP_serviceDesc, srv)
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

func _E2AP_SendControl_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(E2APServer).SendControl(&e2APSendControlServer{stream})
}

type E2AP_SendControlServer interface {
	Send(*sb.ControlUpdate) error
	Recv() (*sb.ControlResponse, error)
	grpc.ServerStream
}

type e2APSendControlServer struct {
	grpc.ServerStream
}

func (x *e2APSendControlServer) Send(m *sb.ControlUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func (x *e2APSendControlServer) Recv() (*sb.ControlResponse, error) {
	m := new(sb.ControlResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _E2AP_SendTelemetry_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(sb.L2MeasConfig)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(E2APServer).SendTelemetry(m, &e2APSendTelemetryServer{stream})
}

type E2AP_SendTelemetryServer interface {
	Send(*sb.TelemetryMessage) error
	grpc.ServerStream
}

type e2APSendTelemetryServer struct {
	grpc.ServerStream
}

func (x *e2APSendTelemetryServer) Send(m *sb.TelemetryMessage) error {
	return x.ServerStream.SendMsg(m)
}

var _E2AP_serviceDesc = grpc.ServiceDesc{
	ServiceName: "interface.e2ap.E2AP",
	HandlerType: (*E2APServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RicControl",
			Handler:       _E2AP_RicControl_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SendControl",
			Handler:       _E2AP_SendControl_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SendTelemetry",
			Handler:       _E2AP_SendTelemetry_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/sb/e2ap/e2ap.proto",
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
