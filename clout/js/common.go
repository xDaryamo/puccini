package js

import (
	"github.com/tliron/kutil/logging"
)

var log = logging.GetLogger("puccini.js")
var logEvaluate = logging.NewSubLogger(log, "evaluate")
var logValidate = logging.NewSubLogger(log, "validate")
var logConvert = logging.NewSubLogger(log, "convert")
