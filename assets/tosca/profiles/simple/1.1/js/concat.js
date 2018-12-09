
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.3.1

function evaluate() {
	a = [];
	length = arguments.length;
	for (var i = 0; i < length; i++) {
		argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join('');
}
