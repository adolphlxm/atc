// Code generated by protoc-gen-go.
// source: pbtest.proto
// DO NOT EDIT!

/*
Package testdata is a generated protocol buffer package.

It is generated from these files:
	pbtest.proto

It has these top-level messages:
	Something
*/
package testdata

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Something struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Age     int32  `protobuf:"varint,2,opt,name=age" json:"age,omitempty"`
	Address string `protobuf:"bytes,3,opt,name=address" json:"address,omitempty"`
}

func (m *Something) Reset()                    { *m = Something{} }
func (m *Something) String() string            { return proto.CompactTextString(m) }
func (*Something) ProtoMessage()               {}
func (*Something) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Something) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Something) GetAge() int32 {
	if m != nil {
		return m.Age
	}
	return 0
}

func (m *Something) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*Something)(nil), "testdata.Something")
}

func init() { proto.RegisterFile("pbtest.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x48, 0x2a, 0x49,
	0x2d, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0xb1, 0x53, 0x12, 0x4b, 0x12,
	0x95, 0xbc, 0xb9, 0x38, 0x83, 0xf3, 0x73, 0x53, 0x4b, 0x32, 0x32, 0xf3, 0xd2, 0x85, 0x84, 0xb8,
	0x58, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x21, 0x01,
	0x2e, 0xe6, 0xc4, 0xf4, 0x54, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x10, 0x53, 0x48, 0x82,
	0x8b, 0x3d, 0x31, 0x25, 0xa5, 0x28, 0xb5, 0xb8, 0x58, 0x82, 0x19, 0xac, 0x10, 0xc6, 0x4d, 0x62,
	0x03, 0x9b, 0x6e, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x62, 0x12, 0x7a, 0xe6, 0x6d, 0x00, 0x00,
	0x00,
}
