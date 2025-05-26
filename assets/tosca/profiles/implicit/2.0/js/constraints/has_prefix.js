// TOSCA 2.0 operator: has_prefix
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let valueToTest, prefix;
    
    if (arguments.length === 2) {
        // Simple case: value and prefix
        valueToTest = arguments[0];
        prefix = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        valueToTest = arguments[arguments.length - 2];
        prefix = arguments[arguments.length - 1];
    } else {
        throw new Error("has_prefix requires at least 2 arguments");
    }
    
    // Validate arguments
    if (valueToTest === undefined || valueToTest === null) {
        return false;
    }
    
    if (prefix === undefined || prefix === null) {
        return false;
    }
    
    // Both arguments must be of the same type (string or list)
    if (typeof valueToTest !== typeof prefix) {
        return false;
    }
    
    // Handle string case
    if (typeof valueToTest === 'string' && typeof prefix === 'string') {
        return valueToTest.startsWith(prefix);
    }
    
    // Handle list case
    if (Array.isArray(valueToTest) && Array.isArray(prefix)) {
        // Check if prefix list is longer than the value list
        if (prefix.length > valueToTest.length) {
            return false;
        }
        
        // Check if the first elements of valueToTest match prefix
        for (let i = 0; i < prefix.length; i++) {
            if (!tosca.deepEqual(valueToTest[i], prefix[i])) {
                return false;
            }
        }
        return true;
    }
    
    // Invalid types
    return false;
};