package clout

import (
	"fmt"
	"io"

	"github.com/tliron/kutil/ard"
)

func Read(reader io.Reader, format string) (*Clout, error) {
	if data, _, err := ard.Read(reader, format, false); err == nil {
		if map_, ok := data.(ard.Map); ok {
			if clout, err := Parse(map_); err == nil {
				if err := clout.Resolve(); err == nil {
					return clout, nil
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("not a map: %T", data)
		}
	} else {
		return nil, err
	}
}
