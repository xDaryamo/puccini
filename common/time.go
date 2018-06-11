package common

import (
	"math/rand"
	"time"
)

// RFC 3339 format
func Timestamp() (string, error) {
	timestamp, err := time.Now().MarshalText()
	if err != nil {
		return "", err
	}
	return BytesToString(timestamp), nil
}

func Delay() {
	time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
}
