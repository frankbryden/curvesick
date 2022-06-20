package main

import (
	"fmt"
	"net/http"
	"server/server/game"
)

func main() {
	fmt.Println("hello world!")
	// s := netcode.NewServer(func(c *netcode.Client) { fmt.Println("Hey there, new client!") })
	g := game.NewGame()
	http.HandleFunc("/", g.GetServer().ServeHome)
	http.HandleFunc("/ws", g.GetServer().ServeWs)

	http.ListenAndServe("0.0.0.0:8090", nil)
}
