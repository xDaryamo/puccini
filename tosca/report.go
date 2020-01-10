package tosca

import (
	"fmt"
	"reflect"

	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/url"
)

//
// Context
//

func (self *Context) Report(message string) bool {
	if self.URL != nil {
		return self.Problems.ReportInSection(message, self.URL.String())
	} else {
		return self.Problems.Report(message)
	}
}

func (self *Context) Reportf(f string, arg ...interface{}) bool {
	return self.Report(fmt.Sprintf(f, arg...))
}

func (self *Context) ReportPath(message string) bool {
	path := self.Path.String()

	location := self.Location()
	if location != "" {
		if path != "" {
			path += " "
		}
		path += "@" + location
	}

	if path != "" {
		message = fmt.Sprintf("%s: %s", terminal.ColorPath(path), message)
	}

	return self.Report(message)
}

func (self *Context) ReportPathf(f string, arg ...interface{}) bool {
	return self.ReportPath(fmt.Sprintf(f, arg...))
}

func (self *Context) ReportError(err error) bool {
	return self.ReportPathf("%s", err)
}

//
// Values
//

func (self *Context) FormatBadData() string {
	return terminal.ColorError(fmt.Sprintf("%+v", self.Data))
}

func (self *Context) ReportValueWrongType(requiredTypeNames ...string) bool {
	return self.ReportPathf("\"%s\" instead of %s", terminal.ColorTypeName(PrimitiveTypeName(self.Data)), terminal.ColoredOptions(requiredTypeNames, terminal.ColorTypeName))
}

func (self *Context) ReportValueWrongFormat(format string) bool {
	return self.ReportPathf("wrong format, must be \"%s\": %s", format, self.FormatBadData())
}

func (self *Context) ReportValueWrongLength(typeName string, length int) bool {
	return self.ReportPathf("\"%s\" does not have %d elements", terminal.ColorTypeName(typeName), length)
}

func (self *Context) ReportValueMalformed(typeName string, reason string) bool {
	if reason == "" {
		return self.ReportPathf("malformed \"%s\": %s", terminal.ColorTypeName(typeName), self.FormatBadData())
	} else {
		return self.ReportPathf("malformed \"%s\", %s: %s", terminal.ColorTypeName(typeName), reason, self.FormatBadData())
	}
}

//
// Read
//

func (self *Context) ReportImportIncompatible(url_ url.URL) bool {
	return self.Reportf("incompatible import \"%s\"", terminal.ColorValue(url_.String()))
}

func (self *Context) ReportImportLoop(url_ url.URL) bool {
	return self.Reportf("endless loop caused by importing \"%s\"", terminal.ColorValue(url_.String()))
}

func (self *Context) ReportRepositoryInaccessible(repositoryName string) bool {
	return self.ReportPathf("inaccessible repository \"%s\"", terminal.ColorValue(repositoryName))
}

func (self *Context) ReportFieldMissing() bool {
	return self.ReportPath("missing required field")
}

func (self *Context) ReportFieldUnsupported() bool {
	return self.ReportPath("unsupported field")
}

func (self *Context) ReportFieldUnsupportedValue() bool {
	return self.ReportPathf("unsupported value for field: %s", self.FormatBadData())
}

func (self *Context) ReportFieldMalformedSequencedList() bool {
	return self.ReportPathf("field must be a \"%s\" of single-key \"%s\" elements", terminal.ColorTypeName("sequenced list"), terminal.ColorTypeName("map"))
}

func (self *Context) ReportPrimitiveType() bool {
	return self.ReportPath("primitive type cannot have properties")
}

func (self *Context) ReportMapKeyReused(key string) bool {
	return self.ReportPathf("reused map key: %s", terminal.ColorValue(key))
}

//
// Namespaces
//

