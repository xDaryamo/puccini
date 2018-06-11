package parser

import (
	"sync"

	"github.com/op/go-logging"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("parser")

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
