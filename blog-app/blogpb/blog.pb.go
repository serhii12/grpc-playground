// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blog-app/blogpb/blog.proto

package blogpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Blog struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AuthorId             string   `protobuf:"bytes,2,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	Title                string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Content              string   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Blog) Reset()         { *m = Blog{} }
func (m *Blog) String() string { return proto.CompactTextString(m) }
func (*Blog) ProtoMessage()    {}
func (*Blog) Descriptor() ([]byte, []int) {
	return fileDescriptor_77490c284db47c9b, []int{0}
}

func (m *Blog) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Blog.Unmarshal(m, b)
}
func (m *Blog) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Blog.Marshal(b, m, deterministic)
}
func (m *Blog) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Blog.Merge(m, src)
}
func (m *Blog) XXX_Size() int {
	return xxx_messageInfo_Blog.Size(m)
}
func (m *Blog) XXX_DiscardUnknown() {
	xxx_messageInfo_Blog.DiscardUnknown(m)
}

var xxx_messageInfo_Blog proto.InternalMessageInfo

func (m *Blog) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Blog) GetAuthorId() string {
	if m != nil {
		return m.AuthorId
	}
	return ""
}

func (m *Blog) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Blog) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*Blog)(nil), "blog.Blog")
}

func init() { proto.RegisterFile("blog-app/blogpb/blog.proto", fileDescriptor_77490c284db47c9b) }

var fileDescriptor_77490c284db47c9b = []byte{
	// 152 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4a, 0xca, 0xc9, 0x4f,
	0xd7, 0x4d, 0x2c, 0x28, 0xd0, 0x07, 0x31, 0x0a, 0x92, 0xc0, 0x94, 0x5e, 0x41, 0x51, 0x7e, 0x49,
	0xbe, 0x10, 0x0b, 0x88, 0xad, 0x94, 0xcc, 0xc5, 0xe2, 0x94, 0x93, 0x9f, 0x2e, 0xc4, 0xc7, 0xc5,
	0x94, 0x99, 0x22, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0xc4, 0x94, 0x99, 0x22, 0x24, 0xcd, 0xc5,
	0x99, 0x58, 0x5a, 0x92, 0x91, 0x5f, 0x14, 0x9f, 0x99, 0x22, 0xc1, 0x04, 0x16, 0xe6, 0x80, 0x08,
	0x78, 0xa6, 0x08, 0x89, 0x70, 0xb1, 0x96, 0x64, 0x96, 0xe4, 0xa4, 0x4a, 0x30, 0x83, 0x25, 0x20,
	0x1c, 0x21, 0x09, 0x2e, 0xf6, 0xe4, 0xfc, 0xbc, 0x92, 0xd4, 0xbc, 0x12, 0x09, 0x16, 0xb0, 0x38,
	0x8c, 0x6b, 0xc4, 0xcb, 0xc5, 0x0d, 0xb2, 0x24, 0x38, 0xb5, 0xa8, 0x2c, 0x33, 0x39, 0xd5, 0x89,
	0x23, 0x8a, 0x0d, 0xe2, 0x9c, 0x24, 0x36, 0xb0, 0x53, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xbc, 0xe9, 0x32, 0x47, 0xa8, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BlogServiceClient is the client API for BlogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BlogServiceClient interface {
}

type blogServiceClient struct {
	cc *grpc.ClientConn
}

func NewBlogServiceClient(cc *grpc.ClientConn) BlogServiceClient {
	return &blogServiceClient{cc}
}

// BlogServiceServer is the server API for BlogService service.
type BlogServiceServer interface {
}

// UnimplementedBlogServiceServer can be embedded to have forward compatible implementations.
type UnimplementedBlogServiceServer struct {
}

func RegisterBlogServiceServer(s *grpc.Server, srv BlogServiceServer) {
	s.RegisterService(&_BlogService_serviceDesc, srv)
}

var _BlogService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blog.BlogService",
	HandlerType: (*BlogServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "blog-app/blogpb/blog.proto",
}
