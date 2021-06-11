package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tapvanvn/go-wspubsub/entity"
	"github.com/tapvanvn/go-wspubsub/runtime"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {

	// The websocket connection.
	conn *websocket.Conn

	publishTopics   []string
	subscribeTopics []string
	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) load() {

}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		for _, topic := range c.subscribeTopics {
			unregisterSubscribe(topic, c)
		}
		for _, topic := range c.subscribeTopics {
			unregisterSubscribe(topic, c)
		}
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println(string(message))
		c.processMessage(message)
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
	}
}

//TODO: secure this
func (c *Client) processMessage(message []byte) {
	raw := &entity.Message{}

	fmt.Println("receive:", string(message))

	err := json.Unmarshal(message, &raw)
	if err == nil {

		topic := strings.TrimSpace(raw.Topic)

		if topic == "control" {
			register := &entity.Register{}
			err := json.Unmarshal([]byte(raw.Message), register)
			if err == nil {
				//if register, ok := raw.Message.(*entity.Register); ok {

				if register.IsPublisher {

					if register.IsUnregister {

						go unregisterPublish(register.Topic, c)

					} else {

						go registerPublish(register.Topic, c)
					}
				} else {
					if register.IsUnregister {

						go unregisterSubscribe(register.Topic, c)

					} else {

						go registerSubscribe(register.Topic, c)
					}
				}
			} else {
				fmt.Println("cannot get control")
			}
		} else {

			if len(topic) > 0 {

				topicHub := GetTopic(topic)
				topicHub.broadcast <- message
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {

	fmt.Println("client come:", r.UserAgent(), r.RemoteAddr)

	if !runtime.Ready {

		w.WriteHeader(http.StatusInternalServerError)
		log.Println("server is not ready")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{subscribeTopics: make([]string, 0), publishTopics: make([]string, 0), conn: conn, send: make(chan []byte, 256)}
	client.load()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.

	go client.readPump()
	go client.writePump()

}
