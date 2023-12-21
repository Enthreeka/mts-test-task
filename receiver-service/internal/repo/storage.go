package repo

import (
	"context"
	"github.com/Entreeka/receiver/internal/entity"
)

type Message interface {
	Create(ctx context.Context, message *entity.Message) error
}
