package random

import (
	"math/rand/v2"
)

func NewRandomString(length int) string {
	chars := []rune("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789")

	alias := make([]rune, length)

	for i := range alias {
		alias[i] = chars[rand.IntN(len(chars))]
	}

	return string(alias)
}
