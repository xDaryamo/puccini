// TOSCA 2.0 logical operator: xor
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // XOR is true if exactly one sub-clause is true
    if (arguments.length <= 1) {
        return false;
    }

    let trueCount = 0;

    // Process each sub-clause
    for (let i = 1; i < arguments.length; i++) {
        const subclauseMap = arguments[i];

        // Extract operator and arguments
        const operatorKey = Object.keys(subclauseMap)[0];
        if (!operatorKey) {
            continue;
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
                const isSubclauseValid = subValidator.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
                
                if (isSubclauseValid) {
                    trueCount++;
                }
            }
        } catch (e) {
            console.warn(`Warning: Error validating ${operatorFunctionName}: ${e.message}`);
            // Treat errors as false results
        }
    }

    return trueCount === 1;
};