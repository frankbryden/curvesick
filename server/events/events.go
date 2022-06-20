package events

import (
	"encoding/json"
)

type Key int
type EventType string
type LobbyEventType string
type GameEventType string

const (
	Left Key = iota
	Right
	Up
	Down
)

const (
	TypeGameEvent               EventType = "game"
	TypeKeyboardUpdateEvent     EventType = "keyboard_state"
	TypeLobbyEvent              EventType = "lobby"
	TypePlayerRegistrationEvent EventType = "register"
	TypePlayerEvent             EventType = "state"
	TypeGenericEvent            EventType = "_generic"
)

const (
	TypeLobbyEventReady      LobbyEventType = "ready"
	TypeLobbyEventUnready    LobbyEventType = "unready"
	TypeLobbyEventUnregister LobbyEventType = "unregister"
)

const (
	TypeGameEventRoundStart GameEventType = "round_start"
	TypeGameEventRoundEnd   GameEventType = "round_end"
)

type KeyState struct {
	ArrowLeft, ArrowRight, ArrowUp, ArrowDown bool
}

type GenericEvent struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
	// Data map[string]interface{} `json:"data"`
}

type EventConsumer interface {
	OnKeystateChange(k KeyState)
}

type ClientExitEventListener interface {
	OnClientExit()
}

type Event interface {
	GetType() EventType
}

type PlayerEvent struct {
	X, Y  int
	Speed int
	Angle float64
}

func (gpe *PlayerEvent) GetType() EventType {
	return TypePlayerEvent
}

type PlayerRegistrationEvent struct {
	Name string `json:"name"`
}

func (gpe PlayerRegistrationEvent) GetType() EventType {
	return TypePlayerRegistrationEvent
}

func NewPlayerRegistrationEvent(name string) *PlayerRegistrationEvent {
	return &PlayerRegistrationEvent{
		Name: name,
	}
}

type LobbyEvent struct {
	SubType LobbyEventType
}

func (ge LobbyEvent) GetType() EventType {
	return TypeLobbyEvent
}

type GameEvent struct {
	SubType GameEventType
}

func (ge GameEvent) GetType() EventType {
	return TypeGameEvent
}

//ClientGameEvent represents a message in the server -> client direction
type ClientGameEvent struct {
	SubType GameEventType `json:"sub_type"`
}

func (cge ClientGameEvent) GetPayload() interface{} {
	return cge
}

func (cge ClientGameEvent) GetType() ClientMessageType {
	return ClientMessageTypeGameEvent
}

type KeyboardUpdateEvent struct {
	ArrowLeft, ArrowRight, ArrowUp, ArrowDown bool
}

func (kue KeyboardUpdateEvent) GetType() EventType {
	return TypeKeyboardUpdateEvent
}

//Abstract errorring methods
func (ge GenericEvent) GetType() EventType {
	// fmt.Println("This is no good! Trying to get event type of generic type")
	return TypeGenericEvent
}
