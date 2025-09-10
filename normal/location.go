package normal

import (
	"strings"
	
	"github.com/tliron/puccini/tosca/parsing"
)

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

func NewLocationForContext(context *parsing.Context) *Location {
	row, column := context.GetLocation()
	return NewLocation(context.Path.String(), row, column)
}

// UpdateNodeTemplatePath updates the path in a location to reflect the correct node template instance name
func (self *Location) UpdateNodeTemplatePath(originalNodeName, instanceNodeName string) {
	if self != nil && self.Path != "" {
		// Replace the original node template name with the instance name in the path
		// e.g., service_template.node_templates["aws_instance.web"] -> service_template.node_templates["aws_instance.web_0"]
		oldPath := `service_template.node_templates["` + originalNodeName + `"]`
		newPath := `service_template.node_templates["` + instanceNodeName + `"]`
		self.Path = strings.Replace(self.Path, oldPath, newPath, 1)
	}
}
