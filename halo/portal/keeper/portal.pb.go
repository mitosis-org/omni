// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: halo/portal/keeper/portal.proto

package keeper

import (
	_ "cosmossdk.io/api/cosmos/orm/v1"
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

// Block groups a set of Msgs adding height and offset.
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                            // Auto-incremented ID (BlockHeight, BlockOffset)
	CreatedHeight uint64 `protobuf:"varint,2,opt,name=created_height,json=createdHeight,proto3" json:"created_height,omitempty"` // Height this block was created at.
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_halo_portal_keeper_portal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_halo_portal_keeper_portal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_halo_portal_keeper_portal_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Block) GetCreatedHeight() uint64 {
	if x != nil {
		return x.CreatedHeight
	}
	return 0
}

// Msg represents a single cross-chain message emitted by the consensus chain portal.
type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                  // Auto-incremented ID (MsgStreamOffset)
	BlockId   uint64 `protobuf:"varint,2,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`         // Block ID to which this msg belongs
	MsgType   uint32 `protobuf:"varint,3,opt,name=msg_type,json=msgType,proto3" json:"msg_type,omitempty"`         // Message type (ValidatorSet, Withdrawal, etc.)
	MsgTypeId uint64 `protobuf:"varint,4,opt,name=msg_type_id,json=msgTypeId,proto3" json:"msg_type_id,omitempty"` // ID of the type referred to be MsgType
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_halo_portal_keeper_portal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_halo_portal_keeper_portal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_halo_portal_keeper_portal_proto_rawDescGZIP(), []int{1}
}

func (x *Msg) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Msg) GetBlockId() uint64 {
	if x != nil {
		return x.BlockId
	}
	return 0
}

func (x *Msg) GetMsgType() uint32 {
	if x != nil {
		return x.MsgType
	}
	return 0
}

func (x *Msg) GetMsgTypeId() uint64 {
	if x != nil {
		return x.MsgTypeId
	}
	return 0
}

var File_halo_portal_keeper_portal_proto protoreflect.FileDescriptor

var file_halo_portal_keeper_portal_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x68, 0x61, 0x6c, 0x6f, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x6b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x68, 0x61, 0x6c, 0x6f, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x1a, 0x17, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x6f, 0x72,
	0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66,
	0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x3a, 0x26,
	0xf2, 0x9e, 0xd3, 0x8e, 0x03, 0x20, 0x0a, 0x06, 0x0a, 0x02, 0x69, 0x64, 0x10, 0x01, 0x12, 0x14,
	0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x10, 0x02, 0x18, 0x01, 0x18, 0x01, 0x22, 0x8b, 0x01, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19,
	0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x73, 0x67,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6d, 0x73, 0x67,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0b, 0x6d, 0x73, 0x67, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x6d, 0x73, 0x67, 0x54, 0x79,
	0x70, 0x65, 0x49, 0x64, 0x3a, 0x1e, 0xf2, 0x9e, 0xd3, 0x8e, 0x03, 0x18, 0x0a, 0x06, 0x0a, 0x02,
	0x69, 0x64, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64,
	0x10, 0x02, 0x18, 0x02, 0x42, 0xc0, 0x01, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x61, 0x6c,
	0x6f, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x42,
	0x0b, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2d,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2f, 0x68, 0x61, 0x6c,
	0x6f, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0xa2,
	0x02, 0x03, 0x48, 0x50, 0x4b, 0xaa, 0x02, 0x12, 0x48, 0x61, 0x6c, 0x6f, 0x2e, 0x50, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0xca, 0x02, 0x12, 0x48, 0x61, 0x6c,
	0x6f, 0x5c, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0xe2,
	0x02, 0x1e, 0x48, 0x61, 0x6c, 0x6f, 0x5c, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x4b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x14, 0x48, 0x61, 0x6c, 0x6f, 0x3a, 0x3a, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x3a,
	0x3a, 0x4b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_halo_portal_keeper_portal_proto_rawDescOnce sync.Once
	file_halo_portal_keeper_portal_proto_rawDescData = file_halo_portal_keeper_portal_proto_rawDesc
)

func file_halo_portal_keeper_portal_proto_rawDescGZIP() []byte {
	file_halo_portal_keeper_portal_proto_rawDescOnce.Do(func() {
		file_halo_portal_keeper_portal_proto_rawDescData = protoimpl.X.CompressGZIP(file_halo_portal_keeper_portal_proto_rawDescData)
	})
	return file_halo_portal_keeper_portal_proto_rawDescData
}

var file_halo_portal_keeper_portal_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_halo_portal_keeper_portal_proto_goTypes = []interface{}{
	(*Block)(nil), // 0: halo.portal.keeper.Block
	(*Msg)(nil),   // 1: halo.portal.keeper.Msg
}
var file_halo_portal_keeper_portal_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_halo_portal_keeper_portal_proto_init() }
func file_halo_portal_keeper_portal_proto_init() {
	if File_halo_portal_keeper_portal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_halo_portal_keeper_portal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_halo_portal_keeper_portal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
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
			RawDescriptor: file_halo_portal_keeper_portal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_halo_portal_keeper_portal_proto_goTypes,
		DependencyIndexes: file_halo_portal_keeper_portal_proto_depIdxs,
		MessageInfos:      file_halo_portal_keeper_portal_proto_msgTypes,
	}.Build()
	File_halo_portal_keeper_portal_proto = out.File
	file_halo_portal_keeper_portal_proto_rawDesc = nil
	file_halo_portal_keeper_portal_proto_goTypes = nil
	file_halo_portal_keeper_portal_proto_depIdxs = nil
}