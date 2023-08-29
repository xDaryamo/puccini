package commands

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

const toolName = "puccini-csar"

var log = commonlog.GetLogger(toolName)

var archiveFormat string

func Bases(urlContext *exturl.Context) []exturl.URL {
	workingDir, err := urlContext.NewWorkingDirFileURL()
	util.FailOnError(err)
	return []exturl.URL{workingDir}
}
