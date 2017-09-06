// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/claim.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Claim_Version int32

const (
	Claim_UNKNOWN_VERSION Claim_Version = 0
	Claim__0_0_1          Claim_Version = 1
)

var Claim_Version_name = map[int32]string{
	0: "UNKNOWN_VERSION",
	1: "_0_0_1",
}
var Claim_Version_value = map[string]int32{
	"UNKNOWN_VERSION": 0,
	"_0_0_1":          1,
}

func (x Claim_Version) Enum() *Claim_Version {
	p := new(Claim_Version)
	*p = x
	return p
}
func (x Claim_Version) String() string {
	return proto.EnumName(Claim_Version_name, int32(x))
}
func (x *Claim_Version) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Claim_Version_value, data, "Claim_Version")
	if err != nil {
		return err
	}
	*x = Claim_Version(value)
	return nil
}
func (Claim_Version) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

type Claim_ClaimType int32

const (
	Claim_UNKNOWN_CLAIM_TYPE Claim_ClaimType = 0
	Claim_streamType         Claim_ClaimType = 1
	Claim_certificateType    Claim_ClaimType = 2
)

var Claim_ClaimType_name = map[int32]string{
	0: "UNKNOWN_CLAIM_TYPE",
	1: "streamType",
	2: "certificateType",
}
var Claim_ClaimType_value = map[string]int32{
	"UNKNOWN_CLAIM_TYPE": 0,
	"streamType":         1,
	"certificateType":    2,
}

func (x Claim_ClaimType) Enum() *Claim_ClaimType {
	p := new(Claim_ClaimType)
	*p = x
	return p
}
func (x Claim_ClaimType) String() string {
	return proto.EnumName(Claim_ClaimType_name, int32(x))
}
func (x *Claim_ClaimType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Claim_ClaimType_value, data, "Claim_ClaimType")
	if err != nil {
		return err
	}
	*x = Claim_ClaimType(value)
	return nil
}
func (Claim_ClaimType) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 1} }

type Claim struct {
	Version            *Claim_Version   `protobuf:"varint,1,req,name=version,enum=pb.Claim_Version" json:"version,omitempty"`
	ClaimType          *Claim_ClaimType `protobuf:"varint,2,req,name=claimType,enum=pb.Claim_ClaimType" json:"claimType,omitempty"`
	Stream             *Stream          `protobuf:"bytes,3,opt,name=stream" json:"stream,omitempty"`
	Certificate        *Certificate     `protobuf:"bytes,4,opt,name=certificate" json:"certificate,omitempty"`
	PublisherSignature *Signature       `protobuf:"bytes,5,opt,name=publisherSignature" json:"publisherSignature,omitempty"`
	XXX_unrecognized   []byte           `json:"-"`
}

func (m *Claim) Reset()                    { *m = Claim{} }
func (m *Claim) String() string            { return proto.CompactTextString(m) }
func (*Claim) ProtoMessage()               {}
func (*Claim) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Claim) GetVersion() Claim_Version {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return Claim_UNKNOWN_VERSION
}

func (m *Claim) GetClaimType() Claim_ClaimType {
	if m != nil && m.ClaimType != nil {
		return *m.ClaimType
	}
	return Claim_UNKNOWN_CLAIM_TYPE
}

func (m *Claim) GetStream() *Stream {
	if m != nil {
		return m.Stream
	}
	return nil
}

func (m *Claim) GetCertificate() *Certificate {
	if m != nil {
		return m.Certificate
	}
	return nil
}

func (m *Claim) GetPublisherSignature() *Signature {
	if m != nil {
		return m.PublisherSignature
	}
	return nil
}

func init() {
	proto.RegisterType((*Claim)(nil), "pb.Claim")
	proto.RegisterEnum("pb.Claim_Version", Claim_Version_name, Claim_Version_value)
	proto.RegisterEnum("pb.Claim_ClaimType", Claim_ClaimType_name, Claim_ClaimType_value)
}

