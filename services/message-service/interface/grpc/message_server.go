package grpc

import (
	"IM/services/message-service/application/service"
	pb "IM/api/gen/message"
	"context"
)

type MessageGrpcServer struct {
	pb.UnimplementedMessageServiceServer
	msgSvc *service.MessageApplicationService
}

func NewMessageGrpcServer(msgSvc *service.MessageApplicationService) *MessageGrpcServer {
	return &MessageGrpcServer{
		msgSvc: msgSvc,
	}
}

func (s *MessageGrpcServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	msg, err := s.msgSvc.SendMessage(ctx, req.SenderId, req.ReceiverId, req.Content, req.Type)
	if err != nil {
		return &pb.SendMessageResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SendMessageResponse{
		Success: true,
		Message: "message sent",
		MessageInfo: &pb.Message{
			Id:         msg.ID,
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			Type:       msg.Type,
			Timestamp:  msg.Timestamp,
		},
	}, nil
}

func (s *MessageGrpcServer) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	msg, err := s.msgSvc.GetMessage(ctx, req.MessageId)
	if err != nil {
		return &pb.GetMessageResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetMessageResponse{
		Success: true,
		Message: &pb.Message{
			Id:         msg.ID,
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			Type:       msg.Type,
			Timestamp:  msg.Timestamp,
		},
	}, nil
}

func (s *MessageGrpcServer) GetOfflineMessages(ctx context.Context, req *pb.GetOfflineMessagesRequest) (*pb.GetOfflineMessagesResponse, error) {
	msgs, err := s.msgSvc.GetOfflineMessages(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return &pb.GetOfflineMessagesResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbMessages := make([]*pb.Message, len(msgs))
	for i, m := range msgs {
		pbMessages[i] = &pb.Message{
			Id:         m.ID,
			SenderId:   m.SenderID,
			ReceiverId: m.ReceiverID,
			Content:    m.Content,
			Type:       m.Type,
			Timestamp:  m.Timestamp,
		}
	}

	return &pb.GetOfflineMessagesResponse{
		Success:  true,
		Messages: pbMessages,
	}, nil
}

func (s *MessageGrpcServer) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	err := s.msgSvc.MarkAsRead(ctx, req.MessageId)
	if err != nil {
		return &pb.MarkAsReadResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.MarkAsReadResponse{
		Success: true,
		Message: "message marked as read",
	}, nil
}

func (s *MessageGrpcServer) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	count, err := s.msgSvc.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		return &pb.GetUnreadCountResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetUnreadCountResponse{
		Success:     true,
		UnreadCount: count,
	}, nil
}