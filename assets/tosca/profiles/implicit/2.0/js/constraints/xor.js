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

        // Process arguments, replacing $value with actual property value
        const processedArgsForSubclause = [];
        for (const arg of originalOperatorArgs) {
            if (arg === '$value') {
                processedArgsForSubclause.push(currentPropertyValue);
            } else {
                processedArgsForSubclause.push(arg);
            }
        }
        
        try {
            // Dynamically load the validation module
            const validatorModule = require('tosca.validation.' + operatorFunctionName);

            if (validatorModule && typeof validatorModule.validate === 'function') {
                let isSubclauseValid = false;
                const subValidator = validatorModule.validate;

                // Call the sub-validator based on its type
                if (operatorFunctionName === 'valid_values' || operatorFunctionName === 'in_range') {
                    isSubclauseValid = subValidator.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
                } else if (operatorFunctionName === 'and' || operatorFunctionName === 'or' || operatorFunctionName === 'not' || operatorFunctionName === 'xor') {
                    isSubclauseValid = subValidator.apply(null, [currentPropertyValue, ...originalOperatorArgs]);
                } else {
                    isSubclauseValid = subValidator.apply(null, processedArgsForSubclause);
                }
                
                if (isSubclauseValid) {
                    trueCount++;
                }
            }
        } catch (e) {
            // Treat errors as false results
        }
    }

    return trueCount === 1;
};