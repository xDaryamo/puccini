// TOSCA 2.0 operator: contains
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const valueToTest = parsed.val1;
    const searchValue = parsed.val2;
    
    // Validate arguments
    if (valueToTest === undefined || valueToTest === null) {
        return false;
    }
    
    if (searchValue === undefined || searchValue === null) {
        return false;
    }
    
    // Handle string case
    if (typeof valueToTest === 'string') {
        if (typeof searchValue === 'string') {
            return valueToTest.includes(searchValue);
        } else {
            return false;
        }
    }
    
    // Handle list case
    if (Array.isArray(valueToTest)) {
        if (Array.isArray(searchValue)) {
            // Check if searchValue (sublist) is contained in valueToTest as an uninterrupted sequence
            if (searchValue.length === 0) {
                return true; // Empty array is contained in any array
            }
            
            if (searchValue.length > valueToTest.length) {
                return false; // Can't contain a longer sequence
            }
            
            // Look for the uninterrupted sequence in the array
            for (let i = 0; i <= valueToTest.length - searchValue.length; i++) {
                let sequenceFound = true;
                
                for (let j = 0; j < searchValue.length; j++) {
                    if (!tosca.deepEqual(valueToTest[i + j], searchValue[j])) {
                        sequenceFound = false;
                        break;
                    }
                }
                
                if (sequenceFound) {
                    return true;
                }
            }
            
            return false;
        } else {
            // Check if single value is contained in list
            for (let i = 0; i < valueToTest.length; i++) {
                if (tosca.deepEqual(valueToTest[i], searchValue)) {
                    return true;
                }
            }
            return false;
        }
    }
    
    // Invalid types
    return false;
};