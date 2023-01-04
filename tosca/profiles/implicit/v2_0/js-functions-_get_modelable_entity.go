// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/functions/_get_modelable_entity.js"] = `

const tosca = require('tosca.lib.utils');

exports.evaluate = function(name) {
	return tosca.getModelableEntity.call(this, name).properties.name;
};
`
}
