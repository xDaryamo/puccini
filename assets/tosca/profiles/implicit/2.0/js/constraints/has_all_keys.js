// TOSCA 2.0 operator: has_all_keys
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let mapToTest, requiredKeys;
    
    if (arguments.length === 2) {
        // Simple case: map and required keys list
        mapToTest = arguments[0];
        requiredKeys = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        mapToTest = arguments[arguments.length - 2];
        requiredKeys = arguments[arguments.length - 1];
    } else {
        throw new Error("has_all_keys requires at least 2 arguments");
    }
    
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
        let found = false;
        for (let key in mapToTest) {
            if (mapToTest.hasOwnProperty(key)) {
                if (tosca.deepEqual(key, requiredKeys[i])) {
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
};