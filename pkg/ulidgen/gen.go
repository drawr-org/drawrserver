package ulidgen

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// New returns a new ULID based on `t`
func New(t time.Time) ulid.ULID {
	entrop := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entrop)
}

// Now returns a new ULID based on current time
func Now() ulid.ULID {
	t := time.Now()
	return New(t)
}
