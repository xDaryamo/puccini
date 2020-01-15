package js

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/beevik/etree"
	"github.com/fatih/color"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/url"
)

//
// PucciniAPI
//

type PucciniAPI struct {
	Log    *Log
	Stdout *os.File
	Stderr *os.File
	Stdin  *os.File
	Output string
	Format string
	Pretty bool

	context *Context
}

func (self *Context) NewPucciniAPI() *PucciniAPI {
	format := self.Format
	if format == "" {
		format = "yaml"
	}
	return &PucciniAPI{
		Log:     self.Log,
		Stdout:  self.Stdout,
		Stderr:  self.Stderr,
		Stdin:   self.Stdin,
		Output:  self.Output,
		Format:  format,
		Pretty:  self.Pretty,
		context: self,
	}
}

func (entry *PucciniAPI) Sprintf(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}

func (entry *PucciniAPI) JoinFilePath(elements ...string) string {
	return filepath.Join(elements...)
}

func (entry *PucciniAPI) ValidateFormat(code string, format_ string) error {
	return format.Validate(code, format_)
}

func (self *PucciniAPI) Timestamp() string {
	return common.Timestamp()
}

func (self *PucciniAPI) NewXMLDocument() *etree.Document {
	return etree.NewDocument()
}

func (self *PucciniAPI) Write(data interface{}, path string, dontOverwrite bool) {
	output := self.context.Output
	if path != "" {
		output = filepath.Join(output, path)
		var err error
		output, err = filepath.Abs(output)
		self.context.FailOnError(err)
	}

	if output == "" {
		if self.context.Quiet {
			return
		}
	} else {
		_, err := os.Stat(output)
		var message string
		var skip bool
		if (err == nil) || os.IsExist(err) {
			if dontOverwrite {
				message = color.RedString("skippping:  ")
				skip = true
			} else {
				message = color.YellowString("overwriting:")
			}
		} else {
			message = color.GreenString("writing:    ")
		}
		if !self.context.Quiet {
			fmt.Fprintln(self.Stdout, fmt.Sprintf("%s %s", message, output))
		}
		if skip {
			return
		}
	}

	self.context.FailOnError(format.WriteOrPrint(data, self.Format, self.Pretty, output))
}

func (self *PucciniAPI) Exec(name string, arguments ...string) (string, error) {
	cmd := exec.Command(name, arguments...)
	if out, err := cmd.Output(); err == nil {
		return common.BytesToString(out), nil
	} else {
		return "", err
	}
}

func (self *PucciniAPI) Download(sourceUrl string, targetPath string) error {
	if sourceUrl_, err := url.NewValidURL(sourceUrl, nil); err == nil {
		return url.DownloadTo(sourceUrl_, targetPath)
	} else {
		return err
	}
}
