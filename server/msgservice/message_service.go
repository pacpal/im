package msgService

import (
	"IM/server/model"
	"IM/server/msgservice/hub"
	repo "IM/server/repository"
	"context"
	"fmt"
	"time"
)

type MessageService struct {
	MsgRepo     repo.MessageRepo
	Friendship  repo.FriendshipRepo
	GroupMember repo.GroupMemberRepo
	Hub         *hub.Hub
}

func NewMessageService(msgRepo repo.MessageRepo, friendship repo.FriendshipRepo, groupMember repo.GroupMemberRepo) *MessageService {
	return &MessageService{
		MsgRepo:     msgRepo,
		Friendship:  friendship,
		GroupMember: groupMember,
	}
}

func (s *MessageService) RouteMessage(ctx context.Context, msg model.Message) error {
	if msg.RcID == "" || msg.SdID == "" || msg.Content == "" {
		return fmt.Errorf("invalid message: missing required fields")
	}

	if msg.Time == 0 {
		msg.Time = time.Now().UnixMilli()
	}

	switch msg.Type {
	case "private":
		return s.routePrivate(ctx, msg)
	case "group":
		return s.routeGroup(ctx, msg)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func (s *MessageService) routePrivate(ctx context.Context, msg model.Message) error {
	isFriend, err := s.Friendship.Exists(ctx, msg.SdID, msg.RcID)
	if err != nil {
		return fmt.Errorf("check friendship failed: %w", err)
	}
	if !isFriend {
		return fmt.Errorf("cannot send message: not friends")
	}

	if client, online := s.Hub.GetOnlineClient(msg.RcID); online {
		select {
		case client.Send <- msg:
			msg.IsRead = true
			return s.MsgRepo.Create(ctx, &msg)
		default:
			return s.CacheOffline(ctx, msg)
		}
	}

	return s.CacheOffline(ctx, msg)
}

func (s *MessageService) routeGroup(ctx context.Context, msg model.Message) error {
	isMember, err := s.GroupMember.IsMember(ctx, msg.RcID, msg.SdID)
	if err != nil {
		return fmt.Errorf("check group membership failed: %w", err)
	}
	if !isMember {
		return fmt.Errorf("sender is not a member of the group")
	}

	memberIDs, err := s.GroupMember.GetMemberIDs(ctx, msg.RcID)
	if err != nil {
		return fmt.Errorf("get group members failed: %w", err)
	}

	var deliveryErrors []error
	for _, memberID := range memberIDs {
		if memberID == msg.SdID {
			continue
		}

		if client, online := s.Hub.GetOnlineClient(memberID); online {
			select {
			case client.Send <- msg:
			default:
				if err := s.cacheOfflineForMember(ctx, msg, memberID); err != nil {
					deliveryErrors = append(deliveryErrors, err)
				}
			}
		} else {
			if err := s.cacheOfflineForMember(ctx, msg, memberID); err != nil {
				deliveryErrors = append(deliveryErrors, err)
			}
		}
	}

	if len(deliveryErrors) > 0 {
		return fmt.Errorf("some messages failed to deliver: %d errors", len(deliveryErrors))
	}

	return nil
}

func (s *MessageService) cacheOfflineForMember(ctx context.Context, msg model.Message, memberID string) error {
	memberMsg := model.Message{
		ID:      fmt.Sprintf("%s_%s", msg.ID, memberID),
		SdID:    msg.SdID,
		RcID:    memberID,
		Content: msg.Content,
		Type:    "group",
		Time:    msg.Time,
		IsRead:  false,
	}
	return s.MsgRepo.Create(ctx, &memberMsg)
}

func (s *MessageService) CacheOffline(ctx context.Context, msg model.Message) error {
	msg.IsRead = false
	return s.MsgRepo.Create(ctx, &msg)
}

func (s *MessageService) GetOfflineMsgs(ctx context.Context, uid string) (*[]model.Message, error) {
	msgs, err := s.MsgRepo.GetOfflineMessages(ctx, uid, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("get offline messages failed: %w", err)
	}

	result := make([]model.Message, len(msgs))
	for i, m := range msgs {
		result[i] = *m
	}

	if len(msgs) > 0 {
		if err := s.MsgRepo.MarkAllAsRead(ctx, uid); err != nil {
			return &result, fmt.Errorf("mark messages as read failed: %w", err)
		}
	}

	return &result, nil
}

func (s *MessageService) GetOnlineStatus(ctx context.Context, uid string) (*[]string, error) {
	if s.Hub == nil {
		empty := make([]string, 0)
		return &empty, nil
	}

	friendIDs, err := s.Friendship.GetFriendIDs(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get friend ids failed: %w", err)
	}

	onlineFriends := make([]string, 0)
	for _, friendID := range friendIDs {
		if _, online := s.Hub.GetOnlineClient(friendID); online {
			onlineFriends = append(onlineFriends, friendID)
		}
	}

	return &onlineFriends, nil
}
