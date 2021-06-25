
exports.validate = function(v) {
	let values = Array.prototype.slice.call(arguments, 1);
	return values.indexOf(v) !== -1;
};
