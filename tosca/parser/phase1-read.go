package parser

import (
	"fmt"
	"sort"
	"sync"

	"github.com/tliron/kutil/reflection"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/csar"
	"github.com/tliron/puccini/tosca/grammars"
	"github.com/tliron/yamlkeys"
)

func (self *Context) ReadRoot(url urlpkg.URL) bool {
	toscaContext := tosca.NewContext(self.Quirks)

	toscaContext.URL = url

	var ok bool

	self.WaitGroup.Add(1)
	self.Root, ok = self.read(nil, toscaContext, nil, nil, "$Root")
	self.WaitGroup.Wait()

	sort.Sort(self.Units)

	return ok
}

var readCache sync.Map // entityPtr or Promise

func (self *Context) read(promise Promise, toscaContext *tosca.Context, container *Unit, nameTransformer tosca.NameTransformer, readerName string) (*Unit, bool) {
	defer self.WaitGroup.Done()
	if promise != nil {
		// For the goroutines waiting for our cached entityPtr
		defer promise.Release()
	}

	logRead.Infof("%s: %s", readerName, toscaContext.URL.Key())

	switch toscaContext.URL.Format() {
	case "csar", "zip":
		var err error
		if toscaContext.URL, err = csar.GetRootURL(toscaContext.URL); err != nil {
			toscaContext.ReportError(err)
			unit := NewUnitNoEntity(toscaContext, container, nameTransformer)
			self.AddUnit(unit)
			return unit, false
		}
	}

	// Read ARD
	var err error
	if toscaContext.Data, toscaContext.Locator, err = urlpkg.ReadARD(toscaContext.URL, true); err != nil {
		if decodeError, ok := err.(*yamlkeys.DecodeError); ok {
			err = NewYAMLDecodeError(decodeError)
		}
		toscaContext.ReportError(err)
		unit := NewUnitNoEntity(toscaContext, container, nameTransformer)
		self.AddUnit(unit)
		return unit, false
	}

	// Detect grammar
	if !grammars.Detect(toscaContext) {
		unit := NewUnitNoEntity(toscaContext, container, nameTransformer)
		self.AddUnit(unit)
		return unit, false
	}

	// Read entityPtr
	read, ok := toscaContext.Grammar.Readers[readerName]
	if !ok {
		panic(fmt.Sprintf("grammar does not support reader %q", readerName))
	}
	entityPtr := read(toscaContext)
	if entityPtr == nil {
		// Even if there are problems, the reader should return an entityPtr
		panic(fmt.Sprintf("reader %q returned a non-entity: %T", reflection.GetFunctionName(read), entityPtr))
	}

	// Validate required fields
	reflection.Traverse(entityPtr, tosca.ValidateRequiredFields)

	readCache.Store(toscaContext.URL.Key(), entityPtr)

	return self.AddImportUnit(entityPtr, container, nameTransformer), true
}

// From tosca.Importer interface
func (self *Context) goReadImports(container *Unit) {
	var importSpecs []*tosca.ImportSpec
	if importer, ok := container.EntityPtr.(tosca.Importer); ok {
		importSpecs = importer.GetImportSpecs()
	}

	// Implicit import
	if !container.GetContext().HasQuirk(tosca.QuirkImportsImplicitDisable) {
		if implicitImportSpec, ok := grammars.GetImplicitImportSpec(container.GetContext()); ok {
			importSpecs = append(importSpecs, implicitImportSpec)
		}
	}

	for _, importSpec := range importSpecs {
		key := importSpec.URL.Key()

		// Skip if causes import loop
		skip := false
		for container_ := container; container_ != nil; container_ = container_.Container {
			url := container_.GetContext().URL
			if url.Key() == key {
				if !importSpec.Implicit {
					// Import loops are considered errors
					container.GetContext().ReportImportLoop(url)
				}
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		promise := NewPromise()
		if cached, inCache := readCache.LoadOrStore(key, promise); inCache {
			switch cached_ := cached.(type) {
			case Promise:
				// Wait for promise
				logRead.Debugf("wait for promise: %s", key)
				self.WaitGroup.Add(1)
				go self.waitForPromise(cached_, key, container, importSpec.NameTransformer)

			default: // entityPtr
				// Cache hit
				logRead.Debugf("cache hit: %s", key)
				self.AddImportUnit(cached, container, importSpec.NameTransformer)
			}
		} else {
			importToscaContext := container.GetContext().NewImportContext(importSpec.URL)

			// Read (concurrently)
			self.WaitGroup.Add(1)
			go self.read(promise, importToscaContext, container, importSpec.NameTransformer, "$Unit")
		}
	}
}

func (self *Context) waitForPromise(promise Promise, key string, container *Unit, nameTransformer tosca.NameTransformer) {
	defer self.WaitGroup.Done()
	promise.Wait()

	if cached, inCache := readCache.Load(key); inCache {
		switch cached.(type) {
		case Promise:
			logRead.Debugf("promise broken: %s", key)

		default: // entityPtr
			// Cache hit
			logRead.Debugf("promise kept: %s", key)
			self.AddImportUnit(cached, container, nameTransformer)
		}
	} else {
		logRead.Debugf("promise broken (empty): %s", key)
	}
}
