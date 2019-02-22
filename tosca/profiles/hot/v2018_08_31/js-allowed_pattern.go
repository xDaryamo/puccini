// This file was auto-generated from YAML files

package v2018_08_31

func init() {
	Profile["/hot/2018-08-31/js/allowed_pattern.js"] = `

function validate(v, re) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	if (v.$string !== undefined)
		v = v.$string;
	return new RegExp('^' + re + '$').test(v);
}
`
}
