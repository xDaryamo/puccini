// [TOSCA-v2.0] @ 10.2.4.2

exports.evaluate = function() {
    if (arguments.length === 0) {
        throw 'The $intersection function expects at least one argument.';
    }
    
    // Validate that all arguments are lists/arrays
    for (let i = 0; i < arguments.length; i++) {
        if (!Array.isArray(arguments[i])) {
            throw 'The $intersection function argument ' + i + ' must be a list; got ' + (typeof arguments[i]);
        }
    }
    
    // Special case: only one list - remove duplicates (same as union for single list)
    if (arguments.length === 1) {
        let result = [];
        let seen = new Set();
        
        for (let j = 0; j < arguments[0].length; j++) {
            let element = arguments[0][j];
            let key = JSON.stringify(element);
            
            if (!seen.has(key)) {
                seen.add(key);
                result.push(element);
            }
        }
        
        return result;
    }
    
    // Multiple lists: find intersection
    let result = [];
    let firstList = arguments[0];
    
    // For each element in the first list
    for (let i = 0; i < firstList.length; i++) {
        let element = firstList[i];
        let elementKey = JSON.stringify(element);
        let foundInAll = true;
        
        // Check if this element exists in all other lists
        for (let listIndex = 1; listIndex < arguments.length; listIndex++) {
            let currentList = arguments[listIndex];
            let foundInCurrentList = false;
            
            for (let j = 0; j < currentList.length; j++) {
                if (JSON.stringify(currentList[j]) === elementKey) {
                    foundInCurrentList = true;
                    break;
                }
            }
            
            if (!foundInCurrentList) {
                foundInAll = false;
                break;
            }
        }
        
        // If element is found in all lists and not already in result
        if (foundInAll) {
            let alreadyInResult = false;
            for (let k = 0; k < result.length; k++) {
                if (JSON.stringify(result[k]) === elementKey) {
                    alreadyInResult = true;
                    break;
                }
            }
            
            if (!alreadyInResult) {
                result.push(element);
            }
        }
    }
    
    return result;
};