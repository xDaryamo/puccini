package common

import (
	"math/rand"
	"time"
)

// RFC 3339 format
func Timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func RandomDelay() {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
}
