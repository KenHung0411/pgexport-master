// Code generated by protoc-gen-go. DO NOT EDIT.
// source: notifications/notifications.proto

package notifications

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Target struct {
	PlatformId           uint32   `protobuf:"varint,1,opt,name=platform_id,json=platformId,proto3" json:"platform_id,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Type                 string   `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Target) Reset()         { *m = Target{} }
func (m *Target) String() string { return proto.CompactTextString(m) }
func (*Target) ProtoMessage()    {}
func (*Target) Descriptor() ([]byte, []int) {
	return fileDescriptor_notifications_aeb75df5379d4b15, []int{0}
}
func (m *Target) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Target.Unmarshal(m, b)
}
func (m *Target) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Target.Marshal(b, m, deterministic)
}
func (dst *Target) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Target.Merge(dst, src)
}
func (m *Target) XXX_Size() int {
	return xxx_messageInfo_Target.Size(m)
}
func (m *Target) XXX_DiscardUnknown() {
	xxx_messageInfo_Target.DiscardUnknown(m)
}

var xxx_messageInfo_Target proto.InternalMessageInfo

func (m *Target) GetPlatformId() uint32 {
	if m != nil {
		return m.PlatformId
	}
	return 0
}

func (m *Target) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Target) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Target) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Sender struct {
	PlatformId           uint32   `protobuf:"varint,1,opt,name=platform_id,json=platformId,proto3" json:"platform_id,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Type                 string   `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Sender) Reset()         { *m = Sender{} }
func (m *Sender) String() string { return proto.CompactTextString(m) }
func (*Sender) ProtoMessage()    {}
func (*Sender) Descriptor() ([]byte, []int) {
	return fileDescriptor_notifications_aeb75df5379d4b15, []int{1}
}
func (m *Sender) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Sender.Unmarshal(m, b)
}
func (m *Sender) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Sender.Marshal(b, m, deterministic)
}
func (dst *Sender) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Sender.Merge(dst, src)
}
func (m *Sender) XXX_Size() int {
	return xxx_messageInfo_Sender.Size(m)
}
func (m *Sender) XXX_DiscardUnknown() {
	xxx_messageInfo_Sender.DiscardUnknown(m)
}

var xxx_messageInfo_Sender proto.InternalMessageInfo

func (m *Sender) GetPlatformId() uint32 {
	if m != nil {
		return m.PlatformId
	}
	return 0
}

