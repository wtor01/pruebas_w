package event

import "encoding/json"

var (
	EventTypeKey = "type"
)

type Message[P any] struct {
	Type    string `json:"type"`
	Payload P      `json:"payload"`
}

func (m Message[P]) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m Message[P]) GetAttributes() map[string]string {
	attributes := make(map[string]string)
	attributes["type"] = m.Type

	return attributes
}
