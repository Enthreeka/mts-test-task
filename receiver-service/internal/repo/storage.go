package repo

import (
	"context"
)

type Message interface {
	Create(ctx context.Context) error
}
