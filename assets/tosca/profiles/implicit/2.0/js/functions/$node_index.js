const tosca = require('tosca.lib.utils');

exports.evaluate = function() {
	// Check if we're in a node template context
	if (!this || !this.site)
		throw '$node_index can only be used in a node template context';
	
	// Get the current node template
	let nodeTemplate = this.site;
	if (!tosca.isNodeTemplate(nodeTemplate))
		throw '$node_index can only be used in a node template context';
	
	// Return the node index for this specific instance
	// For single nodes (count=1), this will be 0
	// For multiple nodes (count>1), this will be 0, 1, 2, etc.
	return nodeTemplate.properties.nodeIndex || 0;
};
