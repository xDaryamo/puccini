exports.isTosca = function (o, kind) {
  if (o.metadata === undefined) return false;
  o = o.metadata['puccini'];
  if (o === undefined) return false;
  if (o.version !== '1.0') return false;
  if (kind !== undefined) return kind === o.kind;
  return true;
};

exports.isNodeTemplate = function (vertex, typeName) {
  if (exports.isTosca(vertex, 'NodeTemplate')) {
    if (typeName !== undefined) return typeName in vertex.properties.types;
    return true;
  }
  return false;
};

exports.setOutputValue = function (name, value) {
  if (clout.properties.tosca === undefined) return false;
  let output = clout.properties.tosca.outputs[name];
  if (output === undefined) return false;

  if (output.$type && output.$type.type)
    switch (output.$type.type.name) {
      case 'boolean':
        value = value === 'true';
        break;
      case 'integer':
        value = parseInt(value);
        break;
      case 'float':
        value = parseFloat(value);
        break;
    }

  output.$value = value;
  return true;
};

exports.getPolicyTargets = function (vertex) {
  let targets = [];

  function addTarget(target) {
    for (let t = 0, l = targets.length; t < l; t++) if (targets[t].name === target.name) return;
    targets.push(target);
  }

  for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
    let edge = vertex.edgesOut[e];
    if (exports.isTosca(edge, 'NodeTemplateTarget')) targets.push(clout.vertexes[edge.targetID].properties);
    else if (toexportssca.isTosca(edge, 'GroupTarget')) {
      let members = exports.getGroupMembers(clout.vertexes[edge.targetID]);
      for (let m = 0, ll = members.length; m < ll; m++) addTarget(members[m]);
    }
  }
  return targets;
};

exports.getGroupMembers = function (vertex) {
  let members = [];
  for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
    let edge = vertex.edgesOut[e];
    if (exports.isTosca(edge, 'Member')) members.push(clout.vertexes[edge.targetID].properties);
  }
  return members;
};

exports.addHistory = function (description) {
  let metadata = clout.metadata;
  if (metadata === undefined) metadata = clout.metadata = {};
  let history = metadata.history;
  if (history === undefined) history = [];
  else history = history.slice(0);
  history.push({
    timestamp: util.now().string(),
    description: description,
  });
  metadata.history = history;
};

exports.getNestedValue = function(singular, plural, args) {
    args = Array.prototype.slice.call(args);
    let length = args.length;
    
    if (length < 2)
        throw 'must have at least 2 arguments';
    
    // Check if this looks like TOSCA 2.0 syntax (contains keywords)
    let isTosca20 = false;
    for (let i = 1; i < args.length; i++) {
        if (typeof args[i] === 'string' && 
            ['RELATIONSHIP', 'TARGET', 'SOURCE', 'CAPABILITY'].includes(args[i])) {
            isTosca20 = true;
            break;
        }
    }
    
    // If not explicitly TOSCA 2.0, check if this could be a relationship access
    if (!isTosca20 && length >= 3) {
        let vertex = exports.getModelableEntity.call(this, args[0]);
        let nodeTemplate = vertex.properties;
        let relationshipName = args[1];
        
        // Check if args[1] is a relationship name rather than a capability or direct property/attribute
        let isRelationship = false;
        if (!(relationshipName in nodeTemplate.capabilities) && !(relationshipName in nodeTemplate[plural])) {
            // Check if it's a relationship
            for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
                let edge = vertex.edgesOut[e];
                
                if (exports.isTosca(edge, 'Relationship')) {
                    let relationship = edge.properties;
                    
                    if (relationship.name === relationshipName) {
                        isRelationship = true;
                        break;
                    }
                }
            }
        }
        
        if (isRelationship) {
            isTosca20 = true;
        }
    }
    
    if (isTosca20) {
        return exports.getNestedValueTosca20.call(this, singular, plural, args);
    }
    
    // Original implementation for backward compatibility
    let vertex = exports.getModelableEntity.call(this, args[0]);
    let nodeTemplate = vertex.properties;
    let value = nodeTemplate[plural];
    let a = 1;
    let arg = args[a];
    let nextArg = args[a+1];
    let count = 0;
    
    if (arg in nodeTemplate.capabilities) {
        value = nodeTemplate.capabilities[arg][plural];
        singular = util.sprintf('capability %q %s', arg, singular);
        arg = args[++a];
    } else {
        for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
            let edge = vertex.edgesOut[e];
            if (!exports.isTosca(edge, 'Relationship'))
                continue;
            let relationship = edge.properties;
            if (relationship.name === arg) {
                if (count++ === nextArg) {
                    value = relationship[plural];
                    singular = util.sprintf('relationship %q %s', arg, singular);
                    a += 2;
                    arg = args[a];
                    break;
                }
            }
        }
    }
    
    if ((typeof value === 'object') && (value !== null) && (arg in value))
        value = value[arg];
    else
        throw util.sprintf('%s %q not found in %q', singular, arg, nodeTemplate.name);
    
    value = clout.coerce(value);
    
    for (let i = a + 1; i < length; i++) {
        arg = args[i];
        if ((typeof value === 'object') && (value !== null) && (arg in value))
            value = value[arg];
        else
            throw util.sprintf('nested %s %q not found in %q', singular, args.slice(a, i+1).join('.'), nodeTemplate.name);
    }
    return value;
};

