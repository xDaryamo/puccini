// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Take the last 2 arguments like other constraints do
    if (arguments.length < 2)
        throw 'must have at least 2 arguments: value and minimum length';
    
    var v = arguments[arguments.length - 2];      // Second-to-last argument
    var length = arguments[arguments.length - 1]; // Last argument
    
    return tosca.getLength(v) >= length;
};
