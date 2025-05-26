// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare (ignoring property value)
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
    } else {
        throw new Error("greater_or_equal requires at least 2 arguments");
    }
    
    return tosca.compare(v1, v2) >= 0;
};