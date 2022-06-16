package events

type Key int

const (
	Left Key = iota
	Right
	Up
	Down
)

type KeyState struct {
	ArrowLeft, ArrowRight, ArrowUp, ArrowDown bool
}

type EventConsumer interface {
	OnKeystateChange(k KeyState)
}

type ClientExitEventListener interface {
	OnClientExit()
}
