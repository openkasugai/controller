// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: api/operator/operator.proto

package operator

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

type OperateResponseStatus int32

const (
	OperateResponseStatus_UNKNOWN OperateResponseStatus = 0
	OperateResponseStatus_OK      OperateResponseStatus = 1
	OperateResponseStatus_NG      OperateResponseStatus = 2
)

// Enum value maps for OperateResponseStatus.
var (
	OperateResponseStatus_name = map[int32]string{
		0: "UNKNOWN",
		1: "OK",
		2: "NG",
	}
	OperateResponseStatus_value = map[string]int32{
		"UNKNOWN": 0,
		"OK":      1,
		"NG":      2,
	}
)

func (x OperateResponseStatus) Enum() *OperateResponseStatus {
	p := new(OperateResponseStatus)
	*p = x
	return p
}

func (x OperateResponseStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OperateResponseStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_operator_operator_proto_enumTypes[0].Descriptor()
}

func (OperateResponseStatus) Type() protoreflect.EnumType {
	return &file_api_operator_operator_proto_enumTypes[0]
}

func (x OperateResponseStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OperateResponseStatus.Descriptor instead.
func (OperateResponseStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_operator_operator_proto_rawDescGZIP(), []int{0}
}

type OperateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inputs  []float64        `protobuf:"fixed64,1,rep,packed,name=inputs,proto3" json:"inputs,omitempty"`
	Results []*OperateResult `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *OperateRequest) Reset() {
	*x = OperateRequest{}
	mi := &file_api_operator_operator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OperateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperateRequest) ProtoMessage() {}

func (x *OperateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_operator_operator_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperateRequest.ProtoReflect.Descriptor instead.
func (*OperateRequest) Descriptor() ([]byte, []int) {
	return file_api_operator_operator_proto_rawDescGZIP(), []int{0}
}

func (x *OperateRequest) GetInputs() []float64 {
	if x != nil {
		return x.Inputs
	}
	return nil
}

func (x *OperateRequest) GetResults() []*OperateResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type OperateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status OperateResponseStatus `protobuf:"varint,1,opt,name=status,proto3,enum=calcapp.grpc.operator.OperateResponseStatus" json:"status,omitempty"`
}

func (x *OperateResponse) Reset() {
	*x = OperateResponse{}
	mi := &file_api_operator_operator_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OperateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperateResponse) ProtoMessage() {}

func (x *OperateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_operator_operator_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperateResponse.ProtoReflect.Descriptor instead.
func (*OperateResponse) Descriptor() ([]byte, []int) {
	return file_api_operator_operator_proto_rawDescGZIP(), []int{1}
}

func (x *OperateResponse) GetStatus() OperateResponseStatus {
	if x != nil {
		return x.Status
	}
	return OperateResponseStatus_UNKNOWN
}

type OperateResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operator string  `protobuf:"bytes,1,opt,name=operator,proto3" json:"operator,omitempty"`
	Value    float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *OperateResult) Reset() {
	*x = OperateResult{}
	mi := &file_api_operator_operator_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OperateResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperateResult) ProtoMessage() {}

func (x *OperateResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_operator_operator_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperateResult.ProtoReflect.Descriptor instead.
func (*OperateResult) Descriptor() ([]byte, []int) {
	return file_api_operator_operator_proto_rawDescGZIP(), []int{2}
}

func (x *OperateResult) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *OperateResult) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_api_operator_operator_proto protoreflect.FileDescriptor

var file_api_operator_operator_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x63,
	0x61, 0x6c, 0x63, 0x61, 0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x22, 0x68, 0x0a, 0x0e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x12, 0x3e,
	0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x61, 0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22, 0x57,
	0x0a, 0x0f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x44, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x2c, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x61, 0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x41, 0x0a, 0x0d, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x34, 0x0a, 0x15, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x4e, 0x47, 0x10, 0x02,
	0x32, 0x6d, 0x0a, 0x0f, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x07, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x65, 0x12, 0x25,
	0x2e, 0x63, 0x61, 0x6c, 0x63, 0x61, 0x70, 0x70, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x61, 0x70, 0x70, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x16, 0x5a, 0x14, 0x63, 0x61, 0x6c, 0x63, 0x61, 0x70, 0x70, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_operator_operator_proto_rawDescOnce sync.Once
	file_api_operator_operator_proto_rawDescData = file_api_operator_operator_proto_rawDesc
)

func file_api_operator_operator_proto_rawDescGZIP() []byte {
	file_api_operator_operator_proto_rawDescOnce.Do(func() {
		file_api_operator_operator_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_operator_operator_proto_rawDescData)
	})
	return file_api_operator_operator_proto_rawDescData
}

var file_api_operator_operator_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_operator_operator_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_operator_operator_proto_goTypes = []any{
	(OperateResponseStatus)(0), // 0: calcapp.grpc.operator.OperateResponseStatus
	(*OperateRequest)(nil),     // 1: calcapp.grpc.operator.OperateRequest
	(*OperateResponse)(nil),    // 2: calcapp.grpc.operator.OperateResponse
	(*OperateResult)(nil),      // 3: calcapp.grpc.operator.OperateResult
}
var file_api_operator_operator_proto_depIdxs = []int32{
	3, // 0: calcapp.grpc.operator.OperateRequest.results:type_name -> calcapp.grpc.operator.OperateResult
	0, // 1: calcapp.grpc.operator.OperateResponse.status:type_name -> calcapp.grpc.operator.OperateResponseStatus
	1, // 2: calcapp.grpc.operator.OperatorService.Operate:input_type -> calcapp.grpc.operator.OperateRequest
	2, // 3: calcapp.grpc.operator.OperatorService.Operate:output_type -> calcapp.grpc.operator.OperateResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_operator_operator_proto_init() }
func file_api_operator_operator_proto_init() {
	if File_api_operator_operator_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_operator_operator_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_operator_operator_proto_goTypes,
		DependencyIndexes: file_api_operator_operator_proto_depIdxs,
		EnumInfos:         file_api_operator_operator_proto_enumTypes,
		MessageInfos:      file_api_operator_operator_proto_msgTypes,
	}.Build()
	File_api_operator_operator_proto = out.File
	file_api_operator_operator_proto_rawDesc = nil
	file_api_operator_operator_proto_goTypes = nil
	file_api_operator_operator_proto_depIdxs = nil
}
