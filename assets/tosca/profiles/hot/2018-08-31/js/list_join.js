
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#list_join]

function evaluate() {
	length = arguments.length;
	if (length < 1)
		throw 'must have at least 1 arguments';
	a = [];
	for (var i = 1; i < length; i++) {
		argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join(arguments[0]);
}
