// TOSCA 2.0 operator: has_any_key
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const mapToTest = parsed.val1;
    const candidateKeys = parsed.val2;
    
    // Validate arguments
    if (mapToTest === undefined || mapToTest === null) {
        return false;
    }
    
    if (candidateKeys === undefined || candidateKeys === null) {
        return false;
    }
    
    // First argument must be a map (object)
    if (typeof mapToTest !== 'object' || Array.isArray(mapToTest)) {
        return false;
    }
    
    // Second argument must be a list
    if (!Array.isArray(candidateKeys)) {
        return false;
    }
    
    // Empty candidate keys list always returns false
    if (candidateKeys.length === 0) {
        return false;
    }
    
    // Check if ANY key in candidateKeys exists in the map
    for (let i = 0; i < candidateKeys.length; i++) {
        if (mapToTest.hasOwnProperty(candidateKeys[i])) {
            return true;
        }
    }
    
    return false;
};