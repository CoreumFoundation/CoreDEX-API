// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: domain/trade/trade.proto

package trade

import (
	decimal "github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	denom "github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	metadata "github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	order_properties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Key in store is TXID-Sequence-Metadata.Network
type Trade struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Account  string           `protobuf:"bytes,1,opt,name=Account,proto3" json:"Account,omitempty"`
	OrderID  string           `protobuf:"bytes,2,opt,name=OrderID,proto3" json:"OrderID,omitempty"`    // User assigned order reference
	Sequence int64            `protobuf:"varint,3,opt,name=Sequence,proto3" json:"Sequence,omitempty"` // The sequence number of the order, assigned by the DEX (guaranteed unique value for the order)
	Amount   *decimal.Decimal `protobuf:"bytes,4,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Price    float64          `protobuf:"fixed64,5,opt,name=Price,proto3" json:"Price,omitempty"`
	Denom1   *denom.Denom     `protobuf:"bytes,6,opt,name=Denom1,proto3" json:"Denom1,omitempty"`
	Denom2   *denom.Denom     `protobuf:"bytes,7,opt,name=Denom2,proto3" json:"Denom2,omitempty"`
	// The buy/sell (e.g. did the user place a buy or sell order)
	Side      order_properties.Side  `protobuf:"varint,8,opt,name=Side,proto3,enum=orderproperties.Side" json:"Side,omitempty"`
	BlockTime *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=BlockTime,proto3" json:"BlockTime,omitempty"` // The time the trade was executed in UTC
	// Standard storage related fields
	MetaData    *metadata.MetaData `protobuf:"bytes,30,opt,name=MetaData,proto3" json:"MetaData,omitempty"`
	TXID        *string            `protobuf:"bytes,31,opt,name=TXID,proto3,oneof" json:"TXID,omitempty"`
	BlockHeight int64              `protobuf:"varint,32,opt,name=BlockHeight,proto3" json:"BlockHeight,omitempty"`
	Enriched    bool               `protobuf:"varint,33,opt,name=Enriched,proto3" json:"Enriched,omitempty"` // If the trade has been enriched with precision data
	// USD representation of the trade values and trading fee (fixed base for easy data comparisson in reports etc)
	USD *float32 `protobuf:"fixed32,40,opt,name=USD,proto3,oneof" json:"USD,omitempty"` // The USD value of the trade, calculated from the USD value of the currencies and the trading fee.
	// Trades get stored in alphabetical order of the denom pair.
	// Data is "uninverted" on retrieval and
	// this flag only indicates that the denoms as seen in the record are not in the original order
	Inverted bool `protobuf:"varint,50,opt,name=Inverted,proto3" json:"Inverted,omitempty"`
}

func (x *Trade) Reset() {
	*x = Trade{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_trade_trade_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Trade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trade) ProtoMessage() {}

func (x *Trade) ProtoReflect() protoreflect.Message {
	mi := &file_domain_trade_trade_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trade.ProtoReflect.Descriptor instead.
func (*Trade) Descriptor() ([]byte, []int) {
	return file_domain_trade_trade_proto_rawDescGZIP(), []int{0}
}

func (x *Trade) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *Trade) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

func (x *Trade) GetSequence() int64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *Trade) GetAmount() *decimal.Decimal {
	if x != nil {
		return x.Amount
	}
	return nil
}

func (x *Trade) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Trade) GetDenom1() *denom.Denom {
	if x != nil {
		return x.Denom1
	}
	return nil
}

func (x *Trade) GetDenom2() *denom.Denom {
	if x != nil {
		return x.Denom2
	}
	return nil
}

func (x *Trade) GetSide() order_properties.Side {
	if x != nil {
		return x.Side
	}
	return order_properties.Side(0)
}

func (x *Trade) GetBlockTime() *timestamppb.Timestamp {
	if x != nil {
		return x.BlockTime
	}
	return nil
}

func (x *Trade) GetMetaData() *metadata.MetaData {
	if x != nil {
		return x.MetaData
	}
	return nil
}

func (x *Trade) GetTXID() string {
	if x != nil && x.TXID != nil {
		return *x.TXID
	}
	return ""
}

func (x *Trade) GetBlockHeight() int64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *Trade) GetEnriched() bool {
	if x != nil {
		return x.Enriched
	}
	return false
}

func (x *Trade) GetUSD() float32 {
	if x != nil && x.USD != nil {
		return *x.USD
	}
	return 0
}

func (x *Trade) GetInverted() bool {
	if x != nil {
		return x.Inverted
	}
	return false
}

type Trades struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Trades []*Trade `protobuf:"bytes,1,rep,name=Trades,proto3" json:"Trades,omitempty"`
}

func (x *Trades) Reset() {
	*x = Trades{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_trade_trade_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Trades) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trades) ProtoMessage() {}

