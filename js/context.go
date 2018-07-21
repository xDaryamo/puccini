package js

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
	"github.com/op/go-logging"
	"github.com/tliron/puccini/format"
)

//
// Context
//

type Context struct {
	Quiet     bool
	ArdFormat string
	Output    string
	Log       *format.Log
	Stdout    *os.File
	Stderr    *os.File
	Stdin     *os.File
}

func NewContext(name string, logger *logging.Logger, quiet bool, ardFormat string, output string) *Context {
	return &Context{
		Quiet:     quiet,
		ArdFormat: ardFormat,
		Output:    output,
		Log:       format.NewLog(logger, name),
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Stdin:     os.Stdin,
	}
}

func (self *Context) NewRuntime() *goja.Runtime {
	runtime := goja.New()
	runtime.SetFieldNameMapper(mapper)
	runtime.Set("puccini", self.NewPucciniApi())
	return runtime
}

func (self *Context) Errorf(f string, args ...interface{}) {
	if !self.Quiet {
		fmt.Fprintf(self.Stderr, f+"\n", args...)
	}
	os.Exit(1)
}

func (self *Context) ValidateError(err error) {
	if err != nil {
		self.Errorf("%s", err)
	}
}
