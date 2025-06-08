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

    // Process and evaluate nested functions in arguments
    const processedArgsForSubclause = [];
    for (const arg of originalOperatorArgs) {
        if (arg === '$value') {
            processedArgsForSubclause.push(currentPropertyValue);
        } else if (typeof arg === 'object' && arg !== null && !Array.isArray(arg)) {
            // Check if the argument is a function to evaluate
            const functionResult = tosca.evaluateNestedFunction(arg, currentPropertyValue);
            processedArgsForSubclause.push(functionResult);
        } else if (typeof arg === 'string' && currentPropertyValue) {
            const parsed = tosca.tryParseScalar(arg, currentPropertyValue);
            processedArgsForSubclause.push(parsed || arg);
        } else {
            processedArgsForSubclause.push(arg);
        }
    }
    
    try {
        // Dynamically load the validation module
        const validatorModule = require('tosca.validation.' + operatorFunctionName);

        if (validatorModule && typeof validatorModule.validate === 'function') {
            const subValidator = validatorModule.validate;

            // Always use the same argument pattern: currentPropertyValue followed by constraint arguments
            const subclauseResult = subValidator.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
            
            // NOT negates the result of the sub-clause
            return !subclauseResult;
        } else {
            return true; // Module or function not found
        }
    } catch (e) {
        console.warn(`Warning: Error validating ${operatorFunctionName}: ${e.message}`);
        return true; // Error during sub-clause validation
    }
};