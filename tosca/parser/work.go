package parser

import (
	"sync"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/reflection"
)

func (self *Context) Traverse(phase string, traverse reflection.Traverser) {
	work := make(EntityWork)
	var traversed tosca.EntityPtrs

	traverseWrapper := func(entityPtr interface{}) bool {
		if work.Start(phase, entityPtr) {
			return false
		}

		// Don't traverse the same entity more than once
		for _, entityPtr_ := range traversed {
			if entityPtr_ == entityPtr {
				return false
			}
		}
		traversed = append(traversed, entityPtr)

		return traverse(entityPtr)
	}

	// Root
	reflection.Traverse(self.Root.EntityPtr, traverseWrapper)

	// Types
	self.Root.GetContext().Namespace.Range(func(forType interface{}, entityPtr interface{}) bool {
		reflection.Traverse(entityPtr, traverseWrapper)
		return true
	})
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
