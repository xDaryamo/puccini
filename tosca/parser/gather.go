package parser

import (
	"regexp"
	"strings"

	"github.com/tliron/puccini/tosca/parsing"
)

func (self *Context) Gather(pattern string) parsing.EntityPtrs {
	var entityPtrs parsing.EntityPtrs

	re := compileGatherPattern(pattern)

	self.TraverseEntities(logGather, nil, func(entityPtr parsing.EntityPtr) bool {
		context := parsing.GetContext(entityPtr)

		if re.MatchString(context.Path.String()) {
			entityPtrs = append(entityPtrs, entityPtr)
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
			reString += `.*`
		}
	}
	return regexp.MustCompile(reString)
}
