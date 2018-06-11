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
	return self.readRootEntity(url_, ServiceTemplateGrammars)
}

func (self *Context) readRootEntity(url_ url.URL, readers Grammars) bool {
	toscaContext := tosca.NewContext(self.Problems)
	toscaContext.URL = url_

	self.WG.Add(1)
	serviceTemplate, ok := self.read(nil, &toscaContext, nil, nil, readers)
	self.WG.Wait()

	self.ServiceTemplate = serviceTemplate
	sort.Sort(self.Imports)

	return ok
}

func (self *Context) read(promise Promise, toscaContext *tosca.Context, container *Import, nameTransfomer tosca.NameTransformer, grammars Grammars) (*Import, bool) {
	defer self.WG.Done()
	if promise != nil {
		// For the goroutines waiting for our cached entityPtr
		defer promise.Release()
	}

	log.Infof("{read} %s", toscaContext.URL.Key())

	switch toscaContext.URL.Format() {
	case "csar", "zip":
		url_, err := csar.GetServiceTemplateURL(toscaContext.URL)
		if err != nil {
			toscaContext.ReportError(err)
			return nil, false
		}
		toscaContext.URL = url_
	}

	// Read ARD
	data, err := ard.ReadURL(toscaContext.URL)
	if err != nil {
		toscaContext.ReportError(err)
		return nil, false
	}

	// Grammar
	toscaContext.Data = data
	read, ok := GetGrammar(toscaContext, grammars)
	if !ok {
		return nil, false
	}

	// Read entityPtr
	entityPtr := read(toscaContext)
	if !ok {
		// Even if there are problems, the reader should return an entityPtr
		panic(fmt.Sprintf("reader \"%s\" returned a non-entity: %T", reflection.GetFunctionName(read), entityPtr))
	}

	// Validate required fields
	reflection.Traverse(entityPtr, tosca.ValidateRequiredFields)

	cache.Store(toscaContext.URL.Key(), entityPtr)

	import_ := NewImport(entityPtr, container, nameTransfomer)
	self.AddImport(import_)

	self.readImports(import_)

	return import_, true
}

// From Importer interface
func (self *Context) readImports(container *Import) {
	hasImportSpecs, ok := container.EntityPtr.(tosca.Importer)
	if !ok {
		return
	}

	for _, importSpec := range hasImportSpecs.GetImportSpecs() {
		key := importSpec.URL.Key()

		// Check for import loop
		for c := container; c != nil; c = c.Container {
			url_ := c.GetContext().URL
			if url_.Key() == key {
				container.GetContext().ReportImportLoop(url_)
				return
			}
		}

		promise := NewPromise()
		cached, inCache := cache.LoadOrStore(key, promise)
		if inCache {
			switch c := cached.(type) {
			case Promise:
				// Wait for promise
				log.Debugf("{read} cache promise: %s", key)
				self.WG.Add(1)
				go self.waitForPromise(c, key, container, importSpec.NameTransformer)
			default: // entityPtr
				// Cache hit
				log.Debugf("{read} cache hit: %s", key)
				self.AddImportFor(cached, container, importSpec.NameTransformer)
			}
		} else {
			importToscaContext := container.GetContext().Import(importSpec.URL)

			// Read (concurrently)
			self.WG.Add(1)
			go self.read(promise, importToscaContext, container, importSpec.NameTransformer, UnitGrammars)
		}
	}
}

func (self *Context) waitForPromise(promise Promise, key string, container *Import, nameTransformer tosca.NameTransformer) {
	defer self.WG.Done()
	promise.Wait()

	cached, inCache := cache.Load(key)
	if inCache {
		switch cached.(type) {
		case Promise:
			log.Debugf("{read} cache promise failed: %s", key)
		default: // entityPtr
			// Cache hit
			log.Debugf("{read} cache promise hit: %s", key)
			self.AddImportFor(cached, container, nameTransformer)
		}
	} else {
		log.Debugf("{read} cache promise failed (empty): %s", key)
	}
}

var cache sync.Map // entityPtr or Promise
