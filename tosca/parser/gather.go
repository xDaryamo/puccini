package parser

import (
	"regexp"
	"strings"

	"github.com/tliron/puccini/tosca"
)

func (self *Context) Gather(path string) tosca.EntityPtrs {
	var entityPtrs tosca.EntityPtrs

	split := strings.Split(path, "*")
	last := len(split) - 1
	var reString string
	for index, s := range split {
		reString += regexp.QuoteMeta(s)
		if index != last {
			reString += ".*"
		}
	}
	re := regexp.MustCompile(reString)

	self.Traverse("gather", func(entityPtr interface{}) bool {
		context := tosca.GetContext(entityPtr)

		if re.MatchString(context.Path.String()) {
			found := false
			for _, e := range entityPtrs {
				if e == entityPtr {
					found = true
					break
				}
			}
			if !found {
				entityPtrs = append(entityPtrs, entityPtr)
			}
		}

		return true
	})

	return entityPtrs
}
