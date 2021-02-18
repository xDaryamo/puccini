package normal

import "github.com/tliron/puccini/tosca"

//
// Location
//

type Location struct {
	Path   string `json:"path" yaml:"path"`
	Row    int    `json:"row" yaml:"row"`
	Column int    `json:"column" yaml:"column"`
}

func NewLocation(path string, row int, column int) *Location {
	return &Location{
		Path:   path,
		Row:    row,
		Column: column,
	}
}

func NewLocationForContext(context *tosca.Context) *Location {
	row, column := context.GetLocation()
	return NewLocation(context.Path.String(), row, column)
}
