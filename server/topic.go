package server

import (
	"fmt"
	"math/rand"
)

var __topic_map map[string]*Topic = map[string]*Topic{}

type Topic struct {
	topic string
	// Registered clients.
	publishers  map[*Client]bool
	subscribers map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte
	pick      chan []byte

	// Register requests from the clients.
	registerSubscribe chan *Client
	registerPublish   chan *Client

	// Unregister requests from clients.
	unregisterSubscribe chan *Client
	unregisterPublish   chan *Client
}

func GetTopic(topic string) *Topic {
	if oldTopic, ok := __topic_map[topic]; ok {
		return oldTopic
	}

	topicHub := &Topic{
		topic:               topic,
		broadcast:           make(chan []byte),
		pick:                make(chan []byte),
		registerSubscribe:   make(chan *Client),
		registerPublish:     make(chan *Client),
		unregisterSubscribe: make(chan *Client),
		unregisterPublish:   make(chan *Client),
		subscribers:         make(map[*Client]bool),
		publishers:          make(map[*Client]bool),
	}
	__topic_map[topic] = topicHub
	go topicHub.Run()
	return topicHub
}

func registerSubscribe(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.registerSubscribe <- client
}
func registerPublish(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.registerPublish <- client
}
func unregisterSubscribe(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.unregisterSubscribe <- client
}
func unregisterPublish(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.unregisterPublish <- client
}

func close(client *Client) {
	for _, hub := range __topic_map {
		delete(hub.subscribers, client)
		delete(hub.publishers, client)
	}
}
func (h *Topic) Run() {
	for {
		select {
		case client := <-h.registerSubscribe:
			fmt.Println("add subscriber to topic:", h.topic)
			h.subscribers[client] = true

		case client := <-h.registerPublish:
			fmt.Println("add publisher to topic:", h.topic)
			h.publishers[client] = true

		case client := <-h.unregisterSubscribe:
			fmt.Println("remove subscriber to topic:", h.topic)
			if _, ok := h.subscribers[client]; ok {
				delete(h.subscribers, client)

			}
		case client := <-h.unregisterPublish:
			fmt.Println("remove publisher to topic:", h.topic)
			if _, ok := h.publishers[client]; ok {
				delete(h.publishers, client)

			}
		case message := <-h.broadcast:
			fmt.Println("broadcast to:", len(h.subscribers), "member")
			for client := range h.subscribers {
				select {
				case client.send <- message:
				default:
					delete(h.subscribers, client)
				}
			}
		case pick := <-h.pick:
			num := len(h.subscribers)
			if num > 0 {
				choice := rand.Intn(num)
				fmt.Println("sent to:", choice, num)
				i := 0
				for client := range h.subscribers {
					if i != choice {
						i++
						continue
					}
					select {
					case client.send <- pick:
					default:
						delete(h.subscribers, client)
					}
				}
			}
		}
	}
}
