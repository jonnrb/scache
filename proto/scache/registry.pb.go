// Code generated by protoc-gen-go. DO NOT EDIT.
// source: scache/registry.proto

package scache

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/any"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ProviderAddress struct {
	Uri string `protobuf:"bytes,1,opt,name=uri" json:"uri,omitempty"`
	// Defaults to gRPC.
	Proto string `protobuf:"bytes,2,opt,name=proto" json:"proto,omitempty"`
	// Options to apply when first adding this provider. Valid until this
	// provider (as indicated by `uri` and `proto`) handles no more source
	// types.
	Options *google_protobuf.Any `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
}

func (m *ProviderAddress) Reset()                    { *m = ProviderAddress{} }
func (m *ProviderAddress) String() string            { return proto.CompactTextString(m) }
func (*ProviderAddress) ProtoMessage()               {}
func (*ProviderAddress) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *ProviderAddress) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

func (m *ProviderAddress) GetProto() string {
	if m != nil {
		return m.Proto
	}
	return ""
}

func (m *ProviderAddress) GetOptions() *google_protobuf.Any {
	if m != nil {
		return m.Options
	}
	return nil
}

type ProviderSpec struct {
	Addr       *ProviderAddress `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	SourceType []string         `protobuf:"bytes,2,rep,name=source_type,json=sourceType" json:"source_type,omitempty"`
}

func (m *ProviderSpec) Reset()                    { *m = ProviderSpec{} }
func (m *ProviderSpec) String() string            { return proto.CompactTextString(m) }
func (*ProviderSpec) ProtoMessage()               {}
func (*ProviderSpec) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ProviderSpec) GetAddr() *ProviderAddress {
	if m != nil {
		return m.Addr
	}
	return nil
}

func (m *ProviderSpec) GetSourceType() []string {
	if m != nil {
		return m.SourceType
	}
	return nil
}

type AddProviderResponse struct {
}

