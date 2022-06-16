package game

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"server/server/events"
)

type Player struct {
	x, y             int
	speed            float64
	angle, turnAngle float64
	history          []Point
	alive            bool
	inGap            bool
	inGapTimer       float64
	gapTimer         float64
	pressedKeys      events.KeyState
}

func NewPlayer(gs *GameState) *Player {
	fmt.Printf("Creating player with rand n bounds %d, %d\n", gs.GetMapWidth()/2, gs.GetMapHeight()/2)
	//Place Player at random x,y with 1/4 margin
	x := rand.Intn(gs.GetMapWidth()/2) + gs.GetMapWidth()/4
	y := rand.Intn(gs.GetMapHeight()/2) + gs.GetMapHeight()/4
	return &Player{
		x:          x,
		y:          y,
		speed:      BaseSpeed,
		angle:      rand.Float64() * 2 * math.Pi,
		turnAngle:  TurnAngle,
		history:    make([]Point, 0),
		alive:      true,
		inGap:      false,
		inGapTimer: GapTime,
		gapTimer:   0,
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
		p.angle -= p.turnAngle * delta
	} else if p.pressedKeys.ArrowRight {
		p.angle += p.turnAngle * delta
	}
	fmt.Printf("We need to move in direction %f at a speed of %f (%f s have passed)\n", p.angle, p.speed, delta)
	p.x += int(p.speed * math.Cos(p.angle) * delta)
	p.y += int(p.speed * math.Sin(p.angle) * delta)

	if p.gapTimer <= 0 {
		p.inGap = true
		p.inGapTimer = GapLength
		p.gapTimer = GapTime
	}

	if !p.inGap {
		p.gapTimer -= delta
		p.history = append(p.history, *NewPoint(p.x, p.y))
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

func (p *Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"x":     p.x,
		"y":     p.y,
		"alive": p.alive,
		"angle": p.angle,
		"speed": p.speed,
	})
}
