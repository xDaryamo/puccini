
const tosca = require('tosca.lib.utils');

exports.evaluate = function() {
	if (!this || !this.target)
		throw 'TARGET cannot be used in this context';
	if (!tosca.isNodeTemplate(this.target))
		throw 'TARGET is not a node template';
	return this.target.properties.name;
};
