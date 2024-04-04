// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: aggregator/aggregator.proto

package aggregator

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

type InitOperatorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The layer1 chain id for operator to use
	Layer1ChainId uint32 `protobuf:"varint,1,opt,name=layer1_chain_id,json=layer1ChainId,proto3" json:"layer1_chain_id,omitempty"`
	// The layer2 chain id for operator to use
	ChainId uint32 `protobuf:"varint,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// The operator 's id
	OperatorId []byte `protobuf:"bytes,3,opt,name=operator_id,json=operatorId,proto3" json:"operator_id,omitempty"`
	// The operator 's ecdsa address
	OperatorAddress string `protobuf:"bytes,4,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty"`
	// The operator_state_retriever_addr
	OperatorStateRetrieverAddr string `protobuf:"bytes,5,opt,name=operator_state_retriever_addr,json=operatorStateRetrieverAddr,proto3" json:"operator_state_retriever_addr,omitempty"`
	// The registry_coordinator_addr
	RegistryCoordinatorAddr string `protobuf:"bytes,6,opt,name=registry_coordinator_addr,json=registryCoordinatorAddr,proto3" json:"registry_coordinator_addr,omitempty"`
}

func (x *InitOperatorRequest) Reset() {
	*x = InitOperatorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_aggregator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitOperatorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitOperatorRequest) ProtoMessage() {}

func (x *InitOperatorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_aggregator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitOperatorRequest.ProtoReflect.Descriptor instead.
func (*InitOperatorRequest) Descriptor() ([]byte, []int) {
	return file_aggregator_aggregator_proto_rawDescGZIP(), []int{0}
}

func (x *InitOperatorRequest) GetLayer1ChainId() uint32 {
	if x != nil {
		return x.Layer1ChainId
	}
	return 0
}

func (x *InitOperatorRequest) GetChainId() uint32 {
	if x != nil {
		return x.ChainId
	}
	return 0
}

func (x *InitOperatorRequest) GetOperatorId() []byte {
	if x != nil {
		return x.OperatorId
	}
	return nil
}

func (x *InitOperatorRequest) GetOperatorAddress() string {
	if x != nil {
		return x.OperatorAddress
	}
	return ""
}

func (x *InitOperatorRequest) GetOperatorStateRetrieverAddr() string {
	if x != nil {
		return x.OperatorStateRetrieverAddr
	}
	return ""
}

func (x *InitOperatorRequest) GetRegistryCoordinatorAddr() string {
	if x != nil {
		return x.RegistryCoordinatorAddr
	}
	return ""
}

type InitOperatorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If the operator 's state is ok
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	// Reason
	Reason string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *InitOperatorResponse) Reset() {
	*x = InitOperatorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_aggregator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitOperatorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitOperatorResponse) ProtoMessage() {}

func (x *InitOperatorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_aggregator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitOperatorResponse.ProtoReflect.Descriptor instead.
func (*InitOperatorResponse) Descriptor() ([]byte, []int) {
	return file_aggregator_aggregator_proto_rawDescGZIP(), []int{1}
}

func (x *InitOperatorResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *InitOperatorResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type CreateTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The hash of alert
	AlertHash []byte `protobuf:"bytes,1,opt,name=alert_hash,json=alertHash,proto3" json:"alert_hash,omitempty"`
}

func (x *CreateTaskRequest) Reset() {
	*x = CreateTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_aggregator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTaskRequest) ProtoMessage() {}

func (x *CreateTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_aggregator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTaskRequest.ProtoReflect.Descriptor instead.
func (*CreateTaskRequest) Descriptor() ([]byte, []int) {
	return file_aggregator_aggregator_proto_rawDescGZIP(), []int{2}
}

func (x *CreateTaskRequest) GetAlertHash() []byte {
	if x != nil {
		return x.AlertHash
	}
	return nil
}

type CreateTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The info of alert
	Info *AlertTaskInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *CreateTaskResponse) Reset() {
	*x = CreateTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_aggregator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTaskResponse) ProtoMessage() {}

func (x *CreateTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_aggregator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTaskResponse.ProtoReflect.Descriptor instead.
func (*CreateTaskResponse) Descriptor() ([]byte, []int) {
	return file_aggregator_aggregator_proto_rawDescGZIP(), []int{3}
}

func (x *CreateTaskResponse) GetInfo() *AlertTaskInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type SignedTaskRespRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The alert
	Alert *AlertTaskInfo `protobuf:"bytes,1,opt,name=alert,proto3" json:"alert,omitempty"`
	// The operator's BLS signature signed on the keccak256 hash
	OperatorRequestSignature []byte `protobuf:"bytes,2,opt,name=operator_request_signature,json=operatorRequestSignature,proto3" json:"operator_request_signature,omitempty"`
	// The operator 's id
	OperatorId []byte `protobuf:"bytes,3,opt,name=operator_id,json=operatorId,proto3" json:"operator_id,omitempty"`
}

func (x *SignedTaskRespRequest) Reset() {
	*x = SignedTaskRespRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_aggregator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedTaskRespRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedTaskRespRequest) ProtoMessage() {}

func (x *SignedTaskRespRequest) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_aggregator_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedTaskRespRequest.ProtoReflect.Descriptor instead.
func (*SignedTaskRespRequest) Descriptor() ([]byte, []int) {
	return file_aggregator_aggregator_proto_rawDescGZIP(), []int{4}
}

func (x *SignedTaskRespRequest) GetAlert() *AlertTaskInfo {
	if x != nil {
		return x.Alert
	}
	return nil
}

func (x *SignedTaskRespRequest) GetOperatorRequestSignature() []byte {
	if x != nil {
		return x.OperatorRequestSignature
	}
	return nil
}

func (x *SignedTaskRespRequest) GetOperatorId() []byte {
	if x != nil {
		return x.OperatorId
	}
	return nil
}

type SignedTaskRespResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If need reply
	Reply bool `protobuf:"varint,1,opt,name=reply,proto3" json:"reply,omitempty"`
	// The tx hash of send
	TxHash []byte `protobuf:"bytes,2,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
}

