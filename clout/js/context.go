package js

import (
	"io"
	"os"
	"sync"

	"github.com/dop251/goja"
	"github.com/tliron/kutil/js"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
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
	Output     string
	Log        logging.Logger
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Writer
	Stylist    *terminal.Stylist
	URLContext *urlpkg.Context

	programCache sync.Map
}

func NewContext(name string, log logging.Logger, arguments map[string]string, quiet bool, format string, strict bool, pretty bool, output string, urlContext *urlpkg.Context) *Context {
	if arguments == nil {
		arguments = make(map[string]string)
	}

	return &Context{
		Arguments:  arguments,
		Quiet:      quiet,
		Format:     format,
		Strict:     strict,
		Pretty:     pretty,
		Output:     output,
		Log:        logging.NewScopeLogger(log, name),
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Stdin:      os.Stdin,
		Stylist:    terminal.DefaultStylist,
		URLContext: urlContext,
	}
}

func (self *Context) NewEnvironment(clout *cloutpkg.Clout, apis map[string]any) *js.Environment {
	environment := js.NewEnvironment(self.URLContext, nil)

	environment.CreateResolver = func(url urlpkg.URL, context *js.Context) js.ResolveFunc {
		return func(id string, raw bool) (urlpkg.URL, error) {
			if scriptlet, err := GetScriptlet(id, clout); err == nil {
				url := urlpkg.NewInternalURL(id, self.URLContext)
				url.SetContent(scriptlet)
				return url, nil
			} else {
				return nil, err
			}
		}
	}

	environment.Extensions = append(environment.Extensions, js.Extension{
		Name: "puccini",
		Create: func(context *js.Context) goja.Value {
			return context.Environment.Runtime.ToValue(self.NewPucciniAPI())
		},
	})

	environment.Extensions = append(environment.Extensions, js.Extension{
		Name: "clout",
		Create: func(context *js.Context) goja.Value {
			return context.Environment.Runtime.ToValue(self.NewCloutAPI(clout, context))
		},
	})

	for name, api := range apis {
		environment.Extensions = append(environment.Extensions, js.Extension{
			Name: name,
			Create: func(context *js.Context) goja.Value {
				return context.Environment.Runtime.ToValue(api)
			},
		})
	}

	return environment
}

func (self *Context) Require(clout *cloutpkg.Clout, scriptletName string, apis map[string]any) (*goja.Object, error) {
	environment := self.NewEnvironment(clout, apis)
	return environment.RequireID(scriptletName)
}
