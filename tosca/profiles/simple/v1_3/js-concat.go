// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/concat.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.3.1

function evaluate() {
	var a = [];
	var length = arguments.length;
	for (var i = 0; i < length; i++) {
		var argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join('');
}
`
}
