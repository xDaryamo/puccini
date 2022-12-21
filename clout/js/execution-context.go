package js

import (
	"fmt"
)

//
// ExecutionContext
//

type ExecutionContext struct {
	CloutContext *CloutContext
	Site         any
	Source       any
	Target       any
}

func (self *CloutContext) NewExecutionContext(site any, source any, target any) *ExecutionContext {
	return &ExecutionContext{
		CloutContext: self,
		Site:         site,
		Source:       source,
		Target:       target,
	}
}

func (self *ExecutionContext) Call(scriptletName string, functionName string, arguments []any) (any, error) {
	if exports, err := self.CloutContext.JSContext.Environment.RequireID(scriptletName); err == nil {
		if function := exports.Get(functionName); function != nil {
			return CallGojaFunction(self.CloutContext.JSContext.Environment.Runtime, function, self, arguments)
		} else {
			return nil, fmt.Errorf("function not exported from module: %s", functionName)
		}
	} else {
		return nil, err
	}
}