func (x *Trades) ProtoReflect() protoreflect.Message {
	mi := &file_domain_trade_trade_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trades.ProtoReflect.Descriptor instead.
func (*Trades) Descriptor() ([]byte, []int) {
	return file_domain_trade_trade_proto_rawDescGZIP(), []int{1}
}

func (x *Trades) GetTrades() []*Trade {
	if x != nil {
		return x.Trades
	}
	return nil
}

type TradePair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Denom1   *denom.Denom       `protobuf:"bytes,1,opt,name=Denom1,proto3" json:"Denom1,omitempty"`
	Denom2   *denom.Denom       `protobuf:"bytes,2,opt,name=Denom2,proto3" json:"Denom2,omitempty"`
	MetaData *metadata.MetaData `protobuf:"bytes,3,opt,name=MetaData,proto3" json:"MetaData,omitempty"`
}

func (x *TradePair) Reset() {
	*x = TradePair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_trade_trade_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradePair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradePair) ProtoMessage() {}

func (x *TradePair) ProtoReflect() protoreflect.Message {
	mi := &file_domain_trade_trade_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradePair.ProtoReflect.Descriptor instead.
func (*TradePair) Descriptor() ([]byte, []int) {
	return file_domain_trade_trade_proto_rawDescGZIP(), []int{2}
}

func (x *TradePair) GetDenom1() *denom.Denom {
	if x != nil {
		return x.Denom1
	}
	return nil
}

func (x *TradePair) GetDenom2() *denom.Denom {
	if x != nil {
		return x.Denom2
	}
	return nil
}

func (x *TradePair) GetMetaData() *metadata.MetaData {
	if x != nil {
		return x.MetaData
	}
	return nil
}

type TradePairs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TradePairs []*TradePair `protobuf:"bytes,1,rep,name=TradePairs,proto3" json:"TradePairs,omitempty"`
	Offset     *int32       `protobuf:"varint,2,opt,name=Offset,proto3,oneof" json:"Offset,omitempty"`
}

func (x *TradePairs) Reset() {
	*x = TradePairs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_trade_trade_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradePairs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradePairs) ProtoMessage() {}

