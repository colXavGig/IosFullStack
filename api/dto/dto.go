package dto

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model

		Username string `gorm:"unique;not null"`
		Password string `gorm:"not null"`
	}

	Conversation struct {
		gorm.Model

		Members  []User `gorm:"many2many:Members"`
		Messages []Message
	}

	Message struct {
		gorm.Model

		ConversationId uint `gorm:"not null;foreignKey:Conversation"`

		SenderID uint   `gorm:"not null;foreignKey:User"`
		Sender   User   `gorm:"not null"`
		Content  string `gorm:"not null"`
	}
)
