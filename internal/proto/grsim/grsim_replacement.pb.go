// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grsim_replacement.proto

package grsim

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type GrSim_RobotReplacement struct {
	X                    *float64 `protobuf:"fixed64,1,req,name=x" json:"x,omitempty"`
	Y                    *float64 `protobuf:"fixed64,2,req,name=y" json:"y,omitempty"`
	Dir                  *float64 `protobuf:"fixed64,3,req,name=dir" json:"dir,omitempty"`
	Id                   *uint32  `protobuf:"varint,4,req,name=id" json:"id,omitempty"`
	Yellowteam           *bool    `protobuf:"varint,5,req,name=yellowteam" json:"yellowteam,omitempty"`
	Turnon               *bool    `protobuf:"varint,6,opt,name=turnon" json:"turnon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrSim_RobotReplacement) Reset()         { *m = GrSim_RobotReplacement{} }
func (m *GrSim_RobotReplacement) String() string { return proto.CompactTextString(m) }
func (*GrSim_RobotReplacement) ProtoMessage()    {}
func (*GrSim_RobotReplacement) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9447b4cd61a47d6, []int{0}
}

func (m *GrSim_RobotReplacement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GrSim_RobotReplacement.Unmarshal(m, b)
}
func (m *GrSim_RobotReplacement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GrSim_RobotReplacement.Marshal(b, m, deterministic)
}
func (m *GrSim_RobotReplacement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrSim_RobotReplacement.Merge(m, src)
}
func (m *GrSim_RobotReplacement) XXX_Size() int {
	return xxx_messageInfo_GrSim_RobotReplacement.Size(m)
}
func (m *GrSim_RobotReplacement) XXX_DiscardUnknown() {
	xxx_messageInfo_GrSim_RobotReplacement.DiscardUnknown(m)
}

var xxx_messageInfo_GrSim_RobotReplacement proto.InternalMessageInfo

func (m *GrSim_RobotReplacement) GetX() float64 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *GrSim_RobotReplacement) GetY() float64 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *GrSim_RobotReplacement) GetDir() float64 {
	if m != nil && m.Dir != nil {
		return *m.Dir
	}
	return 0
}

func (m *GrSim_RobotReplacement) GetId() uint32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *GrSim_RobotReplacement) GetYellowteam() bool {
	if m != nil && m.Yellowteam != nil {
		return *m.Yellowteam
	}
	return false
}

func (m *GrSim_RobotReplacement) GetTurnon() bool {
	if m != nil && m.Turnon != nil {
		return *m.Turnon
	}
	return false
}

type GrSim_BallReplacement struct {
	X                    *float64 `protobuf:"fixed64,1,opt,name=x" json:"x,omitempty"`
	Y                    *float64 `protobuf:"fixed64,2,opt,name=y" json:"y,omitempty"`
	Vx                   *float64 `protobuf:"fixed64,3,opt,name=vx" json:"vx,omitempty"`
	Vy                   *float64 `protobuf:"fixed64,4,opt,name=vy" json:"vy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrSim_BallReplacement) Reset()         { *m = GrSim_BallReplacement{} }
func (m *GrSim_BallReplacement) String() string { return proto.CompactTextString(m) }
func (*GrSim_BallReplacement) ProtoMessage()    {}
func (*GrSim_BallReplacement) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9447b4cd61a47d6, []int{1}
}

func (m *GrSim_BallReplacement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GrSim_BallReplacement.Unmarshal(m, b)
}
func (m *GrSim_BallReplacement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GrSim_BallReplacement.Marshal(b, m, deterministic)
}
func (m *GrSim_BallReplacement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrSim_BallReplacement.Merge(m, src)
}
func (m *GrSim_BallReplacement) XXX_Size() int {
	return xxx_messageInfo_GrSim_BallReplacement.Size(m)
}
func (m *GrSim_BallReplacement) XXX_DiscardUnknown() {
	xxx_messageInfo_GrSim_BallReplacement.DiscardUnknown(m)
}

var xxx_messageInfo_GrSim_BallReplacement proto.InternalMessageInfo

func (m *GrSim_BallReplacement) GetX() float64 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *GrSim_BallReplacement) GetY() float64 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

func (m *GrSim_BallReplacement) GetVx() float64 {
	if m != nil && m.Vx != nil {
		return *m.Vx
	}
	return 0
}

func (m *GrSim_BallReplacement) GetVy() float64 {
	if m != nil && m.Vy != nil {
		return *m.Vy
	}
	return 0
}

