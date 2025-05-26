// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: matches
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let stringToTest, pattern;
    
    if (arguments.length === 2) {
        // Simple case: string value and pattern
        stringToTest = arguments[0];
        pattern = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        stringToTest = arguments[arguments.length - 2];
        pattern = arguments[arguments.length - 1];
    } else {
        throw new Error("matches requires at least 2 arguments");
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
        
        // Perform regex matching
        console.log(`Matching string "${stringToTest}" against pattern "${pattern}"`);
        return new RegExp(pattern).test(stringToTest);
    } catch (e) {
        console.error(`Error in regex test: ${e.message}`);
        return false;
    }
};