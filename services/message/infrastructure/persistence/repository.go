// Package persistence 提供 message 服务使用的 MongoDB 持久化实现。
package persistence

import (
	"IM/pkg/config"
	"IM/services/message/domain/entity"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoDB 包装 mongo client、database 与 collection。
type MongoDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewMongoDB 根据配置创建 MongoDB 连接。
func NewMongoDB(cfg config.MongoDBConfig) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	database := client.Database(cfg.Database)
	collection := database.Collection(cfg.Collection)

	return &MongoDB{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

// Close 断开 MongoDB 连接。
func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}

// GetCollection 返回底层 collection 引用。
func (m *MongoDB) GetCollection() *mongo.Collection {
	return m.collection
}

// MessageRepository 提供对消息集合的 CRUD 操作。
type MessageRepository struct {
	db *MongoDB
}

func NewMessageRepository(db *MongoDB) *MessageRepository {
	return &MessageRepository{db: db}
}

type messageDocument struct {
	ID         string `bson:"_id"`
	SenderID   string `bson:"sender_id"`
	ReceiverID string `bson:"receiver_id"`
	Content    string `bson:"content"`
	MsgType    string `bson:"msg_type"`
	Timestamp  int64  `bson:"timestamp"`
	IsRead     bool   `bson:"is_read"`
	IsRevoked  bool   `bson:"is_revoked"`
	ReadAt     int64  `bson:"read_at"`
	CreatedAt  int64  `bson:"created_at"`
}

func toMessageDocument(m *entity.Message) *messageDocument {
	return &messageDocument{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		MsgType:    string(m.MsgType),
		Timestamp:  m.Timestamp,
		IsRead:     m.IsRead,
		IsRevoked:  m.IsRevoked,
		ReadAt:     m.ReadAt,
		CreatedAt:  m.CreatedAt.Unix(),
	}
}

func toMessageEntity(d *messageDocument) *entity.Message {
	return &entity.Message{
		ID:         d.ID,
		SenderID:   d.SenderID,
		ReceiverID: d.ReceiverID,
		Content:    d.Content,
		MsgType:    entity.MessageType(d.MsgType),
		Timestamp:  d.Timestamp,
		IsRead:     d.IsRead,
		IsRevoked:  d.IsRevoked,
		ReadAt:     d.ReadAt,
	}
}

// Create 插入一条消息文档。
func (r *MessageRepository) Create(ctx context.Context, message *entity.Message) error {
	doc := toMessageDocument(message)
	_, err := r.db.collection.InsertOne(ctx, doc)
	return err
}

// GetByID 按 ID 查询消息。
func (r *MessageRepository) GetByID(ctx context.Context, id string) (*entity.Message, error) {
	var doc messageDocument
	filter := bson.M{"_id": id}
	err := r.db.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return toMessageEntity(&doc), nil
}

// GetByReceiverID 获取某用户的离线消息（分页）。
func (r *MessageRepository) GetByReceiverID(ctx context.Context, receiverID string, limit, offset int) ([]*entity.Message, error) {
	filter := bson.M{"receiver_id": receiverID}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.db.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []messageDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, len(docs))
	for i, doc := range docs {
		messages[i] = toMessageEntity(&doc)
	}
	return messages, nil
}

// GetHistory 获取两者之间的历史消息。
func (r *MessageRepository) GetHistory(ctx context.Context, userID, targetID string, beforeTime int64, limit int) ([]*entity.Message, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": userID, "receiver_id": targetID},
			{"sender_id": targetID, "receiver_id": userID},
		},
	}

	if beforeTime > 0 {
		filter["timestamp"] = bson.M{"$lt": beforeTime}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.db.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []messageDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, len(docs))
	for i, doc := range docs {
		messages[i] = toMessageEntity(&doc)
	}
	return messages, nil
}

// GetUnreadByReceiverID 获取用户所有未读消息。
func (r *MessageRepository) GetUnreadByReceiverID(ctx context.Context, receiverID string) ([]*entity.Message, error) {
	filter := bson.M{
		"receiver_id": receiverID,
		"is_read":     false,
	}

	cursor, err := r.db.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []messageDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, len(docs))
	for i, doc := range docs {
		messages[i] = toMessageEntity(&doc)
	}
	return messages, nil
}

// GetUnreadCount 统计未读消息数。
func (r *MessageRepository) GetUnreadCount(ctx context.Context, receiverID string) (int64, error) {
	filter := bson.M{
		"receiver_id": receiverID,
		"is_read":     false,
	}
	return r.db.collection.CountDocuments(ctx, filter)
}

// MarkAsRead 将消息标记为已读，若未曾标记则设置 read_at。
func (r *MessageRepository) MarkAsRead(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
			"read_at": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$read_at", 0}}, bson.M{"$now": 1}, "$read_at"}},
		},
	}
	_, err := r.db.collection.UpdateOne(ctx, filter, update)
	return err
}

// MarkAllAsRead 将指定 sender 与 receiver 之间的消息全部标为已读。
func (r *MessageRepository) MarkAllAsRead(ctx context.Context, receiverID, senderID string) error {
	filter := bson.M{
		"receiver_id": receiverID,
		"sender_id":   senderID,
		"is_read":     false,
	}
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
		},
	}
	_, err := r.db.collection.UpdateMany(ctx, filter, update)
	return err
}

// Revoke 将消息标为已撤回并替换内容为占位文本。
func (r *MessageRepository) Revoke(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_revoked": true,
			"content":    "[消息已撤回]",
		},
	}
	_, err := r.db.collection.UpdateOne(ctx, filter, update)
	return err
}
