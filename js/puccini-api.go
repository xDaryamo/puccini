package js

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
	"github.com/fatih/color"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
)

//
// PucciniApi
//

type PucciniApi struct {
	Log    *format.Log
	Stdout *os.File
	Stderr *os.File
	Stdin  *os.File
	Format string

	context *Context
}

func (self *Context) NewPucciniApi() *PucciniApi {
	format := self.ArdFormat
	if format == "" {
		format = "yaml"
	}
	return &PucciniApi{
		Log:     self.Log,
		Stdout:  self.Stdout,
		Stderr:  self.Stderr,
		Stdin:   self.Stdin,
		Format:  format,
		context: self,
	}
}

func (self *PucciniApi) Sprintf(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}

func (self *PucciniApi) Timestamp() (string, error) {
	return common.Timestamp()
}

func (self *PucciniApi) NewXmlDocument() *etree.Document {
	return etree.NewDocument()
}

func (self *PucciniApi) Write(data interface{}, path string, dontOverwrite bool) {
	output := self.context.Output
	if path != "" {
		output = filepath.Join(output, path)
		var err error
		output, err = filepath.Abs(output)
		self.context.ValidateError(err)
	}

	if self.context.Quiet && (output == "") {
		return
	}

	if output != "" {
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

	var err error
	if s, ok := data.(string); ok {
		// String is a special case: we just write the string contents
		var f *os.File

		if output != "" {
			var err error
			f, err = format.OpenFileForWrite(output)
			self.context.ValidateError(err)
			defer f.Close()
		} else {
			f = self.Stdout
		}

		_, err = f.WriteString(s)
	} else {
		err = format.WriteOrPrint(data, self.context.ArdFormat, true, output)
	}
	self.context.ValidateError(err)
}
