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
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) WriteMessage(message string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(message))
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

			var state events.KeyState
			err = json.Unmarshal(message, &state)

			utils.CheckErr(err)

			if c.eventConsumer != nil {
				c.eventConsumer.OnKeystateChange(state)
			}
		}
	}()
}
