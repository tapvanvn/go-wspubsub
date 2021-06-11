package entity

type Message struct {
	Topic   string `json:"topic"`
	Message string `json:"message,omitempty"`
}
