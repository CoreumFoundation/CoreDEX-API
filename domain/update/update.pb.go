// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: domain/update/update.proto

package update

import (
	metadata "github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
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

type Action int32

const (
	Action_SUBSCRIBE   Action = 0
	Action_UNSUBSCRIBE Action = 1
	Action_CLOSE       Action = 2
	Action_RESPONSE    Action = 3
)

// Enum value maps for Action.
var (
	Action_name = map[int32]string{
		0: "SUBSCRIBE",
		1: "UNSUBSCRIBE",
		2: "CLOSE",
		3: "RESPONSE",
	}
	Action_value = map[string]int32{
		"SUBSCRIBE":   0,
		"UNSUBSCRIBE": 1,
		"CLOSE":       2,
		"RESPONSE":    3,
	}
)

func (x Action) Enum() *Action {
	p := new(Action)
	*p = x
	return p
}

func (x Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Action) Descriptor() protoreflect.EnumDescriptor {
	return file_domain_update_update_proto_enumTypes[0].Descriptor()
}

func (Action) Type() protoreflect.EnumType {
	return &file_domain_update_update_proto_enumTypes[0]
}

func (x Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Action.Descriptor instead.
func (Action) EnumDescriptor() ([]byte, []int) {
	return file_domain_update_update_proto_rawDescGZIP(), []int{0}
}

type Method int32

const (
	Method_METHOD_DO_NOT_USE                Method = 0
	Method_TRADES_FOR_SYMBOL                Method = 1 // ID: {denom1}_{denom2}
	Method_TRADES_FOR_ACCOUNT               Method = 2 // ID: {account}
	Method_TRADES_FOR_ACCOUNT_AND_SYMBOL    Method = 3 // ID: {account}_{denom1}_{denom2}
	Method_OHLC                             Method = 4 // ID: {denom1}_{denom2}_{interval}
	Method_TICKER                           Method = 5 // ID: {denom1}_{denom2}
	Method_ORDERBOOK                        Method = 6 // ID: {denom1}_{denom2}
	Method_ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT Method = 7 // ID: {account}_{denom1}_{denom2}
	Method_WALLET                           Method = 8 // ID: {account}
)

// Enum value maps for Method.
var (
	Method_name = map[int32]string{
		0: "METHOD_DO_NOT_USE",
		1: "TRADES_FOR_SYMBOL",
		2: "TRADES_FOR_ACCOUNT",
		3: "TRADES_FOR_ACCOUNT_AND_SYMBOL",
		4: "OHLC",
		5: "TICKER",
		6: "ORDERBOOK",
		7: "ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT",
		8: "WALLET",
	}
	Method_value = map[string]int32{
		"METHOD_DO_NOT_USE":                0,
		"TRADES_FOR_SYMBOL":                1,
		"TRADES_FOR_ACCOUNT":               2,
		"TRADES_FOR_ACCOUNT_AND_SYMBOL":    3,
		"OHLC":                             4,
		"TICKER":                           5,
		"ORDERBOOK":                        6,
		"ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT": 7,
		"WALLET":                           8,
	}
)

func (x Method) Enum() *Method {
	p := new(Method)
	*p = x
	return p
}

func (x Method) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Method) Descriptor() protoreflect.EnumDescriptor {
	return file_domain_update_update_proto_enumTypes[1].Descriptor()
}

func (Method) Type() protoreflect.EnumType {
	return &file_domain_update_update_proto_enumTypes[1]
}

func (x Method) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Method.Descriptor instead.
func (Method) EnumDescriptor() ([]byte, []int) {
	return file_domain_update_update_proto_rawDescGZIP(), []int{1}
}

type Subscribe struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action       Action        `protobuf:"varint,1,opt,name=Action,proto3,enum=update.Action" json:"Action,omitempty"`
	Subscription *Subscription `protobuf:"bytes,2,opt,name=Subscription,proto3" json:"Subscription,omitempty"`
}

