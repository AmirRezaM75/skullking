package support

import (
	"math/rand"
	"time"
)

type Str struct{}

func (s Str) Random(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	rand.Seed(time.Now().UnixNano())

	characters := make([]rune, length)

	for i := range characters {
		characters[i] = letters[rand.Intn(len(letters))]
	}

	return string(characters)
}
