package model

import "encoding/json"

// Session holds the session information
type Session struct {
	ID    string   `json:"id"`
	Users []string `json:"users"`
}

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(raw []byte) error {
	return json.Unmarshal(raw, s)
}
