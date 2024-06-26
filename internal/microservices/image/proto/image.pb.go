// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: internal/microservices/image/proto/image.proto

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

type GetImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetImageRequest) Reset() {
	*x = GetImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageRequest) ProtoMessage() {}

func (x *GetImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageRequest.ProtoReflect.Descriptor instead.
func (*GetImageRequest) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{0}
}

func (x *GetImageRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetImageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Image      []byte `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	LocalError int64  `protobuf:"varint,2,opt,name=localError,proto3" json:"localError,omitempty"`
}

func (x *GetImageResponse) Reset() {
	*x = GetImageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageResponse) ProtoMessage() {}

func (x *GetImageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageResponse.ProtoReflect.Descriptor instead.
func (*GetImageResponse) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{1}
}

func (x *GetImageResponse) GetImage() []byte {
	if x != nil {
		return x.Image
	}
	return nil
}

func (x *GetImageResponse) GetLocalError() int64 {
	if x != nil {
		return x.LocalError
	}
	return 0
}

type UploadImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Image    []byte `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	Filename string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *UploadImageRequest) Reset() {
	*x = UploadImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadImageRequest) ProtoMessage() {}

func (x *UploadImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadImageRequest.ProtoReflect.Descriptor instead.
func (*UploadImageRequest) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{2}
}

func (x *UploadImageRequest) GetImage() []byte {
	if x != nil {
		return x.Image
	}
	return nil
}

func (x *UploadImageRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type UploadImageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	LocalError int64  `protobuf:"varint,2,opt,name=localError,proto3" json:"localError,omitempty"`
}

func (x *UploadImageResponse) Reset() {
	*x = UploadImageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadImageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadImageResponse) ProtoMessage() {}

func (x *UploadImageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadImageResponse.ProtoReflect.Descriptor instead.
func (*UploadImageResponse) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{3}
}

func (x *UploadImageResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UploadImageResponse) GetLocalError() int64 {
	if x != nil {
		return x.LocalError
	}
	return 0
}

type FormUrlRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *FormUrlRequest) Reset() {
	*x = FormUrlRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FormUrlRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FormUrlRequest) ProtoMessage() {}

func (x *FormUrlRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FormUrlRequest.ProtoReflect.Descriptor instead.
func (*FormUrlRequest) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{4}
}

func (x *FormUrlRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type FormUrlResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *FormUrlResponse) Reset() {
	*x = FormUrlResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FormUrlResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FormUrlResponse) ProtoMessage() {}

func (x *FormUrlResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FormUrlResponse.ProtoReflect.Descriptor instead.
func (*FormUrlResponse) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{5}
}

func (x *FormUrlResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetImageBoundsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *GetImageBoundsRequest) Reset() {
	*x = GetImageBoundsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageBoundsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageBoundsRequest) ProtoMessage() {}

func (x *GetImageBoundsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageBoundsRequest.ProtoReflect.Descriptor instead.
func (*GetImageBoundsRequest) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{6}
}

func (x *GetImageBoundsRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetImageBoundsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dx         int64 `protobuf:"varint,1,opt,name=dx,proto3" json:"dx,omitempty"`
	Dy         int64 `protobuf:"varint,2,opt,name=dy,proto3" json:"dy,omitempty"`
	LocalError int64 `protobuf:"varint,3,opt,name=localError,proto3" json:"localError,omitempty"`
}

func (x *GetImageBoundsResponse) Reset() {
	*x = GetImageBoundsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_microservices_image_proto_image_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageBoundsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageBoundsResponse) ProtoMessage() {}

