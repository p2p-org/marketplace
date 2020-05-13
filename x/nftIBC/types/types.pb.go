// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.7.1
// source: x/nftIBC/types/types.proto

package types

import (
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type MsgTransferNFT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the port on which the packet will be sent
	SourcePort string `protobuf:"bytes,1,opt,name=source_port,json=sourcePort,proto3" json:"source_port,omitempty"`
	// the channel by which the packet will be sent
	SourceChannel string `protobuf:"bytes,2,opt,name=source_channel,json=sourceChannel,proto3" json:"source_channel,omitempty"`
	// the current height of the destination chain
	DestinationHeight uint64 `protobuf:"varint,3,opt,name=destination_height,json=destinationHeight,proto3" json:"destination_height,omitempty"`
	Id                string `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
	Denom             string `protobuf:"bytes,5,opt,name=denom,proto3" json:"denom,omitempty"`
	// the sender address
	Sender []byte `protobuf:"bytes,6,opt,name=sender,proto3" json:"sender,omitempty"`
	// the recipient address on the destination chain
	Receiver string `protobuf:"bytes,7,opt,name=receiver,proto3" json:"receiver,omitempty"`
}

func (x *MsgTransferNFT) Reset() {
	*x = MsgTransferNFT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_x_nftIBC_types_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgTransferNFT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgTransferNFT) ProtoMessage() {}

func (x *MsgTransferNFT) ProtoReflect() protoreflect.Message {
	mi := &file_x_nftIBC_types_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgTransferNFT.ProtoReflect.Descriptor instead.
func (*MsgTransferNFT) Descriptor() ([]byte, []int) {
	return file_x_nftIBC_types_types_proto_rawDescGZIP(), []int{0}
}

func (x *MsgTransferNFT) GetSourcePort() string {
	if x != nil {
		return x.SourcePort
	}
	return ""
}

func (x *MsgTransferNFT) GetSourceChannel() string {
	if x != nil {
		return x.SourceChannel
	}
	return ""
}

func (x *MsgTransferNFT) GetDestinationHeight() uint64 {
	if x != nil {
		return x.DestinationHeight
	}
	return 0
}

func (x *MsgTransferNFT) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MsgTransferNFT) GetDenom() string {
	if x != nil {
		return x.Denom
	}
	return ""
}

func (x *MsgTransferNFT) GetSender() []byte {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *MsgTransferNFT) GetReceiver() string {
	if x != nil {
		return x.Receiver
	}
	return ""
}

type NFTPacketData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Denom    string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	Owner    []byte `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Receiver []byte `protobuf:"bytes,4,opt,name=receiver,proto3" json:"receiver,omitempty"`
	TokenURI string `protobuf:"bytes,5,opt,name=tokenURI,proto3" json:"tokenURI,omitempty"`
}

func (x *NFTPacketData) Reset() {
	*x = NFTPacketData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_x_nftIBC_types_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NFTPacketData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NFTPacketData) ProtoMessage() {}

func (x *NFTPacketData) ProtoReflect() protoreflect.Message {
	mi := &file_x_nftIBC_types_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NFTPacketData.ProtoReflect.Descriptor instead.
func (*NFTPacketData) Descriptor() ([]byte, []int) {
	return file_x_nftIBC_types_types_proto_rawDescGZIP(), []int{1}
}

func (x *NFTPacketData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *NFTPacketData) GetDenom() string {
	if x != nil {
		return x.Denom
	}
	return ""
}

func (x *NFTPacketData) GetOwner() []byte {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *NFTPacketData) GetReceiver() []byte {
	if x != nil {
		return x.Receiver
	}
	return nil
}

func (x *NFTPacketData) GetTokenURI() string {
	if x != nil {
		return x.TokenURI
	}
	return ""
}

type NFTPacketAcknowledgement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *NFTPacketAcknowledgement) Reset() {
	*x = NFTPacketAcknowledgement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_x_nftIBC_types_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NFTPacketAcknowledgement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NFTPacketAcknowledgement) ProtoMessage() {}

func (x *NFTPacketAcknowledgement) ProtoReflect() protoreflect.Message {
	mi := &file_x_nftIBC_types_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NFTPacketAcknowledgement.ProtoReflect.Descriptor instead.
func (*NFTPacketAcknowledgement) Descriptor() ([]byte, []int) {
	return file_x_nftIBC_types_types_proto_rawDescGZIP(), []int{2}
}

func (x *NFTPacketAcknowledgement) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *NFTPacketAcknowledgement) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_x_nftIBC_types_types_proto protoreflect.FileDescriptor

