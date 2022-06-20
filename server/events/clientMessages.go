package events

import (
	"encoding/json"
	"server/server/utils"
)

type ClientMessageType string

const (
	ClientMessageTypeLobby           = "lobby"
	ClientMessageTypeInitialStateDef = "init_state_def"
	ClientMessageTypeRegularStateDef = "reg_state_def"
	ClientMessageTypeGameEvent       = "game_event"

	// ClientMessageTypeLobby = "lobby"
	// ClientMessageTypeLobby = "lobby"
)

type ClientMessage interface {
	GetType() ClientMessageType
	GetPayload() interface{}
}

func ConstructClientMessage(message ClientMessage) []byte {
	finalMessage := struct {
		Type ClientMessageType `json:"type"`
		// Data []byte            `type:"data"`
		Data interface{} `json:"data"`
	}{
		message.GetType(),
		message.GetPayload(),
	}
	val, err := json.Marshal(finalMessage)
	utils.CheckErr(err)
	return val
}
