package js

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/common"
	formatpkg "github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
	urlpkg "github.com/tliron/puccini/url"
)

//
// PucciniAPI
//

type PucciniAPI struct {
	Log             *Log
	Stdout          io.Writer
	Stderr          io.Writer
	Stdin           io.Writer
	Output          string
	Format          string
	Strict          bool
	AllowTimestamps bool
	Pretty          bool

	context *Context
}

func (self *Context) NewPucciniAPI() *PucciniAPI {
	format := self.Format
	if format == "" {
		format = "yaml"
	}
	return &PucciniAPI{
		Log:             self.Log,
		Stdout:          self.Stdout,
		Stdin:           self.Stdin,
		Output:          self.Output,
		Format:          format,
		Strict:          self.Strict,
		AllowTimestamps: self.AllowTimestamps,
		Pretty:          self.Pretty,
		context:         self,
	}
}

func (entry *PucciniAPI) Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (entry *PucciniAPI) JoinFilePath(elements ...string) string {
	return filepath.Join(elements...)
}

func (entry *PucciniAPI) IsType(value ard.Value, type_ string) (bool, error) {
	// Special case whereby an integer stored as a float type has been optimized to an integer type
	if (type_ == "!!float") && ard.IsInteger(value) {
		return true, nil
	}

	if validate, ok := ard.TypeValidators[type_]; ok {
		return validate(value), nil
	} else {
		return false, fmt.Errorf("unsupported type: %s", type_)
	}
}

func (entry *PucciniAPI) ValidateFormat(code string, format string) error {
	return formatpkg.Validate(code, format)
}

func (self *PucciniAPI) Timestamp() ard.Value {
	return common.Timestamp(!self.AllowTimestamps)
}

func (self *PucciniAPI) NewXMLDocument() *etree.Document {
	return etree.NewDocument()
}

func (self *PucciniAPI) Write(data interface{}, path string, dontOverwrite bool) {
	output := self.context.Output
	if path != "" {
		// Our path is relative to output path
		// (output path is here considered to be a directory)
		output = filepath.Join(output, path)
		var err error
		output, err = filepath.Abs(output)
		self.failOnError(err)
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
				message = terminal.ColorError("skipping:   ")
				skip = true
			} else {
				message = terminal.ColorValue("overwriting:")
			}
		} else {
			message = terminal.ColorHeading("writing:    ")
		}
		if !self.context.Quiet {
			fmt.Fprintf(self.Stdout, "%s %s\n", message, output)
		}
		if skip {
			return
		}
	}

	self.failOnError(formatpkg.WriteOrPrint(data, self.Format, self.Stdout, self.Strict, self.Pretty, output))
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
	if sourceUrl_, err := urlpkg.NewValidURL(sourceUrl, nil); err == nil {
		return urlpkg.DownloadTo(sourceUrl_, targetPath)
	} else {
		return err
	}
}

func (self *PucciniAPI) LoadString(url string) (string, error) {
	if url_, err := urlpkg.NewValidURL(url, nil); err == nil {
		return urlpkg.ReadToString(url_)
	} else {
		return "", err
	}
}

func (self *PucciniAPI) Btoa(from []byte) (string, error) {
	var builder strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &builder)
	if _, err := encoder.Write(from); err == nil {
		if err := encoder.Close(); err == nil {
			return builder.String(), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func (self *PucciniAPI) Atob(from string) ([]byte, error) {
	// Note: if you need a string in JavaScript: String.fromCharCode.apply(null, puccini.atob(...))
	reader := strings.NewReader(from)
	decoder := base64.NewDecoder(base64.StdEncoding, reader)
	return ioutil.ReadAll(decoder)
}

func (self *PucciniAPI) Fail(message string) {
	if !self.context.Quiet {
		fmt.Fprintln(self.Stderr, terminal.ColorError(message))
	}
	atexit.Exit(1)
}

func (self *PucciniAPI) Failf(format string, args ...interface{}) {
	self.Fail(fmt.Sprintf(format, args...))
}

func (self *PucciniAPI) failOnError(err error) {
	if err != nil {
		self.Fail(err.Error())
	}
}
