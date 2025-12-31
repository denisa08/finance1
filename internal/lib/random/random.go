package random

import (
	"time"

	"golang.org/x/exp/rand"
)

var rnd = rand.New(
	rand.NewSource(uint64(time.Now().UnixNano())),
)

func NewRandomString(size int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	b := make([]byte, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}
