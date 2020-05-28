// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/nb/a1/a1-p/tsp/tsp.proto

package tsp

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/onosproject/onos-ric/api/nb/a1/types"
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

// Preference the preference of cell usage
type Preference int32

const (
	Preference_SHALL  Preference = 0
	Preference_PREFER Preference = 1
	Preference_AVOID  Preference = 2
	Preference_FORBID Preference = 3
)

var Preference_name = map[int32]string{
	0: "SHALL",
	1: "PREFER",
	2: "AVOID",
	3: "FORBID",
}

var Preference_value = map[string]int32{
	"SHALL":  0,
	"PREFER": 1,
	"AVOID":  2,
	"FORBID": 3,
}

func (x Preference) String() string {
	return proto.EnumName(Preference_name, int32(x))
}

func (Preference) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_170a5a892e6249bc, []int{0}
}

// TspResources Traffic Steering Preference (TSP) Attributes used to schedule traffic
// on available cells in a different way than what would be through default behavior
type TspResources struct {
	CellIdList []*types.CellID `protobuf:"bytes,1,rep,name=cell_id_list,json=cellIdList,proto3" json:"cell_id_list,omitempty"`
	Preference Preference      `protobuf:"varint,2,opt,name=preference,proto3,enum=a1.tsp.Preference" json:"preference,omitempty"`
	Primary    bool            `protobuf:"varint,3,opt,name=primary,proto3" json:"primary,omitempty"`
}

func (m *TspResources) Reset()         { *m = TspResources{} }
func (m *TspResources) String() string { return proto.CompactTextString(m) }
func (*TspResources) ProtoMessage()    {}
func (*TspResources) Descriptor() ([]byte, []int) {
	return fileDescriptor_170a5a892e6249bc, []int{0}
}
func (m *TspResources) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TspResources) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TspResources.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TspResources) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TspResources.Merge(m, src)
}
func (m *TspResources) XXX_Size() int {
	return m.Size()
}
func (m *TspResources) XXX_DiscardUnknown() {
	xxx_messageInfo_TspResources.DiscardUnknown(m)
}

var xxx_messageInfo_TspResources proto.InternalMessageInfo

func (m *TspResources) GetCellIdList() []*types.CellID {
	if m != nil {
		return m.CellIdList
	}
	return nil
}

func (m *TspResources) GetPreference() Preference {
	if m != nil {
		return m.Preference
	}
	return Preference_SHALL
}

func (m *TspResources) GetPrimary() bool {
	if m != nil {
		return m.Primary
	}
	return false
}

func init() {
	proto.RegisterEnum("a1.tsp.Preference", Preference_name, Preference_value)
	proto.RegisterType((*TspResources)(nil), "a1.tsp.TspResources")
}

func init() { proto.RegisterFile("api/nb/a1/a1-p/tsp/tsp.proto", fileDescriptor_170a5a892e6249bc) }

