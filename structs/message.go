package structs

import "encoding/json"

type MessageType string

const (
	MessageTypeHello = "connection_established"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
