package js

import (
	"io"
	"os"
	"sync"

	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/terminal"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// Context
//

type Context struct {
	Arguments  map[string]string
	Quiet      bool
	Format     string
	Strict     bool
	Pretty     bool
	Base64     bool
	Output     string
	Log        commonlog.Logger
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Writer
	Stylist    *terminal.Stylist
	URLContext *exturl.Context

	programCache sync.Map
}

func NewContext(name string, log commonlog.Logger, arguments map[string]string, quiet bool, format string, strict bool, pretty bool, base64 bool, output string, urlContext *exturl.Context) *Context {
	if arguments == nil {
		arguments = make(map[string]string)
	}

	return &Context{
		Arguments:  arguments,
		Quiet:      quiet,
		Format:     format,
		Strict:     strict,
		Pretty:     pretty,
		Base64:     base64,
		Output:     output,
		Log:        commonlog.NewScopeLogger(log, name),
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Stdin:      os.Stdin,
		Stylist:    terminal.DefaultStylist,
		URLContext: urlContext,
	}
}

func (self *Context) NewEnvironment(clout *cloutpkg.Clout, apis map[string]any) *commonjs.Environment {
	environment := commonjs.NewEnvironment(self.URLContext, nil)

	environment.CreateResolver = func(url exturl.URL, context *commonjs.Context) commonjs.ResolveFunc {
		return func(id string, raw bool) (exturl.URL, error) {
			if scriptlet, err := GetScriptlet(id, clout); err == nil {
				url := self.URLContext.NewInternalURL(id)
				url.SetContent(scriptlet)
				return url, nil
			} else {
				return nil, err
			}
		}
	}

	environment.Extensions = append(environment.Extensions, commonjs.Extension{
		Name: "puccini",
		Create: func(context *commonjs.Context) goja.Value {
			return context.Environment.Runtime.ToValue(self.NewPucciniAPI())
		},
	})

	environment.Extensions = append(environment.Extensions, commonjs.Extension{
		Name: "clout",
		Create: func(context *commonjs.Context) goja.Value {
			return context.Environment.Runtime.ToValue(self.NewCloutAPI(clout, context))
		},
	})

	for name, api := range apis {
		environment.Extensions = append(environment.Extensions, commonjs.Extension{
			Name: name,
			Create: func(context *commonjs.Context) goja.Value {
				return context.Environment.Runtime.ToValue(api)
			},
		})
	}

	return environment
}

func (self *Context) Require(clout *cloutpkg.Clout, scriptletName string, apis map[string]any) (*goja.Object, error) {
	environment := self.NewEnvironment(clout, apis)
	if r, err := environment.RequireID(scriptletName); err == nil {
		return r, nil
	} else {
		return r, UnwrapException(err)
	}
}
