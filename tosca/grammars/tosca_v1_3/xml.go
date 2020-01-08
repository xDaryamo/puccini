package tosca_v1_3

import (
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
)

//
// XML
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 5.3.4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 5.3.4
//

// tosca.Reader signature
func ReadXml(context *tosca.Context) interface{} {
	if content := context.ReadString(); content != nil {
		if err := format.ValidateXml(*content); err != nil {
			context.ReportValueMalformed("XML", err.Error())
		}
	}
	return context.Data
}
