package ulidgen

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// GeneratedULID is a unique lexicographically sortable
// unique identifier
type GeneratedULID [16]byte

// String returns ULID as a string
func (u GeneratedULID) String() string {
	return ulid.ULID(u).String()
}

// New returns a new ULID based on `t`
func New(t time.Time) GeneratedULID {
	entrop := rand.New(rand.NewSource(t.UnixNano()))
	return GeneratedULID(ulid.MustNew(ulid.Timestamp(t), entrop))
}

// Now returns a new ULID based on current time
func Now() GeneratedULID {
	t := time.Now()
	return New(t)
}
