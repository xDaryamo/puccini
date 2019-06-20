// This file was auto-generated from a YAML file

package v1_2

func init() {
	Profile["/tosca/simple/1.2/js/less_than.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2

clout.exec('tosca.helpers');

function validate(v1, v2) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	return tosca.getComparable(v1) < tosca.getComparable(v2);
}
`
}
