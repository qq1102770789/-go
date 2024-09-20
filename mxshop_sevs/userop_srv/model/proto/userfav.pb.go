// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: userfav.proto

package proto

import (
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

type UserFavRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  int32 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	GoodsId int32 `protobuf:"varint,2,opt,name=goodsId,proto3" json:"goodsId,omitempty"`
}

func (x *UserFavRequest) Reset() {
	*x = UserFavRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userfav_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserFavRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserFavRequest) ProtoMessage() {}

func (x *UserFavRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userfav_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserFavRequest.ProtoReflect.Descriptor instead.
func (*UserFavRequest) Descriptor() ([]byte, []int) {
	return file_userfav_proto_rawDescGZIP(), []int{0}
}

func (x *UserFavRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserFavRequest) GetGoodsId() int32 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

type UserFavResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  int32 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	GoodsId int32 `protobuf:"varint,2,opt,name=goodsId,proto3" json:"goodsId,omitempty"`
}

func (x *UserFavResponse) Reset() {
	*x = UserFavResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userfav_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserFavResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserFavResponse) ProtoMessage() {}

func (x *UserFavResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userfav_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserFavResponse.ProtoReflect.Descriptor instead.
func (*UserFavResponse) Descriptor() ([]byte, []int) {
	return file_userfav_proto_rawDescGZIP(), []int{1}
}

func (x *UserFavResponse) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserFavResponse) GetGoodsId() int32 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

type UserFavListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total int32              `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Data  []*UserFavResponse `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *UserFavListResponse) Reset() {
	*x = UserFavListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_userfav_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserFavListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserFavListResponse) ProtoMessage() {}

func (x *UserFavListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userfav_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserFavListResponse.ProtoReflect.Descriptor instead.
func (*UserFavListResponse) Descriptor() ([]byte, []int) {
	return file_userfav_proto_rawDescGZIP(), []int{2}
}

func (x *UserFavListResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *UserFavListResponse) GetData() []*UserFavResponse {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_userfav_proto protoreflect.FileDescriptor

var file_userfav_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x66, 0x61, 0x76, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x42,
	0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x6f, 0x6f, 0x64,
	0x73, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x67, 0x6f, 0x6f, 0x64, 0x73,
	0x49, 0x64, 0x22, 0x43, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x22, 0x51, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72, 0x46,
	0x61, 0x76, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xbc, 0x01, 0x0a, 0x07, 0x55,
	0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x12, 0x33, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x0f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0a, 0x41,
	0x64, 0x64, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x12, 0x0f, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x46, 0x61, 0x76, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x28, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x46, 0x61, 0x76, 0x12, 0x0f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2b, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x12, 0x0f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x46, 0x61, 0x76, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_userfav_proto_rawDescOnce sync.Once
	file_userfav_proto_rawDescData = file_userfav_proto_rawDesc
)

func file_userfav_proto_rawDescGZIP() []byte {
	file_userfav_proto_rawDescOnce.Do(func() {
		file_userfav_proto_rawDescData = protoimpl.X.CompressGZIP(file_userfav_proto_rawDescData)
	})
	return file_userfav_proto_rawDescData
}

var file_userfav_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_userfav_proto_goTypes = []interface{}{
	(*UserFavRequest)(nil),      // 0: UserFavRequest
	(*UserFavResponse)(nil),     // 1: UserFavResponse
	(*UserFavListResponse)(nil), // 2: UserFavListResponse
	(*Empty)(nil),               // 3: Empty
}
var file_userfav_proto_depIdxs = []int32{
	1, // 0: UserFavListResponse.data:type_name -> UserFavResponse
	0, // 1: UserFav.GetFavList:input_type -> UserFavRequest
	0, // 2: UserFav.AddUserFav:input_type -> UserFavRequest
	0, // 3: UserFav.DeleteUserFav:input_type -> UserFavRequest
	0, // 4: UserFav.GetUserFavDetail:input_type -> UserFavRequest
	2, // 5: UserFav.GetFavList:output_type -> UserFavListResponse
	3, // 6: UserFav.AddUserFav:output_type -> Empty
	3, // 7: UserFav.DeleteUserFav:output_type -> Empty
	3, // 8: UserFav.GetUserFavDetail:output_type -> Empty
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_userfav_proto_init() }
func file_userfav_proto_init() {
	if File_userfav_proto != nil {
		return
	}
	file_address_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_userfav_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserFavRequest); i {
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
		file_userfav_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserFavResponse); i {
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
		file_userfav_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserFavListResponse); i {
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
			RawDescriptor: file_userfav_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_userfav_proto_goTypes,
		DependencyIndexes: file_userfav_proto_depIdxs,
		MessageInfos:      file_userfav_proto_msgTypes,
	}.Build()
	File_userfav_proto = out.File
	file_userfav_proto_rawDesc = nil
	file_userfav_proto_goTypes = nil
	file_userfav_proto_depIdxs = nil
}
