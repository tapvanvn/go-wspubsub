package entity

type Config struct {
	MaxMessageSize int64 `json:"max_message_size"`
}

var DefaultConfig *Config = &Config{
	MaxMessageSize: 2048,
}
