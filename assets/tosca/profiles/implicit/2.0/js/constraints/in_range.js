// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: in_range
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // Handle both old (4 args) and new (3 args) calling conventions
    let valueToTest, lowerBound, upperBound;
    
    if (arguments.length === 4) {
        // Old style: arguments[0] = currentPropertyValue, arguments[1] = currentPropertyValue, arguments[2] = lowerBound, arguments[3] = upperBound
        valueToTest = arguments[1];
        lowerBound = arguments[2];
        upperBound = arguments[3];
    } else if (arguments.length === 3) {
        // New style: arguments[0] = currentPropertyValue, arguments[1] = lowerBound, arguments[2] = upperBound
        valueToTest = currentPropertyValue;
        lowerBound = arguments[1];
        upperBound = arguments[2];
    } else {
        return false;
    }

    if (valueToTest === undefined || valueToTest === null ||
        lowerBound === undefined || lowerBound === null ||
        upperBound === undefined || upperBound === null) {
        return false;
    }

    // Parse upperBound if it's a string and we have scalar context
    if (typeof upperBound === 'string' && currentPropertyValue && 
        currentPropertyValue.$originalString !== undefined) {
        const parsedUpperBound = tosca.tryParseScalar(upperBound, currentPropertyValue);
        if (parsedUpperBound) {
            upperBound = parsedUpperBound;
        }
    }

    // Parse lowerBound if it's a string and we have scalar context
    if (typeof lowerBound === 'string' && currentPropertyValue && 
        currentPropertyValue.$originalString !== undefined) {
        const parsedLowerBound = tosca.tryParseScalar(lowerBound, currentPropertyValue);
        if (parsedLowerBound) {
            lowerBound = parsedLowerBound;
        }
    }

    // Use canonical comparison for scalars and other comparable types
    return (tosca.compare(valueToTest, lowerBound) >= 0) && 
           (tosca.compare(valueToTest, upperBound) <= 0);
};
