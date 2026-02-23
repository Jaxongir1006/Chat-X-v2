package chat

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
)

func (r *chatRepo) CreateConversation(ctx context.Context, conv *domain.Conversation) error {
	query := `INSERT INTO conversations (type, title, username, description, is_public, created_by, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at`
	
	err := r.execer().QueryRowContext(ctx, query, conv.Type, conv.Title, conv.Username, conv.Description, conv.IsPublic, conv.CreatedBy).
		Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)
	return err
}

func (r *chatRepo) GetConversationByID(ctx context.Context, id uint64) (*domain.Conversation, error) {
	query := `SELECT id, type, title, username, description, is_public, created_by, last_message_id, created_at, updated_at
			  FROM conversations WHERE id = $1`
	
	var conv domain.Conversation
	err := r.execer().QueryRowContext(ctx, query, id).Scan(
		&conv.ID, &conv.Type, &conv.Title, &conv.Username, &conv.Description, &conv.IsPublic,
		&conv.CreatedBy, &conv.LastMessageID, &conv.CreatedAt, &conv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *chatRepo) ListConversationsByUserID(ctx context.Context, userID uint64) ([]domain.Conversation, error) {
	query := `SELECT c.id, c.type, c.title, c.username, c.description, c.is_public, c.created_by, c.last_message_id, c.created_at, c.updated_at
			  FROM conversations c
			  JOIN conversation_participants cp ON c.id = cp.conversation_id
			  WHERE cp.user_id = $1 AND cp.left_at IS NULL
			  ORDER BY c.updated_at DESC`
	
	rows, err := r.execer().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []domain.Conversation
	for rows.Next() {
		var c domain.Conversation
		if err := rows.Scan(&c.ID, &c.Type, &c.Title, &c.Username, &c.Description, &c.IsPublic,
			&c.CreatedBy, &c.LastMessageID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, nil
}

func (r *chatRepo) AddParticipant(ctx context.Context, part *domain.Participant) error {
	query := `INSERT INTO conversation_participants (conversation_id, user_id, role, joined_at)
			  VALUES ($1, $2, $3, NOW())
			  ON CONFLICT (conversation_id, user_id) DO UPDATE SET role = $3, left_at = NULL, joined_at = NOW()`
	_, err := r.execer().ExecContext(ctx, query, part.ConversationID, part.UserID, part.Role)
	return err
}

func (r *chatRepo) GetParticipants(ctx context.Context, conversationID uint64) ([]domain.Participant, error) {
	query := `SELECT conversation_id, user_id, role, joined_at, left_at, muted_until, is_pinned, last_read_message_id
			  FROM conversation_participants WHERE conversation_id = $1 AND left_at IS NULL`
	
	rows, err := r.execer().QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.Participant
	for rows.Next() {
		var p domain.Participant
		if err := rows.Scan(&p.ConversationID, &p.UserID, &p.Role, &p.JoinedAt, &p.LeftAt, &p.MutedUntil, &p.IsPinned, &p.LastReadMessageID); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *chatRepo) RemoveParticipant(ctx context.Context, conversationID, userID uint64) error {
	query := `UPDATE conversation_participants SET left_at = NOW() WHERE conversation_id = $1 AND user_id = $2`
	_, err := r.execer().ExecContext(ctx, query, conversationID, userID)
	return err
}

func (r *chatRepo) SendMessage(ctx context.Context, msg *domain.Message) error {
	query := `INSERT INTO messages (conversation_id, sender_id, type, text, created_at)
			  VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	
	err := r.execer().QueryRowContext(ctx, query, msg.ConversationID, msg.SenderID, msg.Type, msg.Text).
		Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return err
	}

	// Update last_message_id and updated_at in conversation
	updateQuery := `UPDATE conversations SET last_message_id = $1, updated_at = NOW() WHERE id = $2`
	_, err = r.execer().ExecContext(ctx, updateQuery, msg.ID, msg.ConversationID)
	return err
}

func (r *chatRepo) GetMessages(ctx context.Context, conversationID uint64, limit, offset int) ([]domain.Message, error) {
	query := `SELECT id, conversation_id, sender_id, type, text, created_at, edited_at, reply_to_id, forward_from_id, deleted_at
			  FROM messages WHERE conversation_id = $1 AND deleted_at IS NULL
			  ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.execer().QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Type, &m.Text, &m.CreatedAt, &m.EditedAt, &m.ReplyToID, &m.ForwardFromID, &m.DeletedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *chatRepo) GetMessageByID(ctx context.Context, id uint64) (*domain.Message, error) {
	query := `SELECT id, conversation_id, sender_id, type, text, created_at, edited_at, reply_to_id, forward_from_id, deleted_at
			  FROM messages WHERE id = $1`
	
	var m domain.Message
	err := r.execer().QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.ConversationID, &m.SenderID, &m.Type, &m.Text, &m.CreatedAt, &m.EditedAt, &m.ReplyToID, &m.ForwardFromID, &m.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *chatRepo) UpdateMessage(ctx context.Context, msg *domain.Message) error {
	query := `UPDATE messages SET text = $1, edited_at = NOW() WHERE id = $2 AND sender_id = $3`
	_, err := r.execer().ExecContext(ctx, query, msg.Text, msg.ID, msg.SenderID)
	return err
}

func (r *chatRepo) DeleteMessage(ctx context.Context, id uint64) error {
	query := `UPDATE messages SET deleted_at = NOW() WHERE id = $1`
	_, err := r.execer().ExecContext(ctx, query, id)
	return err
}

func (r *chatRepo) GetDMConversation(ctx context.Context, user1ID, user2ID uint64) (*domain.Conversation, error) {
	// Ensure user1ID < user2ID as per DM pairs check
	u1, u2 := user1ID, user2ID
	if u1 > u2 {
		u1, u2 = u2, u1
	}

	query := `SELECT c.id, c.type, c.title, c.username, c.description, c.is_public, c.created_by, c.last_message_id, c.created_at, c.updated_at
			  FROM conversations c
			  JOIN dm_pairs dm ON c.id = dm.conversation_id
			  WHERE dm.user1_id = $1 AND dm.user2_id = $2`
	
	var c domain.Conversation
	err := r.execer().QueryRowContext(ctx, query, u1, u2).Scan(
		&c.ID, &c.Type, &c.Title, &c.Username, &c.Description, &c.IsPublic,
		&c.CreatedBy, &c.LastMessageID, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *chatRepo) CreateDMConversation(ctx context.Context, user1ID, user2ID uint64, convID uint64) error {
	u1, u2 := user1ID, user2ID
	if u1 > u2 {
		u1, u2 = u2, u1
	}

	query := `INSERT INTO dm_pairs (user1_id, user2_id, conversation_id) VALUES ($1, $2, $3)`
	_, err := r.execer().ExecContext(ctx, query, u1, u2, convID)
	return err
}
