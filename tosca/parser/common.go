package parser

import (
	"sync"
	"time"

	"github.com/op/go-logging"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("parser")

func GetVersion(context *tosca.Context) (*tosca.Context, *string) {
	var versionContext *tosca.Context
	var ok bool

	for keyword := range Grammars {
		if versionContext, ok = context.GetFieldChild(keyword); ok {
			if keyword == "heat_template_version" {
				// Hack to allow HOT to use YAML !!timestamp values

				if versionContext.Is("string") {
					return versionContext, versionContext.ReadString()
				}

				switch versionContext.Data.(type) {
				case time.Time:
					versionContext.Data = versionContext.Data.(time.Time).Format("2006-01-02")
					return versionContext, versionContext.ReadString()
				}

				versionContext.ReportValueWrongType("string", "timestamp")
			} else {
				if versionContext.ValidateType("string") {
					return versionContext, versionContext.ReadString()
				}
			}
		}
	}

	return nil, nil
}

//
// ContextsWork
//

type ContextsWork struct {
	sync.Map
	Phase string
}

func (self *ContextsWork) Start(context *tosca.Context) (Promise, bool) {
	key := context.URL.Key()
	promise := NewPromise()
	existing, loaded := self.LoadOrStore(key, promise)
	if loaded {
		log.Debugf("{%s} wait: %s", self.Phase, key)
		promise = existing.(Promise)
		promise.Wait()
		return nil, false
	}
	log.Debugf("{%s} start: %s", self.Phase, key)
	return promise, true
}

//
// EntitiesDone
//

type EntitiesDone map[interface{}]bool

func (self EntitiesDone) IsDone(phase string, entityPtr interface{}) bool {
	if _, ok := self[entityPtr]; ok {
		log.Debugf("{%s} skip: %s", phase, tosca.GetContext(entityPtr).Path)
		return true
	}
	self[entityPtr] = true
	return false
}
