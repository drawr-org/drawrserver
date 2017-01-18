package ulidgen

import (
	"io"
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

// GeneratorTime returns the arguemnts for a ulid.New() call based on t
func GeneratorTime(t time.Time) (uint64, io.Reader) {
	return ulid.Timestamp(t), rand.New(rand.NewSource(t.UnixNano()))
}

// GeneratorNow returns the arguemnts for a ulid.New() call based on time.Now()
func GeneratorNow() (uint64, io.Reader) {
	t := time.Now()
	return GeneratorTime(t)

}
