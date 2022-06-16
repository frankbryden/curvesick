package netcode

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameController interface {
	RegisterPlayer(conn *websocket.Conn)
}

type Server struct {
	controller       *GameController
	registrationFunc func(*Client)
}

// func NewServer(controller *GameController) *Server {
// 	return &Server{
// 		controller: controller,
// 	}
// }
func NewServer(registrationFunc func(*Client)) *Server {
	return &Server{
		registrationFunc: registrationFunc,
	}
}

func (s *Server) ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("HEY"))
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New conn!")
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		fmt.Printf("Headers: %s\n", r.Header)
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	c := NewClient(ws)
	s.registrationFunc(c)
	fmt.Println("Reading forever...")
	go c.ReadForever()

	// //register ws
	// for {
	// 	message := c.ReadMessage()
	// 	fmt.Println("Received message from client: " + message)
	// 	// c.WriteMessage(message)
	// }
}
