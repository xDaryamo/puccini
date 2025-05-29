// TOSCA $length function
// [TOSCA-v2.0] @ 10.2.3.1

exports.evaluate = function(arg) {
    if (arguments.length !== 1) {
        throw 'The $length function expects exactly one argument.';
    }

    if (arg === null || arg === undefined) {
        throw 'The $length function argument cannot be null or undefined.';
    }

    if (typeof arg === 'string') {
        return arg.length;
    }

    if (Array.isArray(arg)) {
        return arg.length;
    }

    // In Goja, ard.Map becomes a JS object
    if (typeof arg === 'object' && arg !== null && !Array.isArray(arg)) {
        return Object.keys(arg).length;
    }
    
    throw 'The $length function argument must be a string, list, or map; got ' + (typeof arg);
};