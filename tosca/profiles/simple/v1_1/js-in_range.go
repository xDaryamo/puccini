// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/in_range.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2

clout.exec('tosca.helpers');

function validate(v, lower, upper) {
	if (arguments.length !== 3)
		throw 'must have 2 arguments';
	v = tosca.getComparable(v);
	return (v >= tosca.getComparable(lower)) && (v <= tosca.getComparable(upper));
}
`
}
