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
	Pretty    bool
	Output    string
	Log       *format.Log
	Stdout    *os.File
	Stderr    *os.File
	Stdin     *os.File
}

func NewContext(name string, logger *logging.Logger, quiet bool, ardFormat string, pretty bool, output string) *Context {
	return &Context{
		Quiet:     quiet,
		ArdFormat: ardFormat,
		Pretty:    pretty,
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

func (self *Context) Failf(f string, args ...interface{}) {
	if !self.Quiet {
		fmt.Fprintln(self.Stderr, format.ColorError(fmt.Sprintf(f, args...)))
	}
	os.Exit(1)
}

func (self *Context) FailOnError(err error) {
	if err != nil {
		self.Failf("%s", err)
	}
}
