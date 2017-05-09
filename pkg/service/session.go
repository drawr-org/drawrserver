package service

import (
	"github.com/drawr-team/drawrserver/pkg/model"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"

	"github.com/boltdb/bolt"
	"github.com/simia-tech/boltx"

	log "github.com/golang/glog"
)

// ListSessions returns all the sessions in the database
func ListSessions() (sl []model.Session, err error) {
	log.Info("list sessions")

	err = svc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(sessionBucketKey)
		_, _, err := boltx.ForEach(b, &model.Session{}, func(k []byte, v interface{}) (boltx.Action, error) {
			sl = append(sl, *v.(*model.Session))
			return boltx.ActionContinue, nil
		})
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	log.V(2).Infof("payload: %+v", sl)
	return
}

// NewSession returns a new Session
func NewSession(in *model.Session) (s model.Session, err error) {
	log.Info("new session")

	if in == nil {
		log.Warning("deprecated api endpoint!")
	}
	if in != nil {
		log.V(2).Infof("put payload: %+v", in)
		s = *in
	}
	s.ID = ulidgen.Now().String()

	err = svc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(sessionBucketKey)
		return boltx.PutModel(b, []byte(s.ID), &s)
	})

	log.V(2).Infof("payload: %+v", s)
	return
}

// GetSession returns the Session with id
func GetSession(id string) (s model.Session, err error) {
	log.Info("get session")

	err = svc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(sessionBucketKey)
		ok, err := boltx.GetModel(b, []byte(id), &s)
		if !ok {
			err = ErrNotFound
		}
		return err
	})

	log.V(2).Infof("payload: %+v", s)
	return
}

// UpdateSession changes the Session with id
func UpdateSession(id string, s model.Session) (err error) {
	log.Info("update session")
	defer log.V(2).Infof("payload: %+v", s)

	return svc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(sessionBucketKey)
		return boltx.PutModel(b, []byte(id), &s)
	})
}

// DeleteSession removes s
func DeleteSession(s model.Session) error {
	log.Info("delete session")
	log.V(2).Infof("payload: %+v", s)

	return boltx.DeleteFromBucket(svc.db, sessionBucketKey, []byte(s.ID))
}

// Join connects to the websocket of s
func Join(s model.Session) error {
	log.Info("join session")
	log.Error("not implemented!")
	log.V(2).Infof("payload: %+v", s)
	return nil
}
