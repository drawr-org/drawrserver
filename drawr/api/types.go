package api

import "net"

type Session struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Users []User `json:"users"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	IP   net.IP `json:"ip_address"`
}
