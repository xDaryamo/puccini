package js

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/tliron/kutil/ard"
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
	"github.com/tliron/yamlkeys"
)

//
// PucciniAPI
//

type PucciniAPI struct {
	Arguments       map[string]string
	Log             logging.Logger
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
		Arguments:       self.Arguments,
		Log:             self.Log,
		Stdout:          self.Stdout,
		Stderr:          self.Stderr,
		Stdin:           self.Stdin,
		Output:          self.Output,
		Format:          format,
		Strict:          self.Strict,
		AllowTimestamps: self.AllowTimestamps,
		Pretty:          self.Pretty,
		context:         self,
	}
}

func (self *PucciniAPI) Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (self *PucciniAPI) JoinFilePath(elements ...string) string {
	return filepath.Join(elements...)
}

func (self *PucciniAPI) IsType(value ard.Value, type_ string) (bool, error) {
	// Special case whereby an integer stored as a float type has been optimized to an integer type
	if (type_ == "!!float") && ard.IsInteger(value) {
		return true, nil
	}

	if validate, ok := ard.TypeValidators[ard.TypeName(type_)]; ok {
		return validate(value), nil
	} else {
		return false, fmt.Errorf("unsupported type: %s", type_)
	}
}

func (self *PucciniAPI) ValidateFormat(code string, format string) error {
	return formatpkg.Validate(code, format)
}

func (self *PucciniAPI) Timestamp() ard.Value {
	return util.Timestamp(!self.AllowTimestamps)
}

func (self *PucciniAPI) NewXMLDocument() *etree.Document {
	return etree.NewDocument()
}

func (self *PucciniAPI) Decode(code string, format string, all bool) (ard.Value, error) {
	switch format {
	case "yaml", "":
		if all {
			if value, err := yamlkeys.DecodeAll(strings.NewReader(code)); err == nil {
				value_, _ := ard.MapsToStringMaps(value)
				return value_, err
			} else {
				return nil, err
			}
		} else {
			value, _, err := ard.DecodeYAML(code, false)
			value, _ = ard.MapsToStringMaps(value)
			return value, err
		}

	case "json":
		value, _, err := ard.DecodeJSON(code, false)
		value, _ = ard.MapsToStringMaps(value)
		return value, err

	case "cjson":
		value, _, err := ard.DecodeCompatibleJSON(code, false)
		value, _ = ard.MapsToStringMaps(value)
		return value, err

	case "xml":
		value, _, err := ard.DecodeCompatibleXML(code, false)
		value, _ = ard.MapsToStringMaps(value)
		return value, err

	case "cbor":
		value, _, err := ard.DecodeCBOR(code, false)
		value, _ = ard.MapsToStringMaps(value)
		return value, err

	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
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
				message = terminal.StyleError("skipping:   ")
				skip = true
			} else {
				message = terminal.StyleValue("overwriting:")
			}
		} else {
			message = terminal.StyleHeading("writing:    ")
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
		return util.BytesToString(out), nil
	} else if err_, ok := err.(*exec.ExitError); ok {
		return "", fmt.Errorf("%s\n%s", err_.Error(), util.BytesToString(err_.Stderr))
	} else {
		return "", err
	}
}

func (self *PucciniAPI) TemporaryFile(pattern string, directory string) (string, error) {
	if file, err := os.CreateTemp(directory, pattern); err == nil {
		name := file.Name()
		os.Remove(name)
		return name, nil
	} else {
		return "", err
	}
}

func (self *PucciniAPI) TemporaryDirectory(pattern string, directory string) (string, error) {
	return os.MkdirTemp(directory, pattern)
}

func (self *PucciniAPI) Download(sourceUrl string, targetPath string) error {
	if sourceUrl_, err := urlpkg.NewValidURL(sourceUrl, nil, self.context.URLContext); err == nil {
		return urlpkg.DownloadTo(sourceUrl_, targetPath)
	} else {
		return err
	}
}

func (self *PucciniAPI) LoadString(url string) (string, error) {
	if url_, err := urlpkg.NewValidURL(url, nil, self.context.URLContext); err == nil {
		return urlpkg.ReadString(url_)
	} else {
		return "", err
	}
}

// Encode bytes as base64
func (self *PucciniAPI) Btoa(bytes []byte) string {
	return util.ToBase64(bytes)
}

// Decode base64 to bytes
func (self *PucciniAPI) Atob(b64 string) ([]byte, error) {
	// Note: if you need a string in JavaScript: String.fromCharCode.apply(null, puccini.atob(...))
	return util.FromBase64(b64)
}

func (self *PucciniAPI) DeepCopy(value ard.Value) ard.Value {
	return ard.Copy(value)
}

func (self *PucciniAPI) DeepEquals(a ard.Value, b ard.Value) bool {
	return ard.Equals(a, b)
}

func (self *PucciniAPI) Fail(message string) {
	if !self.context.Quiet {
		fmt.Fprintln(self.Stderr, terminal.StyleError(message))
	}
	util.Exit(1)
}

func (self *PucciniAPI) Failf(format string, args ...interface{}) {
	self.Fail(fmt.Sprintf(format, args...))
}

func (self *PucciniAPI) failOnError(err error) {
	if err != nil {
		self.Fail(err.Error())
	}
}
