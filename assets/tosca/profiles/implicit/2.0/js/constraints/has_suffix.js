// TOSCA 2.0 operator: has_suffix
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const valueToTest = parsed.val1;
    const suffix = parsed.val2;
    
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