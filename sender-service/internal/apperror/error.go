package apperror

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
)

var (
	ErrEmptyMessage = NewError("message is empty", errors.New("empty_message"))
)

type AppError struct {
	Msg string `json:"message"`
	Err error  `json:"-"`
}

func (a *AppError) Error() string {
	return fmt.Sprintf("%s", a.Msg)
}

func NewError(msg string, err error) *AppError {
	return &AppError{
		Msg: msg,
		Err: err,
	}
}

func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, ErrEmptyMessage):
		return codes.InvalidArgument

	}

	return codes.Internal
}
