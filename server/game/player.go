package game

import (
	"fmt"
	"math"
	"math/rand"
	"server/server/events"
)

type Player struct {
	Name        string  `json:"name"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	XFloat      float64 `json:"x_float"`
	YFloat      float64 `json:"y_float"`
	Speed       float64 `json:"speed"`
	Angle       float64 `json:"angle"`
	TurnAngle   float64 `json:"turnAngle"`
	history     []Point
	alive       bool
	inGap       bool
	inGapTimer  float64
	gapTimer    float64
	pressedKeys events.KeyState
	ready       bool
}

func NewPlayer(gs *GameState, name string) *Player {
	fmt.Printf("Creating player with rand n bounds %d, %d\n", gs.GetMapWidth()/2, gs.GetMapHeight()/2)
	//Place Player at random x,y with 1/4 margin
	x := rand.Intn(gs.GetMapWidth()/2) + gs.GetMapWidth()/4
	y := rand.Intn(gs.GetMapHeight()/2) + gs.GetMapHeight()/4
	return &Player{
		Name:       name,
		X:          x,
		Y:          y,
		XFloat:     float64(x),
		YFloat:     float64(y),
		Speed:      BaseSpeed / 2,
		Angle:      rand.Float64() * 2 * math.Pi,
		TurnAngle:  TurnAngle,
		history:    make([]Point, 0),
		alive:      true,
		inGap:      false,
		inGapTimer: GapTime,
		gapTimer:   0,
		ready:      false,
	}
}

func (p *Player) kill(name string) {
	p.alive = false
}

func (p *Player) update(delta float64) {
	if !p.alive {
		return
	}
	if p.pressedKeys.ArrowLeft {
		p.Angle -= p.TurnAngle * delta
	} else if p.pressedKeys.ArrowRight {
		p.Angle += p.TurnAngle * delta
	}

	p.X += int(p.Speed * math.Cos(p.Angle) * delta)
	p.Y += int(p.Speed * math.Sin(p.Angle) * delta)
	p.XFloat += p.Speed * math.Cos(p.Angle) * delta
	p.YFloat += p.Speed * math.Sin(p.Angle) * delta

	if p.gapTimer <= 0 {
		p.inGap = true
		p.inGapTimer = GapLength
		p.gapTimer = GapTime
	}

	if !p.inGap {
		p.gapTimer -= delta
		p.history = append(p.history, *NewPoint(p.X, p.Y))
	} else {
		p.inGapTimer -= delta
		if p.inGapTimer <= 0 {
			p.inGap = false
		}
	}
}

func (p *Player) OnKeystateChange(state events.KeyState) {
	p.pressedKeys = state
}

// func (p *Player) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(map[string]interface{}{
// 		"x":     p.X,
// 		"y":     p.Y,
// 		"alive": p.alive,
// 		"angle": p.Angle,
// 		"speed": p.Speed,
// 	})
// }
