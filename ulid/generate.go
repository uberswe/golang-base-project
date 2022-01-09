package ulid

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

// Generate a new ULID string
func Generate() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	res := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return res.String()
}
