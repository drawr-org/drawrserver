package message

import "encoding/json"

// GenericMessage is a message from/to the clients
type GenericMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// GenericAck is ...
type GenericAck struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// NewMessage takes a message type and an interface
// It generates a GenericMessage type
func NewMessage(messageType string, data interface{}) (*GenericMessage, error) {
	var m = &GenericMessage{
		Type: messageType,
		Data: nil,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	m.Data = b
	return m, nil
}
