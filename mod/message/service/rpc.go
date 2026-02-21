package service

import (
	messageV1 "github.com/GoldenSheep402/Hermes/pkg/proto/message/v1"

	"context"

	"go.uber.org/zap"
)

var _ messageV1.MessageServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	messageV1.UnimplementedMessageServiceServer
}

func (s S) SendMessage(ctx context.Context, request *messageV1.SendMessageRequest) (*messageV1.SendMessageResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) GetMessage(ctx context.Context, request *messageV1.GetMessageRequest) (*messageV1.GetMessageResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListInbox(ctx context.Context, request *messageV1.ListInboxRequest) (*messageV1.ListInboxResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListOutbox(ctx context.Context, request *messageV1.ListOutboxRequest) (*messageV1.ListOutboxResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) MarkAsRead(ctx context.Context, request *messageV1.MarkAsReadRequest) (*messageV1.MarkAsReadResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) DeleteMessage(ctx context.Context, request *messageV1.DeleteMessageRequest) (*messageV1.DeleteMessageResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) GetUnreadCount(ctx context.Context, request *messageV1.GetUnreadCountRequest) (*messageV1.GetUnreadCountResponse, error) {
	// TODO implement me
	panic("implement me")
}
