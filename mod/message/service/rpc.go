package service

import (
	messageV1 "github.com/GoldenSheep402/Hermes/pkg/proto/message/v1"

	"context"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/message/dao"
	"github.com/GoldenSheep402/Hermes/mod/message/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
)

var _ messageV1.MessageServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	messageV1.UnimplementedMessageServiceServer
}

func requireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return userID, nil
}

func requireAdmin(ctx context.Context) error {
	userID, err := requireAuth(ctx)
	if err != nil {
		return err
	}
	isAdmin, err := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if err != nil || !isAdmin {
		return status.Error(codes.PermissionDenied, "admin privileges required")
	}
	return nil
}

func (s S) SendMessage(ctx context.Context, req *messageV1.SendMessageRequest) (*messageV1.SendMessageResponse, error) {
	senderID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if req.ReceiverId == "" || req.Subject == "" || req.Body == "" {
		return nil, status.Error(codes.InvalidArgument, "ReceiverId, Subject, and Body are required")
	}

	msg := &model.Message{
		Model:      stdao.Model{ID: ulid.Make().String()},
		SenderID:   senderID,
		ReceiverID: req.ReceiverId,
		Subject:    req.Subject,
		Body:       req.Body,
		IsRead:     false,
		Type:       model.MessageTypePrivate,
	}

	if err := dao.Message.Create(ctx, msg); err != nil {
		s.Log.Errorw("failed to send message", "err", err)
		return nil, status.Error(codes.Internal, "Failed to send message")
	}

	return &messageV1.SendMessageResponse{Id: msg.ID}, nil
}

func (s S) GetMessage(ctx context.Context, req *messageV1.GetMessageRequest) (*messageV1.GetMessageResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Message ID required")
	}

	msg, err := dao.Message.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Message not found")
	}

	// Make sure only the sender or receiver can view the message
	if msg.SenderID != userID && msg.ReceiverID != userID {
		if errAdmin := requireAdmin(ctx); errAdmin != nil {
			return nil, status.Error(codes.PermissionDenied, "Not authorized to view this message")
		}
	}

	// Mark as read if the receiver is opening it for the first time
	if msg.ReceiverID == userID && !msg.IsRead {
		_ = dao.Message.MarkAsRead(ctx, userID, []string{msg.ID})
		msg.IsRead = true
	}

	return &messageV1.GetMessageResponse{
		Message: &messageV1.MessageInfo{
			Id:           msg.ID,
			SenderId:     msg.SenderID,
			SenderName:   "User " + msg.SenderID, // TODO: Interconnect with User module
			ReceiverId:   msg.ReceiverID,
			ReceiverName: "User " + msg.ReceiverID, // TODO: Interconnect with User module
			Subject:      msg.Subject,
			Body:         msg.Body,
			IsRead:       msg.IsRead,
			Type:         int32(msg.Type),
			CreatedAt:    msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (s S) ListInbox(ctx context.Context, req *messageV1.ListInboxRequest) (*messageV1.ListInboxResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	list, count, err := dao.Message.ListInbox(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list inbox", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list inbox")
	}

	var pbMsgs []*messageV1.MessageInfo
	for _, msg := range list {
		pbMsgs = append(pbMsgs, &messageV1.MessageInfo{
			Id:           msg.ID,
			SenderId:     msg.SenderID,
			SenderName:   "User " + msg.SenderID, // Placeholder
			ReceiverId:   msg.ReceiverID,
			ReceiverName: "User " + msg.ReceiverID,
			Subject:      msg.Subject,
			Body:         msg.Body,
			IsRead:       msg.IsRead,
			Type:         int32(msg.Type),
			CreatedAt:    msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &messageV1.ListInboxResponse{Messages: pbMsgs, Total: count}, nil
}

func (s S) ListOutbox(ctx context.Context, req *messageV1.ListOutboxRequest) (*messageV1.ListOutboxResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	list, count, err := dao.Message.ListOutbox(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list outbox", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list outbox")
	}

	var pbMsgs []*messageV1.MessageInfo
	for _, msg := range list {
		pbMsgs = append(pbMsgs, &messageV1.MessageInfo{
			Id:           msg.ID,
			SenderId:     msg.SenderID,
			SenderName:   "User " + msg.SenderID,
			ReceiverId:   msg.ReceiverID,
			ReceiverName: "User " + msg.ReceiverID,
			Subject:      msg.Subject,
			Body:         msg.Body,
			IsRead:       msg.IsRead,
			Type:         int32(msg.Type),
			CreatedAt:    msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &messageV1.ListOutboxResponse{Messages: pbMsgs, Total: count}, nil
}

func (s S) MarkAsRead(ctx context.Context, req *messageV1.MarkAsReadRequest) (*messageV1.MarkAsReadResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "IDs list cannot be empty")
	}

	if err := dao.Message.MarkAsRead(ctx, userID, req.Ids); err != nil {
		s.Log.Errorw("failed to mark messages as read", "err", err)
		return nil, status.Error(codes.Internal, "Failed to mark as read")
	}

	return &messageV1.MarkAsReadResponse{}, nil
}

func (s S) DeleteMessage(ctx context.Context, req *messageV1.DeleteMessageRequest) (*messageV1.DeleteMessageResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Message ID required")
	}

	msg, err := dao.Message.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Message not found")
	}

	if msg.SenderID != userID && msg.ReceiverID != userID {
		if errAdmin := requireAdmin(ctx); errAdmin != nil {
			return nil, status.Error(codes.PermissionDenied, "Not authorized to delete this message")
		}
	}

	if err := dao.Message.Delete(ctx, &model.Message{Model: stdao.Model{ID: req.Id}}).Error; err != nil {
		s.Log.Errorw("failed to delete message", "err", err)
		return nil, status.Error(codes.Internal, "Failed to delete message")
	}

	return &messageV1.DeleteMessageResponse{}, nil
}

func (s S) GetUnreadCount(ctx context.Context, req *messageV1.GetUnreadCountRequest) (*messageV1.GetUnreadCountResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	count, err := dao.Message.GetUnreadCount(ctx, userID)
	if err != nil {
		s.Log.Errorw("failed to get unread count", "err", err)
		return nil, status.Error(codes.Internal, "Failed to fetch unread count")
	}

	return &messageV1.GetUnreadCountResponse{Count: count}, nil
}
