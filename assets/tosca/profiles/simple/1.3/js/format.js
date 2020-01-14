
function validate(v, format) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	puccini.validateFormat(v, format);
	return true;
}