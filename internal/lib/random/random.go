package random

import (
	"math/rand"
)

func NewRandomString(size int) string {
	chars := []rune("qwertyuiopasdfghjklzxcvbnmMNBVCXZASDFGHJKLPOIUYTREWQ0192837465")
	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
