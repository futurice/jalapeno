package util

import (
	"errors"
)

var (
	ErrRequired    = errors.New("value can not be empty")
	ErrRegExFailed = errors.New("validation failed")
	ErrUserAborted = errors.New("user aborted")
)
