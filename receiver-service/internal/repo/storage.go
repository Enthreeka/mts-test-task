package repo

import (
	"context"
	"github.com/Entreeka/receiver/internal/entity"
)

//go:generate mockgen -source storage.go -destination mock/pg_repository_mock.go -package mock
type Message interface {
	Create(ctx context.Context, message *entity.Message) error
}
