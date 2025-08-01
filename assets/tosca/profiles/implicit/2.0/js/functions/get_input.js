// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.4.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.4.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.4.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.4.1

const tosca = require('tosca.lib.utils');

exports.evaluate = function() {
    if (arguments.length < 1 || arguments.length > 2)
        throw 'must have 1 or 2 arguments';
    if (!tosca.isTosca(clout))
        throw 'Clout is not TOSCA';
    
    let input, nestedPath;
    
    // Handle both syntaxes:
    // $get_input: input_name
    // $get_input: [input_name, index/key, ...]
    if (arguments.length === 1) {
        if (Array.isArray(arguments[0])) {
            // Syntax: $get_input: [input_name, nested_elements...]
            input = arguments[0][0];
            nestedPath = arguments[0].slice(1);
        } else {
            // Syntax: $get_input: input_name
            input = arguments[0];
            nestedPath = null;
        }
    } else {
        // Alternative syntax with separate arguments (for compatibility)
        input = arguments[0];
        nestedPath = [arguments[1]];
    }
    
    let inputs = clout.properties.tosca.inputs;
    if (!(input in inputs))
        throw util.sprintf('input %q not found', input);
    let r = inputs[input];
    r = clout.coerce(r);
    
    // Navigate through the structure if a path is specified
    if (nestedPath && nestedPath.length > 0) {
        for (let i = 0; i < nestedPath.length; i++) {
            let key = nestedPath[i];
            
            // If key is "$node_index", evaluate it as a function
            if (key === '$node_index') {
                try {
                    const nodeIndexFunction = require('tosca.function.$node_index');
                    key = nodeIndexFunction.evaluate.call(this);
                } catch (e) {
                    throw util.sprintf('failed to evaluate $node_index: %s', e);
                }
            }
            
            if (Array.isArray(r)) {
                // Array index access
                let index = parseInt(key);
                if (isNaN(index) || index < 0 || index >= r.length) {
                    throw util.sprintf('index %d out of bounds for input %q (length %d)', index, input, r.length);
                }
                r = r[index];
            } else if (r && typeof r === 'object') {
                // Object property access
                if (!(key in r)) {
                    throw util.sprintf('property %q not found in input %q', key, input);
                }
                r = r[key];
            } else {
                throw util.sprintf('cannot access %q in non-object/non-array value', key);
            }
        }
    }
    
    return r;
};