func (self *Context) ReportNameAmbiguous(type_ reflect.Type, name string, entityPtrs ...interface{}) bool {
	url := make([]string, len(entityPtrs))
	for i, entityPtr := range entityPtrs {
		url[i] = GetContext(entityPtr).URL.String()
	}
	return self.Reportf("ambiguous %s name \"%s\", can be in %s", GetEntityTypeName(type_), terminal.ColorName(name), terminal.ColoredOptions(url, terminal.ColorValue))
}

func (self *Context) ReportFieldReferenceNotFound(types ...reflect.Type) bool {
	var entityTypeNames []string
	for _, type_ := range types {
		entityTypeNames = append(entityTypeNames, GetEntityTypeName(type_))
	}
	return self.ReportPathf("reference to unknown %s: %s", terminal.Options(entityTypeNames), self.FormatBadData())
}

//
// Inheritance
//

func (self *Context) ReportInheritanceLoop(parent interface{}) bool {
	return self.ReportPathf("inheritance loop by deriving from \"%s\"", terminal.ColorTypeName(GetContext(parent).Name))
}

func (self *Context) ReportTypeIncomplete(parent interface{}) bool {
	return self.ReportPathf("deriving from incomplete type \"%s\"", terminal.ColorTypeName(GetContext(parent).Name))
}

//
// Render
//

func (self *Context) ReportUndeclared(kind string) bool {
	return self.ReportPathf("undeclared %s", kind)
}

func (self *Context) ReportUnknown(kind string) bool {
	return self.ReportPathf("unknown %s: %s", kind, self.FormatBadData())
}

func (self *Context) ReportReferenceNotFound(kind string, entityPtr interface{}) bool {
	typeName := GetEntityTypeName(reflect.TypeOf(entityPtr).Elem())
	name := GetContext(entityPtr).Name
	return self.ReportPathf("unknown %s reference in %s \"%s\": %s", kind, typeName, terminal.ColorName(name), self.FormatBadData())
}

func (self *Context) ReportReferenceAmbiguous(kind string, entityPtr interface{}) bool {
	typeName := GetEntityTypeName(reflect.TypeOf(entityPtr).Elem())
	name := GetContext(entityPtr).Name
	return self.ReportPathf("ambiguous %s in %s \"%s\": %s", kind, typeName, terminal.ColorName(name), self.FormatBadData())
}

func (self *Context) ReportPropertyRequired(kind string) bool {
	return self.ReportPathf("unassigned required %s", kind)
}

func (self *Context) ReportReservedMetadata() bool {
	return self.ReportPath("reserved for use by Puccini")
}

func (self *Context) ReportUnknownDataType(dataTypeName string) bool {
	return self.ReportPathf("unknown data type \"%s\"", terminal.ColorError(dataTypeName))
}

func (self *Context) ReportMissingEntrySchema(kind string) bool {
	return self.ReportPathf("missing entry schema for %s definition", kind)
}

func (self *Context) ReportUnsupportedType() bool {
	return self.ReportPathf("unsupported puccini.type \"%s\"", terminal.ColorError(self.Name))
}

func (self *Context) ReportIncompatibleType(typeName string, parentTypeName string) bool {
	return self.ReportPathf("type \"%s\" is incompatible with parent type \"%s\"", terminal.ColorTypeName(typeName), terminal.ColorTypeName(parentTypeName))
}

func (self *Context) ReportIncompatible(name string, typeName string, kind string) bool {
	return self.ReportPathf("\"%s\" cannot be %s of %s", terminal.ColorName(name), kind, typeName)
}

func (self *Context) ReportIncompatibleExtension(extension string, requiredExtensions []string) bool {
	return self.ReportPathf("extension \"%s\" is not %s", terminal.ColorValue(extension), terminal.ColoredOptions(requiredExtensions, terminal.ColorValue))
}

func (self *Context) ReportNotInRange(name string, value uint64, lower uint64, upper uint64) bool {
	return self.ReportPathf("%s is %d, must be >= %d and <= %d", name, value, lower, upper)
}
