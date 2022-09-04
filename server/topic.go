package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
)

var __topic_map map[string]*Topic = map[string]*Topic{}
var __topic_mux sync.RWMutex

type Topic struct {
	topic string
	// Registered clients.
	publishers  map[*Client]bool
	subscribers map[*Client]bool

	publisherMux  sync.RWMutex
	subscriberMux sync.RWMutex

	// Inbound messages from the clients.
	broadcast chan *Message
	pick      chan *Message

	// Register requests from the clients.
	registerSubscribe chan *Client
	registerPublish   chan *Client

	// Unregister requests from clients.
	unregisterSubscribe chan *Client
	unregisterPublish   chan *Client
}

//GetTopic return the existed Topic object for a topic title, or create new one if not existed.
func GetTopic(topic string) *Topic {

	__topic_mux.RLock()
	if oldTopic, ok := __topic_map[topic]; ok {
		__topic_mux.RUnlock()
		return oldTopic
	}
	__topic_mux.RUnlock()
	topicHub := &Topic{
		topic:               topic,
		broadcast:           make(chan *Message),
		pick:                make(chan *Message),
		registerSubscribe:   make(chan *Client),
		registerPublish:     make(chan *Client),
		unregisterSubscribe: make(chan *Client),
		unregisterPublish:   make(chan *Client),
		subscribers:         make(map[*Client]bool),
		publishers:          make(map[*Client]bool),
	}
	__topic_mux.Lock()
	__topic_map[topic] = topicHub
	__topic_mux.Unlock()

	go topicHub.Run()
	return topicHub
}

//registerSubscribe register a client as subcriber on a topic
func registerSubscribe(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.registerSubscribe <- client
}

//registerPublish register a client as publisher of a topic
func registerPublish(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.registerPublish <- client
}

//unregisterSubscribe unregister a client from subscriber list of topic
func unregisterSubscribe(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.unregisterSubscribe <- client
}

//unregisterPublish unregister a client from publisher list of topic
func unregisterPublish(topic string, client *Client) {
	topicHub := GetTopic(topic)
	topicHub.unregisterPublish <- client
}

func close(client *Client) {
	client.live = false
	for _, hub := range __topic_map {

		delete(hub.subscribers, client)
		delete(hub.publishers, client)
	}
}

func (h *Topic) Run() {
	for {
		select {
		case client := <-h.registerSubscribe:
			h.subscriberMux.Lock()
			h.subscribers[client] = true
			h.subscriberMux.Unlock()

		case client := <-h.registerPublish:
			h.publisherMux.Lock()
			h.publishers[client] = true
			h.publisherMux.Unlock()
		case client := <-h.unregisterSubscribe:

			if _, ok := h.subscribers[client]; ok {
				h.subscriberMux.Lock()
				delete(h.subscribers, client)
				h.subscriberMux.Unlock()
			}
		case client := <-h.unregisterPublish:

			if _, ok := h.publishers[client]; ok {
				h.publisherMux.Lock()
				delete(h.publishers, client)
				h.publisherMux.Unlock()
			}
		case message := <-h.broadcast:
			fmt.Println("broadcast to:", len(h.subscribers), "member")
			data, err := json.Marshal(message)
			if err == nil {
				if message.NotMe {
					h.subscriberMux.RLock()
					for client := range h.subscribers {
						if client == message.client {
							continue
						}
						select {
						case client.send <- data:
						default:
							client.live = false
							h.unregisterSubscribe <- client
						}
					}
					h.subscriberMux.RUnlock()
				} else {
					h.subscriberMux.RLock()
					for client := range h.subscribers {

						select {
						case client.send <- data:
						default:
							client.live = false
							h.unregisterSubscribe <- client
						}
					}
					h.subscriberMux.RUnlock()
				}
			}
		case pick := <-h.pick:

			if pick.Tier == 1 {

				ObserveTier1(pick)
			}
			data, err := json.Marshal(pick)
			if err == nil {

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
						if pick.Tier == 2 {
							ObserveTier2(pick, client)
						}
						select {
						case client.send <- data:
						default:
							client.live = false
							h.unregisterSubscribe <- client
						}
						break
					}
				}
			}
		}
	}
}
