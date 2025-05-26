// TOSCA 2.0 operator: has_key
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let mapToTest, keyValue;
    
    if (arguments.length === 2) {
        // Simple case: map and key value
        mapToTest = arguments[0];
        keyValue = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        mapToTest = arguments[arguments.length - 2];
        keyValue = arguments[arguments.length - 1];
    } else {
        throw new Error("has_key requires at least 2 arguments");
    }
    
    // Validate arguments
    if (mapToTest === undefined || mapToTest === null) {
        return false;
    }
    
    if (keyValue === undefined || keyValue === null) {
        return false;
    }
    
    // First argument must be a map (object)
    if (typeof mapToTest !== 'object' || Array.isArray(mapToTest)) {
        return false;
    }
    
    // Check if keyValue exists as a key in the map
    for (let key in mapToTest) {
        if (mapToTest.hasOwnProperty(key)) {
            if (tosca.deepEqual(key, keyValue)) {
                return true;
            }
        }
    }
    
    return false;
};