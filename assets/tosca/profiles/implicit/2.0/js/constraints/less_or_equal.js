// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // Handle both old (3 args) and new (2 args) calling conventions
    let compareValue;
    if (arguments.length === 3) {
        compareValue = arguments[2];
    } else if (arguments.length === 2) {
        compareValue = arguments[1];
    } else {
        return false;
    }
    
    // Parse compareValue if it's a string and we have scalar context
    if (typeof compareValue === 'string' && currentPropertyValue && 
        currentPropertyValue.$originalString !== undefined) {
        const parsedCompareValue = tosca.tryParseScalar(compareValue, currentPropertyValue);
        if (parsedCompareValue) {
            compareValue = parsedCompareValue;
        }
    }
    
    return tosca.compare(currentPropertyValue, compareValue) <= 0;  // FIX: <= instead of >=
};