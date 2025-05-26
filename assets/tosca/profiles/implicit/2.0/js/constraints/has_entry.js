// TOSCA 2.0 operator: has_entry
const tosca = require('tosca.lib.utils');

exports.validate = function() {
    // Extract the actual values we need to compare
    let container, entryValue;
    
    if (arguments.length === 2) {
        // Simple case: container and entry value
        container = arguments[0];
        entryValue = arguments[1];
    } else if (arguments.length >= 3) {
        // When function calls are involved, the last two arguments 
        // contain the values we need to compare
        container = arguments[arguments.length - 2];
        entryValue = arguments[arguments.length - 1];
    } else {
        throw new Error("has_entry requires at least 2 arguments");
    }
    
    // Validate arguments
    if (container === undefined || container === null) {
        return false;
    }
    
    if (entryValue === undefined || entryValue === null) {
        return false;
    }
    
    // Handle list case
    if (Array.isArray(container)) {
        // Check if entryValue exists as an element in the list
        for (let i = 0; i < container.length; i++) {
            if (tosca.deepEqual(container[i], entryValue)) {
                return true;
            }
        }
        return false;
    }
    
    // Handle map case
    if (typeof container === 'object' && container !== null && !Array.isArray(container)) {
        // Check if entryValue exists as a value in any of the key-value pairs
        for (let key in container) {
            if (container.hasOwnProperty(key)) {
                if (tosca.deepEqual(container[key], entryValue)) {
                    return true;
                }
            }
        }
        return false;
    }
    
    // Invalid container type
    return false;
};