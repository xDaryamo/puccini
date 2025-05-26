// TOSCA 2.0 operator: contains
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let valueToTest, containedValue;
    
    if (arguments.length === 2) {
        // Simple case: value and contained value
        valueToTest = arguments[0];
        containedValue = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        valueToTest = arguments[arguments.length - 2];
        containedValue = arguments[arguments.length - 1];
    } else {
        throw new Error("contains requires at least 2 arguments");
    }
    
    // Validate arguments
    if (valueToTest === undefined || valueToTest === null) {
        return false;
    }
    
    if (containedValue === undefined || containedValue === null) {
        return false;
    }
    
    // Handle string case
    if (typeof valueToTest === 'string' && typeof containedValue === 'string') {
        return valueToTest.includes(containedValue);
    }
    
    // Handle list cases
    if (Array.isArray(valueToTest)) {
        // Case 1: List contains another list (subsequence)
        if (Array.isArray(containedValue)) {
            // Empty containedValue is always contained in any array
            if (containedValue.length === 0) {
                return true;
            }
            
            // containedValue can't be longer than the valueToTest
            if (containedValue.length > valueToTest.length) {
                return false;
            }
            
            // Look for the uninterrupted sequence in the array
            for (let i = 0; i <= valueToTest.length - containedValue.length; i++) {
                let sequenceFound = true;
                
                for (let j = 0; j < containedValue.length; j++) {
                    if (!tosca.deepEqual(valueToTest[i + j], containedValue[j])) {
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
            // Case 2: List contains a single value
            for (let i = 0; i < valueToTest.length; i++) {
                if (tosca.deepEqual(valueToTest[i], containedValue)) {
                    return true;
                }
            }
            return false;
        }
    }
    
    // Handle map case
    if (typeof valueToTest === 'object' && valueToTest !== null && !Array.isArray(valueToTest)) {
        // Check if containedValue has key/value structure for map validation
        if (typeof containedValue === 'object' && containedValue !== null && 
            'key' in containedValue && 'value' in containedValue) {
            const key = containedValue.key;
            const value = containedValue.value;
            return key in valueToTest && tosca.deepEqual(valueToTest[key], value);
        }
    }
    
    // Fallback: direct equality check
    return tosca.deepEqual(valueToTest, containedValue);
};