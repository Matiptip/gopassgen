package password

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Character sets
const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	symbols   = "!@#$%^&*()-_=+[]{}<>?/|~"
	ambiguous = "O0Il"
)

func Random(length int, useSymbols bool, noAmbiguous bool) (string, error) {
	if length <= 0 {
		return "", errors.New("la longitud debe ser mayor a 0")
	}

	// Construir el alfabeto dinámicamente
	alphabet := letters + digits
	if useSymbols {
		alphabet += symbols
	}
	if noAmbiguous {
		filtered := make([]rune, 0, len(alphabet))
	outer:
		for _, r := range alphabet {
			for _, bad := range ambiguous {
				if r == bad {
					continue outer
				}
			}
			filtered = append(filtered, r)
		}
		alphabet = string(filtered)
	}

	out := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))

	for i := range out {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		out[i] = alphabet[n.Int64()]
	}

	return string(out), nil
}
