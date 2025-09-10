// [TOSCA-v2.0] @ 10.2.5.1

exports.evaluate = function() {
    if (arguments.length === 0) {
        throw 'The $sum function expects at least one argument.';
    }
    
    let result = 0;
    let resultType = 'integer'; // Start with integer, upgrade to float if needed
    
    // Process each argument
    for (let i = 0; i < arguments.length; i++) {
        let arg = arguments[i];
        
        // Validate argument type
        if (typeof arg === 'number') {
            // JavaScript number (can be integer or float)
            result += arg;
            
            // If any argument is a float, result becomes float
            if (!Number.isInteger(arg)) {
                resultType = 'float';
            }
        } else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$scalar')) {
            // TOSCA scalar type with unit
            let scalarValue = arg.$scalar;
            if (typeof scalarValue === 'number') {
                result += scalarValue;
                
                if (!Number.isInteger(scalarValue)) {
                    resultType = 'float';
                }
            } else {
                throw 'The $sum function scalar argument ' + i + ' must have a numeric value; got ' + (typeof scalarValue);
            }
        } else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$integer')) {
            // TOSCA integer type
            let intValue = arg.$integer;
            if (typeof intValue === 'number' && Number.isInteger(intValue)) {
                result += intValue;
            } else {
                throw 'The $sum function integer argument ' + i + ' must be an integer; got ' + (typeof intValue);
            }
        } else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$float')) {
            // TOSCA float type
            let floatValue = arg.$float;
            if (typeof floatValue === 'number') {
                result += floatValue;
                resultType = 'float';
            } else {
                throw 'The $sum function float argument ' + i + ' must be a number; got ' + (typeof floatValue);
            }
        } else {
            throw 'The $sum function argument ' + i + ' must be a number, integer, float, or scalar type; got ' + (typeof arg);
        }
    }
    
    // Return result as appropriate type
    if (resultType === 'float') {
        return parseFloat(result);
    } else {
        return parseInt(result);
    }
};