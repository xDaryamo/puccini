// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.3.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.3.1

// TOSCA 1.3/2.0 compatible $concat function
// TOSCA 1.3: Allows implicit conversion of compatible types to strings
// TOSCA 2.0: Strict type checking with no implicit conversion

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
    const isFirstArgNumber = typeof firstArg === 'number';
    const isFirstArgArray = Array.isArray(firstArg);
    
    // For TOSCA 1.3 compatibility: treat numbers as strings for concatenation
    const isStringConcatenation = isFirstArgString || isFirstArgNumber;
    
    if (!isStringConcatenation && !isFirstArgArray) {
        throw new Error('$concat function arguments must be all of type string/number or all of type list, got: ' + 
                       (typeof firstArg) + ' for first argument');
    }
    
    if (isStringConcatenation) {
        // String concatenation mode (TOSCA 1.3 compatible)
        let result = [];
        
        for (let i = 0; i < arguments.length; i++) {
            let argument = arguments[i];
            
            // Handle scalar objects with $string property
            if (argument && typeof argument === 'object' && argument.$string !== undefined) {
                argument = argument.$string;
            } else if (argument && typeof argument === 'object' && argument.$number !== undefined) {
                // Handle scalar numbers
                argument = argument.$number;
            }
            
            // TOSCA 1.3: Allow implicit conversion of compatible types
            if (typeof argument === 'string') {
                result.push(argument);
            } else if (typeof argument === 'number') {
                // Implicit conversion: number to string
                result.push(String(argument));
            } else if (typeof argument === 'boolean') {
                // Implicit conversion: boolean to string
                result.push(String(argument));
            } else if (argument === null || argument === undefined) {
                // Handle null/undefined as empty string
                result.push('');
            } else {
                throw new Error('$concat function with string arguments cannot accept argument of type: ' + 
                               (typeof argument) + ' at position ' + i + 
                               '. Supported types: string, number, boolean.');
            }
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
