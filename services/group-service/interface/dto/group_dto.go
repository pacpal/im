package dto

import "IM/services/group-service/domain/entity"

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JoinGroupRequest struct {
	GroupID string `json:"group_id"`
	Reason  string `json:"reason"`
}

type GroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	Type        string `json:"type"`
	ImageURL    string `json:"image_url"`
}

type GroupMemberResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Tele     string `json:"tele"`
	Role     int16  `json:"role"`
	Nickname string `json:"nickname"`
}

type JoinRequestActionRequest struct {
	RequestID string `json:"request_id"`
}

type UpdateMemberRoleRequest struct {
	TargetID string `json:"target_id"`
	NewRole  int16  `json:"new_role"`
}

func ToGroupResponse(group *entity.Group) *GroupResponse {
	if group == nil {
		return nil
	}
	return &GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		OwnerID:     group.OwnerID,
		Type:        group.Type,
		ImageURL:    group.ImageURL,
	}
}

func ToGroupResponseList(groups []*entity.Group) []*GroupResponse {
	result := make([]*GroupResponse, len(groups))
	for i, g := range groups {
		result[i] = ToGroupResponse(g)
	}
	return result
}

func ToGroupMemberResponse(user *entity.User) *GroupMemberResponse {
	if user == nil {
		return nil
	}
	return &GroupMemberResponse{
		ID:   user.ID,
		Name: user.Name,
		Tele: user.Tele,
		Role: 0,
	}
}

func ToGroupMemberResponseList(users []*entity.User) []*GroupMemberResponse {
	result := make([]*GroupMemberResponse, len(users))
	for i, u := range users {
		result[i] = ToGroupMemberResponse(u)
	}
	return result
}