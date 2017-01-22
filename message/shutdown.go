package message

import "encoding/json"

const (
	ShutdownMessageType = "server-down"
)

type ShutdownMessageData struct {
	Error string `json:"error"`
}

func ShutdownMessage(msg string) []byte {
	m, _ := NewMessage(ShutdownMessageType, ShutdownMessageData{
		Error: msg,
	})

	b, _ := json.Marshal(m)

	return b
}
