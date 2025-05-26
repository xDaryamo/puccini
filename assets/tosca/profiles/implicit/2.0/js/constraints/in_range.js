// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: in_range
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // Requires exactly 4 arguments (currentPropertyValue + 3 YAML args)
    if (arguments.length !== 4) {
        return false;
    }

    let valueToTest;
    let lowerBound;
    let upperBound;

    // Resolve first YAML argument (value to test)
    if (arguments[1] === '$value') {
        valueToTest = currentPropertyValue;
    } else {
        valueToTest = arguments[1];
    }

    // Resolve second YAML argument (lower bound)
    if (arguments[2] === '$value') {
        lowerBound = currentPropertyValue;
    } else {
        lowerBound = arguments[2];
    }

    // Resolve third YAML argument (upper bound)
    if (arguments[3] === '$value') {
        upperBound = currentPropertyValue;
    } else {
        upperBound = arguments[3];
    }

    if (valueToTest === undefined || valueToTest === null ||
        lowerBound === undefined || lowerBound === null ||
        upperBound === undefined || upperBound === null) {
        return false;
    }

    // Special case for timestamps
    if (valueToTest.$number !== undefined) {
        // Extract numeric value from timestamp objects if needed
        const valueNum = valueToTest.$number;
        const lowerNum = lowerBound.$number !== undefined ? lowerBound.$number : 
                         (typeof lowerBound === 'string' ? Date.parse(lowerBound) * 1000000 : lowerBound);
        const upperNum = upperBound.$number !== undefined ? upperBound.$number : 
                         (typeof upperBound === 'string' ? Date.parse(upperBound) * 1000000 : upperBound);
        
        return valueNum >= lowerNum && valueNum <= upperNum;
    }
    
    // Special case: if valueToTest is itself a range
    if ((valueToTest.lower !== undefined) && (valueToTest.upper !== undefined)) {
        // Check if the range is within the bounds
        return (tosca.compare(valueToTest.lower, lowerBound) >= 0) && 
               (tosca.compare(valueToTest.upper, upperBound) <= 0);
    } else {
        // Regular case: check if value is within bounds
        return (tosca.compare(valueToTest, lowerBound) >= 0) && 
               (tosca.compare(valueToTest, upperBound) <= 0);
    }
};
