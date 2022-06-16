package game

import (
	"encoding/json"
	"fmt"
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

type miniPlayer struct {
	X, Y int
}

func (gs *GameState) Marshall() []byte {
	mp := miniPlayer{
		X: gs.players[0].x,
		Y: gs.players[0].y,
	}
	data, err := json.Marshal(mp)
	utils.CheckErr(err)
	return data
}
