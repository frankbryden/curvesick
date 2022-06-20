package game

import (
	"encoding/json"
	"fmt"
	"server/server/events"
	"server/server/netcode"
	"server/server/utils"
)

type Lobby struct {
	Players []*MiniLobbyPlayer `json:"players"`
}

func NewLobby(clients map[*netcode.Client]*Player) *Lobby {
	fmt.Println("NewLobby")
	fmt.Println(clients)
	players := make([]*MiniLobbyPlayer, len(clients))
	i := 0
	for _, player := range clients {
		players[i] = NewMiniPlayer(player)
		i++
	}
	return &Lobby{
		Players: players,
	}
}

// func (l Lobby) GetPayload() []byte {
func (l Lobby) GetPayload() interface{} {
	// val, err := json.Marshal(l)
	// utils.CheckErr(err)
	// return val
	return l
}

func (l Lobby) GetType() events.ClientMessageType {
	return events.ClientMessageTypeLobby
}

type MiniLobbyPlayer struct {
	Name  string `json:"name"`
	Ready bool   `json:"isReady"`
	//level, powerups...
}

func NewMiniPlayer(player *Player) *MiniLobbyPlayer {
	return &MiniLobbyPlayer{
		Name:  player.Name,
		Ready: player.ready,
	}
}

func (l *Lobby) Marshall() []byte {
	jsonData, err := json.Marshal(l)
	utils.CheckErr(err)
	return jsonData
}
