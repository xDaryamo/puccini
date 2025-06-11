// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: matches
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const stringToTest = parsed.val1;
    const regexPattern = parsed.val2;
    
    // Validate arguments
    if (stringToTest === undefined || stringToTest === null) {
        return false;
    }
    
    if (regexPattern === undefined || regexPattern === null) {
        return false;
    }
    
    // Both arguments must be strings
    if (typeof stringToTest !== 'string' || typeof regexPattern !== 'string') {
        return false;
    }
    
    try {
        const regex = new RegExp(regexPattern);
        return regex.test(stringToTest);
    } catch (e) {
        // Invalid regex pattern
        return false;
    }
};