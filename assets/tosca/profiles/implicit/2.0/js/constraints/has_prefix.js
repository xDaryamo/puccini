// TOSCA 2.0 operator: has_prefix
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const valueToTest = parsed.val1;
    const prefix = parsed.val2;
    
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