var file_x_nftIBC_types_types_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x78, 0x2f, 0x6e, 0x66, 0x74, 0x49, 0x42, 0x43, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x2e, 0x78, 0x2e, 0x6e, 0x66, 0x74, 0x49,
	0x42, 0x43, 0x2e, 0x76, 0x31, 0x1a, 0x26, 0x74, 0x68, 0x69, 0x72, 0x64, 0x5f, 0x70, 0x61, 0x72,
	0x74, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe6, 0x02,
	0x0a, 0x0e, 0x4d, 0x73, 0x67, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x4e, 0x46, 0x54,
	0x12, 0x37, 0x0a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x16, 0xf2, 0xde, 0x1f, 0x12, 0x79, 0x61, 0x6d, 0x6c, 0x3a,
	0x22, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x52, 0x0a, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x40, 0x0a, 0x0e, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x19, 0xf2, 0xde, 0x1f, 0x15, 0x79, 0x61, 0x6d, 0x6c, 0x3a, 0x22, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x22, 0x52, 0x0d, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x4c, 0x0a, 0x12, 0x64,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x1d, 0xf2, 0xde, 0x1f, 0x19, 0x79, 0x61, 0x6d,
	0x6c, 0x3a, 0x22, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x52, 0x11, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6e,
	0x6f, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x12,
	0x49, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x42,
	0x31, 0xfa, 0xde, 0x1f, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2d, 0x73, 0x64,
	0x6b, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x41, 0x63, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x22, 0xe9, 0x01, 0x0a, 0x0d, 0x4e, 0x46, 0x54, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6e, 0x6f,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x12, 0x47,
	0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x31, 0xfa,
	0xde, 0x1f, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f,
	0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2d, 0x73, 0x64, 0x6b, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x41, 0x63, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x4d, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x31, 0xfa, 0xde, 0x1f, 0x2d, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73,
	0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x41, 0x63, 0x63, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x08, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x55,
	0x52, 0x49, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x55,
	0x52, 0x49, 0x22, 0x4a, 0x0a, 0x18, 0x4e, 0x46, 0x54, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x41,
	0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x32,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x73, 0x74, 0x61, 0x72, 0x69, 0x6f, 0x2f, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x2f, 0x78, 0x2f, 0x6e, 0x66, 0x74, 0x49, 0x42, 0x43, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_x_nftIBC_types_types_proto_rawDescOnce sync.Once
	file_x_nftIBC_types_types_proto_rawDescData = file_x_nftIBC_types_types_proto_rawDesc
)

func file_x_nftIBC_types_types_proto_rawDescGZIP() []byte {
	file_x_nftIBC_types_types_proto_rawDescOnce.Do(func() {
		file_x_nftIBC_types_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_x_nftIBC_types_types_proto_rawDescData)
	})
	return file_x_nftIBC_types_types_proto_rawDescData
}

var file_x_nftIBC_types_types_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_x_nftIBC_types_types_proto_goTypes = []interface{}{
	(*MsgTransferNFT)(nil),           // 0: marketplace.x.nftIBC.v1.MsgTransferNFT
	(*NFTPacketData)(nil),            // 1: marketplace.x.nftIBC.v1.NFTPacketData
	(*NFTPacketAcknowledgement)(nil), // 2: marketplace.x.nftIBC.v1.NFTPacketAcknowledgement
}
var file_x_nftIBC_types_types_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_x_nftIBC_types_types_proto_init() }
func file_x_nftIBC_types_types_proto_init() {
	if File_x_nftIBC_types_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_x_nftIBC_types_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgTransferNFT); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_x_nftIBC_types_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NFTPacketData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_x_nftIBC_types_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NFTPacketAcknowledgement); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_x_nftIBC_types_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_x_nftIBC_types_types_proto_goTypes,
		DependencyIndexes: file_x_nftIBC_types_types_proto_depIdxs,
		MessageInfos:      file_x_nftIBC_types_types_proto_msgTypes,
	}.Build()
	File_x_nftIBC_types_types_proto = out.File
	file_x_nftIBC_types_types_proto_rawDesc = nil
	file_x_nftIBC_types_types_proto_goTypes = nil
	file_x_nftIBC_types_types_proto_depIdxs = nil
}
