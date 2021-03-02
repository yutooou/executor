// created by yutooou
package utils

import (
	"math/rand"
	"time"
)

// Get a random string of n bits
func UUID(n int) string {
	chars := [...]rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'g', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'G', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}
	len := len(chars)
	runes := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		runes[i] = chars[rand.Intn(len)]
	}
	return string(runes)
}
