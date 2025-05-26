// TOSCA 2.0 operator: has_any_key
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let mapToTest, candidateKeys;
    
    if (arguments.length === 2) {
        // Simple case: map and candidate keys list
        mapToTest = arguments[0];
        candidateKeys = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        mapToTest = arguments[arguments.length - 2];
        candidateKeys = arguments[arguments.length - 1];
    } else {
        throw new Error("has_any_key requires at least 2 arguments");
    }
    
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
        for (let key in mapToTest) {
            if (mapToTest.hasOwnProperty(key)) {
                if (tosca.deepEqual(key, candidateKeys[i])) {
                    return true;
                }
            }
        }
    }
    
    return false;
};