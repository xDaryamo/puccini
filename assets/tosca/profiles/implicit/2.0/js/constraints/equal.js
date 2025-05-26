// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: equal
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    if (arguments.length < 2) {
        throw new Error("equal requires at least 2 arguments");
    }
    
    // Extract the actual values we need to compare
    let v1, v2;
    
    if (arguments.length === 2) {
        // Simple case: property value and constraint value
        v1 = arguments[0];
        v2 = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        v1 = arguments[arguments.length - 2];
        v2 = arguments[arguments.length - 1];
    }
    
    if (v1 === undefined || v2 === undefined) {
        return false;
    }

    // Handle arrays and objects with deepEqual instead of compare
    if ((Array.isArray(v1) && Array.isArray(v2)) || 
        (typeof v1 === 'object' && v1 !== null && 
         typeof v2 === 'object' && v2 !== null)) {
        return tosca.deepEqual(v1, v2);
    }
    
    // For primitive types, use the standard compare function
    return tosca.compare(v1, v2) === 0;
};
