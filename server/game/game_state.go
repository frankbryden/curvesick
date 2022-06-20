package game

import (
	"encoding/json"
	"fmt"
	"server/server/events"
	"server/server/utils"
)

type GameState struct {
	stateMap map[int](map[int]bool)
	players  []*Player
}

func NewGameState(width, height int) *GameState {
	stateMap := make(map[int](map[int]bool), height)
	for i := 0; i < height; i++ {
		stateMap[i] = make(map[int]bool)
		for j := 0; j < width; j++ {
			stateMap[i][j] = false
		}
	}
	fmt.Printf("Created new state map of size %d x %d\n", len(stateMap[0]), len(stateMap))
	return &GameState{
		stateMap: stateMap,
		players:  make([]*Player, 0),
	}
}

// Pass in some player related data - not sure what yet
func (gs *GameState) RegisterPlayer(player *Player) {
	gs.players = append(gs.players, player)
}

func (gs *GameState) GetMapWidth() int {
	return len(gs.stateMap[0])
}

func (gs *GameState) GetMapHeight() int {
	return len(gs.stateMap)
}

type MiniStateUpdatePlayer struct {
	X     int     `json:"x"`
	Y     int     `json:"y"`
	Speed float64 `json:"speed"`
	Angle float64 `json:"angle"`
}

func NewMiniStateUpdatePlayer(player *Player) *MiniStateUpdatePlayer {
	return &MiniStateUpdatePlayer{
		X:     player.X,
		Y:     player.Y,
		Speed: player.Speed,
		Angle: player.Angle,
	}
}

func (gs *GameState) Marshall2() []byte {
	players := make([]*MiniStateUpdatePlayer, len(gs.players))
	for index, player := range gs.players {
		players[index] = NewMiniStateUpdatePlayer(player)
	}
	data, err := json.Marshal(players)
	utils.CheckErr(err)
	return data
}

type InitialGameDef struct {
	Width   int       `json:"width"`
	Height  int       `json:"height"`
	Players []*Player `json:"players"`
}

type RegularGameStateDef struct {
	Players []*Player `json:"players"`
}

func NewInitialGameDef(gs *GameState) *InitialGameDef {
	return &InitialGameDef{
		Width:   len(gs.stateMap[0]),
		Height:  len(gs.stateMap),
		Players: gs.players,
	}
}

func (igd InitialGameDef) GetPayload() interface{} {
	return igd
}

func (igd InitialGameDef) GetType() events.ClientMessageType {
	return events.ClientMessageTypeInitialStateDef
}

func NewRegularGameStateDef(gs *GameState) *RegularGameStateDef {
	return &RegularGameStateDef{
		Players: gs.players,
	}
}

func (rgd RegularGameStateDef) GetPayload() interface{} {
	return rgd
}

func (igd RegularGameStateDef) GetType() events.ClientMessageType {
	return events.ClientMessageTypeRegularStateDef
}
