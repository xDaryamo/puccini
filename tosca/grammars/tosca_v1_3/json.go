package tosca_v1_3

import (
	"github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/tosca"
)

//
// JSON
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 5.3.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 5.3.2
//

// TODO: implement as constraint, so that we can validate function call results

// tosca.Reader signature
func ReadJSON(context *tosca.Context) interface{} {
	if content := context.ReadString(); content != nil {
		if err := format.ValidateJSON(*content); err != nil {
			context.ReportValueMalformed("JSON", err.Error())
		}
	}
	return context.Data
}
