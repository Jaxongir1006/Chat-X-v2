package chat

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Jaxongir1006/Chat-X-v2/internal/domain"
	apperr "github.com/Jaxongir1006/Chat-X-v2/internal/errors"
)

func (u *ChatUsecase) StartDM(ctx context.Context, currentUserID, targetUserID uint64) (*ConversationResponse, error) {
	if currentUserID == targetUserID {
		return nil, apperr.New(apperr.CodeBadRequest, http.StatusBadRequest, "cannot start DM with yourself")
	}

	existing, err := u.chatStore.GetDMConversation(ctx, currentUserID, targetUserID)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}
	if existing != nil {
		return &ConversationResponse{
			ID:            existing.ID,
			Type:          existing.Type,
			LastMessageID: existing.LastMessageID,
			UpdatedAt:     existing.UpdatedAt,
		}, nil
	}

	var conv domain.Conversation
	err = u.uow.Do(ctx, func(tx *sql.Tx) error {
		chatTx := u.chatStore.WithTx(tx)

		conv = domain.Conversation{
			Type:      domain.ConversationTypeDM,
			CreatedBy: currentUserID,
		}

		if err := chatTx.CreateConversation(ctx, &conv); err != nil {
			return err
		}

		if err := chatTx.CreateDMConversation(ctx, currentUserID, targetUserID, conv.ID); err != nil {
			return err
		}

		// Add both participants
		if err := chatTx.AddParticipant(ctx, &domain.Participant{ConversationID: conv.ID, UserID: currentUserID, Role: domain.ParticipantRoleMember}); err != nil {
			return err
		}
		if err := chatTx.AddParticipant(ctx, &domain.Participant{ConversationID: conv.ID, UserID: targetUserID, Role: domain.ParticipantRoleMember}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "failed to start DM", err)
	}

	return &ConversationResponse{
		ID:        conv.ID,
		Type:      conv.Type,
		UpdatedAt: conv.UpdatedAt,
	}, nil
}

func (u *ChatUsecase) CreateGroup(ctx context.Context, currentUserID uint64, req CreateGroupRequest) (*ConversationResponse, error) {
	var conv domain.Conversation
	err := u.uow.Do(ctx, func(tx *sql.Tx) error {
		chatTx := u.chatStore.WithTx(tx)

		conv = domain.Conversation{
			Type:        domain.ConversationTypeGroup,
			Title:       &req.Title,
			Description: &req.Description,
			CreatedBy:   currentUserID,
		}

		if err := chatTx.CreateConversation(ctx, &conv); err != nil {
			return err
		}

		// Add owner
		if err := chatTx.AddParticipant(ctx, &domain.Participant{ConversationID: conv.ID, UserID: currentUserID, Role: domain.ParticipantRoleOwner}); err != nil {
			return err
		}

		// Add other participants
		for _, userID := range req.UserIDs {
			if userID == currentUserID {
				continue
			}
			if err := chatTx.AddParticipant(ctx, &domain.Participant{ConversationID: conv.ID, UserID: userID, Role: domain.ParticipantRoleMember}); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "failed to create group", err)
	}

	return &ConversationResponse{
		ID:        conv.ID,
		Type:      conv.Type,
		Title:     conv.Title,
		UpdatedAt: conv.UpdatedAt,
	}, nil
}

func (u *ChatUsecase) GetMessages(ctx context.Context, userID, conversationID uint64, limit, offset int) ([]MessageResponse, error) {
	// Optional: Check if user is participant
	
	messages, err := u.chatStore.GetMessages(ctx, conversationID, limit, offset)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	resp := make([]MessageResponse, 0, len(messages))
	for _, m := range messages {
		resp = append(resp, MessageResponse{
			ID:             m.ID,
			ConversationID: m.ConversationID,
			SenderID:       m.SenderID,
			Type:           m.Type,
			Text:           m.Text,
			CreatedAt:      m.CreatedAt,
		})
	}
	return resp, nil
}

func (u *ChatUsecase) GetConversations(ctx context.Context, userID uint64) ([]ConversationResponse, error) {
	convs, err := u.chatStore.ListConversationsByUserID(ctx, userID)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err)
	}

	resp := make([]ConversationResponse, 0, len(convs))
	for _, c := range convs {
		resp = append(resp, ConversationResponse{
			ID:            c.ID,
			Type:          c.Type,
			Title:         c.Title,
			LastMessageID: c.LastMessageID,
			UpdatedAt:     c.UpdatedAt,
		})
	}
	return resp, nil
}
