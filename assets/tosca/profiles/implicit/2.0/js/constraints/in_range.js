// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2

const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    // First handle the case where we're being incorrectly called with a list value in a list context
    if (arguments.length === 3 && Array.isArray(arguments[1]) && Array.isArray(arguments[2])) {
        // This indicates we're incorrectly being called with in_range(firstElement, entireList, bounds)
        // Let's validate each element individually
        const firstElement = currentPropertyValue; // First element already extracted
        const valuesList = arguments[1];
        const bounds = arguments[2];
        
        // Check if the first element is in range
        if (tosca.compare(firstElement, bounds[0]) < 0 || tosca.compare(firstElement, bounds[1]) > 0) {
            return false;
        }
        
        // Check remaining elements in the list
        for (let i = 1; i < valuesList.length; i++) {
            if (tosca.compare(valuesList[i], bounds[0]) < 0 || tosca.compare(valuesList[i], bounds[1]) > 0) {
                return false;
            }
        }
        
        return true;
    }
    
    if (arguments.length === 3) {
        // TOSCA 1.3 syntax: constraint calls in_range(currentValue, lowerBound, upperBound)
        // Or TOSCA 2.0 syntax: $in_range: [ <value_to_test>, [<lower_bound>, <upper_bound>] ]
        let valueToTest = arguments[1];
        let secondArg = arguments[2];
        
        // Handle "$value" substitution for valueToTest
        if (valueToTest === '$value') {
            valueToTest = currentPropertyValue;
        }

        if (valueToTest === undefined || valueToTest === null) {
            return false;
        }
        
        // Check if this is TOSCA 2.0 syntax with bounds array
        if (Array.isArray(secondArg) && secondArg.length === 2) {
            let lowerBound = secondArg[0];
            let upperBound = secondArg[1];

            if (lowerBound === undefined || lowerBound === null ||
                upperBound === undefined || upperBound === null) {
                return false;
            }

            // Parse bounds with scalar context
            lowerBound = tosca.parseScalarOrVersionBound(lowerBound, valueToTest);
            upperBound = tosca.parseScalarOrVersionBound(upperBound, valueToTest);

            return (tosca.compare(valueToTest, lowerBound) >= 0) && 
                   (tosca.compare(valueToTest, upperBound) <= 0);
        } else {
            // TOSCA 1.3 syntax: in_range(currentValue, lowerBound, upperBound)
            // Here valueToTest is actually the current value, secondArg is lowerBound
            // We need to get the third argument as upperBound
            let lowerBound = valueToTest; // arguments[1]
            let upperBound = secondArg;   // arguments[2]
            valueToTest = currentPropertyValue; // The actual value being tested
            
            if (lowerBound === undefined || lowerBound === null ||
                upperBound === undefined || upperBound === null) {
                return false;
            }

            // Parse bounds with scalar context for TOSCA 1.3
            lowerBound = tosca.parseScalarOrVersionBound(lowerBound, valueToTest);
            upperBound = tosca.parseScalarOrVersionBound(upperBound, valueToTest);

            return (tosca.compare(valueToTest, lowerBound) >= 0) && 
                   (tosca.compare(valueToTest, upperBound) <= 0);
        }
               
    } else if (arguments.length === 2) {
        const secondArg = arguments[1];
        
        // Check if this is the nested array syntax: [ $value, [bounds] ]
        if (Array.isArray(secondArg) && secondArg.length === 2) {
            const firstElement = secondArg[0];
            const secondElement = secondArg[1];
            
            // TOSCA 2.0 style: [ $value, [lower, upper] ]
            if ((firstElement === '$value' || typeof firstElement !== 'undefined') && 
                Array.isArray(secondElement) && secondElement.length === 2) {
                
                let valueToTest = firstElement === '$value' ? currentPropertyValue : firstElement;
                let lowerBound = secondElement[0];
                let upperBound = secondElement[1];
                
                if (valueToTest === undefined || valueToTest === null ||
                    lowerBound === undefined || lowerBound === null ||
                    upperBound === undefined || upperBound === null) {
                    return false;
                }

                // Parse bounds with scalar context
                lowerBound = tosca.parseScalarOrVersionBound(lowerBound, valueToTest);
                upperBound = tosca.parseScalarOrVersionBound(upperBound, valueToTest);

                return (tosca.compare(valueToTest, lowerBound) >= 0) && 
                       (tosca.compare(valueToTest, upperBound) <= 0);
            }
            
            // TOSCA 1.3 style: [lower, upper] (two bounds directly)
            if (typeof firstElement === 'number' || typeof firstElement === 'string') {
                let lowerBound = firstElement;
                let upperBound = secondElement;
                
                // Parse bounds with scalar context
                lowerBound = tosca.parseScalarOrVersionBound(lowerBound, currentPropertyValue);
                upperBound = tosca.parseScalarOrVersionBound(upperBound, currentPropertyValue);

                return (tosca.compare(currentPropertyValue, lowerBound) >= 0) && 
                       (tosca.compare(currentPropertyValue, upperBound) <= 0);
            }
        }
        
        return false;
    }
    
    return false;
};