func (x *Subscribe) Reset() {
	*x = Subscribe{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_update_update_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subscribe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subscribe) ProtoMessage() {}

func (x *Subscribe) ProtoReflect() protoreflect.Message {
	mi := &file_domain_update_update_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subscribe.ProtoReflect.Descriptor instead.
func (*Subscribe) Descriptor() ([]byte, []int) {
	return file_domain_update_update_proto_rawDescGZIP(), []int{0}
}

func (x *Subscribe) GetAction() Action {
	if x != nil {
		return x.Action
	}
	return Action_SUBSCRIBE
}

func (x *Subscribe) GetSubscription() *Subscription {
	if x != nil {
		return x.Subscription
	}
	return nil
}

type Subscription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method  Method           `protobuf:"varint,1,opt,name=Method,proto3,enum=update.Method" json:"Method,omitempty"`
	ID      string           `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Network metadata.Network `protobuf:"varint,3,opt,name=Network,proto3,enum=metadata.Network" json:"Network,omitempty"`
	Content string           `protobuf:"bytes,4,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (x *Subscription) Reset() {
	*x = Subscription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_update_update_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subscription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subscription) ProtoMessage() {}

func (x *Subscription) ProtoReflect() protoreflect.Message {
	mi := &file_domain_update_update_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subscription.ProtoReflect.Descriptor instead.
func (*Subscription) Descriptor() ([]byte, []int) {
	return file_domain_update_update_proto_rawDescGZIP(), []int{1}
}

func (x *Subscription) GetMethod() Method {
	if x != nil {
		return x.Method
	}
	return Method_METHOD_DO_NOT_USE
}

func (x *Subscription) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Subscription) GetNetwork() metadata.Network {
	if x != nil {
		return x.Network
	}
	return metadata.Network(0)
}

func (x *Subscription) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_domain_update_update_proto protoreflect.FileDescriptor

var file_domain_update_update_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x1a, 0x1e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x12, 0x26, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0e, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x0c, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x8d, 0x01, 0x0a, 0x0c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x52, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x2b, 0x0a, 0x07,
	0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x52, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x2a, 0x41, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0d, 0x0a,
	0x09, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b,
	0x55, 0x4e, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x10, 0x01, 0x12, 0x09, 0x0a,
	0x05, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x53, 0x50,
	0x4f, 0x4e, 0x53, 0x45, 0x10, 0x03, 0x2a, 0xc8, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x12, 0x15, 0x0a, 0x11, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x44, 0x4f, 0x5f, 0x4e,
	0x4f, 0x54, 0x5f, 0x55, 0x53, 0x45, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x52, 0x41, 0x44,
	0x45, 0x53, 0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x53, 0x59, 0x4d, 0x42, 0x4f, 0x4c, 0x10, 0x01, 0x12,
	0x16, 0x0a, 0x12, 0x54, 0x52, 0x41, 0x44, 0x45, 0x53, 0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x41, 0x43,
	0x43, 0x4f, 0x55, 0x4e, 0x54, 0x10, 0x02, 0x12, 0x21, 0x0a, 0x1d, 0x54, 0x52, 0x41, 0x44, 0x45,
	0x53, 0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x41, 0x4e,
	0x44, 0x5f, 0x53, 0x59, 0x4d, 0x42, 0x4f, 0x4c, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x4f, 0x48,
	0x4c, 0x43, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x54, 0x49, 0x43, 0x4b, 0x45, 0x52, 0x10, 0x05,
	0x12, 0x0d, 0x0a, 0x09, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x42, 0x4f, 0x4f, 0x4b, 0x10, 0x06, 0x12,
	0x24, 0x0a, 0x20, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x42, 0x4f, 0x4f, 0x4b, 0x5f, 0x46, 0x4f, 0x52,
	0x5f, 0x53, 0x59, 0x4d, 0x42, 0x4f, 0x4c, 0x5f, 0x41, 0x4e, 0x44, 0x5f, 0x41, 0x43, 0x43, 0x4f,
	0x55, 0x4e, 0x54, 0x10, 0x07, 0x12, 0x0a, 0x0a, 0x06, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x10,
	0x08, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x43, 0x6f, 0x72, 0x65, 0x75, 0x6d, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x43, 0x6f, 0x72, 0x65, 0x44, 0x45, 0x58, 0x2d, 0x41, 0x50, 0x49, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x3b, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_domain_update_update_proto_rawDescOnce sync.Once
	file_domain_update_update_proto_rawDescData = file_domain_update_update_proto_rawDesc
)

func file_domain_update_update_proto_rawDescGZIP() []byte {
	file_domain_update_update_proto_rawDescOnce.Do(func() {
		file_domain_update_update_proto_rawDescData = protoimpl.X.CompressGZIP(file_domain_update_update_proto_rawDescData)
	})
	return file_domain_update_update_proto_rawDescData
}

var file_domain_update_update_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_domain_update_update_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_domain_update_update_proto_goTypes = []interface{}{
	(Action)(0),           // 0: update.Action
	(Method)(0),           // 1: update.Method
	(*Subscribe)(nil),     // 2: update.Subscribe
	(*Subscription)(nil),  // 3: update.Subscription
	(metadata.Network)(0), // 4: metadata.Network
}
var file_domain_update_update_proto_depIdxs = []int32{
	0, // 0: update.Subscribe.Action:type_name -> update.Action
	3, // 1: update.Subscribe.Subscription:type_name -> update.Subscription
	1, // 2: update.Subscription.Method:type_name -> update.Method
	4, // 3: update.Subscription.Network:type_name -> metadata.Network
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_domain_update_update_proto_init() }
func file_domain_update_update_proto_init() {
	if File_domain_update_update_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_domain_update_update_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subscribe); i {
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
		file_domain_update_update_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subscription); i {
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
			RawDescriptor: file_domain_update_update_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_domain_update_update_proto_goTypes,
		DependencyIndexes: file_domain_update_update_proto_depIdxs,
		EnumInfos:         file_domain_update_update_proto_enumTypes,
		MessageInfos:      file_domain_update_update_proto_msgTypes,
	}.Build()
	File_domain_update_update_proto = out.File
	file_domain_update_update_proto_rawDesc = nil
	file_domain_update_update_proto_goTypes = nil
	file_domain_update_update_proto_depIdxs = nil
}
