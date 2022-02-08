// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/constraints/length.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function(v, length) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	return tosca.getLength(v) == length;
};
`
}
