package game

import (
	"fmt"
	"server/server/netcode"
	"time"
)

const TickRate = 10
const GameTickTimeMs = 1000 / TickRate

type Game struct {
	state        *GameState
	server       *netcode.Server
	players      map[*netcode.Client]*Player
	running      bool
	tickDuration time.Duration
}

func NewGame() *Game {
	game := Game{
		state:        NewGameState(320, 160),
		players:      make(map[*netcode.Client]*Player),
		running:      true,
		tickDuration: time.Duration(GameTickTimeMs * float64(time.Millisecond)),
	}
	game.server = netcode.NewServer(game.RegisterPlayer)
	return &game
}

func (g *Game) RegisterPlayer(client *netcode.Client) {
	fmt.Println("Registering new player!")
	//Register event loop management class which has a websocket and a player
	//any message received on the websocket acts upon the player
	player := NewPlayer(g.state)

	g.state.players = append(g.state.players, player)

	client.SetEventConsumer(player)
	client.SetClientExitEventListener(func() { g.ClientDisconnect(client) })
	g.players[client] = player
}

func (g *Game) GetServer() *netcode.Server {
	return g.server
}

func (g *Game) ClientDisconnect(client *netcode.Client) {
	fmt.Println("Client disconnected!")
	delete(g.players, client)
}

func (g *Game) GameLoop() {
	lastFrame := time.Now()
	for g.running {
		start := time.Now()
		delta := start.Sub(lastFrame)
		//Game loop

		//Update players
		for _, player := range g.players {
			player.update(delta.Seconds())
		}

		//Transmit updated state to players
		for client := range g.players {
			state := g.state.Marshall()
			client.WriteMessage(string(state))
		}

		t := time.Now()
		elapsed := t.Sub(start)

		remainingTime := g.tickDuration - elapsed
		fmt.Printf("Going to sleep for %dms\n", remainingTime.Milliseconds())
		//elapsed.Milliseconds()
		lastFrame = t
		time.Sleep(remainingTime)
	}
}
