// TOSCA 2.0 constraint: has_entry
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const container = parsed.val1;
    const valueToFind = parsed.val2;
    
    if (container === undefined || container === null) {
        return false;
    }
    
    // Handle list/array case
    if (Array.isArray(container)) {
        for (let i = 0; i < container.length; i++) {
            const comparable1 = tosca.getComparable(container[i]);
            const comparable2 = tosca.getComparable(valueToFind);
            if (comparable1 === comparable2) {
                return true;
            }
        }
        return false;
    }
    
    // Handle map/object case
    if (typeof container === 'object' && container !== null) {
        // Check if any value in the map matches
        for (const key in container) {
            if (container.hasOwnProperty(key)) {
                const comparable1 = tosca.getComparable(container[key]);
                const comparable2 = tosca.getComparable(valueToFind);
                if (comparable1 === comparable2) {
                    return true;
                }
            }
        }
        return false;
    }
    
    // Invalid container type
    return false;
};