func (x *SignedTaskRespResponse) Reset() {
	*x = SignedTaskRespResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_aggregator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignedTaskRespResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedTaskRespResponse) ProtoMessage() {}

func (x *SignedTaskRespResponse) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_aggregator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedTaskRespResponse.ProtoReflect.Descriptor instead.
func (*SignedTaskRespResponse) Descriptor() ([]byte, []int) {
	return file_aggregator_aggregator_proto_rawDescGZIP(), []int{5}
}

func (x *SignedTaskRespResponse) GetReply() bool {
	if x != nil {
		return x.Reply
	}
	return false
}

func (x *SignedTaskRespResponse) GetTxHash() []byte {
	if x != nil {
		return x.TxHash
	}
	return nil
}

type AlertTaskInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The hash of alert
	AlertHash []byte `protobuf:"bytes,1,opt,name=alert_hash,json=alertHash,proto3" json:"alert_hash,omitempty"`
	// QuorumNumbers of task
	QuorumNumbers []byte `protobuf:"bytes,2,opt,name=quorum_numbers,json=quorumNumbers,proto3" json:"quorum_numbers,omitempty"`
	// QuorumThresholdPercentages of task
	QuorumThresholdPercentages []byte `protobuf:"bytes,3,opt,name=quorum_threshold_percentages,json=quorumThresholdPercentages,proto3" json:"quorum_threshold_percentages,omitempty"`
	// TaskIndex
	TaskIndex uint32 `protobuf:"varint,4,opt,name=task_index,json=taskIndex,proto3" json:"task_index,omitempty"`
	// ReferenceBlockNumber
	ReferenceBlockNumber uint64 `protobuf:"varint,5,opt,name=reference_block_number,json=referenceBlockNumber,proto3" json:"reference_block_number,omitempty"`
}

func (x *AlertTaskInfo) Reset() {
	*x = AlertTaskInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_aggregator_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlertTaskInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlertTaskInfo) ProtoMessage() {}

func (x *AlertTaskInfo) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_aggregator_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlertTaskInfo.ProtoReflect.Descriptor instead.
func (*AlertTaskInfo) Descriptor() ([]byte, []int) {
	return file_aggregator_aggregator_proto_rawDescGZIP(), []int{6}
}

func (x *AlertTaskInfo) GetAlertHash() []byte {
	if x != nil {
		return x.AlertHash
	}
	return nil
}

func (x *AlertTaskInfo) GetQuorumNumbers() []byte {
	if x != nil {
		return x.QuorumNumbers
	}
	return nil
}

func (x *AlertTaskInfo) GetQuorumThresholdPercentages() []byte {
	if x != nil {
		return x.QuorumThresholdPercentages
	}
	return nil
}

func (x *AlertTaskInfo) GetTaskIndex() uint32 {
	if x != nil {
		return x.TaskIndex
	}
	return 0
}

func (x *AlertTaskInfo) GetReferenceBlockNumber() uint64 {
	if x != nil {
		return x.ReferenceBlockNumber
	}
	return 0
}

var File_aggregator_aggregator_proto protoreflect.FileDescriptor

