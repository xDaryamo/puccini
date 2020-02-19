// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/constraints/_format.js"] = `

function validate(v, format) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	puccini.validateFormat(v, format);
	return true;
}`
}
