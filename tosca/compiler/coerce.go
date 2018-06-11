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

	Transform(PrepareValue, c, p)
	Transform(CoerceValue, c, p)

	return c
}

// Transformer signature
func PrepareValue(value interface{}, site interface{}, source interface{}, target interface{}, c *clout.Clout) (interface{}, bool, error) {
	var err error
	value, err = js.NewCoercible(value, site, source, target, c)
	return value, true, err
}

// Transformer signature
func CoerceValue(value interface{}, site interface{}, source interface{}, target interface{}, c *clout.Clout) (interface{}, bool, error) {
	coercible, ok := value.(js.Coercible)
	if !ok {
		return nil, false, nil
	}
	var err error
	value, err = coercible.Coerce()
	return value, true, err
}
