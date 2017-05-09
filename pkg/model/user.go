package model

import "encoding/json"

// User holds the user information
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsEditor  bool   `json:"isEditor"`
	IsManager bool   `json:"isManager"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(raw []byte) error {
	return json.Unmarshal(raw, u)
}
