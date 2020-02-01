package parser

import (
	"sync"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/reflection"
)

func (self *Context) Traverse(phase string, traverse reflection.Traverser) {
	work := make(EntityWork)

	traverseWrapper := func(entityPtr interface{}) bool {
		if work.Start(phase, entityPtr) {
			return false
		}
		return traverse(entityPtr)
	}

	reflection.Traverse(self.Root.EntityPtr, traverseWrapper)

	for _, forType := range self.Root.GetContext().Namespace {
		for _, entityPtr := range forType {
			reflection.Traverse(entityPtr, traverseWrapper)
		}
	}
}

//
// EntityWork
//

type EntityWork map[interface{}]bool

func (self EntityWork) Start(phase string, entityPtr interface{}) bool {
	if _, ok := self[entityPtr]; ok {
		log.Debugf("{%s} skip: %s", phase, tosca.GetContext(entityPtr).Path)
		return true
	}
	self[entityPtr] = true
	return false
}

//
// ContextualWork
//

type ContextualWork struct {
	sync.Map
	Phase string
}

func NewContextualWork(phase string) *ContextualWork {
	return &ContextualWork{Phase: phase}
}

func (self *ContextualWork) Start(context *tosca.Context) (Promise, bool) {
	key := context.URL.Key()
	promise := NewPromise()
	if existing, loaded := self.LoadOrStore(key, promise); !loaded {
		log.Debugf("{%s} start: %s", self.Phase, key)
		return promise, true
	} else {
		log.Debugf("{%s} wait for: %s", self.Phase, key)
		promise = existing.(Promise)
		promise.Wait()
		return nil, false
	}
}

//
// Promise
//

type Promise chan bool

func NewPromise() Promise {
	return make(Promise)
}

func (self Promise) Release() {
	close(self)
}

func (self Promise) Wait() {
	<-self
}