// TOSCA 2.0 path traversal following BNF grammar
exports.getNestedValueTosca20 = function(singular, plural, args) {
    args = Array.prototype.slice.call(args);
    let length = args.length;
    if (length < 2)
        throw 'must have at least 2 arguments';
    
    let currentVertex = exports.getModelableEntity.call(this, args[0]);
    let pathIndex = 1;
    
    // Parse according to BNF grammar: <tosca_path>
    // <tosca_path> ::= <node_symbolic_name>, <idx>, <node_context> |
    //                  SELF, <node_context> |
    //                  SELF, <rel_context>
    
    while (pathIndex < args.length) {
        let step = args[pathIndex];
        
        // Check if this is an explicit keyword
        if (['RELATIONSHIP', 'TARGET', 'SOURCE', 'CAPABILITY'].includes(step)) {
            switch (step) {
                case 'RELATIONSHIP':
                    // <node_context> ::= RELATIONSHIP, <requirement_name>, <idx>, <rel_context>
                    pathIndex++;
                    if (pathIndex >= args.length)
                        throw 'RELATIONSHIP keyword must be followed by requirement name';
                    
                    let requirementName = args[pathIndex];
                    pathIndex++;
                    
                    // Check for optional index
                    let relationshipIndex = 0;
                    if (pathIndex < args.length && typeof args[pathIndex] === 'number') {
                        relationshipIndex = args[pathIndex];
                        pathIndex++;
                    }
                    
                    currentVertex = exports.traverseToRelationship(currentVertex, requirementName, relationshipIndex);
                    break;
                    
                case 'TARGET':
                    // <rel_context> ::= TARGET, <node_context>
                    if (!exports.isTosca(currentVertex, 'Relationship'))
                        throw 'TARGET can only be used in relationship context';
                    
                    currentVertex = currentVertex.target;
                    pathIndex++;
                    break;
                    
                case 'SOURCE':
                    // <rel_context> ::= SOURCE, <node_context>
                    if (!exports.isTosca(currentVertex, 'Relationship'))
                        throw 'SOURCE can only be used in relationship context';
                    
                    currentVertex = currentVertex.source;
                    pathIndex++;
                    break;
                    
                case 'CAPABILITY':
                    // <node_context> ::= CAPABILITY, <capability_name>, <cap_context>
                    // <rel_context> ::= CAPABILITY, <cap_context>
                    pathIndex++;
                    if (pathIndex >= args.length)
                        throw 'CAPABILITY keyword must be followed by capability name';
                    
                    if (exports.isTosca(currentVertex, 'Relationship')) {
                        // In relationship context, CAPABILITY refers to the target capability
                        // This is handled by the relationship's capability property
                        pathIndex--; // Back up to handle capability access
                        break;
                    } else if (exports.isNodeTemplate(currentVertex)) {
                        let capabilityName = args[pathIndex];
                        pathIndex++;
                        
                        if (!(capabilityName in currentVertex.properties.capabilities))
                            throw util.sprintf('capability %q not found in node %q', capabilityName, currentVertex.properties.name);
                        
                        // Switch to capability context for remaining path
                        let capabilityTarget = currentVertex.properties.capabilities[capabilityName];
                        return exports.getNestedPropertyValue(capabilityTarget, plural, args, pathIndex, singular);
                    } else {
                        throw 'CAPABILITY can only be used in node or relationship context';
                    }
                    break;
            }
        } else {
            // Handle implicit syntax and property access
            if (exports.isNodeTemplate(currentVertex)) {
                let entityName = step;
                
                // First, check if it's a capability
                if (entityName in currentVertex.properties.capabilities) {
                    // This is a capability access
                    let capabilityTarget = currentVertex.properties.capabilities[entityName];
                    pathIndex++;
                    return exports.getNestedPropertyValue(capabilityTarget, plural, args, pathIndex, singular);
                } else {
                    // Check if it's a requirement name (implicit relationship access)
                    let requirementName = entityName;
                    pathIndex++;
                    
                    // Check for optional index
                    let relationshipIndex = 0;
                    if (pathIndex < args.length && typeof args[pathIndex] === 'number') {
                        relationshipIndex = args[pathIndex];
                        pathIndex++;
                    }
                    
                    // Try to find the relationship
                    try {
                        currentVertex = exports.traverseToRelationship(currentVertex, requirementName, relationshipIndex);
                    } catch (e) {
                        throw util.sprintf('neither capability nor relationship %q found in node %q', requirementName, currentVertex.properties.name);
                    }
                }
            } else if (exports.isTosca(currentVertex, 'Relationship')) {
                // We're in a relationship context, so remaining args are property access
                return exports.getNestedPropertyValue(currentVertex.properties, plural, args, pathIndex, singular);
            } else {
                throw util.sprintf('unexpected token in TOSCA path: %s', step);
            }
        }
    }
    
    // If we've consumed all arguments, we're accessing the entity itself
    // This shouldn't happen for property/attribute access
    throw 'invalid TOSCA path: no property specified';
};

