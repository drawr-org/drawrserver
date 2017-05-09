package service

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
	"github.com/drawr-team/drawrserver/pkg/websock"

	log "github.com/golang/glog"
)

var svc service

// Service holds the database client for the session service
type service struct {
	db *bolt.DB

	hubs map[string]*websock.Hub
}

var (
	dbTimeout        = 5 * time.Second
	sessionBucketKey = []byte("sessions")
	userBucketKey    = []byte("users")
)

var (
	ErrNotFound = errors.New("key not in database")
)

// Init takes a database config and returns a Service
func Init(dbPath string, dbTimeout int) error {
	var err error
	svc.db, err = bolt.Open(dbPath, 0666, &bolt.Options{
		Timeout: time.Duration(dbTimeout) * time.Second,
	})
	if err != nil {
		return err
	}

	if err := svc.db.Update(func(tx *bolt.Tx) error {
		log.Info("Create bucket for sessions")
		if _, err := tx.CreateBucketIfNotExists(sessionBucketKey); err != nil {
			log.Error("Error creating bucket sessions:", err)
			return err
		}
		log.Info("Create bucket for users")
		if _, err := tx.CreateBucketIfNotExists(userBucketKey); err != nil {
			log.Error("Error creating bucket users:", err)
			return err
		}
		log.Info("Create bucket for updates")
		return nil
	}); err != nil {
		log.Error("Error creating buckets:", err)
		return err
	}

	svc.hubs = make(map[string]*websock.Hub)
	return nil
}

func Close() {
	log.Warning("closing database")
	if err := svc.db.Close(); err != nil {
		log.Error(err)
	}

	log.Warning("closing websocket hubs")
	for _, h := range svc.hubs {
		h.Close()
	}
}
