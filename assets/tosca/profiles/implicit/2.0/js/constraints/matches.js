// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: matches
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // Extract pattern from arguments
    let stringToTest, pattern;
    
    if (arguments.length === 2) {
        // Standard case: currentPropertyValue and pattern
        stringToTest = arguments[0];
        pattern = arguments[1];
    } else if (arguments.length === 3) {
        // Case with expanded arguments: currentPropertyValue, stringToTest, pattern
        stringToTest = arguments[1];
        pattern = arguments[2];
    } else {
        throw new Error("matches requires 2 or 3 arguments");
    }
    
    // Validate arguments
    if (stringToTest === undefined || stringToTest === null) {
        return false;
    }
    
    if (pattern === undefined || pattern === null) {
        return false;
    }
    
    try {
        // Convert both arguments to strings
        stringToTest = String(stringToTest);
        pattern = String(pattern);
        
        // Handle YAML escaping: convert \\d to \d, \\w to \w, etc.
        pattern = pattern.replace(/\\\\([dwsWDSbBnrtfv])/g, '\\$1');
        
        // Perform regex matching
        return new RegExp(pattern).test(stringToTest);
    } catch (e) {
        console.error(`Error in regex test: ${e.message}`);
        return false;
    }
};