// Helper function to get nested property values from an entity
exports.getNestedPropertyValue = function(entity, plural, args, startIndex, singular) {
    let value = entity[plural];
    let entityName = entity.name || 'entity';
    
    // Coerce the initial value to handle Puccini's Value object structure
    value = clout.coerce(value);
    
    for (let i = startIndex; i < args.length; i++) {
        let arg = args[i];
        
        if ((typeof value === 'object') && (value !== null) && (arg in value)) {
            value = value[arg];
            // Coerce each intermediate value to handle nested Value objects
            value = clout.coerce(value);
        } else {
            // Build a proper error message with the property path
            let propertyPath = [];
            for (let j = startIndex; j <= i; j++) {
                propertyPath.push(args[j]);
            }
            throw util.sprintf('%s %q not found in %q (looking for path: %s)', 
                singular, arg, entityName, propertyPath.join('.'));
        }
    }
    
    return value;
};

// Helper function to traverse to a relationship by requirement name and index
exports.traverseToRelationship = function(vertex, requirementName, relationshipIndex) {
    if (!exports.isNodeTemplate(vertex)) {
        throw util.sprintf('can only traverse relationships from node template context');
    }
    
    let relationshipFound = false;
    let relationshipCount = 0;
    
    for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
        let edge = vertex.edgesOut[e];
        if (!exports.isTosca(edge, 'Relationship'))
            continue;
        
        let relationship = edge.properties;
        if (relationship.name === requirementName) {
            if (relationshipCount === relationshipIndex) {
                return edge; // Return relationship context
            }
            relationshipCount++;
        }
    }
    
    throw util.sprintf('relationship %q (index %d) not found in node %q', requirementName, relationshipIndex, vertex.properties.name);
};

