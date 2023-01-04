
const tosca = require('tosca.lib.utils');

exports.evaluate = function(name) {
	return tosca.getModelableEntity.call(this, name).properties.name;
};