func init() { proto.RegisterFile("pb/claim.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0x86, 0xd7, 0xea, 0x36, 0xf6, 0x0d, 0xdb, 0xfa, 0x4d, 0xa4, 0xec, 0x34, 0x7a, 0x1a, 0x0a,
	0xdd, 0xba, 0xbb, 0x07, 0x29, 0x03, 0x87, 0xda, 0x49, 0x3a, 0x27, 0x9e, 0x4a, 0x3b, 0xa2, 0x06,
	0xe6, 0x1a, 0xd2, 0x4e, 0xf0, 0x77, 0xfb, 0x07, 0x24, 0x69, 0x63, 0x7b, 0xf0, 0xfa, 0xbc, 0x4f,
	0xbe, 0xbc, 0x2f, 0x58, 0x3c, 0x9b, 0xed, 0xf6, 0x29, 0xfb, 0xf4, 0xb9, 0xc8, 0xcb, 0x1c, 0x4d,
	0x9e, 0x8d, 0x6d, 0x9e, 0xcd, 0x8a, 0x52, 0xd0, 0xb4, 0x86, 0xe3, 0x0b, 0x29, 0x51, 0x51, 0xb2,
	0x37, 0xb6, 0x4b, 0x4b, 0x5a, 0x53, 0x94, 0x1a, 0x7b, 0x3f, 0xa4, 0xe5, 0x51, 0xd4, 0xcc, 0xfb,
	0x31, 0xa1, 0x1b, 0xca, 0x73, 0x78, 0x0d, 0xfd, 0x2f, 0x2a, 0x0a, 0x96, 0x1f, 0x5c, 0x63, 0x62,
	0x4e, 0xad, 0xc5, 0xb9, 0xcf, 0x33, 0x5f, 0x65, 0xfe, 0xb6, 0x0a, 0x88, 0x36, 0x30, 0x80, 0x81,
	0x2a, 0xb1, 0xf9, 0xe6, 0xd4, 0x35, 0x95, 0x3e, 0x6a, 0xf4, 0x50, 0x47, 0xa4, 0xb1, 0xd0, 0x83,
	0x5e, 0xd5, 0xd1, 0x3d, 0x99, 0x18, 0xd3, 0xe1, 0x02, 0xa4, 0x1f, 0x2b, 0x42, 0xea, 0x04, 0x03,
	0x18, 0xb6, 0x6a, 0xbb, 0xa7, 0x4a, 0xb4, 0xd5, 0xe1, 0x06, 0x93, 0xb6, 0x83, 0x37, 0x80, 0xfc,
	0x98, 0xed, 0x59, 0xf1, 0x41, 0x45, 0xac, 0xc7, 0xb9, 0x5d, 0xf5, 0xf2, 0x4c, 0x7d, 0xa1, 0x21,
	0xf9, 0x47, 0xf4, 0xae, 0xa0, 0x5f, 0x8f, 0xc3, 0x11, 0xd8, 0xcf, 0xd1, 0x7d, 0xb4, 0x7e, 0x89,
	0x92, 0xed, 0x92, 0xc4, 0xab, 0x75, 0xe4, 0x74, 0x10, 0xa0, 0x97, 0xcc, 0x93, 0x79, 0x12, 0x38,
	0x86, 0x77, 0x07, 0x83, 0xbf, 0x65, 0x78, 0x09, 0xa8, 0xed, 0xf0, 0xe1, 0x76, 0xf5, 0x98, 0x6c,
	0x5e, 0x9f, 0x96, 0x4e, 0x07, 0x2d, 0x80, 0x6a, 0x8c, 0xb4, 0x1c, 0x43, 0x5e, 0x6d, 0xd5, 0x55,
	0xd0, 0xfc, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x9f, 0xc2, 0x3f, 0x9f, 0xc5, 0x01, 0x00, 0x00,
}
