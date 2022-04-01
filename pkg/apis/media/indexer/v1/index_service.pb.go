// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: media/indexer/v1/index_service.proto

package indexerv1

import (
	reflect "reflect"
	sync "sync"

	common "github.com/wenerme/torrenti/pkg/apis/media/common"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/httpbody"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IndexRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *common.File `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	Url  string       `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *IndexRequest) Reset() {
	*x = IndexRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_indexer_v1_index_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexRequest) ProtoMessage() {}

func (x *IndexRequest) ProtoReflect() protoreflect.Message {
	mi := &file_media_indexer_v1_index_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexRequest.ProtoReflect.Descriptor instead.
func (*IndexRequest) Descriptor() ([]byte, []int) {
	return file_media_indexer_v1_index_service_proto_rawDescGZIP(), []int{0}
}

func (x *IndexRequest) GetFile() *common.File {
	if x != nil {
		return x.File
	}
	return nil
}

func (x *IndexRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type IndexResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *IndexResponse) Reset() {
	*x = IndexResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_indexer_v1_index_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexResponse) ProtoMessage() {}

func (x *IndexResponse) ProtoReflect() protoreflect.Message {
	mi := &file_media_indexer_v1_index_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexResponse.ProtoReflect.Descriptor instead.
func (*IndexResponse) Descriptor() ([]byte, []int) {
	return file_media_indexer_v1_index_service_proto_rawDescGZIP(), []int{1}
}

var File_media_indexer_v1_index_service_proto protoreflect.FileDescriptor

var file_media_indexer_v1_index_service_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2f,
	0x76, 0x31, 0x2f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x62, 0x6f, 0x64, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a, 0x0c, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x0f, 0x0a, 0x0d, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x6b, 0x0a, 0x0c, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x1e, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1f, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x22, 0x06, 0x2f, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x3a, 0x01, 0x2a, 0x42, 0xcc, 0x01, 0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65,
	0x6e, 0x65, 0x72, 0x6d, 0x65, 0x2f, 0x74, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x69, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72,
	0x76, 0x31, 0xa2, 0x02, 0x03, 0x4d, 0x49, 0x58, 0xaa, 0x02, 0x10, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x10, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x5c, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x1c, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x5c, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x12,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x3a, 0x3a, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_media_indexer_v1_index_service_proto_rawDescOnce sync.Once
	file_media_indexer_v1_index_service_proto_rawDescData = file_media_indexer_v1_index_service_proto_rawDesc
)

func file_media_indexer_v1_index_service_proto_rawDescGZIP() []byte {
	file_media_indexer_v1_index_service_proto_rawDescOnce.Do(func() {
		file_media_indexer_v1_index_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_media_indexer_v1_index_service_proto_rawDescData)
	})
	return file_media_indexer_v1_index_service_proto_rawDescData
}

var (
	file_media_indexer_v1_index_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
	file_media_indexer_v1_index_service_proto_goTypes  = []interface{}{
		(*IndexRequest)(nil),  // 0: media.indexer.v1.IndexRequest
		(*IndexResponse)(nil), // 1: media.indexer.v1.IndexResponse
		(*common.File)(nil),   // 2: media.common.File
	}
)
var file_media_indexer_v1_index_service_proto_depIdxs = []int32{
	2, // 0: media.indexer.v1.IndexRequest.file:type_name -> media.common.File
	0, // 1: media.indexer.v1.IndexService.Index:input_type -> media.indexer.v1.IndexRequest
	1, // 2: media.indexer.v1.IndexService.Index:output_type -> media.indexer.v1.IndexResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_media_indexer_v1_index_service_proto_init() }
func file_media_indexer_v1_index_service_proto_init() {
	if File_media_indexer_v1_index_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_media_indexer_v1_index_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexRequest); i {
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
		file_media_indexer_v1_index_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexResponse); i {
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
			RawDescriptor: file_media_indexer_v1_index_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_media_indexer_v1_index_service_proto_goTypes,
		DependencyIndexes: file_media_indexer_v1_index_service_proto_depIdxs,
		MessageInfos:      file_media_indexer_v1_index_service_proto_msgTypes,
	}.Build()
	File_media_indexer_v1_index_service_proto = out.File
	file_media_indexer_v1_index_service_proto_rawDesc = nil
	file_media_indexer_v1_index_service_proto_goTypes = nil
	file_media_indexer_v1_index_service_proto_depIdxs = nil
}