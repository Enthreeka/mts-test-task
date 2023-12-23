package service

import (
	"context"
	"github.com/Entreeka/receiver/internal/entity"
	"github.com/Entreeka/receiver/internal/repo/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMessageService(t *testing.T) {

	type mockBehavior func(r *mock.MockMessage, message *entity.Message)
	type args struct {
		message *entity.Message
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				message: &entity.Message{
					Msg:         "test",
					MsgUUID:     uuid.New(),
					CreatedTime: time.Now(),
				},
			},
			mockBehavior: func(r *mock.MockMessage, message *entity.Message) {
				r.EXPECT().Create(context.Background(), gomock.Eq(message)).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockMsgRepo := mock.NewMockMessage(ctrl)
			tt.mockBehavior(mockMsgRepo, tt.args.message)

			msgService := NewMessageService(mockMsgRepo)
			err := msgService.Create(context.Background(), tt.args.message)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
