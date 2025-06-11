// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    if (arguments.length === 3) {
        // TOSCA 2.0 syntax: $in_range: [ <value_to_test>, [<lower_bound>, <upper_bound>] ]
        // Called as: in_range(currentPropertyValue, valueToTest, [lowerBound, upperBound])
        let valueToTest = arguments[1];
        const boundsArray = arguments[2];
        
        // Handle "$value" substitution
        if (valueToTest === '$value') {
            valueToTest = currentPropertyValue;
        }
        
        if (!Array.isArray(boundsArray) || boundsArray.length !== 2) {
            return false;
        }
        
        let lowerBound = boundsArray[0];
        let upperBound = boundsArray[1];
        
        // Parse bounds if they're strings and we have scalar context
        if (typeof lowerBound === 'string' && valueToTest && 
            valueToTest.$number !== undefined) {
            const parsed = tosca.tryParseScalar(lowerBound, valueToTest);
            if (parsed) {
                lowerBound = parsed;
            }
        }

        if (typeof upperBound === 'string' && valueToTest && 
            valueToTest.$number !== undefined) {
            const parsed = tosca.tryParseScalar(upperBound, valueToTest);
            if (parsed) {
                upperBound = parsed;
            }
        }

        // Use canonical comparison for scalars and other comparable types
        return (tosca.compare(valueToTest, lowerBound) >= 0) && 
               (tosca.compare(valueToTest, upperBound) <= 0);
               
    } else if (arguments.length === 4) {
        // Legacy style: arguments[0] = currentPropertyValue, arguments[1] = valueToTest, arguments[2] = lowerBound, arguments[3] = upperBound
        let valueToTest = arguments[1];
        let lowerBound = arguments[2];
        let upperBound = arguments[3];
        
        // Handle "$value" substitution
        if (valueToTest === '$value') {
            valueToTest = currentPropertyValue;
        }

        if (valueToTest === undefined || valueToTest === null ||
            lowerBound === undefined || lowerBound === null ||
            upperBound === undefined || upperBound === null) {
            return false;
        }

        // Parse lowerBound if it's a string and we have scalar context
        if (typeof lowerBound === 'string' && valueToTest && 
            valueToTest.$number !== undefined) {
            const parsed = tosca.tryParseScalar(lowerBound, valueToTest);
            if (parsed) {
                lowerBound = parsed;
            }
        }

        // Parse upperBound if it's a string and we have scalar context
        if (typeof upperBound === 'string' && valueToTest && 
            valueToTest.$number !== undefined) {
            const parsed = tosca.tryParseScalar(upperBound, valueToTest);
            if (parsed) {
                upperBound = parsed;
            }
        }

        // Use canonical comparison for scalars and other comparable types
        return (tosca.compare(valueToTest, lowerBound) >= 0) && 
               (tosca.compare(valueToTest, upperBound) <= 0);
               
    } else if (arguments.length === 2) {
        // Legacy style: in_range(currentPropertyValue, [lowerBound, upperBound])
        const boundsArray = arguments[1];
        
        if (!Array.isArray(boundsArray) || boundsArray.length !== 2) {
            return false;
        }
        
        let lowerBound = boundsArray[0];
        let upperBound = boundsArray[1];
        
        // Parse bounds if they're strings and we have scalar context
        if (typeof lowerBound === 'string' && currentPropertyValue && 
            currentPropertyValue.$number !== undefined) {
            const parsed = tosca.tryParseScalar(lowerBound, currentPropertyValue);
            if (parsed) {
                lowerBound = parsed;
            }
        }

        if (typeof upperBound === 'string' && currentPropertyValue && 
            currentPropertyValue.$number !== undefined) {
            const parsed = tosca.tryParseScalar(upperBound, currentPropertyValue);
            if (parsed) {
                upperBound = parsed;
            }
        }

        // Use canonical comparison for scalars and other comparable types
        return (tosca.compare(currentPropertyValue, lowerBound) >= 0) && 
               (tosca.compare(currentPropertyValue, upperBound) <= 0);
    } else {
        return false;
    }
};