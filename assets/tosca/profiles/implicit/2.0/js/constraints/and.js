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
            return false; // Una sub-clausola vuota o malformata rende l'AND falso
        }

        const operatorFunctionName = operatorKey.startsWith('$') ? operatorKey.substring(1) : operatorKey;
        let originalOperatorArgs = subclauseMap[operatorKey];

        if (!Array.isArray(originalOperatorArgs)) {
            originalOperatorArgs = [originalOperatorArgs];
        }

        const processedArgsForSubclause = [];
        for (const arg of originalOperatorArgs) {
            if (arg === '$value') {
                processedArgsForSubclause.push(currentPropertyValue);
            } else {
                processedArgsForSubclause.push(arg);
            }
        }
        
        try {
            const validatorModule = require('tosca.validation.' + operatorFunctionName);

            if (validatorModule && typeof validatorModule.validate === 'function') {
                let isSubclauseValid = false;
                const subValidator = validatorModule.validate;

                if (operatorFunctionName === 'valid_values' || operatorFunctionName === 'in_range') {
                    isSubclauseValid = subValidator.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
                } else if (operatorFunctionName === 'and' || operatorFunctionName === 'or' || operatorFunctionName === 'not' || operatorFunctionName === 'xor') {
                    isSubclauseValid = subValidator.apply(null, [currentPropertyValue, ...originalOperatorArgs]);
                } else {
                    isSubclauseValid = subValidator.apply(null, [currentPropertyValue, ...processedArgsForSubclause]);
                }
                
                if (!isSubclauseValid) {
                    return false; // Cortocircuito se una sub-clausola è falsa
                }
            } else {
                return false; // Modulo o funzione non trovati
            }
        } catch (e) {
            return false; // Errore durante la validazione della sub-clausola
        }
    }

    return true; // Tutte le sub-clausole sono vere
};