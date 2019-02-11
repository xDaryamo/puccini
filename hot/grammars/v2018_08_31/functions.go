package v2018_08_31

import (
	profile "github.com/tliron/puccini/hot/profiles/v2018_08_31"
)

//
// Built-in functions
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#intrinsic-functions]
//

var FunctionSourceCode = map[string]string{
	"list_concat": profile.Profile["/hot/2018-08-31/js/list_concat.js"],
}
