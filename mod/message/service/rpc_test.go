package service

import (
	"context"
	"testing"

	messageV1 "github.com/GoldenSheep402/Hermes/pkg/proto/message/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMessageService_Actions_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	r1, err1 := s.SendMessage(context.Background(), &messageV1.SendMessageRequest{ReceiverId: "123", Subject: "test", Body: "body"})
	assert.Error(t, err1)
	assert.Nil(t, r1)
	assert.Contains(t, err1.Error(), "unauthenticated")

	r2, err2 := s.GetMessage(context.Background(), &messageV1.GetMessageRequest{Id: "123"})
	assert.Error(t, err2)
	assert.Nil(t, r2)
	assert.Contains(t, err2.Error(), "unauthenticated")

	r3, err3 := s.ListInbox(context.Background(), &messageV1.ListInboxRequest{})
	assert.Error(t, err3)
	assert.Nil(t, r3)
	assert.Contains(t, err3.Error(), "unauthenticated")

	r4, err4 := s.ListOutbox(context.Background(), &messageV1.ListOutboxRequest{})
	assert.Error(t, err4)
	assert.Nil(t, r4)
	assert.Contains(t, err4.Error(), "unauthenticated")

	r5, err5 := s.MarkAsRead(context.Background(), &messageV1.MarkAsReadRequest{Ids: []string{"123"}})
	assert.Error(t, err5)
	assert.Nil(t, r5)
	assert.Contains(t, err5.Error(), "unauthenticated")

	r6, err6 := s.DeleteMessage(context.Background(), &messageV1.DeleteMessageRequest{Id: "123"})
	assert.Error(t, err6)
	assert.Nil(t, r6)
	assert.Contains(t, err6.Error(), "unauthenticated")

	r7, err7 := s.GetUnreadCount(context.Background(), &messageV1.GetUnreadCountRequest{})
	assert.Error(t, err7)
	assert.Nil(t, r7)
	assert.Contains(t, err7.Error(), "unauthenticated")
}
