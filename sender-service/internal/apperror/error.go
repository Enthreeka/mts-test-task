package apperror

import (
	"fmt"
	"google.golang.org/grpc/codes"
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

	}

	return codes.Internal
}
