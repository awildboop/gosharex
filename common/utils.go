package common

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func HasAnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func TrimPrefixesRecursive(s string, prefixes ...string) (r string) {
	r = strconv.Quote(s)
	r = strings.Trim(r, "\"")
	for HasAnyPrefix(r, prefixes...) {
		for _, p := range prefixes {
			r = strings.TrimPrefix(r, p)
		}
	}
	r = "\"" + r + "\""
	r, _ = strconv.Unquote(r)
	return
}

func RandomString(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixMicro())
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
