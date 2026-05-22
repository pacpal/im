package grpc

import (
	"IM/api/gen/common"
	"IM/api/gen/message"
	"IM/services/message/application/service"
	"IM/services/message/domain/entity"
	"context"
)

type MessageServer struct {
	message.UnimplementedMessageServiceServer
	messageSvc *service.MessageService
}

func NewMessageServer(messageSvc *service.MessageService) *MessageServer {
	return &MessageServer{
		messageSvc: messageSvc,
	}
}

func (s *MessageServer) SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	msg, err := s.messageSvc.SendMessage(ctx, req.SenderId, req.ReceiverId, req.Content, req.MsgType)
	if err != nil {
		return nil, err
	}

	return &message.SendMessageResponse{
		MsgId:   msg.ID,
		Success: true,
	}, nil
}

func (s *MessageServer) GetOfflineMessages(ctx context.Context, req *message.GetOfflineMessagesRequest) (*message.GetOfflineMessagesResponse, error) {
	messages, err := s.messageSvc.GetOfflineMessages(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	pbMessages := make([]*message.MessageInfo, len(messages))
	for i, msg := range messages {
		pbMessages[i] = toMessageInfo(msg)
	}

	return &message.GetOfflineMessagesResponse{
		Messages: pbMessages,
	}, nil
}

func (s *MessageServer) MarkAsRead(ctx context.Context, req *message.MarkAsReadRequest) (*common.Response, error) {
	err := s.messageSvc.MarkAsRead(ctx, req.MsgId, "")
	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "marked as read",
	}, nil
}

func (s *MessageServer) MarkAllAsRead(ctx context.Context, req *message.MarkAllAsReadRequest) (*common.Response, error) {
	err := s.messageSvc.MarkAllAsRead(ctx, req.UserId, "")
	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "all messages marked as read",
	}, nil
}

func (s *MessageServer) GetUnreadCount(ctx context.Context, req *message.GetUnreadCountRequest) (*message.GetUnreadCountResponse, error) {
	count, err := s.messageSvc.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &message.GetUnreadCountResponse{
		Count: count,
	}, nil
}

func (s *MessageServer) GetOnlineStatus(ctx context.Context, req *message.GetOnlineStatusRequest) (*message.GetOnlineStatusResponse, error) {
	statusMap, err := s.messageSvc.GetOnlineStatus(ctx, req.GetUserIds())
	if err != nil {
		return nil, err
	}

	return &message.GetOnlineStatusResponse{
		OnlineStatus: statusMap,
	}, nil
}

func toMessageInfo(msg *entity.Message) *message.MessageInfo {
	return &message.MessageInfo{
		MsgId:      msg.ID,
		SenderId:   msg.SenderID,
		ReceiverId: msg.ReceiverID,
		Content:    msg.Content,
		MsgType:    string(msg.MsgType),
		Timestamp:  msg.Timestamp,
		IsRead:     msg.IsRead,
	}
}
