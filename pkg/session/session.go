package session

import (
	"encoding/json"
	"log"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
)

// List returns all the sessions in the database
func List() (sl []Session, err error) {
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

// New returns a new Session
func New(in *Session) (s Session, err error) {
	if in != nil {
		s = *in
	}
	s.ID = ulidgen.Now().String()
	err = svc.db.Put(bolt.SessionBucket, s.ID, s)
	svc.log.Printf("new: %+v\n", s)
	return
}

// Get returns the Session with id
func Get(id string) (s Session, err error) {
	raw, err := svc.db.Get(bolt.SessionBucket, id)
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &s)
	svc.log.Printf("get: %+v\n", s)
	return
}

// Update changes the Session with id
func Update(id string, s Session) (err error) {
	svc.log.Printf("update: %+v\n", s)
	return svc.db.Update(bolt.SessionBucket, id, s)
}

// Delete removes s
func Delete(s Session) error {
	svc.log.Printf("delete: %+v\n", s)
	return svc.db.Remove(bolt.SessionBucket, s.ID)
}

// Join connects to the websocket of s
func Join(s Session) error {
	log.Println("[session] Join Session logic not implemented yet!")
	return nil
}
