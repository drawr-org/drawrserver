package api

import "encoding/json"

// Session is a drawr session
// It is used to manage the canvas
type Session struct {
	ID        string `json:"id"`
	StartTime string `json:"start_time"`
	Users     []User `json:"users,omitempty"`
}

// NewSession returns the session from provided JSON data
func NewSession(data []byte) (*Session, error) {
	var s Session
	err := json.Unmarshal(data, &s)

	return &s, err
}
