// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/concat.js"] = `

function evaluate() {
	a = [];
	length = arguments.length;
	for (i = 0; i < length; i++) {
		argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join('');
}
`
}
