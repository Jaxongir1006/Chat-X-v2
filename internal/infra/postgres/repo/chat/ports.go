package chat

import (
	"context"
	"database/sql"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

type ChatStore interface {
	WithTx(tx *sql.Tx) *chatRepo

	// Conversations
	CreateConversation(ctx context.Context, conv *domain.Conversation) error
	GetConversationByID(ctx context.Context, id uint64) (*domain.Conversation, error)
	ListConversationsByUserID(ctx context.Context, userID uint64) ([]domain.Conversation, error)
	
	// Participants
	AddParticipant(ctx context.Context, part *domain.Participant) error
	GetParticipants(ctx context.Context, conversationID uint64) ([]domain.Participant, error)
	RemoveParticipant(ctx context.Context, conversationID, userID uint64) error

	// Messages
	SendMessage(ctx context.Context, msg *domain.Message) error
	GetMessages(ctx context.Context, conversationID uint64, limit, offset int) ([]domain.Message, error)
	GetMessageByID(ctx context.Context, id uint64) (*domain.Message, error)
	UpdateMessage(ctx context.Context, msg *domain.Message) error
	DeleteMessage(ctx context.Context, id uint64) error

	// DM specific
	GetDMConversation(ctx context.Context, user1ID, user2ID uint64) (*domain.Conversation, error)
	CreateDMConversation(ctx context.Context, user1ID, user2ID uint64, convID uint64) error
}
