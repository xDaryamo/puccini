package js

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/dop251/goja"
	"github.com/op/go-logging"
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common/terminal"
)

//
// Context
//

type Context struct {
	Quiet  bool
	Format string
	Pretty bool
	Output string
	Log    *Log
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Writer

	programCache sync.Map
}

func NewContext(name string, logger *logging.Logger, quiet bool, format string, pretty bool, output string) *Context {
	return &Context{
		Quiet:  quiet,
		Format: format,
		Pretty: pretty,
		Output: output,
		Log:    NewLog(logger, name),
		Stdout: terminal.Stdout,
		Stderr: terminal.Stderr,
		Stdin:  os.Stdin,
	}
}

func (self *Context) NewCloutRuntime(clout_ *clout.Clout, apis map[string]interface{}) *goja.Runtime {
	runtime := goja.New()
	runtime.SetFieldNameMapper(mapper)

	runtime.Set("puccini", self.NewPucciniAPI())

	runtime.Set("clout", self.NewCloutAPI(clout_, runtime))

	for name, api := range apis {
		runtime.Set(name, api)
	}

	return runtime
}

func (self *Context) GetProgram(name string, scriptlet string) (*goja.Program, error) {
	p, ok := self.programCache.Load(scriptlet)
	if !ok {
		program, err := goja.Compile(name, scriptlet, true)
		if err != nil {
			return nil, err
		}
		p, _ = self.programCache.LoadOrStore(scriptlet, program)
	}

	return p.(*goja.Program), nil
}

func (self *Context) Exec(clout_ *clout.Clout, scriptletName string, apis map[string]interface{}) error {
	scriptlet, err := GetScriptlet(scriptletName, clout_)
	if err != nil {
		return err
	}

	program, err := self.GetProgram(scriptletName, scriptlet)
	if err != nil {
		return err
	}

	runtime := self.NewCloutRuntime(clout_, apis)

	_, err = runtime.RunProgram(program)
	return UnwrapException(err)
}

func (self *Context) Failf(f string, args ...interface{}) {
	if !self.Quiet {
		fmt.Fprintln(self.Stderr, terminal.ColorError(fmt.Sprintf(f, args...)))
	}
	atexit.Exit(1)
}

func (self *Context) FailOnError(err error) {
	if err != nil {
		self.Failf("%s", err)
	}
}