exports.getModelableEntity = function(entity) {
	let vertex;
	switch (entity) {
	case 'SELF':
		if (!this || !this.site)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = this.site;
		break;
	case 'SOURCE':
		if (!this || !this.source)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = this.source;
		break;
	case 'TARGET':
		if (!this || !this.target)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = this.target;
		break;
	case 'HOST':
		if (!this || !this.site)
			throw util.sprintf('%q cannot be used in this context', entity);
		vertex = exports.getHost(this.site);
		break;
	default:
		for (let vertexId in clout.vertexes) {
			let vertex = clout.vertexes[vertexId];
			if (exports.isNodeTemplate(vertex) && (vertex.properties.name === entity))
				return vertex;
		}
		vertex = {};
	}
	if (exports.isNodeTemplate(vertex))
		return vertex;
	else
		throw util.sprintf('%q node template not found', entity);
};

exports.getHost = function(vertex) {
    for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
        let edge = vertex.edgesOut[e];
        if (exports.isTosca(edge, 'Relationship')) {
            for (let typeName in edge.properties.types) {
                let type = edge.properties.types[typeName];
                if (type && type.metadata && type.metadata.role === 'host') {
                    return edge.target;
                }
            }
        }
    }
    if (exports.isNodeTemplate(vertex))
        throw util.sprintf('"HOST" not found for node template %q', vertex.properties.name);
    else
        throw '"HOST" not found';
};

exports.getComparable = function (v) {
  if (v === undefined || v === null) return null;
  let c = v.$number;
  if (c !== undefined) return c;
  c = v.$string;
  if (c !== undefined) return c;
  return v;
};

exports.getLength = function (v) {
  if (v.$string !== undefined) v = v.$string;
  let length = v.length;
  if (length === undefined) length = Object.keys(v).length;
  return length;
};

exports.compare = function (v1, v2) {
  let c = v1.$comparer;
  if (c === undefined) c = v2.$comparer;
  if (c !== undefined) return clout.call(c, 'compare', [v1, v2]);
  v1 = exports.getComparable(v1);
  v2 = exports.getComparable(v2);
  if (v1 == v2) return 0;
  else if (v1 < v2) return -1;
  else return 1;
};

// See: https://stackoverflow.com/a/45683145
exports.deepEqual = function (v1, v2) {
  if (v1 === v2) return true;

  if (exports.isPrimitive(v1) && exports.isPrimitive(v2)) return v1 === v2;

  if (Object.keys(v1).length !== Object.keys(v2).length) return false;

  for (let key in v1) {
    if (!(key in v2)) return false;
    if (!exports.deepEqual(v1[key], v2[key])) return false;
  }

  return true;
};

exports.isPrimitive = function(obj) {
    return obj !== Object(obj);
};

