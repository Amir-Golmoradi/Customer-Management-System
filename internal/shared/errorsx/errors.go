package errorsx

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindNotFound   Kind = "NOT_FOUND"
	KindConflict   Kind = "CONFLICT"
	KindValidation Kind = "VALIDATION_ERROR"
	KindInternal   Kind = "INTERNAL"
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(kind Kind, code, message string, err error) *Error {
	return &Error{Kind: kind, Code: code, Message: message, Err: err}
}

func IsKind(err error, kind Kind) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind == kind
	}
	return false
}

func ToError(err error) *Error {
	if err == nil {
		return nil
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{
		Kind:    KindInternal,
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		Err:     err,
	}
}
