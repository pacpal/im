package ws

import (
	"IM/pkg/event"
	"IM/pkg/logger"
	"encoding/json"
)

type WSEventType string

const (
	WSEventMessageSent            WSEventType = "message.sent"
	WSEventMessageRead            WSEventType = "message.read"
	WSEventMessageRevoked         WSEventType = "message.revoked"
	WSEventFriendRequestCreated   WSEventType = "friend_request.created"
	WSEventFriendRequestAccepted  WSEventType = "friend_request.accepted"
	WSEventFriendRequestRejected  WSEventType = "friend_request.rejected"
	WSEventFriendshipDeleted      WSEventType = "friendship.deleted"
	WSEventGroupCreated           WSEventType = "group.created"
	WSEventGroupDeleted           WSEventType = "group.deleted"
	WSEventMemberJoined           WSEventType = "group_join_request.accepted"
	WSEventMemberLeft             WSEventType = "group.member_left"
	WSEventMemberRemoved          WSEventType = "group.member_removed"
	WSEventJoinRequestCreated     WSEventType = "group_join_request.created"
	WSEventJoinRequestRejected    WSEventType = "group_join_request.rejected"
	WSEventOwnerTransferred       WSEventType = "group.owner_transferred"
	WSEventUserOnline             WSEventType = "user.online"
	WSEventUserOffline            WSEventType = "user.offline"
)

