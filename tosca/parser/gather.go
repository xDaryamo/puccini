package parser

import (
	"regexp"
	"strings"

	"github.com/tliron/puccini/tosca"
)

func (self *Context) Gather(pattern string) tosca.EntityPtrs {
	var entityPtrs tosca.EntityPtrs

	re := compileGatherPattern(pattern)

	self.Traverse("gather", func(entityPtr interface{}) bool {
		context := tosca.GetContext(entityPtr)

		if re.MatchString(context.Path.String()) {
			found := false
			for _, entityPtr_ := range entityPtrs {
				if entityPtr_ == entityPtr {
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

func compileGatherPattern(pattern string) *regexp.Regexp {
	split := strings.Split(pattern, "*")
	last := len(split) - 1
	var reString string
	for index, s := range split {
		reString += regexp.QuoteMeta(s)
		if index != last {
			reString += ".*"
		}
	}
	return regexp.MustCompile(reString)
}
