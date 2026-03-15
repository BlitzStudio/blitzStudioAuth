package utils

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateUlid() string {
	ms := ulid.Timestamp(time.Now())
	ulidToken, _ := ulid.New(ms, rand.Reader)
	return ulidToken.String()
}
