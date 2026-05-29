package handler

import (
	"IM/api/gen/common"
	"IM/api/gen/group"
	"IM/api/gen/message"
	"IM/api/gen/user"
)

func toRegisterResponse(resp *user.RegisterResponse) RegisterResponse {
	return RegisterResponse{
		UserId: resp.UserId,
		Name:   resp.Name,
	}
}

func toLoginResponse(resp *user.LoginResponse) LoginResponse {
	return LoginResponse{
		UserId: resp.UserId,
		Name:   resp.Name,
		Token:  resp.Token,
	}
}

func toUserInfo(resp *user.UserInfo) UserInfo {
	return UserInfo{
		ID:        resp.Id,
		Name:      resp.Name,
		Tele:      resp.Tele,
		AvatarURL: resp.AvatarUrl,
		Status:    resp.Status,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func toFriendsListResponse(resp *user.GetFriendsResponse) FriendsListResponse {
	friends := make([]FriendInfo, len(resp.Friends))
	for i, f := range resp.Friends {
		friends[i] = FriendInfo{
			ID:        f.Id,
			Name:      f.Name,
			Tele:      f.Tele,
			AvatarURL: f.AvatarUrl,
			Status:    f.Status,
		}
	}
	return FriendsListResponse{Friends: friends}
}

func toCommonResponse(resp *common.Response) CommonResponse {
	return CommonResponse{
		Success: resp.Success,
		Message: resp.Message,
		Code:    resp.Code,
	}
}

func toFriendRequestsResponse(resp *user.GetPendingFriendRequestsResponse) FriendRequestsResponse {
	requests := make([]FriendRequestInfo, len(resp.Requests))
	for i, r := range resp.Requests {
		var fromUser *UserInfo
		if r.FromUser != nil {
			u := toUserInfo(r.FromUser)
			fromUser = &u
		}
		requests[i] = FriendRequestInfo{
			ID:        r.Id,
			FromUID:   r.FromUid,
			ToUID:     r.ToUid,
			Reason:    r.Reason,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
			FromUser:  fromUser,
		}
	}
	return FriendRequestsResponse{Requests: requests}
}

func toCreateGroupResponse(resp *group.CreateGroupResponse) CreateGroupResponse {
	return CreateGroupResponse{
		GroupId: resp.GroupId,
		Name:    resp.Name,
	}
}

func toGroupInfo(resp *group.GroupInfo) GroupInfo {
	return GroupInfo{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		OwnerID:     resp.OwnerId,
		Type:        resp.Type,
		ImageURL:    resp.ImageUrl,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		MemberCount: resp.MemberCount,
	}
}

func toGroupMembersResponse(resp *group.GetMembersResponse) GroupMembersResponse {
	members := make([]MemberInfo, len(resp.Members))
	for i, m := range resp.Members {
		role := "member"
		switch m.Role {
		case group.MemberRole_MEMBER_ROLE_ADMIN:
			role = "admin"
		case group.MemberRole_MEMBER_ROLE_OWNER:
			role = "owner"
		}
		members[i] = MemberInfo{
			UserID:    m.UserId,
			Name:      m.Name,
			Nickname:  m.Nickname,
			AvatarURL: m.AvatarUrl,
			Role:      role,
			JoinedAt:  m.JoinedAt,
		}
	}
	return GroupMembersResponse{Members: members}
}

func toUserGroupsResponse(resp *group.GetUserGroupsResponse) UserGroupsResponse {
	groups := make([]GroupInfo, len(resp.Groups))
	for i, g := range resp.Groups {
		groups[i] = toGroupInfo(g)
	}
	return UserGroupsResponse{Groups: groups}
}

func toGroupJoinRequestsResponse(resp *group.GetPendingGroupRequestsResponse) GroupJoinRequestsResponse {
	requests := make([]GroupJoinRequestInfo, len(resp.Requests))
	for i, r := range resp.Requests {
		requests[i] = GroupJoinRequestInfo{
			ID:         r.Id,
			UserID:     r.UserId,
			GroupID:    r.GroupId,
			Reason:     r.Reason,
			Status:     r.Status,
			CreatedAt:  r.CreatedAt,
			UserName:   r.UserName,
			UserAvatar: r.UserAvatar,
		}
	}
	return GroupJoinRequestsResponse{Requests: requests}
}

func toSendMessageResponse(resp *message.SendMessageResponse) SendMessageResponse {
	return SendMessageResponse{
		MsgId:     resp.MsgId,
		Success:   resp.Success,
		Timestamp: resp.Timestamp,
	}
}

func toOfflineMessagesResponse(resp *message.GetOfflineMessagesResponse) OfflineMessagesResponse {
	msgs := make([]MessageInfo, len(resp.Messages))
	for i, m := range resp.Messages {
		msgs[i] = MessageInfo{
			MsgId:      m.MsgId,
			SenderID:   m.SenderId,
			ReceiverID: m.ReceiverId,
			Content:    m.Content,
			MsgType:    m.MsgType,
			Timestamp:  m.Timestamp,
			IsRead:     m.IsRead,
			IsRevoked:  m.IsRevoked,
			ReadAt:     m.ReadAt,
		}
	}
	return OfflineMessagesResponse{
		Messages: msgs,
		Total:    resp.Total,
	}
}

func toUnreadCountResponse(resp *message.GetUnreadCountResponse) UnreadCountResponse {
	return UnreadCountResponse{
		Count:       resp.Count,
		GroupCounts: resp.GroupCounts,
	}
}
