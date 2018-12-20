package common

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/op/go-logging"
)

const FILE_WRITE_PERMISSIONS = 0600

var plainFormatter = logging.MustStringFormatter(
	`%{time:2006/01/02 15:04:05.000} %{level:8.8s} [%{module}] %{message}`,
)

var colorFormatter = logging.MustStringFormatter(
	`%{color}%{time:2006/01/02 15:04:05.000} %{level:8.8s} [%{module}] %{message}%{color:reset}`,
)

func ConfigureLogging(verbosity int, file *string) {
	verbosity += 3 // 0 verbosity is NOTICE
	if verbosity > 5 {
		verbosity = 5
	}
	level := logging.Level(verbosity)

	var backend *logging.LogBackend
	if file != nil {
		f, err := os.OpenFile(*file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, FILE_WRITE_PERMISSIONS)
		if err != nil {
			message := fmt.Sprintf("log file error: %s", err)
			fmt.Fprintln(color.Error, color.RedString(message))
			os.Exit(1)
		}
		// defer f.Close() ???
		backend = logging.NewLogBackend(f, "", 0)
		logging.SetFormatter(plainFormatter)
	} else {
		backend = logging.NewLogBackend(color.Error, "", 0)
		logging.SetFormatter(colorFormatter)
	}

	leveledBackend := logging.AddModuleLevel(backend)
	leveledBackend.SetLevel(level, "")

	logging.SetBackend(leveledBackend)
}
