package dto

import (
	"context"

	"gorm.io/gorm"
)

func (msg Message) AfterCreate(tx *gorm.DB) error {
	update, err := gorm.G[Conversation](tx).Where("id = ?", msg.ConversationId).Update(
		context.Background(),
		"updated_at",
		msg.CreatedAt,
	)
	if err != nil {
		return err
	}
	if update < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
