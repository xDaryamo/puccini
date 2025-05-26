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
            let subclauseResult = false;
            
            // Call the sub-validator based on its type
            if (operatorFunctionName === 'valid_values' || operatorFunctionName === 'in_range') {
                subclauseResult = validatorModule.validate.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
            } else if (operatorFunctionName === 'and' || operatorFunctionName === 'or' || operatorFunctionName === 'not' || operatorFunctionName === 'xor') {
                subclauseResult = validatorModule.validate.apply(null, [currentPropertyValue, ...originalOperatorArgs]);
            } else {
                subclauseResult = validatorModule.validate.apply(null, processedArgsForSubclause);
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