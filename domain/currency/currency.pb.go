// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: domain/currency/currency.proto

package currency

import (
	decimal "github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	denom "github.com/CoreumFoundation/CoreDEX-API/domain/denom"
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

type Currencies struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Currencies []*Currency `protobuf:"bytes,1,rep,name=Currencies,proto3" json:"Currencies,omitempty"`
	Offset     *int32      `protobuf:"varint,2,opt,name=Offset,proto3,oneof" json:"Offset,omitempty"`
}

func (x *Currencies) Reset() {
	*x = Currencies{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_currency_currency_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Currencies) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Currencies) ProtoMessage() {}

func (x *Currencies) ProtoReflect() protoreflect.Message {
	mi := &file_domain_currency_currency_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Currencies.ProtoReflect.Descriptor instead.
func (*Currencies) Descriptor() ([]byte, []int) {
	return file_domain_currency_currency_proto_rawDescGZIP(), []int{0}
}

func (x *Currencies) GetCurrencies() []*Currency {
	if x != nil {
		return x.Currencies
	}
	return nil
}

func (x *Currencies) GetOffset() int32 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

type Currency struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Denom          *denom.Denom       `protobuf:"bytes,1,opt,name=Denom,proto3" json:"Denom,omitempty"`
	SendCommission *decimal.Decimal   `protobuf:"bytes,2,opt,name=SendCommission,proto3" json:"SendCommission,omitempty"`
	BurnRate       *decimal.Decimal   `protobuf:"bytes,3,opt,name=BurnRate,proto3" json:"BurnRate,omitempty"`
	InitialAmount  *decimal.Decimal   `protobuf:"bytes,4,opt,name=InitialAmount,proto3" json:"InitialAmount,omitempty"`
	Chain          string             `protobuf:"bytes,10,opt,name=Chain,proto3" json:"Chain,omitempty"`              // The chain the currency is on (used for IBC tokens, else you can not distinguish between currencies with the same name)
	OriginChain    string             `protobuf:"bytes,11,opt,name=OriginChain,proto3" json:"OriginChain,omitempty"`  // The chain the currency is on (The actual chain which the currency originates from, used for IBC tokens)
	ChainSupply    string             `protobuf:"bytes,12,opt,name=ChainSupply,proto3" json:"ChainSupply,omitempty"`  // The total supply of the currency on the chain (used for IBC tokens)
	Description    string             `protobuf:"bytes,13,opt,name=Description,proto3" json:"Description,omitempty"`  // The description of the currency (used for IBC tokens)
	SkipDisplay    bool               `protobuf:"varint,20,opt,name=SkipDisplay,proto3" json:"SkipDisplay,omitempty"` // Indicates if the currency should be skipped in the display (mainly used to disable 13k+ IBC tokens from being loaded)
	MetaData       *metadata.MetaData `protobuf:"bytes,30,opt,name=MetaData,proto3" json:"MetaData,omitempty"`
}

func (x *Currency) Reset() {
	*x = Currency{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_currency_currency_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Currency) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Currency) ProtoMessage() {}

func (x *Currency) ProtoReflect() protoreflect.Message {
	mi := &file_domain_currency_currency_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Currency.ProtoReflect.Descriptor instead.
func (*Currency) Descriptor() ([]byte, []int) {
	return file_domain_currency_currency_proto_rawDescGZIP(), []int{1}
}

func (x *Currency) GetDenom() *denom.Denom {
	if x != nil {
		return x.Denom
	}
	return nil
}

func (x *Currency) GetSendCommission() *decimal.Decimal {
	if x != nil {
		return x.SendCommission
	}
	return nil
}

func (x *Currency) GetBurnRate() *decimal.Decimal {
	if x != nil {
		return x.BurnRate
	}
	return nil
}

func (x *Currency) GetInitialAmount() *decimal.Decimal {
	if x != nil {
		return x.InitialAmount
	}
	return nil
}

func (x *Currency) GetChain() string {
	if x != nil {
		return x.Chain
	}
	return ""
}

func (x *Currency) GetOriginChain() string {
	if x != nil {
		return x.OriginChain
	}
	return ""
}

func (x *Currency) GetChainSupply() string {
	if x != nil {
		return x.ChainSupply
	}
	return ""
}

func (x *Currency) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Currency) GetSkipDisplay() bool {
	if x != nil {
		return x.SkipDisplay
	}
	return false
}

