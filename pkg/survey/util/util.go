package util

type FocusMsg struct{}

func Focus() FocusMsg {
	return FocusMsg{}
}

type BlurMsg struct{}

func Blur() BlurMsg {
	return BlurMsg{}
}
