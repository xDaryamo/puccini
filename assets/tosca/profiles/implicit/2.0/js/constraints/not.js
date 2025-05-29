// TOSCA 2.0 logical operator: not
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // NOT requires exactly one sub-clause
    if (arguments.length !== 2) {
        return false;
    }

    const subclauseMap = arguments[1];

    // Extract operator and arguments
    const operatorKey = Object.keys(subclauseMap)[0];
    if (!operatorKey) {
        return true;
    }

    const operatorFunctionName = operatorKey.startsWith('$') ? operatorKey.substring(1) : operatorKey;
    let originalOperatorArgs = subclauseMap[operatorKey];

    if (!Array.isArray(originalOperatorArgs)) {
        originalOperatorArgs = [originalOperatorArgs];
    }

    // NEW: Process and evaluate nested functions in arguments
    const processedArgsForSubclause = [];
    for (const arg of originalOperatorArgs) {
        if (arg === '$value') {
            processedArgsForSubclause.push(currentPropertyValue);
        } else if (typeof arg === 'object' && arg !== null && !Array.isArray(arg)) {
            // Check if the argument is a function to evaluate
            const functionResult = evaluateNestedFunction(arg, currentPropertyValue);
            processedArgsForSubclause.push(functionResult);
        } else {
            processedArgsForSubclause.push(arg);
        }
    }
    
    try {
        // Dynamically load the validation module
        const validatorModule = require('tosca.validation.' + operatorFunctionName);

        if (validatorModule && typeof validatorModule.validate === 'function') {
            let subclauseResult = false;
            
            // CORRECTION: Always use processedArgsForSubclause for all validation types
            if (operatorFunctionName === 'valid_values' || operatorFunctionName === 'in_range') {
                subclauseResult = validatorModule.validate.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
            } else {
                // For all other operators (including and, or, not, xor), use processed arguments
                subclauseResult = validatorModule.validate.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
            }
            
            // NOT negates the result of the sub-clause
            return !subclauseResult;
        } else {
            return true;
        }
    } catch (e) {
        return true;
    }
};

// NEW FUNCTION: Evaluate nested functions

// Helper function to dereference a path within an object
function dereferencePathHelper(obj, pathArray) {
    let current = obj;
    // Ensure pathArray is an array, supporting {$value: "property"} as {$value: ["property"]}
    const path = Array.isArray(pathArray) ? pathArray : [pathArray];

    for (const key of path) {
        if (current === null || typeof current !== 'object' || !current.hasOwnProperty(key)) {
            return undefined; // Path not found or invalid intermediate value
        }
        current = current[key];
    }
    return current;
}

function evaluateNestedFunction(arg, currentPropertyValue) {
    if (typeof arg !== 'object' || arg === null || Array.isArray(arg)) {
        return arg;
    }
    
    const keys = Object.keys(arg);
    if (keys.length === 1) {
        const key = keys[0];
        if (key.startsWith('$')) {
            const functionName = key.substring(1);
            const functionArgs = arg[key];
            
            // SPECIAL HANDLING: For $value or {$value: ["path", ...]}
            if (functionName === 'value') {
                // functionArgs contains the path, e.g., ["a"] or "a"
                // currentPropertyValue is the context for dereferencing
                return dereferencePathHelper(currentPropertyValue, functionArgs);
            }
            
            // Check if it's a validation operator (to handle recursive structure)
            if (['and', 'or', 'not', 'xor', 'equal', 'greater_than', 'less_than', 'pattern', 'min_length', 'max_length', 'in_range', 'valid_values'].includes(functionName)) {
                // Recursively process operator arguments
                if (Array.isArray(functionArgs)) {
                    const processedArgs = functionArgs.map(subArg => evaluateNestedFunction(subArg, currentPropertyValue));
                    return { [key]: processedArgs };
                } else {
                    return { [key]: evaluateNestedFunction(functionArgs, currentPropertyValue) };
                }
            }
            
            // Try to evaluate it as a function (length, concat, etc.)
            try {
                const functionModule = require('tosca.function.' + functionName);
                if (functionModule && typeof functionModule.evaluate === 'function') {
                    const processedFunctionArgs = [];
                    const argsArray = Array.isArray(functionArgs) ? functionArgs : [functionArgs];
                    
                    for (const fnArg of argsArray) {
                        if (fnArg === '$value') {
                            processedFunctionArgs.push(currentPropertyValue);
                        } else {
                            processedFunctionArgs.push(evaluateNestedFunction(fnArg, currentPropertyValue));
                        }
                    }
                    
                    return functionModule.evaluate.apply(null, processedFunctionArgs);
                }
            } catch (e) {
                // If the function can't be evaluated, return the processed argument
                if (Array.isArray(functionArgs)) {
                    const processedArgs = functionArgs.map(subArg => evaluateNestedFunction(subArg, currentPropertyValue));
                    return { [key]: processedArgs };
                } else {
                    return { [key]: evaluateNestedFunction(functionArgs, currentPropertyValue) };
                }
            }
        }
    }
    
    return arg;
}