var fileDescriptor_170a5a892e6249bc = []byte{
	// 289 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x4f, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0xdd, 0xb4, 0x58, 0x75, 0xac, 0xb2, 0xe4, 0xb4, 0x88, 0x84, 0xe2, 0xa9, 0x08, 0x4d, 0x68,
	0xbd, 0xe9, 0xa9, 0xb5, 0x2d, 0x2e, 0x14, 0x5a, 0xa2, 0x78, 0x2d, 0xdb, 0x34, 0x6a, 0x64, 0xdb,
	0x84, 0x24, 0x3d, 0xf4, 0x23, 0x04, 0x3f, 0xcb, 0x63, 0x8f, 0x1e, 0x65, 0xf7, 0x47, 0x64, 0x77,
	0xd1, 0x7a, 0x98, 0x61, 0xe6, 0xcd, 0x7b, 0x8f, 0x79, 0x70, 0x91, 0x18, 0xc5, 0xd6, 0x0b, 0x96,
	0x74, 0x59, 0xd2, 0xed, 0x18, 0xe6, 0x5d, 0x59, 0xd4, 0x58, 0xed, 0x35, 0x6e, 0x24, 0x5d, 0xea,
	0x9d, 0x39, 0xbf, 0x7d, 0x51, 0xfe, 0x75, 0xb3, 0xa0, 0x42, 0xaf, 0x98, 0x5e, 0x6b, 0x67, 0xac,
	0x7e, 0x93, 0xc2, 0x97, 0x73, 0xc7, 0x2a, 0xc1, 0xf6, 0x2e, 0x7e, 0x6b, 0xa4, 0xab, 0x7a, 0x65,
	0x72, 0xf9, 0x8e, 0xa0, 0xf9, 0xe8, 0x0c, 0x97, 0x4e, 0x6f, 0xac, 0x90, 0x0e, 0x33, 0x68, 0x0a,
	0x99, 0xa6, 0x73, 0xb5, 0x9c, 0xa7, 0xca, 0xf9, 0x08, 0xb5, 0xea, 0xed, 0x93, 0xde, 0x29, 0xad,
	0x44, 0x77, 0x32, 0x4d, 0xe3, 0x21, 0x87, 0x82, 0x12, 0x2f, 0x27, 0xca, 0x79, 0xdc, 0x03, 0x30,
	0x56, 0x3e, 0x4b, 0x2b, 0xd7, 0x42, 0x46, 0xb5, 0x16, 0x6a, 0x9f, 0xf5, 0x30, 0xad, 0x7e, 0xa3,
	0xb3, 0xbf, 0x0b, 0xff, 0xc7, 0xc2, 0x11, 0x1c, 0x1a, 0xab, 0x56, 0x89, 0xdd, 0x46, 0xf5, 0x16,
	0x6a, 0x1f, 0xf1, 0xdf, 0xf5, 0xea, 0x06, 0x60, 0xaf, 0xc1, 0xc7, 0x70, 0xf0, 0x70, 0xdf, 0x9f,
	0x4c, 0xc2, 0x00, 0x03, 0x34, 0x66, 0x7c, 0x34, 0x1e, 0xf1, 0x10, 0x15, 0x70, 0xff, 0x69, 0x1a,
	0x0f, 0xc3, 0x5a, 0x01, 0x8f, 0xa7, 0x7c, 0x10, 0x0f, 0xc3, 0xfa, 0x20, 0xfa, 0xcc, 0x08, 0xda,
	0x65, 0x04, 0x7d, 0x67, 0x04, 0x7d, 0xe4, 0x24, 0xd8, 0xe5, 0x24, 0xf8, 0xca, 0x49, 0xb0, 0x68,
	0x94, 0x61, 0xaf, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x26, 0x06, 0x55, 0x36, 0x51, 0x01, 0x00,
	0x00,
}

func (m *TspResources) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TspResources) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TspResources) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Primary {
		i--
		if m.Primary {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.Preference != 0 {
		i = encodeVarintTsp(dAtA, i, uint64(m.Preference))
		i--
		dAtA[i] = 0x10
	}
	if len(m.CellIdList) > 0 {
		for iNdEx := len(m.CellIdList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.CellIdList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTsp(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTsp(dAtA []byte, offset int, v uint64) int {
	offset -= sovTsp(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TspResources) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.CellIdList) > 0 {
		for _, e := range m.CellIdList {
			l = e.Size()
			n += 1 + l + sovTsp(uint64(l))
		}
	}
	if m.Preference != 0 {
		n += 1 + sovTsp(uint64(m.Preference))
	}
	if m.Primary {
		n += 2
	}
	return n
}

func sovTsp(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTsp(x uint64) (n int) {
	return sovTsp(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TspResources) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTsp
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
			return fmt.Errorf("proto: TspResources: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TspResources: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CellIdList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTsp
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
				return ErrInvalidLengthTsp
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTsp
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CellIdList = append(m.CellIdList, &types.CellID{})
			if err := m.CellIdList[len(m.CellIdList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Preference", wireType)
			}
			m.Preference = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTsp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Preference |= Preference(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Primary", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTsp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Primary = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipTsp(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTsp
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTsp
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
func skipTsp(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTsp
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
					return 0, ErrIntOverflowTsp
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
					return 0, ErrIntOverflowTsp
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
				return 0, ErrInvalidLengthTsp
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTsp
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTsp
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTsp        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTsp          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTsp = fmt.Errorf("proto: unexpected end of group")
)
