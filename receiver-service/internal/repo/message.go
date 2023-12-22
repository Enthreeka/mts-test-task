package repo

import (
	"context"
	"github.com/Entreeka/receiver/internal/entity"
	"github.com/Entreeka/receiver/pkg/postgres"
	"github.com/jackc/pgx/v5"
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

	tx, err := m.Pool.BeginTx(ctx, pgx.TxOptions{})

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, message.Msg, message.CreatedTime)
	return err
}
