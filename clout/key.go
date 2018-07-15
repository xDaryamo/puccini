package clout

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func NewKey() string {
	return fmt.Sprintf("_%s", ksuid.New().String())
}
