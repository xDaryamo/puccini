package parser

import (
	"sync"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) Traverse(log logging.Logger, traverse reflection.Traverser) {
	work := make(EntityWork)
	var traversed tosca.EntityPtrs

	traverseWrapper := func(entityPtr tosca.EntityPtr) bool {
		if work.Start(log, entityPtr) {
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
	self.Root.GetContext().Namespace.Range(func(forType tosca.EntityPtr, entityPtr tosca.EntityPtr) bool {
		reflection.Traverse(entityPtr, traverseWrapper)
		return true
	})
}

//
// EntityWork
//

type EntityWork map[tosca.EntityPtr]bool

func (self EntityWork) Start(log logging.Logger, entityPtr tosca.EntityPtr) bool {
	if _, ok := self[entityPtr]; ok {
		log.Debugf("skip: %s", tosca.GetContext(entityPtr).Path)
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
	Log logging.Logger
}

func NewContextualWork(log logging.Logger) *ContextualWork {
	return &ContextualWork{
		Log: log,
	}
}

func (self *ContextualWork) Start(context *tosca.Context) (Promise, bool) {
	key := context.URL.Key()
	promise := NewPromise()
	if existing, loaded := self.LoadOrStore(key, promise); !loaded {
		self.Log.Debugf("start: %s", key)
		return promise, true
	} else {
		self.Log.Debugf("wait for: %s", key)
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
