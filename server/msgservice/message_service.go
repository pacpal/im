package msgService

import (
	"IM/server/model"
	"IM/server/msgservice/hub"
	repo "IM/server/repository"
	"fmt"
)

type MessageService struct {
	UserRepo repo.UserRepo
	Hub      *hub.Hub
}

func NewMessageService(ur repo.UserRepo) *MessageService {
	return &MessageService{UserRepo: ur}
}
func (s *MessageService) RouteMessage(msg model.Message) error {
	if msg.RcID == "" || msg.SdID == "" || msg.Content == "" {
		return fmt.Errorf("invalid message")
	}
	switch msg.Type {
	case "private":
		return s.routePrivate(msg)
	case "Group":
		return s.routeGroup(msg)
	default:
	}
	return fmt.Errorf("unknown message type: %s", msg.Type)
}
func (s *MessageService) routePrivate(msg model.Message) error {
	//CheckFriendship
	if client, online := s.Hub.GetOnlineClient(msg.RcID); online {
		select {
		case client.Send <- msg:
		default:
			s.CacheOffline(msg)
		}
		s.CacheOffline(msg)
	}
	return nil
}
func (s *MessageService) routeGroup(msg model.Message) error {
	return nil
}
func (s *MessageService) CacheOffline(msg model.Message) {

}
