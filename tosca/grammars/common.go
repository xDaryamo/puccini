package grammars

import (
	"github.com/tliron/puccini/tosca/parsing"
)

// Map of keyword -> version -> grammar
var Grammars = make(map[string]map[string]*parsing.Grammar)

// Map of keyword -> version -> internal URL path
var ImplicitProfilePaths = make(map[string]map[string]string)
