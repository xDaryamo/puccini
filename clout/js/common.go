package js

import (
	"github.com/tliron/kutil/logging"
)

var log = logging.GetLogger("puccini.js")
var logEvaluate = logging.NewScopeLogger(log, "evaluate")
var logValidate = logging.NewScopeLogger(log, "validate")
var logConvert = logging.NewScopeLogger(log, "convert")
