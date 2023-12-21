package apperror

import (
	"errors"
	"fmt"
)

var (
	ErrUniqueViolation     = NewError("Violation must be unique", errors.New("non_unique_value"))
	ErrForeignKeyViolation = NewError("Foreign Key Violation", errors.New("foreign_key_violation "))
	ErrNoRows              = NewError("No rows in result set", errors.New("no_rows"))
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
