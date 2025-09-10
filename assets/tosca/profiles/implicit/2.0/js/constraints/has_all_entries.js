// TOSCA 2.0 operator: has_all_entries
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const container = parsed.val1;
    const requiredEntries = parsed.val2;
    
    // Validate arguments
    if (container === undefined || container === null) {
        return false;
    }
    
    if (requiredEntries === undefined || requiredEntries === null) {
        return false;
    }
    
    // Second argument must be a list
    if (!Array.isArray(requiredEntries)) {
        return false;
    }
    
    // Empty required entries list is always satisfied
    if (requiredEntries.length === 0) {
        return true;
    }
    
    // Handle list case
    if (Array.isArray(container)) {
        // Check if ALL entries in requiredEntries exist in the container list
        for (let i = 0; i < requiredEntries.length; i++) {
            let found = false;
            for (let j = 0; j < container.length; j++) {
                if (tosca.deepEqual(container[j], requiredEntries[i])) {
                    found = true;
                    break;
                }
            }
            if (!found) {
                return false;
            }
        }
        return true;
    }
    
    // Handle map case
    if (typeof container === 'object' && container !== null && !Array.isArray(container)) {
        // Check if ALL entries in requiredEntries exist as values in the container map
        for (let i = 0; i < requiredEntries.length; i++) {
            let found = false;
            for (let key in container) {
                if (container.hasOwnProperty(key)) {
                    if (tosca.deepEqual(container[key], requiredEntries[i])) {
                        found = true;
                        break;
                    }
                }
            }
            if (!found) {
                return false;
            }
        }
        return true;
    }
    
    // Invalid container type
    return false;
};