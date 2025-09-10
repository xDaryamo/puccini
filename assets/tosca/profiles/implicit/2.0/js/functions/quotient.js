// [TOSCA-v2.0] @ 10.2.5.4

exports.evaluate = function() {
    if (arguments.length !== 2) {
        throw 'The $quotient function expects exactly two arguments.';
    }
    
    let arg1 = arguments[0];
    let arg2 = arguments[1];
    
    // Process second argument (divisor) first - must be number, integer, or float
    let divisor = 0;
    if (typeof arg2 === 'number') {
        divisor = arg2;
    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$integer')) {
        let intValue = arg2.$integer;
        if (typeof intValue === 'number' && Number.isInteger(intValue)) {
            divisor = intValue;
        } else {
            throw 'The $quotient function second argument integer must be an integer; got ' + (typeof intValue);
        }
    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$float')) {
        let floatValue = arg2.$float;
        if (typeof floatValue === 'number') {
            divisor = floatValue;
        } else {
            throw 'The $quotient function second argument float must be a number; got ' + (typeof floatValue);
        }
    } else {
        throw 'The $quotient function second argument must be a number, integer, or float type; got ' + (typeof arg2);
    }
    
    // Check for division by zero
    if (divisor === 0) {
        throw 'The $quotient function cannot divide by zero.';
    }
    
    // Case 1: First argument is a scalar string (like "1024 MiB")
    if (typeof arg1 === 'string') {
        // Parse scalar string format: "value unit"
        let parts = arg1.trim().split(/\s+/);
        if (parts.length >= 2) {
            let scalarValue = parseFloat(parts[0]);
            let unit = parts.slice(1).join(' ');
            
            if (!isNaN(scalarValue)) {
                // Return scalar string with divided value (truncated if necessary)
                let resultValue = scalarValue / divisor;
                // Truncate towards zero (implementation decision)
                resultValue = Math.trunc(resultValue);
                return resultValue + ' ' + unit;
            }
        }
        
        throw 'The $quotient function first argument scalar string is malformed: ' + arg1;
    }
    
    // Case 2: First argument is a TOSCA scalar object
    if (arg1 !== null && typeof arg1 === 'object' && arg1.hasOwnProperty('$scalar')) {
        let scalarValue = arg1.$scalar;
        
        // Validate scalar value
        if (typeof scalarValue !== 'number') {
            throw 'The $quotient function scalar argument must have a numeric value; got ' + (typeof scalarValue);
        }
        
        // Return scalar with divided value (truncated if necessary)
        let resultValue = scalarValue / divisor;
        // Truncate towards zero (implementation decision)
        resultValue = Math.trunc(resultValue);
        
        return {
            $scalar: resultValue,
            unit: arg1.unit // Preserve unit if present
        };
    }
    
    // Case 3: First argument is a number (integer or float)
    let dividend = 0;
    if (typeof arg1 === 'number') {
        dividend = arg1;
    } else if (arg1 !== null && typeof arg1 === 'object' && arg1.hasOwnProperty('$integer')) {
        // TOSCA integer type
        let intValue = arg1.$integer;
        if (typeof intValue === 'number' && Number.isInteger(intValue)) {
            dividend = intValue;
        } else {
            throw 'The $quotient function first argument integer must be an integer; got ' + (typeof intValue);
        }
    } else if (arg1 !== null && typeof arg1 === 'object' && arg1.hasOwnProperty('$float')) {
        // TOSCA float type
        let floatValue = arg1.$float;
        if (typeof floatValue === 'number') {
            dividend = floatValue;
        } else {
            throw 'The $quotient function first argument float must be a number; got ' + (typeof floatValue);
        }
    } else {
        throw 'The $quotient function first argument must be a number, integer, float, or scalar type; got ' + (typeof arg1);
    }
    
    // Calculate quotient - always returns float for numeric division
    let result = dividend / divisor;
    return parseFloat(result);
};