package proxy

import (
	"IM/services/api-gateway/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"IM/api/gen/group"
	"IM/api/gen/message"
	"IM/api/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceProxy struct {
	userSvc    user.UserServiceClient
	groupSvc   group.GroupServiceClient
	messageSvc message.MessageServiceClient
	cfg        *config.Config
	httpClient *http.Client
}

func NewServiceProxy(cfg *config.Config) (*ServiceProxy, error) {
	connUser, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", cfg.Services.UserService.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	connGroup, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", cfg.Services.GroupService.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to group service: %w", err)
	}

	connMessage, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", cfg.Services.MessageService.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to message service: %w", err)
	}

	return &ServiceProxy{
		userSvc:    user.NewUserServiceClient(connUser),
		groupSvc:   group.NewGroupServiceClient(connGroup),
		messageSvc: message.NewMessageServiceClient(connMessage),
		cfg:        cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (p *ServiceProxy) UserService() user.UserServiceClient {
	return p.userSvc
}

func (p *ServiceProxy) GroupService() group.GroupServiceClient {
	return p.groupSvc
}

func (p *ServiceProxy) MessageService() message.MessageServiceClient {
	return p.messageSvc
}

func (p *ServiceProxy) Register(ctx context.Context, tele, name, password string) (*user.RegisterResponse, error) {
	return p.userSvc.Register(ctx, &user.RegisterRequest{
		Tele:     tele,
		Name:     name,
		Password: password,
	})
}

func (p *ServiceProxy) Login(ctx context.Context, tele, password string) (*user.LoginResponse, error) {
	return p.userSvc.Login(ctx, &user.LoginRequest{
		Tele:     tele,
		Password: password,
	})
}

func (p *ServiceProxy) GetUser(ctx context.Context, userID string) (*user.GetUserResponse, error) {
	return p.userSvc.GetUser(ctx, &user.GetUserRequest{
		UserId: userID,
	})
}

func (p *ServiceProxy) GetFriends(ctx context.Context, userID string) (*user.GetFriendsResponse, error) {
	return p.userSvc.GetFriends(ctx, &user.GetFriendsRequest{
		UserId: userID,
	})
}

func (p *ServiceProxy) AddFriend(ctx context.Context, userID, friendID, reason string) (*user.AddFriendResponse, error) {
	return p.userSvc.AddFriend(ctx, &user.AddFriendRequest{
		UserId:   userID,
		FriendId: friendID,
		Reason:   reason,
	})
}

func (p *ServiceProxy) AcceptFriendRequest(ctx context.Context, requestID, userID string) (*user.AcceptFriendRequestResponse, error) {
	return p.userSvc.AcceptFriendRequest(ctx, &user.AcceptFriendRequestRequest{
		RequestId: requestID,
		UserId:    userID,
	})
}

func (p *ServiceProxy) CreateGroup(ctx context.Context, ownerID, name, description string) (*group.CreateGroupResponse, error) {
	return p.groupSvc.CreateGroup(ctx, &group.CreateGroupRequest{
		OwnerId:     ownerID,
		Name:        name,
		Description: description,
	})
}

func (p *ServiceProxy) GetGroup(ctx context.Context, groupID string) (*group.GetGroupResponse, error) {
	return p.groupSvc.GetGroup(ctx, &group.GetGroupRequest{
		GroupId: groupID,
	})
}

func (p *ServiceProxy) GetGroupMembers(ctx context.Context, groupID string) (*group.GetGroupMembersResponse, error) {
	return p.groupSvc.GetGroupMembers(ctx, &group.GetGroupMembersRequest{
		GroupId: groupID,
	})
}

func (p *ServiceProxy) JoinGroup(ctx context.Context, userID, groupID, reason string) (*group.JoinGroupResponse, error) {
	return p.groupSvc.JoinGroup(ctx, &group.JoinGroupRequest{
		UserId:   userID,
		GroupId:  groupID,
		Reason:   reason,
	})
}

func (p *ServiceProxy) LeaveGroup(ctx context.Context, userID, groupID string) (*group.LeaveGroupResponse, error) {
	return p.groupSvc.LeaveGroup(ctx, &group.LeaveGroupRequest{
		UserId:  userID,
		GroupId: groupID,
	})
}

func (p *ServiceProxy) SendMessage(ctx context.Context, senderID, receiverID, content, msgType string) (*message.SendMessageResponse, error) {
	return p.messageSvc.SendMessage(ctx, &message.SendMessageRequest{
		SenderId:   senderID,
		ReceiverId: receiverID,
		Content:    content,
		Type:       msgType,
	})
}

func (p *ServiceProxy) GetOfflineMessages(ctx context.Context, userID string, limit, offset int) (*message.GetOfflineMessagesResponse, error) {
	return p.messageSvc.GetOfflineMessages(ctx, &message.GetOfflineMessagesRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
}

func (p *ServiceProxy) GetUnreadCount(ctx context.Context, userID string) (*message.GetUnreadCountResponse, error) {
	return p.messageSvc.GetUnreadCount(ctx, &message.GetUnreadCountRequest{
		UserId: userID,
	})
}

func (p *ServiceProxy) MarkAsRead(ctx context.Context, messageID string) (*message.MarkAsReadResponse, error) {
	return p.messageSvc.MarkAsRead(ctx, &message.MarkAsReadRequest{
		MessageId: messageID,
	})
}

func (p *ServiceProxy) ForwardToUserService(method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%s%s", p.cfg.Services.UserService.Host, p.cfg.Services.UserService.Port, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return p.httpClient.Do(req)
}

func (p *ServiceProxy) ForwardToGroupService(method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%s%s", p.cfg.Services.GroupService.Host, p.cfg.Services.GroupService.Port, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return p.httpClient.Do(req)
}

func (p *ServiceProxy) ForwardToMessageService(method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%s%s", p.cfg.Services.MessageService.Host, p.cfg.Services.MessageService.Port, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return p.httpClient.Do(req)
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   message,
	}
}

func ParseResponse(resp *http.Response) (APIResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var result APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return APIResponse{}, err
	}
	return result, nil
}

func GetUserIDFromToken(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}
	return strings.TrimPrefix(token, "Bearer "), nil
}