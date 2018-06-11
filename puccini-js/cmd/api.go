package cmd

import (
	"github.com/dop251/goja"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/js"
)

//
// Puccini
//

type Puccini struct {
	c *clout.Clout
}

func NewPuccini(c *clout.Clout) *Puccini {
	return &Puccini{c}
}

func (self *Puccini) GetPlugins(name string) []goja.Value {
	plugins, err := js.GetPlugins(name, self.c)
	common.ValidateError(err)
	return plugins
}

func (self *Puccini) Write(data interface{}) {
	if !common.Quiet || (output != "") {
		err := format.WriteOrPrint(data, ardFormat, true, output)
		common.ValidateError(err)
	}
}
