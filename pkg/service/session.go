package service

import (
	"encoding/json"
	"log"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
)

// ListSessions returns all the sessions in the database
func ListSessions() (sl []Session, err error) {
	rawl, err := svc.db.List(bolt.SessionBucket)
	if err != nil {
		return
	}
	for _, b := range rawl {
		var s Session
		err = json.Unmarshal(b, &s)
		sl = append(sl, s)
	}
	svc.log.Printf("list: %+v\n", sl)
	return
}

// NewSession returns a new Session
func NewSession(in *Session) (s Session, err error) {
	if in != nil {
		s = *in
	}
	s.ID = ulidgen.Now().String()
	err = svc.db.Put(bolt.SessionBucket, s.ID, s)
	svc.log.Printf("new: %+v\n", s)
	return
}

// GetSession returns the Session with id
func GetSession(id string) (s Session, err error) {
	raw, err := svc.db.Get(bolt.SessionBucket, id)
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &s)
	svc.log.Printf("get: %+v\n", s)
	return
}

// UpdateSession changes the Session with id
func UpdateSession(id string, s Session) (err error) {
	svc.log.Printf("update: %+v\n", s)
	return svc.db.Update(bolt.SessionBucket, id, s)
}

// DeleteSession removes s
func DeleteSession(s Session) error {
	svc.log.Printf("delete: %+v\n", s)
	return svc.db.Remove(bolt.SessionBucket, s.ID)
}

// Join connects to the websocket of s
func Join(s Session) error {
	log.Println("[session] Join Session logic not implemented yet!")
	return nil
}
