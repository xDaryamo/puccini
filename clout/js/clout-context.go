package js

import (
	"github.com/tliron/kutil/js"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// CloutContext
//

type CloutContext struct {
	Context   *Context
	Clout     *cloutpkg.Clout
	JSContext *js.Context
}

func (self *Context) NewCloutContext(clout *cloutpkg.Clout, jsContext *js.Context) *CloutContext {
	return &CloutContext{
		Context:   self,
		Clout:     clout,
		JSContext: jsContext,
	}
}