// Parse scalar value using context from another scalar object
exports.tryParseScalar = function(valueString, scalarObject) {
    if (typeof valueString !== 'string' || !scalarObject) {
        return null;
    }

    // Read scalar type information from various possible sources
    const units = scalarObject.units || 
                  scalarObject.Units || 
                  (scalarObject.scalarType && scalarObject.scalarType.Units) ||
                  (scalarObject.$scalarTypeInfo && scalarObject.$scalarTypeInfo.units);
    
    const canonicalUnit = scalarObject.canonicalUnit || 
                         scalarObject.CanonicalUnit ||
                         (scalarObject.scalarType && scalarObject.scalarType.CanonicalUnit) ||
                         (scalarObject.$scalarTypeInfo && scalarObject.$scalarTypeInfo.canonicalUnit);
    
    const baseType = scalarObject.baseType || 
                    scalarObject.BaseType || 
                    scalarObject.dataType ||
                    (scalarObject.scalarType && scalarObject.scalarType.DataTypeName) ||
                    (scalarObject.$scalarTypeInfo && scalarObject.$scalarTypeInfo.baseType);

    const dataTypeName = scalarObject.dataTypeName || 
                        scalarObject.DataTypeName || 
                        (scalarObject.scalarType && scalarObject.scalarType.Name) ||
                        (scalarObject.$scalarTypeInfo && scalarObject.$scalarTypeInfo.name) || 
                        '';

    const prefixes = scalarObject.prefixes ||
                    scalarObject.Prefixes ||
                    (scalarObject.scalarType && scalarObject.scalarType.Prefixes) ||
                    (scalarObject.$scalarTypeInfo && scalarObject.$scalarTypeInfo.prefixes);

    // Verify that we have all necessary information
    if (!units || !canonicalUnit || !baseType) {
        return null;
    }

    // Parse the scalar value string - handle both with and without spaces (TOSCA 1.3)
    let match = valueString.match(/^([+-]?[0-9]*\.?[0-9]+(?:[eE][+-]?[0-9]+)?)\s*(.+)$/);
    if (!match) {
        // Try without any space requirement for TOSCA 1.3 compatibility
        match = valueString.match(/^([+-]?[0-9]*\.?[0-9]+(?:[eE][+-]?[0-9]+)?)([a-zA-Z].*)$/);
        if (!match) {
            return null;
        }
    }

    const numberPart = parseFloat(match[1]);
    if (isNaN(numberPart)) {
        return null;
    }

    const unitPart = match[2];
    let multiplier = exports.findUnitMultiplier(unitPart, units, prefixes);

    if (multiplier === null) {
        return null;
    }

    // Calculate canonical value
    let canonicalNumber = numberPart * multiplier;
    let canonicalString;

    if (baseType === 'integer') {
        canonicalNumber = Math.round(canonicalNumber);
        canonicalString = String(canonicalNumber) + ' ' + canonicalUnit;
    } else {
        canonicalString = String(canonicalNumber) + ' ' + canonicalUnit;
    }

    // Create scalar object with all necessary properties
    return {
        $originalString: valueString,
        $number: canonicalNumber,
        $string: canonicalString,
        scalar: numberPart,
        unit: unitPart,
        baseType: baseType,
        canonicalUnit: canonicalUnit,
        dataTypeName: dataTypeName,
        units: units,
        prefixes: prefixes || {}
    };
};

// Find unit multiplier considering prefixes
exports.findUnitMultiplier = function(unitPart, units, prefixes) {
    // First try direct unit match
    let multiplier = units[unitPart];
    
    // Handle case-insensitive matching for TOSCA 1.3
    if (multiplier === undefined) {
        for (const [unit, mult] of Object.entries(units)) {
            if (unit.toLowerCase() === unitPart.toLowerCase()) {
                multiplier = mult;
                break;
            }
        }
    }

    // Handle prefixes for TOSCA 2.0
    if (multiplier === undefined && prefixes) {
        // Find the longest match: prefix+unit, then unit only
        let bestMatch = null;
        let bestPrefixLength = 0;
        
        for (const [unit, unitMultiplier] of Object.entries(units)) {
            // Check if unitPart ends with this unit
            if (unitPart.endsWith(unit)) {
                const potentialPrefix = unitPart.substring(0, unitPart.length - unit.length);
                
                if (potentialPrefix === '') {
                    // No prefix, direct unit
                    if (bestMatch === null || potentialPrefix.length > bestPrefixLength) {
                        bestMatch = { unit, unitMultiplier, prefix: '', prefixMultiplier: 1 };
                        bestPrefixLength = potentialPrefix.length;
                    }
                } else if (prefixes[potentialPrefix] !== undefined) {
                    // Prefix found
                    if (potentialPrefix.length > bestPrefixLength) {
                        bestMatch = { 
                            unit, 
                            unitMultiplier, 
                            prefix: potentialPrefix, 
                            prefixMultiplier: prefixes[potentialPrefix] 
                        };
                        bestPrefixLength = potentialPrefix.length;
                    }
                }
            }
        }
        
        if (bestMatch) {
            multiplier = bestMatch.unitMultiplier * bestMatch.prefixMultiplier;
        }
    }

    return multiplier !== undefined ? multiplier : null;
};

