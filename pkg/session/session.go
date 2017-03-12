package session

import (
	"encoding/json"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
)

// New returns a new Session
func New() (s Session, err error) {
	s.ID = ulidgen.Now().String()
	err = svc.db.Put(bolt.SessionBucket, s.ID, s)
	return
}

// Get returns the Session with id
func Get(id string) (s Session, err error) {
	raw, err := svc.db.Get(bolt.SessionBucket, id)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &s)
	return
}

// Join connects to the websocket of s
func Join(s Session) error {
	return nil
}

// Delete removes s
func Delete(s Session) error {
	return svc.db.Remove(bolt.SessionBucket, s.ID)
}
