package repo

import (
	"context"
	"github.com/Entreeka/receiver/internal/entity"
	"github.com/Entreeka/receiver/pkg/postgres"
)

type messageRepo struct {
	*postgres.Postgres
}

func NewMessageRepo(pg *postgres.Postgres) Message {
	return &messageRepo{
		pg,
	}
}

func (m *messageRepo) Create(ctx context.Context, message *entity.Message) error {
	query := `INSERT INTO message (msg,created_time) VALUES ($1,$2)`

	_, err := m.Pool.Exec(ctx, query, message.Msg, message.CreatedTime)
	return err
}
