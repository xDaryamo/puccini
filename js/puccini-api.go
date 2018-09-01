package js

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
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

func (self *PucciniApi) Write(data interface{}, path string) {
	output := self.context.Output
	if path != "" {
		output = filepath.Join(output, path)
		err := os.MkdirAll(filepath.Dir(output), os.ModePerm)
		self.context.ValidateError(err)
	}

	if !self.context.Quiet || (output != "") {
		if output != "" {
			fmt.Fprintf(self.Stdout, "writing %s\n", output)
		}
		err := format.WriteOrPrint(data, self.context.ArdFormat, true, output)
		self.context.ValidateError(err)
	}
}
