// Package bolt provides a client interface to one bolt database on disk
package bolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const (
	// DBPath is the location of the bolt DB
	DBPath = "data.db"
	// UserBucket is the name of the bolt bucket that stores the users
	UserBucket = "users"
	// SessionBucket is the name of the bolt bucket that stores sessions
	SessionBucket = "sessions"
)

var (
	ErrExists   = errors.New("Data exists already")
	ErrNotFound = errors.New("Key not found")
	ErrReadOnly = errors.New("Database is read-only")
)

// Client is a client to the bolt DB
type Client struct {
	Path    string
	Timeout time.Duration

	dataBuckets [][]byte
	log         *log.Logger
	db          *bolt.DB
}

// NewClient sets up Client
func NewClient(opts *Options) *Client {
	c := new(Client)
	c.log = log.New(ioutil.Discard, "[db]\t", log.LstdFlags)
	if opts != nil {
		c.Path = opts.Path
		c.Timeout = time.Duration(opts.Timeout)
		if opts.Verbose {
			c.log.SetOutput(os.Stderr)
		}
	} else {
		c.Path = DBPath
		c.Timeout = 3 * time.Second
	}
	return c
}

// Open opens the DB
func (c *Client) Open() error {
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

	c.log.Println("opened database at", c.Path)

	return tx.Commit()
}

// Close closes the DB
func (c *Client) Close() {
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			panic(err)
		}
	}
	c.log.Println("closed database:", c.Path)
}

// Get takes a bucket name and a key
// and returns the value in the DB
func (c *Client) Get(bucket, key string) ([]byte, error) {
	// Open a read-only connection
	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(bucket))
	if err != nil {
		return nil, err
	}

	c.log.Printf("get from: %v::%v\n", bucket, key)

	data := b.Get([]byte(key))
	if data == nil {
		log.Println(ErrNotFound.Error(), "->", key)
		return nil, ErrNotFound
	}

	return data, nil
}

// Put takes a bucket name, a key and a value
// It stores the value in the bucket
func (c *Client) Put(bucket, key string, value interface{}) error {
	if c.db.IsReadOnly() {
		return ErrReadOnly
	}

	// Open a write connection
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// b, err := tx.CreateBucketIfNotExists([]byte(bucket))
	b, err := tx.CreateBucket([]byte(bucket))
	if err != nil {
		if err != bolt.ErrBucketExists {
			return err
		}
		b = tx.Bucket([]byte(bucket))
	}
	c.dataBuckets = append(c.dataBuckets, []byte(bucket))

	c.log.Printf("put to: %v::%v\n", bucket, key)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := b.Put([]byte(key), data); err != nil {
		return err
	}

	return tx.Commit()
}

// Remove takes a bucket name and a key
// It removes the value at key from the bucket
func (c *Client) Remove(bucket, key string) error {
	if c.db.IsReadOnly() {
		return ErrReadOnly
	}

	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucket([]byte(bucket))
	if err != nil {
		if err != bolt.ErrBucketExists {
			return err
		}
		b = tx.Bucket([]byte(bucket))
	}
	c.dataBuckets = append(c.dataBuckets, []byte(bucket))

	c.log.Printf("remove %v from %v\n", key, bucket)

	if err := b.Delete([]byte(key)); err != nil {
		return err
	}

	return tx.Commit()
}

func (c *Client) Update(bucket, key string, value interface{}) error {
	if c.db.IsReadOnly() {
		return ErrReadOnly
	}

	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucket([]byte(bucket))
	if err != nil {
		if err != bolt.ErrBucketExists {
			return err
		}
		b = tx.Bucket([]byte(bucket))
	}
	c.dataBuckets = append(c.dataBuckets, []byte(bucket))

	c.log.Printf("update %v in %v\n", key, bucket)

	// remove old value
	if err := b.Delete([]byte(key)); err != nil {
		return err
	}

	// put new value
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := b.Put([]byte(key), data); err != nil {
		return err
	}

	return tx.Commit()
}

func (c *Client) List(bucket string) ([][]byte, error) {
	// Open a read-only connection
	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(bucket))
	if err != nil {
		return nil, err
	}

	c.log.Printf("list: %v\n", bucket)

	var data [][]byte
	if err := b.ForEach(func(k, v []byte) error {
		c.log.Printf("%+v : %+v\n", string(k), string(v))
		if k == nil {
			return nil
		}
		data = append(data, v)
		return nil
	}); err != nil {
		panic(err)
	}

	return data, nil
}

// Stats returns useful information about the database
func (c *Client) Stats() string {
	txNum := c.db.Stats().TxN
	return fmt.Sprintf("tx_num: %v\n", txNum)
}

// String returns the string representation of the database
func (c *Client) String() string {
	var out string

	var dbReader = func(k, v []byte) error {
		out = out + "> " + string(k) + " => " + string(v) + "\n"
		return nil
	}

	if err := c.db.View(func(tx *bolt.Tx) error {
		var err error

		out = out + "user bucket:\n\n"
		ub := tx.Bucket([]byte(UserBucket))
		err = ub.ForEach(dbReader)

		out = out + "session bucket:\n\n"
		sb := tx.Bucket([]byte(SessionBucket))
		err = sb.ForEach(dbReader)

		return err
	}); err != nil {
		panic(err)
	}

	return out
}
