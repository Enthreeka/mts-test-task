package entity

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Msg         string    `json:"message"`
	MsgUUID     uuid.UUID `json:"msg_id"`
	CreatedTime time.Time `json:"created_time"`
}
