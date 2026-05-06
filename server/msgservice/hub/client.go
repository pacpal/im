package hub

import (
	"IM/server/model"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID    string
	Conn      *websocket.Conn
	Send      chan model.Message
	hub       *Hub
	OnMessage func(model.Message)error
}

const (
	maxMsgSize = 65535
	pongWait   = 60 * time.Second
	pingTime   = 55 * time.Second
	writeWait  = 10 * time.Second
)

func (c *Client) Start() {
	go c.readPump()
	go c.writePump()
}
func (c *Client) readPump() {
	defer func() {
		c.Conn.Close()
		c.hub.Unregister(c)
	}()
	c.Conn.SetReadLimit(maxMsgSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("read error [%s]: %v", c.UserID, err)
			}
			break
		}
		var msg model.Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}
		msg.SdID = c.UserID
		msg.Time = time.Now().UnixMilli()
		if c.OnMessage != nil {
			c.OnMessage(msg) // 交由 service 处理
		}
	}
}
func (c *Client) writePump() {
	ticker := time.NewTicker(pingTime)
	defer func() {
		ticker.Stop()
		c.Conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
