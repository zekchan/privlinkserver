package keyGenerators

import (
	"math/rand"
)

const DICTIONARY = "1234567890qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM"

func RandomKeyGenerator() string {
	b := make([]byte, 4)
	dicLen := len(DICTIONARY)
	for i := range b {
		b[i] = DICTIONARY[rand.Intn(dicLen)]
	}

	return string(b)
}
