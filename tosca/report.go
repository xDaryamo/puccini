package tosca

import (
	"fmt"
	"reflect"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/url"
)

//
// Context
//

func (self *Context) Report(message string) {
	if self.URL != nil {
		self.Problems.ReportInSection(message, self.URL.String())
	} else {
		self.Problems.Report(message)
	}
}

func (self *Context) Reportf(f string, arg ...interface{}) {
	self.Report(fmt.Sprintf(f, arg...))
}

func (self *Context) ReportPath(message string) {
	path := self.Path.String()
	if path != "" {
		message = fmt.Sprintf("%s: %s", format.ColorPath(path), message)
	}

	location := self.Location()
	if location != "" {
		if message != "" {
			message += " "
		}
		message += format.ColorValue("@" + location)
	}

	self.Report(message)
}

func (self *Context) ReportPathf(f string, arg ...interface{}) {
	self.ReportPath(fmt.Sprintf(f, arg...))
}

func (self *Context) ReportError(err error) {
	self.ReportPathf("%s", err)
}

//
// Values
//

func (self *Context) FormatBadData() string {
	return format.ColorError(fmt.Sprintf("%+v", self.Data))
}

func (self *Context) ReportValueWrongType(requiredTypeNames ...string) {
	self.ReportPathf("\"%s\" instead of %s", format.ColorTypeName(PrimitiveTypeName(self.Data)), format.ColoredOptions(requiredTypeNames, format.ColorTypeName))
}

func (self *Context) ReportValueWrongFormat(f string) {
	self.ReportPathf("wrong format, must be \"%s\": %s", f, self.FormatBadData())
}

func (self *Context) ReportValueWrongLength(typeName string, length int) {
	self.ReportPathf("\"%s\" does not have %d elements", format.ColorTypeName(typeName), length)
}

func (self *Context) ReportValueMalformed(typeName string, reason string) {
	if reason == "" {
		self.ReportPathf("malformed \"%s\": %s", format.ColorTypeName(typeName), self.FormatBadData())
	} else {
		self.ReportPathf("malformed \"%s\", %s: %s", format.ColorTypeName(typeName), reason, self.FormatBadData())
	}
}

//
// Read
//

func (self *Context) ReportImportLoop(url_ url.URL) {
	self.Reportf("endless loop caused by importing \"%s\"", format.ColorValue(url_.String()))
}

func (self *Context) ReportRepositoryInaccessible(repositoryName string) {
	self.ReportPathf("inaccessible repository \"%s\"", format.ColorValue(repositoryName))
}

func (self *Context) ReportFieldMissing() {
	self.ReportPath("missing required field")
}

func (self *Context) ReportFieldUnsupported() {
	self.ReportPath("unsupported field")
}

func (self *Context) ReportFieldUnsupportedValue() {
	self.ReportPathf("unsupported value for field: %s", self.FormatBadData())
}

func (self *Context) ReportFieldMalformedSequencedList() {
	self.ReportPathf("field must be a \"%s\" of single-key \"%s\" elements", format.ColorTypeName("sequenced list"), format.ColorTypeName("map"))
}

func (self *Context) ReportPrimitiveType() {
	self.ReportPath("primitive type cannot have properties")
}

func (self *Context) ReportMapKeyReused(key string) {
	self.ReportPathf("reused map key: %s", format.ColorValue(key))
}

//
// Namespaces
//

func (self *Context) ReportNameAmbiguous(type_ reflect.Type, name string, entityPtrs ...interface{}) {
	url := make([]string, len(entityPtrs))
	for i, entityPtr := range entityPtrs {
		url[i] = GetContext(entityPtr).URL.String()
	}
	self.Reportf("ambiguous %s name \"%s\", can be in %s", GetEntityTypeName(type_), format.ColorName(name), format.ColoredOptions(url, format.ColorValue))
}

func (self *Context) ReportFieldReferenceNotFound(types ...reflect.Type) {
	var entityTypeNames []string
	for _, type_ := range types {
		entityTypeNames = append(entityTypeNames, GetEntityTypeName(type_))
	}
	self.ReportPathf("reference to unknown %s: %s", format.Options(entityTypeNames), self.FormatBadData())
}

//
// Inheritance
//

func (self *Context) ReportInheritanceLoop(parent interface{}) {
	self.ReportPathf("inheritance loop by deriving from \"%s\"", format.ColorTypeName(GetContext(parent).Name))
}

func (self *Context) ReportTypeIncomplete(parent interface{}) {
	self.ReportPathf("deriving from incomplete type \"%s\"", format.ColorTypeName(GetContext(parent).Name))
}

//
// Render
//

func (self *Context) ReportUndefined(kind string) {
	self.ReportPathf("undefined %s", kind)
}

func (self *Context) ReportUnknown(kind string) {
	self.ReportPathf("unknown %s: %s", kind, self.FormatBadData())
}

func (self *Context) ReportReferenceNotFound(kind string, entityPtr interface{}) {
	typeName := GetEntityTypeName(reflect.TypeOf(entityPtr).Elem())
	name := GetContext(entityPtr).Name
	self.ReportPathf("unknown %s reference in %s \"%s\": %s", kind, typeName, format.ColorName(name), self.FormatBadData())
}

func (self *Context) ReportReferenceAmbiguous(kind string, entityPtr interface{}) {
	typeName := GetEntityTypeName(reflect.TypeOf(entityPtr).Elem())
	name := GetContext(entityPtr).Name
	self.ReportPathf("ambiguous %s in %s \"%s\": %s", kind, typeName, format.ColorName(name), self.FormatBadData())
}

func (self *Context) ReportPropertyRequired(kind string) {
	self.ReportPathf("unassigned required %s", kind)
}

func (self *Context) ReportReservedMetadata() {
	self.ReportPath("reserved for use by Puccini")
}

func (self *Context) ReportUnknownDataType(dataTypeName string) {
	self.ReportPathf("unknown data type \"%s\"", format.ColorError(dataTypeName))
}

func (self *Context) ReportMissingEntrySchema(kind string) {
	self.ReportPathf("missing entry schema for %s definition", kind)
}

func (self *Context) ReportUnsupportedType() {
	self.ReportPathf("unsupported puccini-tosca.type \"%s\"", format.ColorError(self.Name))
}

func (self *Context) ReportIncompatibleType(typeName string, parentTypeName string) {
	self.ReportPathf("type \"%s\" is incompatible with parent type \"%s\"", format.ColorTypeName(typeName), format.ColorTypeName(parentTypeName))
}

func (self *Context) ReportIncompatible(name string, typeName string, kind string) {
	self.ReportPathf("\"%s\" cannot be %s of %s", format.ColorName(name), kind, typeName)
}

func (self *Context) ReportIncompatibleExtension(extension string, requiredExtensions []string) {
	self.ReportPathf("extension \"%s\" is not %s", format.ColorValue(extension), format.ColoredOptions(requiredExtensions, format.ColorValue))
}

func (self *Context) ReportNotInRange(name string, value uint64, lower uint64, upper uint64) {
	self.ReportPathf("%s is %d, must be >= %d and <= %d", name, value, lower, upper)
}
