package repository

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type ChatRepository interface {
	CreateChatRoom(couple *model.Couple) (*model.ChatRoom, error)
}

type chatRepositoryImpl struct {
	db *sqlx.DB
}

func (c *chatRepositoryImpl) CreateChatRoom(couple *model.Couple) (*model.ChatRoom, error) {
	now := time.Now()
	sql := `INSERT INTO chat_rooms (couple_id, created_at) VALUES (?, ?)`
	result, err := c.db.Exec(sql, couple.ID, now)
	if err != nil {
		return nil, err
	}
	roomId, err := result.LastInsertId()
	return &model.ChatRoom{
		ID:        roomId,
		Couple:    couple,
		CreatedAt: now,
	}, nil
}

func NewChatRepository(db *sqlx.DB) ChatRepository {
	return &chatRepositoryImpl{db: db}
}
