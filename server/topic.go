package server

var __topic_map map[string]*Topic = map[string]*Topic{}

type Topic struct {
	// Registered clients.
	publishers  map[*Client]bool
	subscribers map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func GetTopic(path string) *Topic {
	if oldTopic, ok := __topic_map[path]; ok {
		return oldTopic
	}

	topic := &Topic{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		subscribers: make(map[*Client]bool),
		publishers:  make(map[*Client]bool),
	}
	__topic_map[path] = topic
	topic.Run()
	return topic
}

func (h *Topic) Run() {
	for {
		select {
		case client := <-h.register:
			if client.publisher {
				h.publishers[client] = true
			} else {
				h.subscribers[client] = true
			}

		case client := <-h.unregister:
			if _, ok := h.subscribers[client]; ok {
				delete(h.subscribers, client)
				//close(client.send)
			}
			if _, ok := h.publishers[client]; ok {
				delete(h.publishers, client)
				//close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.subscribers {
				select {
				case client.send <- message:
				default:
					//close(client.send)
					delete(h.subscribers, client)
				}
			}
		}
	}
}
