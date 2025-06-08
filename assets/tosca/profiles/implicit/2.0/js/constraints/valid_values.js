// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: valid_values
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // Handle both old (3+ args) and new (2+ args) calling conventions
    let validValuesArray;
    if (arguments.length === 3) {
        validValuesArray = arguments[2];
    } else if (arguments.length === 2) {
        validValuesArray = arguments[1];
    } else {
        return false;
    }
    
    if (!Array.isArray(validValuesArray)) {
        return false;
    }
    
    // Check if currentPropertyValue matches any of the valid values
    for (const validValue of validValuesArray) {
        // If validValue is a string and currentPropertyValue is a scalar, parse the string
        if (typeof validValue === 'string' && currentPropertyValue && 
            currentPropertyValue.$originalString !== undefined) {
            const parsedValidValue = tosca.tryParseScalar(validValue, currentPropertyValue);
            if (parsedValidValue) {
                // For scalars, compare canonical values ($number)
                const currentComparable = tosca.getComparable(currentPropertyValue);
                const validComparable = tosca.getComparable(parsedValidValue);
                if (currentComparable === validComparable) {
                    return true;
                }
            }
        }
        
        // Direct comparison (for non-scalar types or if parsing fails)
        if (tosca.deepEqual(currentPropertyValue, validValue)) {
            return true;
        }
    }
    
    return false;
};
