package parser

import (
	"github.com/tliron/kutil/logging"
)

var log = logging.GetLogger("puccini.parser")
var logRead = logging.NewScopeLogger(log, "read")
var logNamespaces = logging.NewScopeLogger(log, "namespaces")
var logLookup = logging.NewScopeLogger(log, "lookup")
var logHierarchies = logging.NewScopeLogger(log, "hierarchies")
var logInheritance = logging.NewScopeLogger(log, "inheritance")
var logTasks = logging.NewScopeLogger(log, "tasks")
var logGather = logging.NewScopeLogger(log, "gather")
