package parser

import (
	"sync"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/puccini/tosca"
)

func (self *ServiceContext) TraverseEntities(log logging.Logger, work tosca.EntityWork, traverse reflection.EntityTraverser) {
	if work == nil {
		work = make(tosca.EntityWork)
	}

	// Root
	work.TraverseEntities(self.Root.EntityPtr, traverse)

	// Types
	self.Root.GetContext().Namespace.Range(func(entityPtr tosca.EntityPtr) bool {
		work.TraverseEntities(entityPtr, traverse)
		return true
	})
}

//
// Promise
//

type Promise chan struct{}

func NewPromise() Promise {
	return make(Promise)
}

func (self Promise) Release() {
	close(self)
}

func (self Promise) Wait() {
	<-self
}

// TODO: unused
//
// CoordinatedWork
//

type CoordinatedWork struct {
	sync.Map
	Log logging.Logger
}

func NewCoordinatedWork(log logging.Logger) *CoordinatedWork {
	return &CoordinatedWork{
		Log: log,
	}
}

func (self *CoordinatedWork) Start(key string) (Promise, bool) {
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
