package compiler

import (
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca/problems"
)

// This has the same result as tosca.coerce in JavaScript
func Coerce(clout_ *clout.Clout, problems_ *problems.Problems) *clout.Clout {
	context := js.NewContext("tosca.coerce", log, false, "yaml", "")
	err := context.Exec(clout_, "tosca.coerce")
	if err != nil {
		problems_.ReportError(err)
	}

	// Convert all values to JavaScript coercibles
	//TransformValues(ValueToCoercible, context, problems_)

	// Coerce all coercibles
	//TransformValues(CoerceValue, context, problems_)

	return clout_
}

// TODO: remove

// ValueTransformer signature
func ValueToCoercible(value interface{}, site interface{}, source interface{}, target interface{}, context *js.CloutContext) (interface{}, bool, error) {
	var err error
	value, err = context.NewCoercible(value, site, source, target)
	return value, true, err
}

// ValueTransformer signature
func CoerceValue(value interface{}, site interface{}, source interface{}, target interface{}, context *js.CloutContext) (interface{}, bool, error) {
	coercible, ok := value.(js.Coercible)
	if !ok {
		return nil, false, nil
	}
	var err error
	value, err = coercible.Coerce()
	return value, true, err
}
