// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

// TOSCA 2.0 operator: valid_values
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    if (arguments.length < 2) {
        throw new Error("valid_values requires at least 2 arguments");
    }
    
    // First argument is always the value to check
    const value = arguments[0];
    
    // All remaining arguments are valid values to compare against
    for (let i = 1; i < arguments.length; i++) {
        let validValue = arguments[i];
        
        if (tosca.deepEqual(value, validValue)) {
            return true;
        }
    }
    
    return false;
};
