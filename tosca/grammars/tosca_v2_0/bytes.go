package tosca_v2_0

import (
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
)

//
// Bytes
//

// tosca.Reader signature
func ReadBytes(context *tosca.Context) tosca.EntityPtr {
	var bytes []byte

	if b64 := context.ReadString(); b64 != nil {
		var err error
		if bytes, err = util.FromBase64(*b64); err != nil {
			context.ReportValueMalformed("bytes", "invalid base64")
		}
	}

	return bytes
}
