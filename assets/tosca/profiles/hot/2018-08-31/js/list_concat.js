
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#list-concat]

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
