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
	query := `INSERT INTO message (msg,created_time,msg_uuid) VALUES ($1,$2,$3)`

	tx, err := m.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, message.Msg, message.CreatedTime, message.MsgUUID)
	return err
}