// Parse comparison arguments for constraint validators
exports.parseComparisonArguments = function(currentPropertyValue, argumentsArray) {
    if (argumentsArray.length === 3) {
        // Check if this is a map key validation scenario
        // In this case: args = [keyValue, entireMap, pattern]
        // We should only use: [keyValue, pattern]
        let val1 = argumentsArray[1];
        let val2 = argumentsArray[2];
        
        // If val1 is an object (the entire map) and we have a string currentPropertyValue,
        // this indicates map key validation - skip the map object
        if (typeof currentPropertyValue === 'string' && 
            typeof val1 === 'object' && val1 !== null && 
            typeof val2 === 'string') {
            return { val1: currentPropertyValue, val2: val2 };
        }
        
        // Handle "$value" substitution for val1
        if (val1 === '$value') {
            val1 = currentPropertyValue;
        }
        
        // Handle "$value" substitution for val2
        if (val2 === '$value') {
            val2 = currentPropertyValue;
        }
        
        // Handle scalar parsing: if one value is a scalar object and the other is a string
        if (val1 && typeof val1 === 'object' && val1.$number !== undefined && typeof val2 === 'string') {
            const parsed = exports.tryParseScalar(val2, val1);
            if (parsed) {
                val2 = parsed;
            }
        } else if (val2 && typeof val2 === 'object' && val2.$number !== undefined && typeof val1 === 'string') {
            const parsed = exports.tryParseScalar(val1, val2);
            if (parsed) {
                val1 = parsed;
            }
        }
        
        return { val1, val2 };
    } else if (argumentsArray.length === 2) {
        // Legacy calling convention: currentPropertyValue, compareValue
        const compareValue = argumentsArray[1];
        
        // Handle "$value" substitution
        let val2 = compareValue;
        if (compareValue === '$value') {
            val2 = currentPropertyValue;
        }
        
        // Handle scalar parsing: if currentPropertyValue is a scalar object and val2 is a string
        if (currentPropertyValue && typeof currentPropertyValue === 'object' && currentPropertyValue.$number !== undefined && typeof val2 === 'string') {
            const parsed = exports.tryParseScalar(val2, currentPropertyValue);
            if (parsed) {
                val2 = parsed;
            }
        } else if (val2 && typeof val2 === 'object' && val2.$number !== undefined && typeof currentPropertyValue === 'string') {
            const parsed = exports.tryParseScalar(currentPropertyValue, val2);
            if (parsed) {
                currentPropertyValue = parsed;
            }
        }
        
        return { val1: currentPropertyValue, val2: val2 };
    } else {
        return null;
    }
};

// Evaluate constraint argument with potential nested functions
exports.evaluateConstraintArgument = function(arg, currentPropertyValue) {
    // Handle string literal "$value"
    if (arg === '$value') {
        return currentPropertyValue;
    }
    
    if (typeof arg !== 'object' || arg === null || Array.isArray(arg)) {
        return arg;
    }
    
    const keys = Object.keys(arg);
    if (keys.length === 1) {
        const key = keys[0];
        if (key.startsWith('$')) {
            const functionName = key.substring(1);
            const functionArgs = arg[key];
            
            // SPECIAL HANDLING: For $value or {$value: ["path", ...]}
            if (functionName === 'value') {
                if (Array.isArray(functionArgs)) {
                    return exports.dereferencePathHelper(currentPropertyValue, functionArgs);
                } else if (functionArgs === null || functionArgs === undefined) {
                    return currentPropertyValue;
                } else {
                    return exports.dereferencePathHelper(currentPropertyValue, [functionArgs]);
                }
            }
            
            // Check if it's a validation operator (to handle recursive structure)
            if (exports.isValidationOperator(functionName)) {
                // Return the processed constraint for later evaluation
                if (Array.isArray(functionArgs)) {
                    const processedArgs = functionArgs.map(subArg => exports.evaluateConstraintArgument.call(this, subArg, currentPropertyValue));
                    return { [key]: processedArgs };
                } else {
                    return { [key]: exports.evaluateConstraintArgument.call(this, functionArgs, currentPropertyValue) };
                }
            }
            
            // Try to evaluate it as a function (length, concat, etc.)
            if (exports.isToscaFunction(functionName)) {
                try {
                    const functionModule = require('tosca.function.' + functionName);
                    const processedFunctionArgs = [];
                    const argsArray = Array.isArray(functionArgs) ? functionArgs : [functionArgs];
                    
                    for (const fnArg of argsArray) {
                        processedFunctionArgs.push(exports.evaluateConstraintArgument.call(this, fnArg, currentPropertyValue));
                    }
                    
                    return functionModule.evaluate.apply(this, processedFunctionArgs);
                } catch (e) {
                    // If the function can't be evaluated, return the processed argument
                    if (Array.isArray(functionArgs)) {
                        const processedArgs = functionArgs.map(subArg => exports.evaluateConstraintArgument.call(this, subArg, currentPropertyValue));
                        return { [key]: processedArgs };
                    } else {
                        return { [key]: exports.evaluateConstraintArgument.call(this, functionArgs, currentPropertyValue) };
                    }
                }
            }
        }
    }
    
    return arg;
};