func (x *GetImageBoundsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_microservices_image_proto_image_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageBoundsResponse.ProtoReflect.Descriptor instead.
func (*GetImageBoundsResponse) Descriptor() ([]byte, []int) {
	return file_internal_microservices_image_proto_image_proto_rawDescGZIP(), []int{7}
}

func (x *GetImageBoundsResponse) GetDx() int64 {
	if x != nil {
		return x.Dx
	}
	return 0
}

func (x *GetImageBoundsResponse) GetDy() int64 {
	if x != nil {
		return x.Dy
	}
	return 0
}

func (x *GetImageBoundsResponse) GetLocalError() int64 {
	if x != nil {
		return x.LocalError
	}
	return 0
}

var File_internal_microservices_image_proto_image_proto protoreflect.FileDescriptor

var file_internal_microservices_image_proto_image_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x69, 0x63, 0x72, 0x6f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x25, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x48, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6c, 0x6f, 0x63,
	0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x46, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x49, 0x0a, 0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x24, 0x0a, 0x0e, 0x46, 0x6f,
	0x72, 0x6d, 0x55, 0x72, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0x23, 0x0a, 0x0f, 0x46, 0x6f, 0x72, 0x6d, 0x55, 0x72, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x29, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c,
	0x22, 0x58, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x42, 0x6f, 0x75, 0x6e,
	0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x64, 0x78,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x64, 0x78, 0x12, 0x0e, 0x0a, 0x02, 0x64, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x64, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x93, 0x02, 0x0a, 0x05, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x12, 0x3b, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x12, 0x15, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47,
	0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x44, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x12, 0x18, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x07, 0x46, 0x6f, 0x72, 0x6d, 0x55,
	0x72, 0x6c, 0x12, 0x14, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x46, 0x6f, 0x72, 0x6d, 0x55, 0x72,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x46, 0x6f, 0x72, 0x6d, 0x55, 0x72, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x4d, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x42, 0x6f, 0x75,
	0x6e, 0x64, 0x73, 0x12, 0x1b, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x42, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x1b, 0x5a, 0x19, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_microservices_image_proto_image_proto_rawDescOnce sync.Once
	file_internal_microservices_image_proto_image_proto_rawDescData = file_internal_microservices_image_proto_image_proto_rawDesc
)

func file_internal_microservices_image_proto_image_proto_rawDescGZIP() []byte {
	file_internal_microservices_image_proto_image_proto_rawDescOnce.Do(func() {
		file_internal_microservices_image_proto_image_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_microservices_image_proto_image_proto_rawDescData)
	})
	return file_internal_microservices_image_proto_image_proto_rawDescData
}

var file_internal_microservices_image_proto_image_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_internal_microservices_image_proto_image_proto_goTypes = []interface{}{
	(*GetImageRequest)(nil),        // 0: auth.GetImageRequest
	(*GetImageResponse)(nil),       // 1: auth.GetImageResponse
	(*UploadImageRequest)(nil),     // 2: auth.UploadImageRequest
	(*UploadImageResponse)(nil),    // 3: auth.UploadImageResponse
	(*FormUrlRequest)(nil),         // 4: auth.FormUrlRequest
	(*FormUrlResponse)(nil),        // 5: auth.FormUrlResponse
	(*GetImageBoundsRequest)(nil),  // 6: auth.GetImageBoundsRequest
	(*GetImageBoundsResponse)(nil), // 7: auth.GetImageBoundsResponse
}
var file_internal_microservices_image_proto_image_proto_depIdxs = []int32{
	0, // 0: auth.Image.GetImage:input_type -> auth.GetImageRequest
	2, // 1: auth.Image.UploadImage:input_type -> auth.UploadImageRequest
	4, // 2: auth.Image.FormUrl:input_type -> auth.FormUrlRequest
	6, // 3: auth.Image.GetImageBounds:input_type -> auth.GetImageBoundsRequest
	1, // 4: auth.Image.GetImage:output_type -> auth.GetImageResponse
	3, // 5: auth.Image.UploadImage:output_type -> auth.UploadImageResponse
	5, // 6: auth.Image.FormUrl:output_type -> auth.FormUrlResponse
	7, // 7: auth.Image.GetImageBounds:output_type -> auth.GetImageBoundsResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_microservices_image_proto_image_proto_init() }
func file_internal_microservices_image_proto_image_proto_init() {
	if File_internal_microservices_image_proto_image_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_microservices_image_proto_image_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageRequest); i {
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
		file_internal_microservices_image_proto_image_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageResponse); i {
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
		file_internal_microservices_image_proto_image_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadImageRequest); i {
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
		file_internal_microservices_image_proto_image_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadImageResponse); i {
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
		file_internal_microservices_image_proto_image_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FormUrlRequest); i {
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
		file_internal_microservices_image_proto_image_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FormUrlResponse); i {
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
		file_internal_microservices_image_proto_image_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageBoundsRequest); i {
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
		file_internal_microservices_image_proto_image_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageBoundsResponse); i {
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
			RawDescriptor: file_internal_microservices_image_proto_image_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_microservices_image_proto_image_proto_goTypes,
		DependencyIndexes: file_internal_microservices_image_proto_image_proto_depIdxs,
		MessageInfos:      file_internal_microservices_image_proto_image_proto_msgTypes,
	}.Build()
	File_internal_microservices_image_proto_image_proto = out.File
	file_internal_microservices_image_proto_image_proto_rawDesc = nil
	file_internal_microservices_image_proto_image_proto_goTypes = nil
	file_internal_microservices_image_proto_image_proto_depIdxs = nil
}
