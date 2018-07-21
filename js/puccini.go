package js

import (
	"fmt"
	"os"

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

	context *Context
}

func (self *Context) NewPucciniApi() *PucciniApi {
	return &PucciniApi{
		Log:     self.Log,
		Stdout:  self.Stdout,
		Stderr:  self.Stderr,
		Stdin:   self.Stdin,
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

func (self *PucciniApi) Write(data interface{}) {
	if !self.context.Quiet || (self.context.Output != "") {
		err := format.WriteOrPrint(data, self.context.ArdFormat, true, self.context.Output)
		self.context.ValidateError(err)
	}
}
