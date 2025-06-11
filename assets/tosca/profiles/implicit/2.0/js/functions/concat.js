// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.3.1

// TOSCA 2.0 $concat function
// Section 10.2.3.2: The $concat function takes one or more arguments that must be 
// all of type string or all of type list (with the same entry_schema of the list)

exports.evaluate = function() {
    if (arguments.length === 0) {
        throw new Error('$concat function requires at least one argument');
    }

    // Check first argument to determine expected type
    let firstArg = arguments[0];
    
    // Handle scalar objects with $string property
    if (firstArg && typeof firstArg === 'object' && firstArg.$string !== undefined) {
        firstArg = firstArg.$string;
    }
    
    const isFirstArgString = typeof firstArg === 'string';
    const isFirstArgArray = Array.isArray(firstArg);
    
    if (!isFirstArgString && !isFirstArgArray) {
        throw new Error('$concat function arguments must be all of type string or all of type list, got: ' + 
                       (typeof firstArg) + ' for first argument');
    }
    
    if (isFirstArgString) {
        // String concatenation mode
        let result = [];
        
        for (let i = 0; i < arguments.length; i++) {
            let argument = arguments[i];
            
            // Handle scalar objects with $string property
            if (argument && typeof argument === 'object' && argument.$string !== undefined) {
                argument = argument.$string;
            }
            
            // TOSCA 2.0 Rule: NO implicit type conversion
            // Section 9.1.1.1: A TOSCA parser "MUST NOT attempt to automatically 
            // convert other primitive types to strings if a string type is required"
            if (typeof argument !== 'string') {
                throw new Error('$concat function with string arguments cannot accept argument of type: ' + 
                               (typeof argument) + ' at position ' + i + 
                               '. All arguments must be strings. No implicit type conversion is allowed.');
            }
            
            result.push(argument);
        }
        
        return result.join('');
        
    } else if (isFirstArgArray) {
        // List concatenation mode
        let result = [];
        
        for (let i = 0; i < arguments.length; i++) {
            let argument = arguments[i];
            
            if (!Array.isArray(argument)) {
                throw new Error('$concat function with list arguments cannot accept argument of type: ' + 
                               (typeof argument) + ' at position ' + i + 
                               '. All arguments must be lists.');
            }
            
            // Concatenate array elements
            result = result.concat(argument);
        }
        
        return result;
    }
    
    // This should never be reached, but added for completeness
    throw new Error('$concat function: invalid argument types detected');
};
