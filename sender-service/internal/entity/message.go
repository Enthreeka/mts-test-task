package entity

import "time"

type Message struct {
	Msg         string    `json:"message"`
	CreatedTime time.Time `json:"creeate_time"`
}
