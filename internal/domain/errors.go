package domain

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type ExitError struct {
	Code int64
}

func (e ExitError) Error() string {
	return fmt.Sprintf("exit code: %d", e.Code)
}
