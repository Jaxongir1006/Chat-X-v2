package domain

import "time"

type ConversationType string

const (
	ConversationTypeDM      ConversationType = "dm"
	ConversationTypeGroup   ConversationType = "group"
	ConversationTypeChannel ConversationType = "channel"
)

type ParticipantRole string

const (
	ParticipantRoleOwner      ParticipantRole = "owner"
	ParticipantRoleAdmin      ParticipantRole = "admin"
	ParticipantRoleMember     ParticipantRole = "member"
	ParticipantRoleRestricted ParticipantRole = "restricted"
	ParticipantRoleBanned     ParticipantRole = "banned"
	ParticipantRoleLeft       ParticipantRole = "left"
)

type MessageType string

const (
	MessageTypeText    MessageType = "text"
	MessageTypePhoto   MessageType = "photo"
	MessageTypeVideo   MessageType = "video"
	MessageTypeFile    MessageType = "file"
	MessageTypeVoice   MessageType = "voice"
	MessageTypeSticker MessageType = "sticker"
	MessageTypeSystem  MessageType = "system"
)

type Conversation struct {
	ID            uint64           `json:"id"`
	Type          ConversationType `json:"type"`
	Title         *string          `json:"title,omitempty"`
	Username      *string          `json:"username,omitempty"`
	Description   *string          `json:"description,omitempty"`
	IsPublic      bool             `json:"is_public"`
	CreatedBy     uint64           `json:"created_by"`
	LastMessageID *uint64          `json:"last_message_id,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type Message struct {
	ID             uint64      `json:"id"`
	ConversationID uint64      `json:"conversation_id"`
	SenderID       *uint64     `json:"sender_id,omitempty"`
	Type           MessageType `json:"type"`
	Text           *string     `json:"text,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	EditedAt       *time.Time  `json:"edited_at,omitempty"`
	ReplyToID      *uint64     `json:"reply_to_id,omitempty"`
	ForwardFromID  *uint64     `json:"forward_from_id,omitempty"`
	DeletedAt      *time.Time  `json:"deleted_at,omitempty"`
}

type Participant struct {
	ConversationID    uint64          `json:"conversation_id"`
	UserID            uint64          `json:"user_id"`
	Role              ParticipantRole `json:"role"`
	JoinedAt          time.Time       `json:"joined_at"`
	LeftAt            *time.Time      `json:"left_at,omitempty"`
	MutedUntil        *time.Time      `json:"muted_until,omitempty"`
	IsPinned          bool            `json:"is_pinned"`
	LastReadMessageID *uint64         `json:"last_read_message_id,omitempty"`
}
