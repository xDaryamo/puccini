// TOSCA 2.0 operator: has_any_entry
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const container = parsed.val1;
    const candidateEntries = parsed.val2;
    
    // Validate arguments
    if (container === undefined || container === null) {
        return false;
    }
    
    if (candidateEntries === undefined || candidateEntries === null) {
        return false;
    }
    
    // Second argument must be a list
    if (!Array.isArray(candidateEntries)) {
        return false;
    }
    
    // Empty candidate entries list always returns false
    if (candidateEntries.length === 0) {
        return false;
    }
    
    // Handle list case
    if (Array.isArray(container)) {
        // Check if ANY entry in candidateEntries exists in the container list
        for (let i = 0; i < candidateEntries.length; i++) {
            for (let j = 0; j < container.length; j++) {
                if (tosca.deepEqual(container[j], candidateEntries[i])) {
                    return true;
                }
            }
        }
        return false;
    }
    
    // Handle map case
    if (typeof container === 'object' && container !== null && !Array.isArray(container)) {
        // Check if ANY entry in candidateEntries exists as a value in the container map
        for (let i = 0; i < candidateEntries.length; i++) {
            for (let key in container) {
                if (container.hasOwnProperty(key)) {
                    if (tosca.deepEqual(container[key], candidateEntries[i])) {
                        return true;
                    }
                }
            }
        }
        return false;
    }
    
    // Invalid container type
    return false;
};