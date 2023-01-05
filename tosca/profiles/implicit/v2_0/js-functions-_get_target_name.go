// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/functions/_get_target_name.js"] = `

const tosca = require('tosca.lib.utils');

exports.evaluate = function() {
	if (!this || !this.target)
		throw 'TARGET cannot be used in this context';
	if (!tosca.isNodeTemplate(this.target))
		throw 'TARGET is not a node template';
	return this.target.properties.name;
};
`
}