// Validate a constraint subclause
exports.validateConstraintSubclause = function(operatorFunctionName, originalArgs, currentPropertyValue) {
    try {
        const validatorModule = require('tosca.validation.' + operatorFunctionName);

        if (validatorModule && typeof validatorModule.validate === 'function') {
            const subValidator = validatorModule.validate;
            
            // For constraints like $greater_or_equal: ["$value", "4 GB"]
            // We need to call: greater_or_equal.validate(currentPropertyValue, "$value", "4 GB")
            // The individual constraint will handle "$value" replacement and scalar parsing
            
            const result = subValidator.call(this, currentPropertyValue, ...originalArgs);
            return result;
        } else {
            return false; // Module or function not found
        }
    } catch (e) {
        return false; // Error during sub-clause validation
    }
};

// Helper function to dereference a path within an object
exports.dereferencePathHelper = function(obj, pathArray) {
    let current = obj;
    // Handle both array paths and single property access
    if (pathArray === undefined || pathArray === null) {
        return current;
    }
    
    // Ensure pathArray is an array, supporting {$value: "property"} as {$value: ["property"]}
    const path = Array.isArray(pathArray) ? pathArray : [pathArray];

    for (const key of path) {
        if (current === null || typeof current !== 'object') {
            return undefined; // Path not found or invalid intermediate value
        }
        
        // Handle both object property access and array index access
        if (Array.isArray(current)) {
            const index = parseInt(key);
            if (isNaN(index) || index < 0 || index >= current.length) {
                return undefined;
            }
            current = current[index];
        } else {
            if (!current.hasOwnProperty(key)) {
                return undefined;
            }
            current = current[key];
        }
    }
    return current;
};

// Check if a function name corresponds to a validation operator
exports.isValidationOperator = function(functionName) {
    try {
        const validatorModule = require('tosca.validation.' + functionName);
        return validatorModule && typeof validatorModule.validate === 'function';
    } catch (e) {
        return false;
    }
};

// Check if a function name corresponds to a TOSCA function
exports.isToscaFunction = function(functionName) {
    try {
        const functionModule = require('tosca.function.' + functionName);
        return functionModule && typeof functionModule.evaluate === 'function';
    } catch (e) {
        return false;
    }
};

// Parse constraint subclause from map
exports.parseConstraintSubclause = function(subclauseMap) {
    const operatorKey = Object.keys(subclauseMap)[0];
    if (!operatorKey) {
        return null; // Empty or malformed sub-clause
    }

    const operatorFunctionName = operatorKey.startsWith('$') ? operatorKey.substring(1) : operatorKey;
    let originalOperatorArgs = subclauseMap[operatorKey];

    if (!Array.isArray(originalOperatorArgs)) {
        originalOperatorArgs = [originalOperatorArgs];
    }

    return {
        operatorFunctionName,
        originalOperatorArgs
    };
};