type GrSim_Replacement struct {
	Ball                 *GrSim_BallReplacement    `protobuf:"bytes,1,opt,name=ball" json:"ball,omitempty"`
	Robots               []*GrSim_RobotReplacement `protobuf:"bytes,2,rep,name=robots" json:"robots,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *GrSim_Replacement) Reset()         { *m = GrSim_Replacement{} }
func (m *GrSim_Replacement) String() string { return proto.CompactTextString(m) }
func (*GrSim_Replacement) ProtoMessage()    {}
func (*GrSim_Replacement) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9447b4cd61a47d6, []int{2}
}

func (m *GrSim_Replacement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GrSim_Replacement.Unmarshal(m, b)
}
func (m *GrSim_Replacement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GrSim_Replacement.Marshal(b, m, deterministic)
}
func (m *GrSim_Replacement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrSim_Replacement.Merge(m, src)
}
func (m *GrSim_Replacement) XXX_Size() int {
	return xxx_messageInfo_GrSim_Replacement.Size(m)
}
func (m *GrSim_Replacement) XXX_DiscardUnknown() {
	xxx_messageInfo_GrSim_Replacement.DiscardUnknown(m)
}

var xxx_messageInfo_GrSim_Replacement proto.InternalMessageInfo

func (m *GrSim_Replacement) GetBall() *GrSim_BallReplacement {
	if m != nil {
		return m.Ball
	}
	return nil
}

func (m *GrSim_Replacement) GetRobots() []*GrSim_RobotReplacement {
	if m != nil {
		return m.Robots
	}
	return nil
}

func init() {
	proto.RegisterType((*GrSim_RobotReplacement)(nil), "grSim_RobotReplacement")
	proto.RegisterType((*GrSim_BallReplacement)(nil), "grSim_BallReplacement")
	proto.RegisterType((*GrSim_Replacement)(nil), "grSim_Replacement")
}

func init() { proto.RegisterFile("grsim_replacement.proto", fileDescriptor_e9447b4cd61a47d6) }

var fileDescriptor_e9447b4cd61a47d6 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xd1, 0x6a, 0xc3, 0x20,
	0x18, 0x85, 0xd1, 0xa4, 0x59, 0xf9, 0xbb, 0x95, 0x4d, 0x58, 0xfa, 0x5f, 0x0d, 0xc9, 0x95, 0xec,
	0x22, 0x83, 0x3e, 0x42, 0x1f, 0xc1, 0xde, 0xed, 0xa6, 0xa4, 0x8b, 0x14, 0xc1, 0xc4, 0x60, 0x5d,
	0x17, 0x5f, 0x62, 0xcf, 0x3c, 0x62, 0x02, 0x4b, 0x21, 0x77, 0x7e, 0x72, 0x38, 0x7e, 0x1e, 0xd8,
	0x5d, 0xdc, 0x55, 0x37, 0x27, 0xa7, 0x3a, 0x53, 0x7d, 0xa9, 0x46, 0xb5, 0xbe, 0xec, 0x9c, 0xf5,
	0xb6, 0xf8, 0x25, 0x90, 0x5f, 0xdc, 0x51, 0x37, 0x27, 0x69, 0xcf, 0xd6, 0xcb, 0xff, 0x00, 0x7b,
	0x04, 0xd2, 0x23, 0xe1, 0x54, 0x10, 0x49, 0xfa, 0x81, 0x02, 0xd2, 0x91, 0x02, 0x7b, 0x86, 0xa4,
	0xd6, 0x0e, 0x93, 0xc8, 0xc3, 0x91, 0x6d, 0x81, 0xea, 0x1a, 0x53, 0x4e, 0xc5, 0x93, 0xa4, 0xba,
	0x66, 0x6f, 0x00, 0x41, 0x19, 0x63, 0x7f, 0xbc, 0xaa, 0x1a, 0x5c, 0x71, 0x2a, 0xd6, 0x72, 0x76,
	0xc3, 0x72, 0xc8, 0xfc, 0xb7, 0x6b, 0x6d, 0x8b, 0x19, 0x27, 0x62, 0x2d, 0x27, 0x2a, 0x8e, 0xf0,
	0x3a, 0xfa, 0x1c, 0x2a, 0x63, 0x16, 0x74, 0xc8, 0x9d, 0x0e, 0x19, 0x75, 0xb6, 0x40, 0x6f, 0x3d,
	0x26, 0x11, 0xe9, 0xad, 0x8f, 0x1c, 0x30, 0x9d, 0x38, 0x14, 0x1d, 0xbc, 0x4c, 0x9f, 0x9c, 0x15,
	0xbe, 0x43, 0x7a, 0xae, 0x8c, 0x89, 0x9d, 0x9b, 0x7d, 0x5e, 0x2e, 0x3e, 0x2b, 0x63, 0x86, 0x7d,
	0x40, 0xe6, 0x86, 0x7d, 0xae, 0x48, 0x79, 0x22, 0x36, 0xfb, 0x5d, 0xb9, 0x3c, 0x9a, 0x9c, 0x62,
	0x87, 0x87, 0xcf, 0x55, 0x9c, 0xfc, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xe2, 0xe6, 0xca, 0xa0, 0x7a,
	0x01, 0x00, 0x00,
}