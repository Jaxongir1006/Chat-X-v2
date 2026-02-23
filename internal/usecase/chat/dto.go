package chat

import (
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type CreateGroupRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	UserIDs     []uint64 `json:"user_ids"`
}

type StartDMRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
}

type MessageResponse struct {
	ID             uint64             `json:"id"`
	ConversationID uint64             `json:"conversation_id"`
	SenderID       *uint64            `json:"sender_id"`
	Type           domain.MessageType `json:"type"`
	Text           *string            `json:"text"`
	CreatedAt      time.Time          `json:"created_at"`
}

type ConversationResponse struct {
	ID            uint64                  `json:"id"`
	Type          domain.ConversationType `json:"type"`
	Title         *string                 `json:"title"`
	LastMessageID *uint64                 `json:"last_message_id"`
	UpdatedAt     time.Time               `json:"updated_at"`
}
