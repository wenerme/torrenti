// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: media/scraper/v1/scrape_service.proto

package scraperv1

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/wenerme/torrenti/pkg/apis/media/common"
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

type StateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *StateRequest) Reset() {
	*x = StateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateRequest) ProtoMessage() {}

func (x *StateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateRequest.ProtoReflect.Descriptor instead.
func (*StateRequest) Descriptor() ([]byte, []int) {
	return file_media_scraper_v1_scrape_service_proto_rawDescGZIP(), []int{0}
}

func (x *StateRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type StateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url      string  `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Visiting bool    `protobuf:"varint,2,opt,name=visiting,proto3" json:"visiting,omitempty"`
	Scraped  bool    `protobuf:"varint,3,opt,name=scraped,proto3" json:"scraped,omitempty"`
	Error    *string `protobuf:"bytes,4,opt,name=error,proto3,oneof" json:"error,omitempty"`
}

func (x *StateResponse) Reset() {
	*x = StateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateResponse) ProtoMessage() {}

func (x *StateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateResponse.ProtoReflect.Descriptor instead.
func (*StateResponse) Descriptor() ([]byte, []int) {
	return file_media_scraper_v1_scrape_service_proto_rawDescGZIP(), []int{1}
}

func (x *StateResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *StateResponse) GetVisiting() bool {
	if x != nil {
		return x.Visiting
	}
	return false
}

func (x *StateResponse) GetScraped() bool {
	if x != nil {
		return x.Scraped
	}
	return false
}

func (x *StateResponse) GetError() string {
	if x != nil && x.Error != nil {
		return *x.Error
	}
	return ""
}

type ScrapeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url     string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Referer string `protobuf:"bytes,2,opt,name=referer,proto3" json:"referer,omitempty"`
}

func (x *ScrapeRequest) Reset() {
	*x = ScrapeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScrapeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScrapeRequest) ProtoMessage() {}

func (x *ScrapeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScrapeRequest.ProtoReflect.Descriptor instead.
func (*ScrapeRequest) Descriptor() ([]byte, []int) {
	return file_media_scraper_v1_scrape_service_proto_rawDescGZIP(), []int{2}
}

func (x *ScrapeRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ScrapeRequest) GetReferer() string {
	if x != nil {
		return x.Referer
	}
	return ""
}

type ScrapeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ScrapeResponse) Reset() {
	*x = ScrapeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScrapeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScrapeResponse) ProtoMessage() {}

func (x *ScrapeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_media_scraper_v1_scrape_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScrapeResponse.ProtoReflect.Descriptor instead.
func (*ScrapeResponse) Descriptor() ([]byte, []int) {
	return file_media_scraper_v1_scrape_service_proto_rawDescGZIP(), []int{3}
}

func (x *ScrapeResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_media_scraper_v1_scrape_service_proto protoreflect.FileDescriptor

var file_media_scraper_v1_scrape_service_proto_rawDesc = []byte{
	0x0a, 0x25, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2f,
	0x76, 0x31, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x73,
	0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x62, 0x6f, 0x64, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x20, 0x0a, 0x0c,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x7c,
	0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x69, 0x73, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x76, 0x69, 0x73, 0x69, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x88,
	0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3b, 0x0a, 0x0d,
	0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x18, 0x0a, 0x07, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x72, 0x22, 0x22, 0x0a, 0x0e, 0x53, 0x63, 0x72,
	0x61, 0x70, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x32, 0xca, 0x01,
	0x0a, 0x0d, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x5f, 0x0a, 0x06, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x12, 0x1f, 0x2e, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x72,
	0x61, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63,
	0x72, 0x61, 0x70, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0c, 0x22, 0x07, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x3a, 0x01, 0x2a,
	0x12, 0x58, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0e, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x08, 0x12, 0x06, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0xcd, 0x01, 0x0a, 0x14, 0x63,
	0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x42, 0x12, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x6e, 0x65, 0x72, 0x6d, 0x65, 0x2f, 0x74, 0x6f,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2f, 0x76, 0x31,
	0x3b, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x4d, 0x53, 0x58,
	0xaa, 0x02, 0x10, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x10, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x5c, 0x53, 0x63, 0x72, 0x61,
	0x70, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1c, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x5c, 0x53,
	0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x12, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x3a, 0x3a, 0x53,
	0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_media_scraper_v1_scrape_service_proto_rawDescOnce sync.Once
	file_media_scraper_v1_scrape_service_proto_rawDescData = file_media_scraper_v1_scrape_service_proto_rawDesc
)

func file_media_scraper_v1_scrape_service_proto_rawDescGZIP() []byte {
	file_media_scraper_v1_scrape_service_proto_rawDescOnce.Do(func() {
		file_media_scraper_v1_scrape_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_media_scraper_v1_scrape_service_proto_rawDescData)
	})
	return file_media_scraper_v1_scrape_service_proto_rawDescData
}

var (
	file_media_scraper_v1_scrape_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
	file_media_scraper_v1_scrape_service_proto_goTypes  = []interface{}{
		(*StateRequest)(nil),   // 0: media.scraper.v1.StateRequest
		(*StateResponse)(nil),  // 1: media.scraper.v1.StateResponse
		(*ScrapeRequest)(nil),  // 2: media.scraper.v1.ScrapeRequest
		(*ScrapeResponse)(nil), // 3: media.scraper.v1.ScrapeResponse
	}
)
var file_media_scraper_v1_scrape_service_proto_depIdxs = []int32{
	2, // 0: media.scraper.v1.ScrapeService.Scrape:input_type -> media.scraper.v1.ScrapeRequest
	0, // 1: media.scraper.v1.ScrapeService.State:input_type -> media.scraper.v1.StateRequest
	3, // 2: media.scraper.v1.ScrapeService.Scrape:output_type -> media.scraper.v1.ScrapeResponse
	1, // 3: media.scraper.v1.ScrapeService.State:output_type -> media.scraper.v1.StateResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_media_scraper_v1_scrape_service_proto_init() }
func file_media_scraper_v1_scrape_service_proto_init() {
	if File_media_scraper_v1_scrape_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_media_scraper_v1_scrape_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateRequest); i {
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
		file_media_scraper_v1_scrape_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateResponse); i {
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
		file_media_scraper_v1_scrape_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScrapeRequest); i {
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
		file_media_scraper_v1_scrape_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScrapeResponse); i {
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
	file_media_scraper_v1_scrape_service_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_media_scraper_v1_scrape_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_media_scraper_v1_scrape_service_proto_goTypes,
		DependencyIndexes: file_media_scraper_v1_scrape_service_proto_depIdxs,
		MessageInfos:      file_media_scraper_v1_scrape_service_proto_msgTypes,
	}.Build()
	File_media_scraper_v1_scrape_service_proto = out.File
	file_media_scraper_v1_scrape_service_proto_rawDesc = nil
	file_media_scraper_v1_scrape_service_proto_goTypes = nil
	file_media_scraper_v1_scrape_service_proto_depIdxs = nil
}