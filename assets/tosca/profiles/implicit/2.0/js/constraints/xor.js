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
        const subclauseArg = arguments[i];
        
        // Handle structured objects from parser
        if (subclauseArg && typeof subclauseArg === 'object' && subclauseArg.Operator && subclauseArg.Arguments) {
            // This is a structured constraint object from the parser
            const operatorName = subclauseArg.Operator;
            const constraintArgs = subclauseArg.Arguments || [];
            
            // Validate the structured subclause
            const isSubclauseValid = tosca.validateConstraintSubclause.call(this, operatorName, constraintArgs, currentPropertyValue);
            
            if (isSubclauseValid) {
                trueCount++;
            }
            continue;
        }
        
        // Handle simple constraint map format (legacy)
        const subclauseMap = subclauseArg;
        
        // Parse the subclause
        const parsed = tosca.parseConstraintSubclause(subclauseMap);
        if (!parsed) {
            continue; // Skip empty or malformed sub-clause
        }

        const { operatorFunctionName, originalOperatorArgs } = parsed;

        // Process arguments - handle complex objects but preserve simple "$value" strings
        const processedArgs = [];
        for (const arg of originalOperatorArgs) {
            let processed;
            if (arg === '$value') {
                // Handle string literal "$value" - don't process, let individual constraints handle it
                processed = arg;
            } else if (typeof arg === 'object' && arg !== null && !Array.isArray(arg)) {
                // Handle complex objects like {"$value":["count"]} - these need processing
                processed = tosca.evaluateConstraintArgument.call(this, arg, currentPropertyValue);
            } else {
                // Handle arrays and primitives - pass through
                processed = arg;
            }
            processedArgs.push(processed);
        }
        
        // Validate the subclause
        const isSubclauseValid = tosca.validateConstraintSubclause.call(this, operatorFunctionName, processedArgs, currentPropertyValue);
        
        if (isSubclauseValid) {
            trueCount++;
        }
    }

    return trueCount === 1; // XOR succeeds if exactly one clause is true
};