package netcode

import (
	"encoding/json"
	"fmt"
	"server/server/events"
	"server/server/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn          *websocket.Conn
	eventConsumer events.EventConsumer
	// clientExitEventListener events.ClientExitEventListener
	clientExitEventListener func()
	eventChannel            chan events.Event
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:         conn,
		eventChannel: make(chan events.Event),
	}
}

func (c *Client) WriteMessage(message []byte) {
	c.conn.WriteMessage(websocket.TextMessage, message)
}

func (c *Client) ReadMessage() string {
	messageType, message, err := c.conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	if messageType == websocket.TextMessage {
		return string(message)
	} else {
		return "Binary message"
	}
}

func (c *Client) SetEventConsumer(consumer events.EventConsumer) {
	c.eventConsumer = consumer
}

func (c *Client) SetClientExitEventListener(eventListener func()) {
	c.clientExitEventListener = eventListener
}

func (c *Client) ReadForever() {
	go func() {
		for {
			_, message, err := c.conn.ReadMessage()

			if err != nil {
				if c.clientExitEventListener != nil {
					c.clientExitEventListener()
					return
				}
			}

			var event events.GenericEvent
			err = json.Unmarshal(message, &event)

			c.eventChannel <- event

			switch event.Type {
			case events.TypePlayerRegistrationEvent:
				var ev events.PlayerRegistrationEvent
				err := json.Unmarshal(event.Data, &ev)
				utils.CheckErr(err)
				c.eventChannel <- ev
			case events.TypeLobbyEvent:
				var ev events.LobbyEvent
				err := json.Unmarshal(event.Data, &ev)
				utils.CheckErr(err)
				c.eventChannel <- ev
			case events.TypeKeyboardUpdateEvent:
				fmt.Println(event)
				var ev events.KeyboardUpdateEvent
				err := json.Unmarshal(event.Data, &ev)
				utils.CheckErr(err)
				c.eventChannel <- ev
			}

			utils.CheckErr(err)

			// if c.eventConsumer != nil {
			// 	c.eventConsumer.OnKeystateChange(state)
			// }
		}
	}()
}

func (c *Client) GetEventsChannel() <-chan events.Event {
	return c.eventChannel
}
