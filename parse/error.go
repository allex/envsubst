package parse

import (
	"errors"
)

type interErr struct {
	error
	code string
}

func (e *interErr) Is(err error) bool {
	var ref *interErr
	if errors.As(err, &ref) {
		return e.code == ref.code
	}
	return false
}

// envsubst internal error, with error code wrapped.
func Error(err string, code string) *interErr {
	return &interErr{
		error: errors.New(err),
		code:  code,
	}
}
