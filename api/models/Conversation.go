package models

import (
	"context"
	"errors"
	"ios_full_stack/data"
	"ios_full_stack/dto"
	"log"

	"gorm.io/gorm"
)

const (
	messagePageSize      uint = 10
	conversationPageSize uint = 10
)

var (
	ErrConversationNotFound error = errors.New("conversation not found")
)

type (
	AppConversation struct {
		dto.Conversation
	}
)

func GetAllConversationOfUser(ctx context.Context, user dto.User, page uint) ([]dto.Conversation, error) {
	var (
		tx            = data.GetTransaction(ctx)
		conversations []dto.Conversation

		err error
	)

	err = tx.
		Preload("Members").
		Preload("Messages").
		Joins("JOIN members ON members.conversation_id = conversations.id").
		Where("members.user_id = ?", user.ID).
		Order("updated_at desc").
		Scopes(data.Paginate(page, conversationPageSize)).
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	return conversations, nil
}

func CreateConversation(ctx context.Context, conversation dto.Conversation) error {
	var (
		tx = data.GetTransaction(ctx)
	)
	return gorm.G[dto.Conversation](tx).Create(ctx, &conversation)
}

func TryFindConversationById(ctx context.Context, id uint) (*AppConversation, error) {
	var (
		tx = data.GetTransaction(ctx)
	)

	convo, err := gorm.G[dto.Conversation](tx).Where("id = ?", id).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrConversationNotFound
	} else if err != nil {
		log.Panicln(err)
	}

	return &AppConversation{Conversation: convo}, nil
}

func (ac *AppConversation) AppendMessage(ctx context.Context, message dto.Message) error {
	var (
		tx = data.GetTransaction(ctx)
	)

	message.ConversationId = ac.ID

	err := gorm.G[dto.Message](tx).Create(ctx, &message)
	if err != nil {
		log.Panicln(err)
		return err
	}
	return nil
}

func (ac *AppConversation) GetMessagePage(ctx context.Context, page uint) ([]dto.Message, error) {
	var (
		tx = data.GetTransaction(ctx)
	)

	messages, err := gorm.G[dto.Message](tx).
		Order("created_at desc").
		Where("conversation_id = ?", ac.ID).
		Scopes(data.GenericPaginate(page, messagePageSize)).
		Find(ctx)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	return messages, nil
}
