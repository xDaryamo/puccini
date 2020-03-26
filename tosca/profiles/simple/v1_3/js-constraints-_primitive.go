// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/constraints/_primitive.js"] = `

function validate(v, type) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	if (v === null)
		return true;
	return puccini.validateType(v, type);
}`
}
