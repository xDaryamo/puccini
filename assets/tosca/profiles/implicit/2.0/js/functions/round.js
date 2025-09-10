// [TOSCA-v2.0] @ 10.2.5.6

exports.evaluate = function() {
    if (arguments.length !== 1) {
        throw 'The $round function expects exactly one argument.';
    }
    
    let arg = arguments[0];
    
    // Process the argument - must be float or number
    let floatValue = 0;
    
    // Case 1: JavaScript number
    if (typeof arg === 'number') {
        floatValue = arg;
    }
    // Case 2: TOSCA float object
    else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$float')) {
        let value = arg.$float;
        if (typeof value === 'number') {
            floatValue = value;
        } else {
            throw 'The $round function float argument must be a number; got ' + (typeof value);
        }
    }
    // Case 3: TOSCA integer object (can be converted to float)
    else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$integer')) {
        let intValue = arg.$integer;
        if (typeof intValue === 'number' && Number.isInteger(intValue)) {
            floatValue = parseFloat(intValue);
        } else {
            throw 'The $round function integer argument must be an integer; got ' + (typeof intValue);
        }
    }
    // Case 4: Puccini's scalar representation with numeric value
    else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('scalar') && typeof arg.scalar === 'number') {
        floatValue = arg.scalar;
    }
    // Case 5: Puccini's alternate format with $number
    else if (arg !== null && typeof arg === 'object' && arg.hasOwnProperty('$number') && typeof arg.$number === 'number') {
        floatValue = arg.$number;
    }
    else {
        throw 'The $round function argument must be a float or number type; got ' + (typeof arg);
    }
    
    // Implement TOSCA rounding rule:
    // Equal value distance is rounded down (e.g. 3.5 -> 3, 3.53 -> 4)
    // This means: if fractional part is exactly 0.5, round down; otherwise use normal rounding
    
    let fractionalPart = Math.abs(floatValue - Math.trunc(floatValue));
    
    // Calculate result and normalize -0 to 0
    let result;
    if (Math.abs(fractionalPart - 0.5) < Number.EPSILON) {
        result = Math.trunc(floatValue);
    } else {
        result = Math.round(floatValue);
    }
    return result === 0 ? 0 : result;  // Converte -0 in 0
};