var file_aggregator_aggregator_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x61, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x22, 0xa3, 0x02, 0x0a, 0x13, 0x49, 0x6e,
	0x69, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x26, 0x0a, 0x0f, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x31, 0x5f, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x31, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x41, 0x0a, 0x1d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x5f, 0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64,
	0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1a, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x72, 0x41,
	0x64, 0x64, 0x72, 0x12, 0x3a, 0x0a, 0x19, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5f,
	0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x22,
	0x3e, 0x0a, 0x14, 0x49, 0x6e, 0x69, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22,
	0x32, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x48,
	0x61, 0x73, 0x68, 0x22, 0x43, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x04, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0xa7, 0x01, 0x0a, 0x15, 0x53, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2f, 0x0a, 0x05, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x41,
	0x6c, 0x65, 0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x61, 0x6c,
	0x65, 0x72, 0x74, 0x12, 0x3c, 0x0a, 0x1a, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x5f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x18, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x49, 0x64, 0x22, 0x47, 0x0a, 0x16, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x54, 0x61, 0x73, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x72, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x22, 0xec, 0x01, 0x0a, 0x0d,
	0x41, 0x6c, 0x65, 0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1d, 0x0a,
	0x0a, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x25, 0x0a, 0x0e,
	0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x12, 0x40, 0x0a, 0x1c, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x5f, 0x74, 0x68,
	0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61,
	0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x1a, 0x71, 0x75, 0x6f, 0x72, 0x75,
	0x6d, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e,
	0x74, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x74, 0x61, 0x73, 0x6b, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x34, 0x0a, 0x16, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x14, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x32, 0x96, 0x02, 0x0a, 0x0a, 0x41,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x53, 0x0a, 0x0c, 0x49, 0x6e, 0x69,
	0x74, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x1f, 0x2e, 0x61, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4d,
	0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x1d, 0x2e, 0x61,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x64, 0x0a,
	0x19, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x54, 0x61,
	0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x2e, 0x61, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x54, 0x61,
	0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e,
	0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x6c, 0x74, 0x2d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x61,
	0x76, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aggregator_aggregator_proto_rawDescOnce sync.Once
	file_aggregator_aggregator_proto_rawDescData = file_aggregator_aggregator_proto_rawDesc
)

func file_aggregator_aggregator_proto_rawDescGZIP() []byte {
	file_aggregator_aggregator_proto_rawDescOnce.Do(func() {
		file_aggregator_aggregator_proto_rawDescData = protoimpl.X.CompressGZIP(file_aggregator_aggregator_proto_rawDescData)
	})
	return file_aggregator_aggregator_proto_rawDescData
}

var file_aggregator_aggregator_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_aggregator_aggregator_proto_goTypes = []interface{}{
	(*InitOperatorRequest)(nil),    // 0: aggregator.InitOperatorRequest
	(*InitOperatorResponse)(nil),   // 1: aggregator.InitOperatorResponse
	(*CreateTaskRequest)(nil),      // 2: aggregator.CreateTaskRequest
	(*CreateTaskResponse)(nil),     // 3: aggregator.CreateTaskResponse
	(*SignedTaskRespRequest)(nil),  // 4: aggregator.SignedTaskRespRequest
	(*SignedTaskRespResponse)(nil), // 5: aggregator.SignedTaskRespResponse
	(*AlertTaskInfo)(nil),          // 6: aggregator.AlertTaskInfo
}
var file_aggregator_aggregator_proto_depIdxs = []int32{
	6, // 0: aggregator.CreateTaskResponse.info:type_name -> aggregator.AlertTaskInfo
	6, // 1: aggregator.SignedTaskRespRequest.alert:type_name -> aggregator.AlertTaskInfo
	0, // 2: aggregator.Aggregator.InitOperator:input_type -> aggregator.InitOperatorRequest
	2, // 3: aggregator.Aggregator.CreateTask:input_type -> aggregator.CreateTaskRequest
	4, // 4: aggregator.Aggregator.ProcessSignedTaskResponse:input_type -> aggregator.SignedTaskRespRequest
	1, // 5: aggregator.Aggregator.InitOperator:output_type -> aggregator.InitOperatorResponse
	3, // 6: aggregator.Aggregator.CreateTask:output_type -> aggregator.CreateTaskResponse
	5, // 7: aggregator.Aggregator.ProcessSignedTaskResponse:output_type -> aggregator.SignedTaskRespResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_aggregator_aggregator_proto_init() }
func file_aggregator_aggregator_proto_init() {
	if File_aggregator_aggregator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aggregator_aggregator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitOperatorRequest); i {
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
		file_aggregator_aggregator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitOperatorResponse); i {
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
		file_aggregator_aggregator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTaskRequest); i {
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
		file_aggregator_aggregator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTaskResponse); i {
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
		file_aggregator_aggregator_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedTaskRespRequest); i {
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
		file_aggregator_aggregator_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignedTaskRespResponse); i {
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
		file_aggregator_aggregator_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlertTaskInfo); i {
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
			RawDescriptor: file_aggregator_aggregator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_aggregator_aggregator_proto_goTypes,
		DependencyIndexes: file_aggregator_aggregator_proto_depIdxs,
		MessageInfos:      file_aggregator_aggregator_proto_msgTypes,
	}.Build()
	File_aggregator_aggregator_proto = out.File
	file_aggregator_aggregator_proto_rawDesc = nil
	file_aggregator_aggregator_proto_goTypes = nil
	file_aggregator_aggregator_proto_depIdxs = nil
}
