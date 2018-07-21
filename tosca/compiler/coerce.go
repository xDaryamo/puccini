package compiler

import (
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca/problems"
)

func Coerce(c *clout.Clout, p *problems.Problems) *clout.Clout {
	var err error
	c, err = c.Normalize()
	if err != nil {
		p.ReportError(err)
		return c
	}

	context := js.NewContext("coerce", log, false, "yaml", "")
	cloutContext, _ := context.NewCloutContext(c)

	Transform(PrepareValue, cloutContext, p)
	Transform(CoerceValue, cloutContext, p)

	return c
}

// Transformer signature
func PrepareValue(value interface{}, site interface{}, source interface{}, target interface{}, context *js.CloutContext) (interface{}, bool, error) {
	var err error
	value, err = context.NewCoercible(value, site, source, target)
	return value, true, err
}

// Transformer signature
func CoerceValue(value interface{}, site interface{}, source interface{}, target interface{}, context *js.CloutContext) (interface{}, bool, error) {
	coercible, ok := value.(js.Coercible)
	if !ok {
		return nil, false, nil
	}
	var err error
	value, err = coercible.Coerce()
	return value, true, err
}
