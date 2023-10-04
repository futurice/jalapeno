package util

import "errors"

type FocusMsg struct{}

func Focus() FocusMsg {
	return FocusMsg{}
}

type BlurMsg struct{}

func Blur() BlurMsg {
	return BlurMsg{}
}

var (
	ErrRequired    = errors.New("value can not be empty")
	ErrRegExFailed = errors.New("validation failed")
)
