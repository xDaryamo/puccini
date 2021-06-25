// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/functions/concat.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.3.1

exports.evaluate = function() {
	let a = [];
	let length = arguments.length;
	for (let i = 0; i < length; i++) {
		let argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join('');
};
`
}
