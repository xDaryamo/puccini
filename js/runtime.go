package js

import (
	"fmt"
	"os"
	"reflect"

	"github.com/beevik/etree"
	"github.com/dop251/goja"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
)

func NewRuntime(name string) *goja.Runtime {
	runtime := goja.New()
	runtime.SetFieldNameMapper(mapper)
	runtime.Set("stdout", os.Stdout)
	runtime.Set("stderr", os.Stderr)
	runtime.Set("stdin", os.Stdin)
	runtime.Set("log", format.NewLog(log, name))
	runtime.Set("sprintf", fmt.Sprintf)
	runtime.Set("timestamp", common.Timestamp)
	runtime.Set("newKey", clout.NewKey)
	runtime.Set("newXmlDocument", etree.NewDocument)
	return runtime
}

func NewCloutRuntime(name string, c *clout.Clout) *goja.Runtime {
	runtime := NewRuntime(name)
	runtime.Set("clout", NewClout(c, runtime))
	return runtime
}

//
// Mapper
//

var mapper Mapper

type Mapper struct{}

func (self Mapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return ToJavaScriptStyle(f.Name)
}

func (self Mapper) MethodName(t reflect.Type, m reflect.Method) string {
	return ToJavaScriptStyle(m.Name)
}
