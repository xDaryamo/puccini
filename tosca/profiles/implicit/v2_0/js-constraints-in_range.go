// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/constraints/in_range.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function(v, lower, upper) {
	if (arguments.length !== 3)
		throw 'must have 2 arguments';
	if ((v.lower !== undefined) && (v.upper !== undefined))
		// Special case: is the range in range?
		return (tosca.compare(v.lower, lower) >= 0) && (tosca.compare(v.upper, upper) <= 0);
	else
		return (tosca.compare(v, lower) >= 0) && (tosca.compare(v, upper) <= 0);
};
`
}
