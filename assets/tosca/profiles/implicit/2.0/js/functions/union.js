// [TOSCA-v2.0] @ 10.2.4.1

exports.evaluate = function() {
    if (arguments.length === 0) {
        throw 'The $union function expects at least one argument.';
    }
    
    let result = [];
    let seen = new Set();
    
    // Process each list argument
    for (let i = 0; i < arguments.length; i++) {
        let arg = arguments[i];
        
        // Validate that the argument is a list/array
        if (!Array.isArray(arg)) {
            throw 'The $union function argument ' + i + ' must be a list; got ' + (typeof arg);
        }
        
        // Add unique elements to result
        for (let j = 0; j < arg.length; j++) {
            let element = arg[j];
            
            // Create a string key for comparison (handles objects and primitives)
            let key = JSON.stringify(element);
            
            if (!seen.has(key)) {
                seen.add(key);
                result.push(element);
            }
        }
    }
    
    return result;
};