func (x *Currency) GetMetaData() *metadata.MetaData {
	if x != nil {
		return x.MetaData
	}
	return nil
}

var File_domain_currency_currency_proto protoreflect.FileDescriptor

var file_domain_currency_currency_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63,
	0x79, 0x2f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x1a, 0x18, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x2f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x2f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x64, 0x65, 0x63,
	0x69, 0x6d, 0x61, 0x6c, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x68, 0x0a, 0x0a, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73,
	0x12, 0x32, 0x0a, 0x0a, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x0a, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x69, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x06, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x06, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x88, 0x01,
	0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x9c, 0x03, 0x0a,
	0x08, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x22, 0x0a, 0x05, 0x44, 0x65, 0x6e,
	0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6e, 0x6f, 0x6d,
	0x2e, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x52, 0x05, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x12, 0x38, 0x0a,
	0x0e, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e,
	0x44, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x52, 0x0e, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x08, 0x42, 0x75, 0x72, 0x6e, 0x52,
	0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x65, 0x63, 0x69,
	0x6d, 0x61, 0x6c, 0x2e, 0x44, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x52, 0x08, 0x42, 0x75, 0x72,
	0x6e, 0x52, 0x61, 0x74, 0x65, 0x12, 0x36, 0x0a, 0x0d, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c,
	0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64,
	0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x44, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x52, 0x0d,
	0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x53, 0x75,
	0x70, 0x70, 0x6c, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x68, 0x61, 0x69,
	0x6e, 0x53, 0x75, 0x70, 0x70, 0x6c, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x6b, 0x69,
	0x70, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x14, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x53, 0x6b, 0x69, 0x70, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x12, 0x2e, 0x0a, 0x08, 0x4d,
	0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x42, 0x42, 0x5a, 0x40, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x6f, 0x72, 0x65, 0x75, 0x6d,
	0x46, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x43, 0x6f, 0x72, 0x65, 0x44,
	0x45, 0x58, 0x2d, 0x41, 0x50, 0x49, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x3b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_domain_currency_currency_proto_rawDescOnce sync.Once
	file_domain_currency_currency_proto_rawDescData = file_domain_currency_currency_proto_rawDesc
)

func file_domain_currency_currency_proto_rawDescGZIP() []byte {
	file_domain_currency_currency_proto_rawDescOnce.Do(func() {
		file_domain_currency_currency_proto_rawDescData = protoimpl.X.CompressGZIP(file_domain_currency_currency_proto_rawDescData)
	})
	return file_domain_currency_currency_proto_rawDescData
}

var file_domain_currency_currency_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_domain_currency_currency_proto_goTypes = []interface{}{
	(*Currencies)(nil),        // 0: currency.Currencies
	(*Currency)(nil),          // 1: currency.Currency
	(*denom.Denom)(nil),       // 2: denom.Denom
	(*decimal.Decimal)(nil),   // 3: decimal.Decimal
	(*metadata.MetaData)(nil), // 4: metadata.MetaData
}
var file_domain_currency_currency_proto_depIdxs = []int32{
	1, // 0: currency.Currencies.Currencies:type_name -> currency.Currency
	2, // 1: currency.Currency.Denom:type_name -> denom.Denom
	3, // 2: currency.Currency.SendCommission:type_name -> decimal.Decimal
	3, // 3: currency.Currency.BurnRate:type_name -> decimal.Decimal
	3, // 4: currency.Currency.InitialAmount:type_name -> decimal.Decimal
	4, // 5: currency.Currency.MetaData:type_name -> metadata.MetaData
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_domain_currency_currency_proto_init() }
func file_domain_currency_currency_proto_init() {
	if File_domain_currency_currency_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_domain_currency_currency_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Currencies); i {
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
		file_domain_currency_currency_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Currency); i {
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
	file_domain_currency_currency_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_domain_currency_currency_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_domain_currency_currency_proto_goTypes,
		DependencyIndexes: file_domain_currency_currency_proto_depIdxs,
		MessageInfos:      file_domain_currency_currency_proto_msgTypes,
	}.Build()
	File_domain_currency_currency_proto = out.File
	file_domain_currency_currency_proto_rawDesc = nil
	file_domain_currency_currency_proto_goTypes = nil
	file_domain_currency_currency_proto_depIdxs = nil
}
