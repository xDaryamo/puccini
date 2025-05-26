// [TOSCA-v2.0] @ 10.2.5.2

exports.evaluate = function() {
    if (arguments.length !== 2) {
        throw 'The $difference function expects exactly two arguments.';
    }
    
    let arg1 = arguments[0];
    let arg2 = arguments[1];
    let resultType = 'integer'; // Start with integer, upgrade to float if needed
    
    // Process first argument
    let value1 = 0;
    if (typeof arg1 === 'number') {
        value1 = arg1;
        if (!Number.isInteger(arg1)) {
            resultType = 'float';
        }
    } else if (arg1 !== null && typeof arg1 === 'object' && arg1.hasOwnProperty('$scalar')) {
        // TOSCA scalar type with unit
        let scalarValue = arg1.$scalar;
        if (typeof scalarValue === 'number') {
            value1 = scalarValue;
            if (!Number.isInteger(scalarValue)) {
                resultType = 'float';
            }
        } else {
            throw 'The $difference function first argument scalar must have a numeric value; got ' + (typeof scalarValue);
        }
    } else if (arg1 !== null && typeof arg1 === 'object' && arg1.hasOwnProperty('$integer')) {
        // TOSCA integer type
        let intValue = arg1.$integer;
        if (typeof intValue === 'number' && Number.isInteger(intValue)) {
            value1 = intValue;
        } else {
            throw 'The $difference function first argument integer must be an integer; got ' + (typeof intValue);
        }
    } else if (arg1 !== null && typeof arg1 === 'object' && arg1.hasOwnProperty('$float')) {
        // TOSCA float type
        let floatValue = arg1.$float;
        if (typeof floatValue === 'number') {
            value1 = floatValue;
            resultType = 'float';
        } else {
            throw 'The $difference function first argument float must be a number; got ' + (typeof floatValue);
        }
    } else {
        throw 'The $difference function first argument must be a number, integer, float, or scalar type; got ' + (typeof arg1);
    }
    
    // Process second argument
    let value2 = 0;
    if (typeof arg2 === 'number') {
        value2 = arg2;
        if (!Number.isInteger(arg2)) {
            resultType = 'float';
        }
    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$scalar')) {
        // TOSCA scalar type with unit
        let scalarValue = arg2.$scalar;
        if (typeof scalarValue === 'number') {
            value2 = scalarValue;
            if (!Number.isInteger(scalarValue)) {
                resultType = 'float';
            }
        } else {
            throw 'The $difference function second argument scalar must have a numeric value; got ' + (typeof scalarValue);
        }
    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$integer')) {
        // TOSCA integer type
        let intValue = arg2.$integer;
        if (typeof intValue === 'number' && Number.isInteger(intValue)) {
            value2 = intValue;
        } else {
            throw 'The $difference function second argument integer must be an integer; got ' + (typeof intValue);
        }
    } else if (arg2 !== null && typeof arg2 === 'object' && arg2.hasOwnProperty('$float')) {
        // TOSCA float type
        let floatValue = arg2.$float;
        if (typeof floatValue === 'number') {
            value2 = floatValue;
            resultType = 'float';
        } else {
            throw 'The $difference function second argument float must be a number; got ' + (typeof floatValue);
        }
    } else {
        throw 'The $difference function second argument must be a number, integer, float, or scalar type; got ' + (typeof arg2);
    }
    
    // Calculate difference: value1 - value2
    let result = value1 - value2;
    
    // Return result as appropriate type
    if (resultType === 'float') {
        return parseFloat(result);
    } else {
        return parseInt(result);
    }
};