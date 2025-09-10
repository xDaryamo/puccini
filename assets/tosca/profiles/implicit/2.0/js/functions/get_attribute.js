// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.5.1

const tosca = require('tosca.lib.utils');

exports.evaluate = function(entity, first) {
    // Validate relationship attributes during evaluation
    const args = Array.prototype.slice.call(arguments);
    
    if (args.length >= 4) {
        // This looks like a relationship attribute access: [node, requirement, attribute, index]
        const vertex = tosca.getModelableEntity.call(this, args[0]);
        const requirementName = args[1];
        const attributeName = args[2];
        const relationshipIndex = args[3] || 0;
        
        // Validate that the relationship and attribute exist, and index is valid
        tosca.validateRelationshipAttributeAccess(vertex, requirementName, attributeName, relationshipIndex);
    }
    
    return tosca.getNestedValue.call(this, 'attribute', 'attributes', arguments);
};
