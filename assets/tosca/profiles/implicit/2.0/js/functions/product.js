// [TOSCA-v2.0] @ 10.2.5.3

exports.evaluate = function() {
    if (arguments.length === 0) {
        throw 'The $product function expects at least one argument.';
    }
    
    // Case 1: Two arguments - check for scalar multiplication
    if (arguments.length === 2) {
        let arg1 = arguments[0];
        let arg2 = arguments[1];
        
        // Check if first argument is a scalar string (like "512 MiB")
        if (typeof arg1 === 'string') {
            // Parse scalar string format: "value unit"
            let parts = arg1.trim().split(/\s+/);
            if (parts.length >= 2) {
                let scalarValue = parseFloat(parts[0]);
                let unit = parts.slice(1).join(' ');
                
                if (!isNaN(scalarValue)) {
                    // Process second argument (multiplier)
                    let multiplier = 0;
                    if (typeof arg2 === 'number') {
                        multiplier = arg2;
                    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$integer')) {
                        let intValue = arg2.$integer;
                        if (typeof intValue === 'number' && Number.isInteger(intValue)) {
                            multiplier = intValue;
                        } else {
                            throw 'The $product function integer argument must be an integer; got ' + (typeof intValue);
                        }
                    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$float')) {
                        let floatValue = arg2.$float;
                        if (typeof floatValue === 'number') {
                            multiplier = floatValue;
                        } else {
                            throw 'The $product function float argument must be a number; got ' + (typeof floatValue);
                        }
                    } else {
                        throw 'The $product function second argument must be a number, integer, or float type; got ' + (typeof arg2);
                    }
                    
                    // Return scalar string with multiplied value
                    let resultValue = scalarValue * multiplier;
                    return resultValue + ' ' + unit;
                }
            }
        }
        
        // Check if first argument is a TOSCA scalar object
        if (arg1 !== null && typeof arg1 === 'object' && arg1.hasOwnProperty('$scalar')) {
            let scalarValue = arg1.$scalar;
            let multiplier = 0;
            
            // Validate scalar value
            if (typeof scalarValue !== 'number') {
                throw 'The $product function scalar argument must have a numeric value; got ' + (typeof scalarValue);
            }
            
            // Process second argument (multiplier)
            if (typeof arg2 === 'number') {
                multiplier = arg2;
            } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$integer')) {
                let intValue = arg2.$integer;
                if (typeof intValue === 'number' && Number.isInteger(intValue)) {
                    multiplier = intValue;
                } else {
                    throw 'The $product function integer argument must be an integer; got ' + (typeof intValue);
                }
            } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$float')) {
                let floatValue = arg2.$float;
                if (typeof floatValue === 'number') {
                    multiplier = floatValue;
                } else {
                    throw 'The $product function float argument must be a number; got ' + (typeof floatValue);
                }
            } else {
                throw 'The $product function second argument must be a number, integer, or float type; got ' + (typeof arg2);
            }
            
            // Return scalar with multiplied value
            return {
                $scalar: scalarValue * multiplier,
                unit: arg1.unit // Preserve unit if present
            };
        }
    }
    
    // Case 2: Multiple arguments - numeric multiplication
    let result = 1;
    let resultType = 'integer'; // Start with integer, upgrade to float if needed
    
    // Process each argument
    for (let i = 0; i < arguments.length; i++) {
        let arg = arguments[i];
        
        // Validate argument type
        if (typeof arg === 'number') {
            // JavaScript number (can be integer or float)
            result *= arg;
            
            // If any argument is a float, result becomes float
            if (!Number.isInteger(arg)) {
                resultType = 'float';
            }
        } else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$integer')) {
            // TOSCA integer type
            let intValue = arg.$integer;
            if (typeof intValue === 'number' && Number.isInteger(intValue)) {
                result *= intValue;
            } else {
                throw 'The $product function integer argument ' + i + ' must be an integer; got ' + (typeof intValue);
            }
        } else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$float')) {
            // TOSCA float type
            let floatValue = arg.$float;
            if (typeof floatValue === 'number') {
                result *= floatValue;
                resultType = 'float';
            } else {
                throw 'The $product function float argument ' + i + ' must be a number; got ' + (typeof floatValue);
            }
        } else {
            throw 'The $product function argument ' + i + ' must be a number, integer, or float type; got ' + (typeof arg);
        }
    }
    
    // Return result as appropriate type
    if (resultType === 'float') {
        return parseFloat(result);
    } else {
        return parseInt(result);
    }
};