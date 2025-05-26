// TOSCA 2.0 operator: has_any_entry
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let container, candidateEntries;
    
    if (arguments.length === 2) {
        // Simple case: container and candidate entries list
        container = arguments[0];
        candidateEntries = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        container = arguments[arguments.length - 2];
        candidateEntries = arguments[arguments.length - 1];
    } else {
        throw new Error("has_any_entry requires at least 2 arguments");
    }
    
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