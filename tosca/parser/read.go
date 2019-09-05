package parser

import (
	"fmt"
	"sort"
	"sync"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/csar"
	"github.com/tliron/puccini/tosca/reflection"
	"github.com/tliron/puccini/url"
)

func (self *Context) ReadServiceTemplate(url_ url.URL) bool {
	toscaContext := tosca.NewContext(&self.Problems, self.Quirks)
	toscaContext.URL = url_

	self.WG.Add(1)
	serviceTemplate, ok := self.read(nil, toscaContext, nil, nil, "ServiceTemplate")
	self.WG.Wait()

	self.ServiceTemplate = serviceTemplate
	sort.Sort(self.Units)

	return ok
}

func (self *Context) read(promise Promise, toscaContext *tosca.Context, container *Unit, nameTransfomer tosca.NameTransformer, readerName string) (*Unit, bool) {
	defer self.WG.Done()
	if promise != nil {
		// For the goroutines waiting for our cached entityPtr
		defer promise.Release()
	}

	log.Infof("{read} %s: %s", readerName, toscaContext.URL.Key())

	switch toscaContext.URL.Format() {
	case "csar", "zip":
		var err error
		if toscaContext.URL, err = csar.GetServiceTemplateURL(toscaContext.URL); err != nil {
			toscaContext.ReportError(err)
			return nil, false
		}
	}

	// Read ARD
	var err error
	if toscaContext.Data, toscaContext.Locator, err = ard.ReadURL(toscaContext.URL, true); err != nil {
		toscaContext.ReportError(err)
		return nil, false
	}

	// Detect grammar
	if !DetectGrammar(toscaContext) {
		return nil, false
	}

	// Read entityPtr
	read, ok := toscaContext.Grammar[readerName]
	if !ok {
		panic(fmt.Sprintf("grammar does not support reader \"%s\"", readerName))
	}
	entityPtr := read(toscaContext)
	if entityPtr == nil {
		// Even if there are problems, the reader should return an entityPtr
		panic(fmt.Sprintf("reader \"%s\" returned a non-entity: %T", reflection.GetFunctionName(read), entityPtr))
	}

	// Validate required fields
	reflection.Traverse(entityPtr, tosca.ValidateRequiredFields)

	cache.Store(toscaContext.URL.Key(), entityPtr)

	unit := NewUnit(entityPtr, container, nameTransfomer)
	self.AddUnit(unit)

	self.goReadImports(unit)

	return unit, true
}

// From Importer interface
func (self *Context) goReadImports(container *Unit) {
	var importSpecs []*tosca.ImportSpec
	if importer, ok := container.EntityPtr.(tosca.Importer); ok {
		importSpecs = importer.GetImportSpecs()
	}

	// Implicit import
	if implicitImportSpec, ok := GetProfileImportSpec(container.GetContext()); ok {
		importSpecs = append(importSpecs, implicitImportSpec)
	}

	for _, importSpec := range importSpecs {
		key := importSpec.URL.Key()

		// Skip if causes import loop
		skip := false
		for c := container; c != nil; c = c.Container {
			url_ := c.GetContext().URL
			if url_.Key() == key {
				if !importSpec.Implicit {
					// Explicit import loops are considered errors
					container.GetContext().ReportImportLoop(url_)
				}
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		promise := NewPromise()
		if cached, inCache := cache.LoadOrStore(key, promise); inCache {
			switch c := cached.(type) {
			case Promise:
				// Wait for promise
				log.Debugf("{read} cache promise: %s", key)
				self.WG.Add(1)
				go self.waitForPromise(c, key, container, importSpec.NameTransformer)
			default: // entityPtr
				// Cache hit
				log.Debugf("{read} cache hit: %s", key)
				self.AddUnitFor(cached, container, importSpec.NameTransformer)
			}
		} else {
			importToscaContext := container.GetContext().NewImportContext(importSpec.URL)

			// Read (concurrently)
			self.WG.Add(1)
			go self.read(promise, importToscaContext, container, importSpec.NameTransformer, "Unit")
		}
	}
}

func (self *Context) waitForPromise(promise Promise, key string, container *Unit, nameTransformer tosca.NameTransformer) {
	defer self.WG.Done()
	promise.Wait()

	if cached, inCache := cache.Load(key); inCache {
		switch cached.(type) {
		case Promise:
			log.Debugf("{read} cache promise failed: %s", key)
		default: // entityPtr
			// Cache hit
			log.Debugf("{read} cache promise hit: %s", key)
			self.AddUnitFor(cached, container, nameTransformer)
		}
	} else {
		log.Debugf("{read} cache promise failed (empty): %s", key)
	}
}

var cache sync.Map // entityPtr or Promise
