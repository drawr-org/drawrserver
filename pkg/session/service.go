package session

import (
	"log"
	"os"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/oklog/ulid"
)

var svc service

// Service holds the database client for the session service
type service struct {
	db  *bolt.Client
	id  ulid.ULID
	log *log.Logger
}

// Session holds the session information
type Session struct {
	ID    string `json:"id"`
	Users []User `json:"users"`
}

// User holds the user information
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsEditor  bool   `json:"isEditor"`
	IsManager bool   `json:"isManager"`
}

// Init takes a database client and returns a Service
func Init(client *bolt.Client) {
	svc.db = client
	svc.log = log.New(os.Stderr, "[svc]\t", log.LstdFlags)
	return
}
