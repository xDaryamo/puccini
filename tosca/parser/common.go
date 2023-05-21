package parser

import (
	"github.com/tliron/commonlog"
)

var log = commonlog.GetLogger("puccini.parser")
var logRead = commonlog.NewScopeLogger(log, "read")
var logNamespaces = commonlog.NewScopeLogger(log, "namespaces")
var logLookup = commonlog.NewScopeLogger(log, "lookup")
var logHierarchies = commonlog.NewScopeLogger(log, "hierarchies")
var logInheritance = commonlog.NewScopeLogger(log, "inheritance")
var logTasks = commonlog.NewScopeLogger(log, "tasks")
var logGather = commonlog.NewScopeLogger(log, "gather")
