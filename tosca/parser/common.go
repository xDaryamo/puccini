package parser

import (
	"github.com/tliron/kutil/logging"
)

var log = logging.GetLogger("puccini.parser")
var logRead = logging.NewSubLogger(log, "read")
var logNamespaces = logging.NewSubLogger(log, "namespaces")
var logLookup = logging.NewSubLogger(log, "lookup")
var logHierarchies = logging.NewSubLogger(log, "hierarchies")
var logInheritance = logging.NewSubLogger(log, "inheritance")
var logTasks = logging.NewSubLogger(log, "tasks")
var logGather = logging.NewSubLogger(log, "gather")
