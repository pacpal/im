package msgService

import (
	"IM/server/model"
	"IM/server/msgservice/hub"
	repo "IM/server/repository"
	"context"
	"fmt"
)

type MessageService struct {
	MsgRepo repo.MsgRepo
	Hub     *hub.Hub
}

func NewMessageService(ur repo.MsgRepo) *MessageService {
	return &MessageService{MsgRepo: ur}
}
func (s *MessageService) RouteMessage(ctx context.Context, msg model.Message) error {
	if msg.RcID == "" || msg.SdID == "" || msg.Content == "" {
		return fmt.Errorf("invalid message")
	}
	switch msg.Type {
	case "private":
		return s.routePrivate(ctx, msg)
	case "Group":
		return s.routeGroup(ctx, msg)
	default:
	}
	return fmt.Errorf("unknown message type: %s", msg.Type)
}
func (s *MessageService) routePrivate(ctx context.Context, msg model.Message) error {
	//CheckFriendship
	if client, online := s.Hub.GetOnlineClient(msg.RcID); online {
		select {
		case client.Send <- msg:
		default:
			s.CacheOffline(ctx, msg)
		}
		s.CacheOffline(ctx, msg)
	}
	return nil
}
func (s *MessageService) routeGroup(ctx context.Context, msg model.Message) error {
	//CheckGroup
	members, _ := s.Hub.GetGroupMembers(ctx, msg.RcID)
	for member := range members {
		if client, online := s.Hub.GetOnlineClient(member); online {
			select {
			case client.Send <- msg:
			default:
				s.CacheOffline(ctx, msg)
			}
			s.CacheOffline(ctx, msg)
		}
	}
	return nil
}
func (s *MessageService) CacheOffline(ctx context.Context, msg model.Message) {
}
func (s *MessageService) GetOfflineMsgs(ctx context.Context, uid string) (*[]model.Message, error) {
	msgs, err := s.MsgRepo.GetOfflineMsgs(ctx, uid)
	if err != nil {
		return nil, err
	}
	s.MsgRepo.ClearOfflineMsgs(ctx, uid)

	return msgs, nil
}
func (s *MessageService) GetOnlineStatus(ctx context.Context, uid string) (*[]string, error) {

}
