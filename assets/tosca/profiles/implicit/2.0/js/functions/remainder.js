// [TOSCA-v2.0] @ 10.2.5.5

exports.evaluate = function() {
    if (arguments.length !== 2) {
        throw 'The $remainder function expects exactly two arguments.';
    }
    
    let arg1 = arguments[0];
    let arg2 = arguments[1];
    
    // Process second argument (divisor) - must be integer
    let divisor = 0;
    if (typeof arg2 === 'number' && Number.isInteger(arg2)) {
        divisor = arg2;
    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$integer')) {
        let intValue = arg2.$integer;
        if (typeof intValue === 'number' && Number.isInteger(intValue)) {
            divisor = intValue;
        } else {
            throw 'The $remainder function second argument integer must be an integer; got ' + (typeof intValue);
        }
    } else {
        throw 'The $remainder function second argument must be an integer; got ' + (typeof arg2);
    }
    
    // Check for division by zero
    if (divisor === 0) {
        throw 'The $remainder function cannot divide by zero.';
    }
    
    // Process first argument
    let dividend = 0;
    let unit = '';
    let isScalar = false;
    
    // Case 1: Scalar string (like "90 s")
    if (typeof arg1 === 'string') {
        // Parse scalar string format: "value unit"
        let parts = arg1.trim().split(/\s+/);
        if (parts.length >= 2) {
            let scalarValue = parseFloat(parts[0]);
            unit = parts.slice(1).join(' ');
            
            if (!isNaN(scalarValue)) {
                dividend = Math.floor(scalarValue); // Ensure integer for remainder
                isScalar = true;
            } else {
                throw 'The $remainder function first argument scalar must have a numeric value; got ' + parts[0];
            }
        } else {
            throw 'The $remainder function first argument scalar string is malformed: ' + arg1;
        }
    }
    // Case 2: Puccini's complex scalar object format
    else if (arg1 !== null && typeof arg1 === 'object') {
        // Handle Puccini's scalar representation
        if (arg1.hasOwnProperty('scalar') && typeof arg1.scalar === 'number') {
            dividend = Math.floor(arg1.scalar);
            unit = arg1.unit || '';
            isScalar = true;
        } else if (arg1.hasOwnProperty('$number') && typeof arg1.$number === 'number') {
            dividend = Math.floor(arg1.$number);
            unit = arg1.unit || '';
            isScalar = true;
        } else if (arg1.hasOwnProperty('$scalar')) {
            let scalarValue = arg1.$scalar;
            if (typeof scalarValue === 'number') {
                dividend = Math.floor(scalarValue);
                unit = arg1.unit || '';
                isScalar = true;
            } else {
                throw 'The $remainder function scalar argument must have a numeric value; got ' + (typeof scalarValue);
            }
        } else if (arg1.hasOwnProperty('$integer')) {
            let intValue = arg1.$integer;
            if (typeof intValue === 'number' && Number.isInteger(intValue)) {
                dividend = intValue;
            } else {
                throw 'The $remainder function first argument integer must be an integer; got ' + (typeof intValue);
            }
        } else {
            throw 'The $remainder function first argument format not recognized: ' + JSON.stringify(arg1);
        }
    }
    // Case 3: Plain integer
    else if (typeof arg1 === 'number' && Number.isInteger(arg1)) {
        dividend = arg1;
    }
    else {
        throw 'The $remainder function first argument must be an integer or scalar type; got ' + (typeof arg1);
    }
    
    // Calculate remainder
    let result = dividend % divisor;
    
    // Return appropriate format
    if (isScalar && unit) {
        // Return as scalar string to match TOSCA format
        return result + ' ' + unit;
    } else {
        return parseInt(result);
    }
};