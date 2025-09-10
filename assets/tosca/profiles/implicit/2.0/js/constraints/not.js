// TOSCA 2.0 logical operator: not
const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // NOT requires exactly one sub-clause
    if (arguments.length !== 2) {
        return false;
    }

    const subclauseArg = arguments[1];
    
    // Handle structured objects from parser
    if (subclauseArg && typeof subclauseArg === 'object' && subclauseArg.Operator && subclauseArg.Arguments) {
        // This is a structured constraint object from the parser
        const operatorName = subclauseArg.Operator;
        const constraintArgs = subclauseArg.Arguments || [];
        
        // Validate the structured subclause and return its negation
        const isSubclauseValid = tosca.validateConstraintSubclause.call(this, operatorName, constraintArgs, currentPropertyValue);
        
        return !isSubclauseValid;
    }
    
    // Handle simple constraint map format (legacy)
    const subclauseMap = subclauseArg;
    
    // Parse the subclause
    const parsed = tosca.parseConstraintSubclause(subclauseMap);
    if (!parsed) {
        return true; // Empty or malformed sub-clause -> NOT of nothing is true
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
    
    // Validate the subclause and return its negation
    const isSubclauseValid = tosca.validateConstraintSubclause.call(this, operatorFunctionName, processedArgs, currentPropertyValue);
    
    return !isSubclauseValid; // Return opposite of sub-clause result
};