type WSEventMessage struct {
	Type      WSEventType `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type EventConsumer struct {
	subscriber *event.EventSubscriber
	hub        *Hub
}

func NewEventConsumer(hub *Hub, subscriber *event.EventSubscriber) *EventConsumer {
	c := &EventConsumer{
		subscriber: subscriber,
		hub:        hub,
	}

	c.registerHandlers()
	return c
}

func (c *EventConsumer) registerHandlers() {
	c.subscriber.Register("message.sent", c.handleMessageSent)
	c.subscriber.Register("message.read", c.handleMessageRead)
	c.subscriber.Register("message.revoked", c.handleMessageRevoked)
	c.subscriber.Register("friend_request.created", c.handleFriendRequestCreated)
	c.subscriber.Register("friend_request.accepted", c.handleFriendRequestAccepted)
	c.subscriber.Register("friend_request.rejected", c.handleFriendRequestRejected)
	c.subscriber.Register("friendship.deleted", c.handleFriendshipDeleted)
	c.subscriber.Register("group.created", c.handleGroupCreated)
	c.subscriber.Register("group.deleted", c.handleGroupDeleted)
	c.subscriber.Register("group_join_request.accepted", c.handleMemberJoined)
	c.subscriber.Register("group.member_left", c.handleMemberLeft)
	c.subscriber.Register("group.member_removed", c.handleMemberRemoved)
	c.subscriber.Register("group_join_request.created", c.handleJoinRequestCreated)
	c.subscriber.Register("group_join_request.rejected", c.handleJoinRequestRejected)
	c.subscriber.Register("group.owner_transferred", c.handleOwnerTransferred)
	c.subscriber.Register("user.online", c.handleUserOnline)
	c.subscriber.Register("user.offline", c.handleUserOffline)
}

func (c *EventConsumer) Start() error {
	return c.subscriber.Start()
}

func (c *EventConsumer) Close() error {
	return c.subscriber.Close()
}

func (c *EventConsumer) pushToUser(userID string, msgType WSEventType, data interface{}, timestamp int64) {
	wsMsg := WSEventMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: timestamp,
	}

	payload, err := json.Marshal(wsMsg)
	if err != nil {
		logger.Errorw("Failed to marshal WS event", "component", "event_consumer", "err", err)
		return
	}

	if c.hub.IsOnline(userID) {
		c.hub.SendToUser(userID, payload)
	}
}

func (c *EventConsumer) handleMessageSent(data []byte) error {
	var evt struct {
		event.BaseEvent
		MessageID  string `json:"MessageID"`
		SenderID   string `json:"SenderID"`
		ReceiverID string `json:"ReceiverID"`
		MsgType    string `json:"MsgType"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.ReceiverID, WSEventMessageSent, map[string]interface{}{
		"message_id":  evt.MessageID,
		"sender_id":   evt.SenderID,
		"receiver_id": evt.ReceiverID,
		"msg_type":    evt.MsgType,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleMessageRead(data []byte) error {
	var evt struct {
		event.BaseEvent
		MessageID string `json:"MessageID"`
		UserID    string `json:"UserID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.UserID, WSEventMessageRead, map[string]interface{}{
		"message_id": evt.MessageID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleMessageRevoked(data []byte) error {
	var evt struct {
		event.BaseEvent
		MessageID string `json:"MessageID"`
		UserID    string `json:"UserID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.UserID, WSEventMessageRevoked, map[string]interface{}{
		"message_id": evt.MessageID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleFriendRequestCreated(data []byte) error {
	var evt struct {
		event.BaseEvent
		RequestID string `json:"RequestID"`
		FromUID   string `json:"FromUID"`
		ToUID     string `json:"ToUID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.ToUID, WSEventFriendRequestCreated, map[string]interface{}{
		"request_id": evt.RequestID,
		"from_uid":   evt.FromUID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleFriendRequestAccepted(data []byte) error {
	var evt struct {
		event.BaseEvent
		RequestID string `json:"RequestID"`
		FromUID   string `json:"FromUID"`
		ToUID     string `json:"ToUID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.FromUID, WSEventFriendRequestAccepted, map[string]interface{}{
		"request_id": evt.RequestID,
		"to_uid":     evt.ToUID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleFriendRequestRejected(data []byte) error {
	var evt struct {
		event.BaseEvent
		RequestID string `json:"RequestID"`
		FromUID   string `json:"FromUID"`
		ToUID     string `json:"ToUID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.FromUID, WSEventFriendRequestRejected, map[string]interface{}{
		"request_id": evt.RequestID,
		"to_uid":     evt.ToUID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleFriendshipDeleted(data []byte) error {
	var evt struct {
		event.BaseEvent
		UserID   string `json:"UserID"`
		FriendID string `json:"FriendID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.FriendID, WSEventFriendshipDeleted, map[string]interface{}{
		"user_id": evt.UserID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleGroupCreated(data []byte) error {
	var evt struct {
		event.BaseEvent
		GroupID string `json:"GroupID"`
		Name    string `json:"Name"`
		OwnerID string `json:"OwnerID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.OwnerID, WSEventGroupCreated, map[string]interface{}{
		"group_id": evt.GroupID,
		"name":     evt.Name,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleGroupDeleted(data []byte) error {
	var evt struct {
		event.BaseEvent
		GroupID string `json:"GroupID"`
		OwnerID string `json:"OwnerID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.OwnerID, WSEventGroupDeleted, map[string]interface{}{
		"group_id": evt.GroupID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleMemberJoined(data []byte) error {
	var evt struct {
		event.BaseEvent
		RequestID string `json:"RequestID"`
		UserID    string `json:"UserID"`
		GroupID   string `json:"GroupID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.UserID, WSEventMemberJoined, map[string]interface{}{
		"group_id": evt.GroupID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleMemberLeft(data []byte) error {
	var evt struct {
		event.BaseEvent
		GroupID string `json:"GroupID"`
		UserID  string `json:"UserID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.UserID, WSEventMemberLeft, map[string]interface{}{
		"group_id": evt.GroupID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleMemberRemoved(data []byte) error {
	var evt struct {
		event.BaseEvent
		GroupID string `json:"GroupID"`
		UserID  string `json:"UserID"`
		AdminID string `json:"AdminID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.UserID, WSEventMemberRemoved, map[string]interface{}{
		"group_id": evt.GroupID,
		"admin_id": evt.AdminID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleJoinRequestCreated(data []byte) error {
	var evt struct {
		event.BaseEvent
		RequestID string `json:"RequestID"`
		UserID    string `json:"UserID"`
		GroupID   string `json:"GroupID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.UserID, WSEventJoinRequestCreated, map[string]interface{}{
		"request_id": evt.RequestID,
		"group_id":   evt.GroupID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleJoinRequestRejected(data []byte) error {
	var evt struct {
		event.BaseEvent
		RequestID string `json:"RequestID"`
		UserID    string `json:"UserID"`
		GroupID   string `json:"GroupID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.UserID, WSEventJoinRequestRejected, map[string]interface{}{
		"request_id": evt.RequestID,
		"group_id":   evt.GroupID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleOwnerTransferred(data []byte) error {
	var evt struct {
		event.BaseEvent
		GroupID    string `json:"GroupID"`
		OldOwnerID string `json:"OldOwnerID"`
		NewOwnerID string `json:"NewOwnerID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.pushToUser(evt.NewOwnerID, WSEventOwnerTransferred, map[string]interface{}{
		"group_id":     evt.GroupID,
		"old_owner_id": evt.OldOwnerID,
	}, evt.OccurredAt.Unix())
	return nil
}

func (c *EventConsumer) handleUserOnline(data []byte) error {
	var evt struct {
		event.BaseEvent
		UserID string `json:"UserID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	c.hub.Broadcast(nil)
	_ = evt.UserID
	return nil
}

func (c *EventConsumer) handleUserOffline(data []byte) error {
	var evt struct {
		event.BaseEvent
		UserID string `json:"UserID"`
	}
	if err := json.Unmarshal(data, &evt); err != nil {
		return err
	}

	_ = evt.UserID
	return nil
}