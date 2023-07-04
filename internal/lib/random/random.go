package random

import (
	"math/rand"
	"time"
)

func NewRandomString(aliasLength int) string {
	chars := []rune("ABCDEFGIHJKLMNOPQRSTVUWXYZ" +
		"abcdefghijklmnopqrstvuwxyz" +
		"0123456789")

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, aliasLength)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
