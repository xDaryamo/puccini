// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: valid_values
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    if (arguments.length === 3) {
        // TOSCA 2.0 syntax: $valid_values: [ <value_to_test>, <list_of_valid_values> ]
        // Called as: valid_values(currentPropertyValue, valueToTest, validValuesArray)
        let valueToTest = arguments[1];
        const validValuesArray = arguments[2];
        
        // Handle "$value" substitution
        if (valueToTest === '$value') {
            valueToTest = currentPropertyValue;
        }
        
        if (!Array.isArray(validValuesArray)) {
            // Fallback: if third argument is not an array, treat all arguments from 1 onwards as valid values
            const validValuesArrayFromArgs = Array.prototype.slice.call(arguments, 1);
            return exports.checkValueInList(currentPropertyValue, validValuesArrayFromArgs);
        }
        
        return exports.checkValueInList(valueToTest, validValuesArray);
        
    } else if (arguments.length === 2) {
        // Legacy syntax: $valid_values: <list_of_valid_values>
        // Called as: valid_values(currentPropertyValue, validValuesArray)
        const validValuesArray = arguments[1];
        
        if (!Array.isArray(validValuesArray)) {
            return false;
        }
        
        return exports.checkValueInList(currentPropertyValue, validValuesArray);
        
    } else if (arguments.length > 3) {
        // Alternative: $valid_values: [ val1, val2, val3, ... ]
        // Called as: valid_values(currentPropertyValue, val1, val2, val3, ...)
        const validValuesArray = Array.prototype.slice.call(arguments, 1);
        return exports.checkValueInList(currentPropertyValue, validValuesArray);
        
    } else {
        return false;
    }
};

// Helper function to check if a value is in a list, with scalar parsing support
exports.checkValueInList = function(valueToTest, validValuesArray) {
    for (const validValue of validValuesArray) {
        if (exports.compareValues(valueToTest, validValue)) {
            return true;
        }
    }
    return false;
};

// Enhanced comparison that handles scalar parsing
exports.compareValues = function(val1, val2) {
    // If both are scalar objects, use canonical comparison
    if (val1 && val1.$number !== undefined && val2 && val2.$number !== undefined) {
        return tosca.compare(val1, val2) === 0;
    }
    
    // If val1 is scalar and val2 is string, try to parse val2 as scalar
    if (val1 && val1.$number !== undefined && typeof val2 === 'string') {
        const parsedVal2 = tosca.tryParseScalar(val2, val1);
        if (parsedVal2) {
            return tosca.compare(val1, parsedVal2) === 0;
        }
    }
    
    // If val2 is scalar and val1 is string, try to parse val1 as scalar
    if (val2 && val2.$number !== undefined && typeof val1 === 'string') {
        const parsedVal1 = tosca.tryParseScalar(val1, val2);
        if (parsedVal1) {
            return tosca.compare(parsedVal1, val2) === 0;
        }
    }
    
    // Fallback to standard comparison
    const comparable1 = tosca.getComparable(val1);
    const comparable2 = tosca.getComparable(val2);
    return comparable1 === comparable2;
};
