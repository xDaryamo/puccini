// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: equal
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // Handle both old (3 args) and new (2 args) calling conventions
    let expectedValue;
    if (arguments.length === 3) {
        expectedValue = arguments[2];
    } else if (arguments.length === 2) {
        expectedValue = arguments[1];
    } else {
        return false;
    }
    
    // Parse expectedValue if it's a string and we have scalar context
    if (typeof expectedValue === 'string' && currentPropertyValue && 
        currentPropertyValue.$originalString !== undefined) {
        const parsedExpectedValue = tosca.tryParseScalar(expectedValue, currentPropertyValue);
        if (parsedExpectedValue) {
            expectedValue = parsedExpectedValue;
        }
    }
    
    const currentComparable = tosca.getComparable(currentPropertyValue);
    const expectedComparable = tosca.getComparable(expectedValue);
    
    return currentComparable === expectedComparable;
};
