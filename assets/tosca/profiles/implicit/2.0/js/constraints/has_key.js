// TOSCA 2.0 constraint: has_key

const tosca = require('tosca.lib.utils');

exports.validate = function(currentPropertyValue) {
    const parsed = tosca.parseComparisonArguments(currentPropertyValue, arguments);
    if (!parsed) {
        return false;
    }
    
    const mapValue = parsed.val1;
    const keyToFind = parsed.val2;
    
    if (!mapValue || typeof mapValue !== 'object' || Array.isArray(mapValue)) {
        return false;
    }
    
    return mapValue.hasOwnProperty(keyToFind);
};