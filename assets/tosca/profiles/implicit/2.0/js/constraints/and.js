// and.js
// TOSCA 2.0 logical operator: and
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // Empty AND is vacuously true
    if (arguments.length <= 1) {
        return true;
    }

    // Process each sub-clause
    for (let i = 1; i < arguments.length; i++) {
        const subclauseMap = arguments[i];

        const operatorKey = Object.keys(subclauseMap)[0];
        if (!operatorKey) {
            return false; // Empty or malformed sub-clause makes AND false
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
            const validatorModule = require('tosca.validation.' + operatorFunctionName);

            if (validatorModule && typeof validatorModule.validate === 'function') {
                const subValidator = validatorModule.validate;

                // Always use the same argument pattern: currentPropertyValue followed by constraint arguments
                const isSubclauseValid = subValidator.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
                
                if (!isSubclauseValid) {
                    return false; // Short-circuit if a sub-clause is false
                }
            } else {
                return false; // Module or function not found
            }
        } catch (e) {
            console.warn(`Warning: Error validating ${operatorFunctionName}: ${e.message}`);
            return false; // Error during sub-clause validation
        }
    }

    return true; // All sub-clauses are true
};