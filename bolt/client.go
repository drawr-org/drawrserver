// Package bolt provides a client interface to one bolt database on disk
package bolt

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// Client is a client to the bolt DB
type Client struct {
	Path    string
	Timeout time.Duration

	db *bolt.DB
}

const (
	// DBPath is the location of the bolt DB
	DBPath = "data.db"
	// UserBucket is the name of the bolt bucket that stores the users
	UserBucket = "users"
	// SessionBucket is the name of the bolt bucket that stores sessions
	SessionBucket = "sessions"
)

var (
	// ErrExists means the value was found in the DB
	ErrExists = errors.New("data already exists")
	// ErrNotFound means the requested key could not be found in DB
	ErrNotFound = errors.New("key not found")
)

// NewClient sets up Client
func NewClient() Client {
	return Client{Path: DBPath, Timeout: 1 * time.Second}
}

// Open opens the DB
func (c Client) Open() error {
	// Open the database
	db, err := bolt.Open(c.Path, 0666, &bolt.Options{Timeout: c.Timeout})
	if err != nil {
		return err
	}
	c.db = db

	// Initialize the major buckets
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte(UserBucket)); err != nil {
		return err
	}
	if _, err := tx.CreateBucketIfNotExists([]byte(SessionBucket)); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes the DB
func (c Client) Close() {
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			panic(err)
		}
	}
	log.Println("[bolt] closed database:", c.Path)
}

// Get takes a bucket name and a key
// and returns the value in the DB
func (c Client) Get(bucket, key string) ([]byte, error) {
	// Open a read-only connection
	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(bucket))
	if err != nil {
		return nil, err
	}

	data := b.Get([]byte(key))
	if data == nil {
		log.Println("[bolt] not found:", key)
		return nil, ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return data, nil
}

// Put takes a bucket name, a key and a value
// It stores the value in the bucket
func (c Client) Put(bucket, key string, value interface{}) error {
	// Open a write connection
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(bucket))
	if err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := b.Put([]byte(key), data); err != nil {
		return err
	}
	log.Println("[bolt] stored at key:", key)

	return tx.Commit()
}
