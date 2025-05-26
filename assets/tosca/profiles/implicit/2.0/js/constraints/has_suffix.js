// TOSCA 2.0 operator: has_suffix
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let valueToTest, suffix;
    
    if (arguments.length === 2) {
        // Simple case: value and suffix
        valueToTest = arguments[0];
        suffix = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        valueToTest = arguments[arguments.length - 2];
        suffix = arguments[arguments.length - 1];
    } else {
        throw new Error("has_suffix requires at least 2 arguments");
    }
    
    // Validate arguments
    if (valueToTest === undefined || valueToTest === null) {
        return false;
    }
    
    if (suffix === undefined || suffix === null) {
        return false;
    }
    
    // Both arguments must be of the same type (string or list)
    if (typeof valueToTest !== typeof suffix) {
        return false;
    }
    
    // Handle string case
    if (typeof valueToTest === 'string' && typeof suffix === 'string') {
        return valueToTest.endsWith(suffix);
    }
    
    // Handle list case
    if (Array.isArray(valueToTest) && Array.isArray(suffix)) {
        // Check if suffix list is longer than the value list
        if (suffix.length > valueToTest.length) {
            return false;
        }
        
        // Check if the last elements of valueToTest match suffix
        const startIndex = valueToTest.length - suffix.length;
        for (let i = 0; i < suffix.length; i++) {
            if (!tosca.deepEqual(valueToTest[startIndex + i], suffix[i])) {
                return false;
            }
        }
        return true;
    }
    
    // Invalid types
    return false;
};