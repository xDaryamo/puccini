
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#list_join]

exports.evaluate = function() {
	let length = arguments.length;
	if (length < 1)
		throw 'must have at least 1 arguments';
	let a = [];
	for (let i = 1; i < length; i++) {
		let argument = arguments[i];
		if (argument.$string !== undefined)
			argument = argument.$string;
		a.push(argument);
	}
	return a.join(arguments[0]);
};
