package game

import (
	"fmt"
	"os"
	"server/server/events"
	"server/server/netcode"
	"time"
)

const TickRate = 64
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
		state:        NewGameState(200, 200),
		players:      make(map[*netcode.Client]*Player),
		running:      false,
		tickDuration: time.Duration(GameTickTimeMs * float64(time.Millisecond)),
	}
	game.server = netcode.NewServer(game.ClientConnected)
	return &game
}

func (g *Game) ClientConnected(client *netcode.Client) {
	fmt.Println("Client connected!")
	//We probably want this to happen later, once we're ready to actually create the player.
	// Client connects -> we wait for name message -> then we create player
	// To "wait", we launch a go routine running `clientEventLoop`
	// smth like go clientEventLoop(client)
	go g.clientEventLoop(client)
}

func (g *Game) RegisterPlayer(client *netcode.Client, name string) {
	fmt.Printf("Registering new player with name %s!", name)
	//Register event loop management class which has a websocket and a player
	//any message received on the websocket acts upon the player

	//We probably want this to happen later, once we're ready to actually create the player.
	// Client connects -> we wait for name message -> then we create player
	player := NewPlayer(g.state, name)

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
		// fmt.Println("Updating players")

		//Transmit updated state to players
		g.broadcastRegularGameDefinition()

		t := time.Now()
		elapsed := t.Sub(start)

		remainingTime := g.tickDuration - elapsed
		// fmt.Printf("Going to sleep for %dms\n", remainingTime.Milliseconds())
		//elapsed.Milliseconds()
		lastFrame = t
		time.Sleep(remainingTime)
	}
}

func (g *Game) broadcastLobby() {
	g.broadcastMessage(NewLobby(g.players))
}

func (g *Game) broadcastMessage(message events.ClientMessage) {
	// fmt.Printf("Input: %#v\n", message.GetPayload())
	msg := events.ConstructClientMessage(message)
	for client := range g.players {
		client.WriteMessage(msg)
	}
}

func (g *Game) broadcastInitialGameDefinition() {
	g.broadcastMessage(NewInitialGameDef(g.state))
}

func (g *Game) broadcastRegularGameDefinition() {
	g.broadcastMessage(NewRegularGameStateDef(g.state))
}

func (g *Game) isEveryoneReady() bool {
	fmt.Println("is everyone ready?")
	for _, player := range g.players {
		if !player.ready {
			fmt.Println("No!")
			return false
		}
	}
	fmt.Println("Yes!")
	return true
}

func (g *Game) delayedRoundStart(delay time.Duration) {
	go func() {
		//another thread (goroutine)
		time.Sleep(delay)
		g.startRound()
	}()
}

func (g *Game) startRound() {
	roundStart := events.ClientGameEvent{
		SubType: events.TypeGameEventRoundStart,
	}
	g.broadcastMessage(roundStart)
	g.running = true
	go g.GameLoop()
}

func (g *Game) clientEventLoop(client *netcode.Client) {
	for event := range client.GetEventsChannel() {
		// fmt.Printf("New event of type %s\n", event.GetType())
		switch event.GetType() {
		case events.TypePlayerRegistrationEvent:
			ev, ok := event.(events.PlayerRegistrationEvent)
			if !ok {
				fmt.Println("Failed to assert type of event")
				os.Exit(1)
			}
			g.RegisterPlayer(client, ev.Name)
			g.broadcastLobby()
		case events.TypeLobbyEvent:
			ev, ok := event.(events.LobbyEvent)
			if !ok {
				fmt.Println("Failed to assert type of event")
				os.Exit(1)
			}
			switch ev.SubType {
			case events.TypeLobbyEventReady:
				g.players[client].ready = true
				if g.isEveryoneReady() {
					g.broadcastInitialGameDefinition()
					g.delayedRoundStart(3 * time.Second)
				}
			case events.TypeLobbyEventUnready:
				g.players[client].ready = false
			case events.TypeLobbyEventUnregister:
				delete(g.players, client)
			}
			g.broadcastLobby()
		case events.TypeKeyboardUpdateEvent:
			ev, ok := event.(events.KeyboardUpdateEvent)
			if !ok {
				fmt.Println("Failed to assert type of event")
				os.Exit(1)
			}
			g.players[client].pressedKeys.ArrowDown = ev.ArrowDown
			g.players[client].pressedKeys.ArrowUp = ev.ArrowUp
			g.players[client].pressedKeys.ArrowLeft = ev.ArrowLeft
			g.players[client].pressedKeys.ArrowRight = ev.ArrowRight
		}
	}
}

//A, B (B: A)
//A (hint: I'm B)
//ok, I want to consider c(A) as c(B)
//c(*B)
//c(B) --> B != *B

//ok, I want to consider c(*A) as c(*B)

//a(A)
// b, ok := a.(B)
