// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/common/1.0/js/outputs.js"] = `

const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

traversal.coerce();

if (tosca.isTosca(clout))
    puccini.write(clout.properties.tosca.outputs);
`
}