func (m *Sender) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Sender) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Sender) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Notification struct {
	Uuid                 string               `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Subject              string               `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Target               *Target              `protobuf:"bytes,3,opt,name=target,proto3" json:"target,omitempty"`
	Sender               *Sender              `protobuf:"bytes,4,opt,name=sender,proto3" json:"sender,omitempty"`
	Source               string               `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`
	Action               string               `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`
	Notes                map[string]string    `protobuf:"bytes,8,rep,name=notes,proto3" json:"notes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Notification) Reset()         { *m = Notification{} }
func (m *Notification) String() string { return proto.CompactTextString(m) }
func (*Notification) ProtoMessage()    {}
func (*Notification) Descriptor() ([]byte, []int) {
	return fileDescriptor_notifications_aeb75df5379d4b15, []int{2}
}
func (m *Notification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Notification.Unmarshal(m, b)
}
func (m *Notification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Notification.Marshal(b, m, deterministic)
}
func (dst *Notification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Notification.Merge(dst, src)
}
func (m *Notification) XXX_Size() int {
	return xxx_messageInfo_Notification.Size(m)
}
func (m *Notification) XXX_DiscardUnknown() {
	xxx_messageInfo_Notification.DiscardUnknown(m)
}

var xxx_messageInfo_Notification proto.InternalMessageInfo

func (m *Notification) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Notification) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Notification) GetTarget() *Target {
	if m != nil {
		return m.Target
	}
	return nil
}

func (m *Notification) GetSender() *Sender {
	if m != nil {
		return m.Sender
	}
	return nil
}

func (m *Notification) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Notification) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *Notification) GetNotes() map[string]string {
	if m != nil {
		return m.Notes
	}
	return nil
}

func (m *Notification) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

type AddNotificationRequest struct {
	Notification         *Notification `protobuf:"bytes,1,opt,name=notification,proto3" json:"notification,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *AddNotificationRequest) Reset()         { *m = AddNotificationRequest{} }
func (m *AddNotificationRequest) String() string { return proto.CompactTextString(m) }
func (*AddNotificationRequest) ProtoMessage()    {}
func (*AddNotificationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_notifications_aeb75df5379d4b15, []int{3}
}
func (m *AddNotificationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddNotificationRequest.Unmarshal(m, b)
}
func (m *AddNotificationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddNotificationRequest.Marshal(b, m, deterministic)
}
func (dst *AddNotificationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddNotificationRequest.Merge(dst, src)
}
func (m *AddNotificationRequest) XXX_Size() int {
	return xxx_messageInfo_AddNotificationRequest.Size(m)
}
func (m *AddNotificationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddNotificationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddNotificationRequest proto.InternalMessageInfo

func (m *AddNotificationRequest) GetNotification() *Notification {
	if m != nil {
		return m.Notification
	}
	return nil
}

type AddNotificationResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddNotificationResponse) Reset()         { *m = AddNotificationResponse{} }
func (m *AddNotificationResponse) String() string { return proto.CompactTextString(m) }
func (*AddNotificationResponse) ProtoMessage()    {}
func (*AddNotificationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_notifications_aeb75df5379d4b15, []int{4}
}
func (m *AddNotificationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddNotificationResponse.Unmarshal(m, b)
}
func (m *AddNotificationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddNotificationResponse.Marshal(b, m, deterministic)
}
func (dst *AddNotificationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddNotificationResponse.Merge(dst, src)
}
func (m *AddNotificationResponse) XXX_Size() int {
	return xxx_messageInfo_AddNotificationResponse.Size(m)
}
func (m *AddNotificationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddNotificationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddNotificationResponse proto.InternalMessageInfo

type QueryNotificationsRequest struct {
	LastSync             int64    `protobuf:"varint,1,opt,name=last_sync,json=lastSync,proto3" json:"last_sync,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryNotificationsRequest) Reset()         { *m = QueryNotificationsRequest{} }
func (m *QueryNotificationsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryNotificationsRequest) ProtoMessage()    {}
func (*QueryNotificationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_notifications_aeb75df5379d4b15, []int{5}
}
func (m *QueryNotificationsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNotificationsRequest.Unmarshal(m, b)
}
func (m *QueryNotificationsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNotificationsRequest.Marshal(b, m, deterministic)
}
func (dst *QueryNotificationsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNotificationsRequest.Merge(dst, src)
}
func (m *QueryNotificationsRequest) XXX_Size() int {
	return xxx_messageInfo_QueryNotificationsRequest.Size(m)
}
func (m *QueryNotificationsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryNotificationsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryNotificationsRequest proto.InternalMessageInfo

func (m *QueryNotificationsRequest) GetLastSync() int64 {
	if m != nil {
		return m.LastSync
	}
	return 0
}

type QueryNotificationsResponse struct {
	LastSync             int64           `protobuf:"varint,1,opt,name=last_sync,json=lastSync,proto3" json:"last_sync,omitempty"`
	Notifications        []*Notification `protobuf:"bytes,2,rep,name=notifications,proto3" json:"notifications,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *QueryNotificationsResponse) Reset()         { *m = QueryNotificationsResponse{} }
func (m *QueryNotificationsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryNotificationsResponse) ProtoMessage()    {}
func (*QueryNotificationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_notifications_aeb75df5379d4b15, []int{6}
}
func (m *QueryNotificationsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryNotificationsResponse.Unmarshal(m, b)
}
func (m *QueryNotificationsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryNotificationsResponse.Marshal(b, m, deterministic)
}
func (dst *QueryNotificationsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryNotificationsResponse.Merge(dst, src)
}
func (m *QueryNotificationsResponse) XXX_Size() int {
	return xxx_messageInfo_QueryNotificationsResponse.Size(m)
}
func (m *QueryNotificationsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryNotificationsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryNotificationsResponse proto.InternalMessageInfo

func (m *QueryNotificationsResponse) GetLastSync() int64 {
	if m != nil {
		return m.LastSync
	}
	return 0
}

func (m *QueryNotificationsResponse) GetNotifications() []*Notification {
	if m != nil {
		return m.Notifications
	}
	return nil
}

func init() {
	proto.RegisterType((*Target)(nil), "notifications.Target")
	proto.RegisterType((*Sender)(nil), "notifications.Sender")
	proto.RegisterType((*Notification)(nil), "notifications.Notification")
	proto.RegisterMapType((map[string]string)(nil), "notifications.Notification.NotesEntry")
	proto.RegisterType((*AddNotificationRequest)(nil), "notifications.AddNotificationRequest")
	proto.RegisterType((*AddNotificationResponse)(nil), "notifications.AddNotificationResponse")
	proto.RegisterType((*QueryNotificationsRequest)(nil), "notifications.QueryNotificationsRequest")
	proto.RegisterType((*QueryNotificationsResponse)(nil), "notifications.QueryNotificationsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NotificationHubClient is the client API for NotificationHub service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NotificationHubClient interface {
	AddNotification(ctx context.Context, in *AddNotificationRequest, opts ...grpc.CallOption) (*AddNotificationResponse, error)
	QueryNotifications(ctx context.Context, in *QueryNotificationsRequest, opts ...grpc.CallOption) (*QueryNotificationsResponse, error)
}

type notificationHubClient struct {
	cc *grpc.ClientConn
}

func NewNotificationHubClient(cc *grpc.ClientConn) NotificationHubClient {
	return &notificationHubClient{cc}
}

func (c *notificationHubClient) AddNotification(ctx context.Context, in *AddNotificationRequest, opts ...grpc.CallOption) (*AddNotificationResponse, error) {
	out := new(AddNotificationResponse)
	err := c.cc.Invoke(ctx, "/notifications.NotificationHub/add_notification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationHubClient) QueryNotifications(ctx context.Context, in *QueryNotificationsRequest, opts ...grpc.CallOption) (*QueryNotificationsResponse, error) {
	out := new(QueryNotificationsResponse)
	err := c.cc.Invoke(ctx, "/notifications.NotificationHub/query_notifications", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationHubServer is the server API for NotificationHub service.
type NotificationHubServer interface {
	AddNotification(context.Context, *AddNotificationRequest) (*AddNotificationResponse, error)
	QueryNotifications(context.Context, *QueryNotificationsRequest) (*QueryNotificationsResponse, error)
}

func RegisterNotificationHubServer(s *grpc.Server, srv NotificationHubServer) {
	s.RegisterService(&_NotificationHub_serviceDesc, srv)
}

func _NotificationHub_AddNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationHubServer).AddNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notifications.NotificationHub/AddNotification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationHubServer).AddNotification(ctx, req.(*AddNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationHub_QueryNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryNotificationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationHubServer).QueryNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notifications.NotificationHub/QueryNotifications",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationHubServer).QueryNotifications(ctx, req.(*QueryNotificationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NotificationHub_serviceDesc = grpc.ServiceDesc{
	ServiceName: "notifications.NotificationHub",
	HandlerType: (*NotificationHubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "add_notification",
			Handler:    _NotificationHub_AddNotification_Handler,
		},
		{
			MethodName: "query_notifications",
			Handler:    _NotificationHub_QueryNotifications_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notifications/notifications.proto",
}

func init() {
	proto.RegisterFile("notifications/notifications.proto", fileDescriptor_notifications_aeb75df5379d4b15)
}

var fileDescriptor_notifications_aeb75df5379d4b15 = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0x56, 0xda, 0x35, 0xac, 0xaf, 0x9b, 0x98, 0x0c, 0x0c, 0x2f, 0x3b, 0xac, 0x44, 0x62, 0x2a,
	0x07, 0x52, 0xa9, 0x5c, 0x0a, 0x42, 0x42, 0x3d, 0x20, 0xc1, 0x05, 0x89, 0x6c, 0x17, 0x4e, 0x95,
	0x9b, 0xbc, 0x56, 0x19, 0xad, 0x9d, 0xc5, 0x36, 0x52, 0x24, 0x7e, 0x2c, 0x7f, 0x04, 0x09, 0xd9,
	0x4e, 0x20, 0xe9, 0x4a, 0xe1, 0xb2, 0x9b, 0xdf, 0xe7, 0xef, 0x7d, 0x9f, 0xdf, 0x97, 0x17, 0x78,
	0xc6, 0x85, 0xca, 0x96, 0x59, 0xc2, 0x54, 0x26, 0xb8, 0x1c, 0xb7, 0xaa, 0x28, 0x2f, 0x84, 0x12,
	0xe4, 0xb8, 0x05, 0x06, 0x17, 0x2b, 0x21, 0x56, 0x6b, 0x1c, 0xdb, 0xcb, 0x85, 0x5e, 0x8e, 0x55,
	0xb6, 0x41, 0xa9, 0xd8, 0x26, 0x77, 0xfc, 0x10, 0xc1, 0xbf, 0x66, 0xc5, 0x0a, 0x15, 0xb9, 0x80,
	0x41, 0xbe, 0x66, 0x6a, 0x29, 0x8a, 0xcd, 0x3c, 0x4b, 0xa9, 0x37, 0xf4, 0x46, 0xc7, 0x31, 0xd4,
	0xd0, 0xc7, 0x94, 0x10, 0x38, 0xd0, 0x3a, 0x4b, 0x69, 0x67, 0xe8, 0x8d, 0xfa, 0xb1, 0x3d, 0x1b,
	0x4c, 0x95, 0x39, 0xd2, 0xae, 0xc3, 0xcc, 0xd9, 0x60, 0x9c, 0x6d, 0x90, 0x1e, 0x38, 0xcc, 0x9c,
	0x8d, 0xcd, 0x15, 0xf2, 0x14, 0x8b, 0xfb, 0xb5, 0xf9, 0xd9, 0x81, 0xa3, 0x4f, 0x8d, 0x00, 0x7e,
	0x8b, 0x79, 0x0d, 0x31, 0x0a, 0x0f, 0xa4, 0x5e, 0xdc, 0x60, 0xa2, 0x2a, 0x8f, 0xba, 0x24, 0x2f,
	0xc1, 0x57, 0x36, 0x0c, 0x6b, 0x34, 0x98, 0x3c, 0x89, 0xda, 0x11, 0xbb, 0xa4, 0xe2, 0x8a, 0x64,
	0xe8, 0xd2, 0x0e, 0x65, 0xdf, 0x70, 0x97, 0xee, 0x26, 0x8e, 0x2b, 0x12, 0x39, 0x05, 0x5f, 0x0a,
	0x5d, 0x24, 0x48, 0x7b, 0xd6, 0xb6, 0xaa, 0x0c, 0xce, 0x12, 0xd3, 0x41, 0x7d, 0x87, 0xbb, 0x8a,
	0xbc, 0x85, 0x1e, 0x17, 0x0a, 0x25, 0x3d, 0x1c, 0x76, 0x47, 0x83, 0xc9, 0xe5, 0x96, 0x7a, 0x73,
	0x4e, 0x53, 0xa0, 0x7c, 0xcf, 0x55, 0x51, 0xc6, 0xae, 0x89, 0xbc, 0x06, 0x48, 0x0a, 0x64, 0x0a,
	0xd3, 0x39, 0x53, 0x14, 0xec, 0x03, 0x83, 0xc8, 0xad, 0x43, 0x54, 0xaf, 0x43, 0x74, 0x5d, 0xaf,
	0x43, 0xdc, 0xaf, 0xd8, 0x33, 0x15, 0x4c, 0x01, 0xfe, 0xe8, 0x91, 0x13, 0xe8, 0x7e, 0xc5, 0xb2,
	0x4a, 0xd0, 0x1c, 0xc9, 0x63, 0xe8, 0x7d, 0x63, 0x6b, 0x8d, 0x55, 0x7c, 0xae, 0x78, 0xd3, 0x99,
	0x7a, 0xe1, 0x17, 0x38, 0x9d, 0xa5, 0x69, 0xf3, 0x65, 0x31, 0xde, 0x6a, 0x94, 0x8a, 0xbc, 0x83,
	0xa3, 0xe6, 0xf3, 0xad, 0xdc, 0x60, 0x72, 0xbe, 0x67, 0xa6, 0xb8, 0xd5, 0x10, 0x9e, 0xc1, 0xd3,
	0x3b, 0xd2, 0x32, 0x17, 0x5c, 0x62, 0x38, 0x85, 0xb3, 0xcf, 0x1a, 0x8b, 0xb2, 0x79, 0x29, 0x6b,
	0xe3, 0x73, 0xe8, 0xaf, 0x99, 0x54, 0x73, 0x59, 0xf2, 0xc4, 0xba, 0x76, 0xe3, 0x43, 0x03, 0x5c,
	0x95, 0x3c, 0x09, 0xbf, 0x43, 0xb0, 0xab, 0xd3, 0xe9, 0xee, 0x6d, 0x25, 0x33, 0x68, 0xff, 0x6a,
	0xb4, 0x63, 0xbf, 0xd2, 0xde, 0x89, 0xda, 0x1d, 0x93, 0x1f, 0x1e, 0x3c, 0x6c, 0xde, 0x7f, 0xd0,
	0x0b, 0xc2, 0xe0, 0x84, 0xa5, 0xe9, 0xbc, 0x49, 0x24, 0xcf, 0xb7, 0x34, 0x77, 0x47, 0x1c, 0x5c,
	0xfe, 0x8b, 0x56, 0x8d, 0x75, 0x03, 0x8f, 0x6e, 0xcd, 0xd0, 0x2d, 0x13, 0x49, 0x46, 0x5b, 0xed,
	0x7f, 0x8d, 0x34, 0x78, 0xf1, 0x1f, 0x4c, 0xe7, 0xb5, 0xf0, 0xed, 0xa6, 0xbd, 0xfa, 0x15, 0x00,
	0x00, 0xff, 0xff, 0xb0, 0xa0, 0xdf, 0x0d, 0xba, 0x04, 0x00, 0x00,
}
