package util

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func GenerateULID() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	res := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return res.String()
}
