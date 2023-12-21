package service

import (
	"context"
	"github.com/Entreeka/receiver/internal/apperror"
	pgxError "github.com/Entreeka/receiver/internal/apperror/pgx_error"
	"github.com/Entreeka/receiver/internal/entity"
	"github.com/Entreeka/receiver/internal/repo"
	"github.com/jackc/pgx/v5"
)

type messageService struct {
	messageRepo repo.Message
}

func NewMessageService(messageRepo repo.Message) Message {
	return &messageService{
		messageRepo: messageRepo,
	}
}

func (m *messageService) Create(ctx context.Context, message *entity.Message) error {
	err := m.messageRepo.Create(ctx, message)
	if err == pgx.ErrNoRows {
		return apperror.ErrNoRows
	}
	errCode := pgxError.ErrorCode(err)
	if errCode == pgxError.ForeignKeyViolation {
		return apperror.ErrForeignKeyViolation
	}
	if errCode == pgxError.UniqueViolation {
		return apperror.ErrUniqueViolation
	}

	return nil
}
