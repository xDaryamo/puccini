
function validate(v, format) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	if (!puccini.isType(v, '!!str'))
		return 'not a string';
	try {
		puccini.validateFormat(v, format);
	} catch (x) {
		if (x.value && x.value.error)
			// Unwrap Go error
			return x.value.error();
		else
			throw x;
	}
	return true;
}