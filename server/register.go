package server

type Register struct {
	IsUnregister bool   `json:"is_unregister"`
	Topic        string `json:"topic"`
	IsPublisher  bool   `json:"is_publisher"`
}
