// TOSCA 2.0 operator: has_all_keys
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const mapToTest = parsed.val1;
    const requiredKeys = parsed.val2;
    
    // Validate arguments
    if (mapToTest === undefined || mapToTest === null) {
        return false;
    }
    
    if (requiredKeys === undefined || requiredKeys === null) {
        return false;
    }
    
    // First argument must be a map (object)
    if (typeof mapToTest !== 'object' || Array.isArray(mapToTest)) {
        return false;
    }
    
    // Second argument must be a list
    if (!Array.isArray(requiredKeys)) {
        return false;
    }
    
    // Empty required keys list is always satisfied
    if (requiredKeys.length === 0) {
        return true;
    }
    
    // Check if ALL keys in requiredKeys exist in the map
    for (let i = 0; i < requiredKeys.length; i++) {
        if (!mapToTest.hasOwnProperty(requiredKeys[i])) {
            return false;
        }
    }
    
    return true;
};