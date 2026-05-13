package persistence

import (
	"IM/services/message-service/domain/entity"
	repo "IM/services/message-service/infrastructure/repository"
	"IM/server/model"
	"IM/server/repository/postgres"
	"context"
)

type PostgresMessageRepo struct {
	delegate *postgres.MessageRepo
}

func NewPostgresMessageRepo(delegate *postgres.MessageRepo) repo.MessageRepository {
	return &PostgresMessageRepo{delegate: delegate}
}

func (r *PostgresMessageRepo) toDomain(m *model.Message) *entity.Message {
	if m == nil {
		return nil
	}
	return &entity.Message{
		ID:         m.ID,
		SenderID:   m.SdID,
		ReceiverID: m.RcID,
		Content:    m.Content,
		Type:       m.Type,
		Timestamp:  m.Time,
		IsRead:     m.IsRead,
		CreatedAt:  m.CreatedAt,
	}
}

func (r *PostgresMessageRepo) toDomainList(msgs []*model.Message) []*entity.Message {
	result := make([]*entity.Message, len(msgs))
	for i, m := range msgs {
		result[i] = r.toDomain(m)
	}
	return result
}

func (r *PostgresMessageRepo) toModel(e *entity.Message) *model.Message {
	return &model.Message{
		ID:      e.ID,
		SdID:    e.SenderID,
		RcID:    e.ReceiverID,
		Content: e.Content,
		Type:    e.Type,
		Time:    e.Timestamp,
		IsRead:  e.IsRead,
	}
}

func (r *PostgresMessageRepo) Create(ctx context.Context, msg *entity.Message) error {
	return r.delegate.Create(ctx, r.toModel(msg))
}

func (r *PostgresMessageRepo) GetByID(ctx context.Context, id string) (*entity.Message, error) {
	m, err := r.delegate.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(m), nil
}

func (r *PostgresMessageRepo) GetOfflineMessages(ctx context.Context, userID string, limit, offset int) ([]*entity.Message, error) {
	msgs, err := r.delegate.GetOfflineMessages(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return r.toDomainList(msgs), nil
}

func (r *PostgresMessageRepo) MarkAsRead(ctx context.Context, id string) error {
	return r.delegate.MarkAsRead(ctx, id)
}

func (r *PostgresMessageRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.delegate.MarkAllAsRead(ctx, userID)
}

func (r *PostgresMessageRepo) GetBySender(ctx context.Context, senderID string, limit, offset int) ([]*entity.Message, error) {
	msgs, err := r.delegate.GetBySender(ctx, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	return r.toDomainList(msgs), nil
}

func (r *PostgresMessageRepo) GetByReceiver(ctx context.Context, receiverID string, limit, offset int) ([]*entity.Message, error) {
	msgs, err := r.delegate.GetByReceiver(ctx, receiverID, limit, offset)
	if err != nil {
		return nil, err
	}
	return r.toDomainList(msgs), nil
}

func (r *PostgresMessageRepo) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return r.delegate.GetUnreadCount(ctx, userID)
}

type PostgresFriendshipRepo struct {
	delegate *postgres.FriendshipRepo
}

func NewPostgresFriendshipRepo(delegate *postgres.FriendshipRepo) repo.FriendshipRepository {
	return &PostgresFriendshipRepo{delegate: delegate}
}

func (r *PostgresFriendshipRepo) Exists(ctx context.Context, userID, friendID string) (bool, error) {
	return r.delegate.Exists(ctx, userID, friendID)
}

type PostgresGroupMemberRepo struct {
	delegate *postgres.GroupMemberRepo
}

func NewPostgresGroupMemberRepo(delegate *postgres.GroupMemberRepo) repo.GroupMemberRepository {
	return &PostgresGroupMemberRepo{delegate: delegate}
}

func (r *PostgresGroupMemberRepo) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	return r.delegate.IsMember(ctx, groupID, userID)
}

func (r *PostgresGroupMemberRepo) GetMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	return r.delegate.GetMemberIDs(ctx, groupID)
}