exports.validateRelationshipAttributeAccess = function(vertex, requirementName, attributeName, relationshipIndex) {
    if (!exports.isNodeTemplate(vertex)) {
        throw util.sprintf('node template %q not found', vertex.properties ? vertex.properties.name : 'unknown');
    }
    
    // Count and find relationships by requirement name
    let relationshipCount = 0;
    let targetRelationship = null;
    let relationships = [];
    
    for (let e = 0, l = vertex.edgesOut.size(); e < l; e++) {
        let edge = vertex.edgesOut[e];
        if (!exports.isTosca(edge, 'Relationship'))
            continue;
        
        let relationship = edge.properties;
        if (relationship.name === requirementName) {
            relationships.push(relationship);
            if (relationshipCount === relationshipIndex) {
                targetRelationship = relationship;
            }
            relationshipCount++;
        }
    }
    
    // Validate that we have relationships for this requirement
    if (relationshipCount === 0) {
        throw util.sprintf('no relationships found for requirement %q in node template %q', 
            requirementName, vertex.properties.name);
    }
    
    // Validate relationship index
    if (relationshipIndex >= relationshipCount) {
        throw util.sprintf('relationship index %d exceeds available relationships (%d) for requirement %q in node template %q', 
            relationshipIndex, relationshipCount, requirementName, vertex.properties.name);
    }
    
    if (!targetRelationship) {
        throw util.sprintf('relationship %q (index %d) not found in node template %q', 
            requirementName, relationshipIndex, vertex.properties.name);
    }
    
    // Validate that the attribute exists in the relationship
    if (!targetRelationship.attributes || !(attributeName in targetRelationship.attributes)) {
        // Get relationship type name for better error message
        let relationshipTypeName = 'unknown';
        if (targetRelationship.types) {
            relationshipTypeName = Object.keys(targetRelationship.types)[0] || 'unknown';
        }
        throw util.sprintf('attribute %q not found in relationship %q of type %q', 
            attributeName, requirementName, relationshipTypeName);
    }
    
    return true;
};

// Parse scalar unit or version bound, with proper type conversion for TOSCA 1.3
exports.parseScalarOrVersionBound = function(bound, contextValue) {
    if (typeof bound !== 'string') {
        return bound; // Already parsed or not a string
    }
    
    // Try to parse as scalar unit using context from the value being tested
    if (contextValue && (contextValue.$number !== undefined || contextValue.units)) {
        const parsed = exports.tryParseScalar(bound, contextValue);
        if (parsed) {
            return parsed;
        }
    }
    
    // Try to parse as version if context value is a version
    if (contextValue && contextValue.$comparer === 'tosca.comparer.version') {
        // Check if bound looks like a version string
        if (/^\d+(\.\d+)*(\.\w+(-\d+)?)?$/.test(bound)) {
            try {
                // Create a version object similar to contextValue
                return {
                    $comparer: 'tosca.comparer.version',
                    $originalString: bound,
                    $string: bound,
                    // Parse version components
                    ...exports.parseVersionString(bound)
                };
            } catch (e) {
                // Fall back to string comparison
                return bound;
            }
        }
    }
    
    return bound; // Return as-is if no parsing applies
};

// Helper to parse version string components
exports.parseVersionString = function(versionStr) {
    const parts = versionStr.split('.');
    const result = {
        major: parseInt(parts[0]) || 0,
        minor: parseInt(parts[1]) || 0,
        fix: parseInt(parts[2]) || 0
    };
    
    // Handle qualifier and build number (e.g., "beta-4")
    if (parts.length > 3) {
        const qualifierPart = parts[3];
        const dashIndex = qualifierPart.indexOf('-');
        if (dashIndex !== -1) {
            result.qualifier = qualifierPart.substring(0, dashIndex);
            result.build = parseInt(qualifierPart.substring(dashIndex + 1)) || 0;
        } else {
            result.qualifier = qualifierPart;
            result.build = 0;
        }
    }
    
    return result;
};