func (m *AddProviderResponse) Reset()                    { *m = AddProviderResponse{} }
func (m *AddProviderResponse) String() string            { return proto.CompactTextString(m) }
func (*AddProviderResponse) ProtoMessage()               {}
func (*AddProviderResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

type RemoveProviderResponse struct {
}

func (m *RemoveProviderResponse) Reset()                    { *m = RemoveProviderResponse{} }
func (m *RemoveProviderResponse) String() string            { return proto.CompactTextString(m) }
func (*RemoveProviderResponse) ProtoMessage()               {}
func (*RemoveProviderResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

type ListProvidersRequest struct {
}

func (m *ListProvidersRequest) Reset()                    { *m = ListProvidersRequest{} }
func (m *ListProvidersRequest) String() string            { return proto.CompactTextString(m) }
func (*ListProvidersRequest) ProtoMessage()               {}
func (*ListProvidersRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

type ListProvidersResponse struct {
	Provider []*ProviderSpec `protobuf:"bytes,1,rep,name=provider" json:"provider,omitempty"`
}

func (m *ListProvidersResponse) Reset()                    { *m = ListProvidersResponse{} }
func (m *ListProvidersResponse) String() string            { return proto.CompactTextString(m) }
func (*ListProvidersResponse) ProtoMessage()               {}
func (*ListProvidersResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *ListProvidersResponse) GetProvider() []*ProviderSpec {
	if m != nil {
		return m.Provider
	}
	return nil
}

type GRPCProviderOpts struct {
	DialInsecure bool `protobuf:"varint,1,opt,name=dial_insecure,json=dialInsecure" json:"dial_insecure,omitempty"`
}

func (m *GRPCProviderOpts) Reset()                    { *m = GRPCProviderOpts{} }
func (m *GRPCProviderOpts) String() string            { return proto.CompactTextString(m) }
func (*GRPCProviderOpts) ProtoMessage()               {}
func (*GRPCProviderOpts) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

func (m *GRPCProviderOpts) GetDialInsecure() bool {
	if m != nil {
		return m.DialInsecure
	}
	return false
}

func init() {
	proto.RegisterType((*ProviderAddress)(nil), "scache.ProviderAddress")
	proto.RegisterType((*ProviderSpec)(nil), "scache.ProviderSpec")
	proto.RegisterType((*AddProviderResponse)(nil), "scache.AddProviderResponse")
	proto.RegisterType((*RemoveProviderResponse)(nil), "scache.RemoveProviderResponse")
	proto.RegisterType((*ListProvidersRequest)(nil), "scache.ListProvidersRequest")
	proto.RegisterType((*ListProvidersResponse)(nil), "scache.ListProvidersResponse")
	proto.RegisterType((*GRPCProviderOpts)(nil), "scache.gRPCProviderOpts")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ProviderRegistry service

type ProviderRegistryClient interface {
	AddProvider(ctx context.Context, in *ProviderSpec, opts ...grpc.CallOption) (*AddProviderResponse, error)
	RemoveProvider(ctx context.Context, in *ProviderSpec, opts ...grpc.CallOption) (*RemoveProviderResponse, error)
	ListProviders(ctx context.Context, in *ListProvidersRequest, opts ...grpc.CallOption) (*ListProvidersResponse, error)
}

type providerRegistryClient struct {
	cc *grpc.ClientConn
}

func NewProviderRegistryClient(cc *grpc.ClientConn) ProviderRegistryClient {
	return &providerRegistryClient{cc}
}

func (c *providerRegistryClient) AddProvider(ctx context.Context, in *ProviderSpec, opts ...grpc.CallOption) (*AddProviderResponse, error) {
	out := new(AddProviderResponse)
	err := grpc.Invoke(ctx, "/scache.ProviderRegistry/AddProvider", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerRegistryClient) RemoveProvider(ctx context.Context, in *ProviderSpec, opts ...grpc.CallOption) (*RemoveProviderResponse, error) {
	out := new(RemoveProviderResponse)
	err := grpc.Invoke(ctx, "/scache.ProviderRegistry/RemoveProvider", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerRegistryClient) ListProviders(ctx context.Context, in *ListProvidersRequest, opts ...grpc.CallOption) (*ListProvidersResponse, error) {
	out := new(ListProvidersResponse)
	err := grpc.Invoke(ctx, "/scache.ProviderRegistry/ListProviders", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ProviderRegistry service

type ProviderRegistryServer interface {
	AddProvider(context.Context, *ProviderSpec) (*AddProviderResponse, error)
	RemoveProvider(context.Context, *ProviderSpec) (*RemoveProviderResponse, error)
	ListProviders(context.Context, *ListProvidersRequest) (*ListProvidersResponse, error)
}

func RegisterProviderRegistryServer(s *grpc.Server, srv ProviderRegistryServer) {
	s.RegisterService(&_ProviderRegistry_serviceDesc, srv)
}

func _ProviderRegistry_AddProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProviderSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegistryServer).AddProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scache.ProviderRegistry/AddProvider",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegistryServer).AddProvider(ctx, req.(*ProviderSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderRegistry_RemoveProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProviderSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegistryServer).RemoveProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scache.ProviderRegistry/RemoveProvider",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegistryServer).RemoveProvider(ctx, req.(*ProviderSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderRegistry_ListProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProvidersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegistryServer).ListProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scache.ProviderRegistry/ListProviders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegistryServer).ListProviders(ctx, req.(*ListProvidersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProviderRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "scache.ProviderRegistry",
	HandlerType: (*ProviderRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddProvider",
			Handler:    _ProviderRegistry_AddProvider_Handler,
		},
		{
			MethodName: "RemoveProvider",
			Handler:    _ProviderRegistry_RemoveProvider_Handler,
		},
		{
			MethodName: "ListProviders",
			Handler:    _ProviderRegistry_ListProviders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scache/registry.proto",
}

func init() { proto.RegisterFile("scache/registry.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x4f, 0x6f, 0xda, 0x40,
	0x10, 0xc5, 0x31, 0x6e, 0x29, 0x8c, 0xa1, 0x45, 0x5b, 0xa0, 0x2e, 0xfd, 0x67, 0xb9, 0x17, 0xa4,
	0x4a, 0xa6, 0xa2, 0x87, 0x9e, 0x49, 0x2e, 0x41, 0x8a, 0x12, 0xb4, 0xc9, 0x31, 0x12, 0x32, 0xde,
	0x89, 0xb3, 0x12, 0xf1, 0x6e, 0x76, 0x6d, 0x24, 0x7f, 0xeb, 0x7c, 0x84, 0x88, 0x35, 0x9b, 0x04,
	0x02, 0x37, 0xfb, 0x37, 0x6f, 0xe6, 0xcd, 0x3c, 0x1b, 0xfa, 0x3a, 0x89, 0x93, 0x3b, 0x1c, 0x2b,
	0x4c, 0xb9, 0xce, 0x55, 0x19, 0x49, 0x25, 0x72, 0x41, 0x1a, 0x15, 0x1e, 0x7e, 0x4d, 0x85, 0x48,
	0x57, 0x38, 0x36, 0x74, 0x59, 0xdc, 0x8e, 0xe3, 0x6c, 0x2b, 0x09, 0x39, 0x7c, 0x9a, 0x2b, 0xb1,
	0xe6, 0x0c, 0xd5, 0x94, 0x31, 0x85, 0x5a, 0x93, 0x2e, 0xb8, 0x85, 0xe2, 0xbe, 0x13, 0x38, 0xa3,
	0x16, 0xdd, 0x3c, 0x92, 0x1e, 0xbc, 0x37, 0x6a, 0xbf, 0x6e, 0x58, 0xf5, 0x42, 0x22, 0xf8, 0x20,
	0x64, 0xce, 0x45, 0xa6, 0x7d, 0x37, 0x70, 0x46, 0xde, 0xa4, 0x17, 0x55, 0x3e, 0x91, 0xf5, 0x89,
	0xa6, 0x59, 0x49, 0xad, 0x28, 0xbc, 0x81, 0xb6, 0xb5, 0xba, 0x92, 0x98, 0x90, 0x3f, 0xf0, 0x2e,
	0x66, 0x4c, 0x19, 0x23, 0x6f, 0xf2, 0x25, 0xaa, 0x96, 0x8d, 0xf6, 0xd6, 0xa1, 0x46, 0x44, 0x7e,
	0x81, 0xa7, 0x45, 0xa1, 0x12, 0x5c, 0xe4, 0xa5, 0x44, 0xbf, 0x1e, 0xb8, 0xa3, 0x16, 0x85, 0x0a,
	0x5d, 0x97, 0x12, 0xc3, 0x3e, 0x7c, 0x9e, 0x32, 0x66, 0x9b, 0x29, 0x6a, 0x29, 0x32, 0x8d, 0xa1,
	0x0f, 0x03, 0x8a, 0xf7, 0x62, 0x8d, 0x6f, 0x2a, 0x03, 0xe8, 0x9d, 0x73, 0x9d, 0x5b, 0xae, 0x29,
	0x3e, 0x14, 0xa8, 0xf3, 0x70, 0x06, 0xfd, 0x3d, 0x5e, 0x35, 0x90, 0xbf, 0xd0, 0x94, 0x5b, 0xe8,
	0x3b, 0x81, 0x6b, 0x0e, 0xde, 0xdb, 0x79, 0x73, 0x17, 0x7d, 0x56, 0x85, 0xff, 0xa1, 0x9b, 0xd2,
	0xf9, 0xa9, 0xad, 0x5e, 0xca, 0x5c, 0x93, 0xdf, 0xd0, 0x61, 0x3c, 0x5e, 0x2d, 0x78, 0xa6, 0x31,
	0x29, 0x14, 0x9a, 0xf3, 0x9b, 0xb4, 0xbd, 0x81, 0xb3, 0x2d, 0x9b, 0x3c, 0x3a, 0xd0, 0x7d, 0x59,
	0xb8, 0xfa, 0xa6, 0xe4, 0x04, 0xbc, 0x57, 0x17, 0x92, 0x83, 0xe6, 0xc3, 0x6f, 0x96, 0x1e, 0x0a,
	0xa3, 0x46, 0xce, 0xe0, 0xe3, 0x6e, 0x1c, 0x47, 0xc6, 0xfc, 0xb4, 0xf4, 0x48, 0x78, 0x35, 0x72,
	0x01, 0x9d, 0x9d, 0x98, 0xc8, 0x77, 0xdb, 0x72, 0x28, 0xd5, 0xe1, 0x8f, 0x23, 0x55, 0x3b, 0x6f,
	0xd9, 0x30, 0x3f, 0xcd, 0xbf, 0xa7, 0x00, 0x00, 0x00, 0xff, 0xff, 0x25, 0xcd, 0x46, 0xee, 0xcb,
	0x02, 0x00, 0x00,
}