func (x *TradePairs) ProtoReflect() protoreflect.Message {
	mi := &file_domain_trade_trade_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradePairs.ProtoReflect.Descriptor instead.
func (*TradePairs) Descriptor() ([]byte, []int) {
	return file_domain_trade_trade_proto_rawDescGZIP(), []int{3}
}

func (x *TradePairs) GetTradePairs() []*TradePair {
	if x != nil {
		return x.TradePairs
	}
	return nil
}

func (x *TradePairs) GetOffset() int32 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

var File_domain_trade_trade_proto protoreflect.FileDescriptor

var file_domain_trade_trade_proto_rawDesc = []byte{
	0x0a, 0x18, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x65, 0x2f, 0x74,
	0x72, 0x61, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x72, 0x61, 0x64,
	0x65, 0x1a, 0x18, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x2f,
	0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2f, 0x64, 0x65, 0x63, 0x69,
	0x6d, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74,
	0x69, 0x65, 0x73, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72,
	0x74, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x04, 0x0a, 0x05, 0x54,
	0x72, 0x61, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x44,
	0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x52, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x31, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x2e, 0x44, 0x65, 0x6e,
	0x6f, 0x6d, 0x52, 0x06, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x31, 0x12, 0x24, 0x0a, 0x06, 0x44, 0x65,
	0x6e, 0x6f, 0x6d, 0x32, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6e,
	0x6f, 0x6d, 0x2e, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x52, 0x06, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x32,
	0x12, 0x29, 0x0a, 0x04, 0x53, 0x69, 0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15,
	0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73,
	0x2e, 0x53, 0x69, 0x64, 0x65, 0x52, 0x04, 0x53, 0x69, 0x64, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x52, 0x08, 0x4d, 0x65, 0x74,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x17, 0x0a, 0x04, 0x54, 0x58, 0x49, 0x44, 0x18, 0x1f, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x54, 0x58, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x20,
	0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x20, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x45, 0x6e, 0x72, 0x69, 0x63, 0x68, 0x65, 0x64, 0x18, 0x21, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x45, 0x6e, 0x72, 0x69, 0x63, 0x68, 0x65, 0x64, 0x12, 0x15, 0x0a, 0x03,
	0x55, 0x53, 0x44, 0x18, 0x28, 0x20, 0x01, 0x28, 0x02, 0x48, 0x01, 0x52, 0x03, 0x55, 0x53, 0x44,
	0x88, 0x01, 0x01, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x64, 0x18,
	0x32, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x49, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x64, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x54, 0x58, 0x49, 0x44, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x55, 0x53, 0x44,
	0x22, 0x2e, 0x0a, 0x06, 0x54, 0x72, 0x61, 0x64, 0x65, 0x73, 0x12, 0x24, 0x0a, 0x06, 0x54, 0x72,
	0x61, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x72, 0x61,
	0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x52, 0x06, 0x54, 0x72, 0x61, 0x64, 0x65, 0x73,
	0x22, 0x87, 0x01, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x64, 0x65, 0x50, 0x61, 0x69, 0x72, 0x12, 0x24,
	0x0a, 0x06, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x2e, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x52, 0x06, 0x44, 0x65,
	0x6e, 0x6f, 0x6d, 0x31, 0x12, 0x24, 0x0a, 0x06, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x32, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x2e, 0x44, 0x65, 0x6e,
	0x6f, 0x6d, 0x52, 0x06, 0x44, 0x65, 0x6e, 0x6f, 0x6d, 0x32, 0x12, 0x2e, 0x0a, 0x08, 0x4d, 0x65,
	0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x22, 0x66, 0x0a, 0x0a, 0x54, 0x72,
	0x61, 0x64, 0x65, 0x50, 0x61, 0x69, 0x72, 0x73, 0x12, 0x30, 0x0a, 0x0a, 0x54, 0x72, 0x61, 0x64,
	0x65, 0x50, 0x61, 0x69, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x74,
	0x72, 0x61, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x50, 0x61, 0x69, 0x72, 0x52, 0x0a,
	0x54, 0x72, 0x61, 0x64, 0x65, 0x50, 0x61, 0x69, 0x72, 0x73, 0x12, 0x1b, 0x0a, 0x06, 0x4f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x06, 0x4f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x43, 0x6f, 0x72, 0x65, 0x75, 0x6d, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x43, 0x6f, 0x72, 0x65, 0x44, 0x45, 0x58, 0x2d, 0x41, 0x50, 0x49, 0x2f, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x65, 0x3b, 0x74, 0x72, 0x61, 0x64, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_domain_trade_trade_proto_rawDescOnce sync.Once
	file_domain_trade_trade_proto_rawDescData = file_domain_trade_trade_proto_rawDesc
)

func file_domain_trade_trade_proto_rawDescGZIP() []byte {
	file_domain_trade_trade_proto_rawDescOnce.Do(func() {
		file_domain_trade_trade_proto_rawDescData = protoimpl.X.CompressGZIP(file_domain_trade_trade_proto_rawDescData)
	})
	return file_domain_trade_trade_proto_rawDescData
}

var file_domain_trade_trade_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_domain_trade_trade_proto_goTypes = []interface{}{
	(*Trade)(nil),                 // 0: trade.Trade
	(*Trades)(nil),                // 1: trade.Trades
	(*TradePair)(nil),             // 2: trade.TradePair
	(*TradePairs)(nil),            // 3: trade.TradePairs
	(*decimal.Decimal)(nil),       // 4: decimal.Decimal
	(*denom.Denom)(nil),           // 5: denom.Denom
	(order_properties.Side)(0),    // 6: orderproperties.Side
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
	(*metadata.MetaData)(nil),     // 8: metadata.MetaData
}
var file_domain_trade_trade_proto_depIdxs = []int32{
	4,  // 0: trade.Trade.Amount:type_name -> decimal.Decimal
	5,  // 1: trade.Trade.Denom1:type_name -> denom.Denom
	5,  // 2: trade.Trade.Denom2:type_name -> denom.Denom
	6,  // 3: trade.Trade.Side:type_name -> orderproperties.Side
	7,  // 4: trade.Trade.BlockTime:type_name -> google.protobuf.Timestamp
	8,  // 5: trade.Trade.MetaData:type_name -> metadata.MetaData
	0,  // 6: trade.Trades.Trades:type_name -> trade.Trade
	5,  // 7: trade.TradePair.Denom1:type_name -> denom.Denom
	5,  // 8: trade.TradePair.Denom2:type_name -> denom.Denom
	8,  // 9: trade.TradePair.MetaData:type_name -> metadata.MetaData
	2,  // 10: trade.TradePairs.TradePairs:type_name -> trade.TradePair
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_domain_trade_trade_proto_init() }
func file_domain_trade_trade_proto_init() {
	if File_domain_trade_trade_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_domain_trade_trade_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Trade); i {
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
		file_domain_trade_trade_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Trades); i {
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
		file_domain_trade_trade_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TradePair); i {
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
		file_domain_trade_trade_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TradePairs); i {
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
	file_domain_trade_trade_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_domain_trade_trade_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_domain_trade_trade_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_domain_trade_trade_proto_goTypes,
		DependencyIndexes: file_domain_trade_trade_proto_depIdxs,
		MessageInfos:      file_domain_trade_trade_proto_msgTypes,
	}.Build()
	File_domain_trade_trade_proto = out.File
	file_domain_trade_trade_proto_rawDesc = nil
	file_domain_trade_trade_proto_goTypes = nil
	file_domain_trade_trade_proto_depIdxs = nil
}
