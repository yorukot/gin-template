package encryption

import (
	"crypto/rand"
	"math/big"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+-=[]{}|;:,.<>?")
var urlSafeRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_.~")

// Generate random string this can provide 88^n possible
func RandStringRunes(n int, urlSafe bool) (string, error) {
	b := make([]rune, n)
	for i := range b {
		runesLen := 0
		if urlSafe {
			runesLen = len(urlSafeRunes)
		} else {
			runesLen = len(letterRunes)
		}
		index, err := rand.Int(rand.Reader, big.NewInt(int64(runesLen)))
		if err != nil {
			return "", err
		}
		if urlSafe {
			b[i] = urlSafeRunes[index.Int64()]
		} else {
			b[i] = letterRunes[index.Int64()]
		}
	}
	return string(b), nil
}
