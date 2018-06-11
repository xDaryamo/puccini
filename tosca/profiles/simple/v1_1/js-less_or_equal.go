// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/less_or_equal.js"] = `

clout.exec('tosca.helpers');

function validate(v1, v2) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	v1 = tosca.getComparable(v1);
	v2 = tosca.getComparable(v2);
	return v1 <= v2;
}
`
}
