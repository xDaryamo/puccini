package tosca

import (
	"fmt"
	"reflect"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca/problems"
	"github.com/tliron/puccini/url"
)

//
// Context
//

func (self *Context) Report(message string) {
	*self.Problems = append(*self.Problems, problems.Problem{Message: message, URL: self.URL.String()})
}

func (self *Context) Reportf(format string, arg ...interface{}) {
	self.Report(fmt.Sprintf(format, arg...))
}

func (self *Context) ReportPath(message string) {
	self.Report(fmt.Sprintf("%s: %s", format.ColorPath(self.Path), message))
}

func (self *Context) ReportPathf(format string, arg ...interface{}) {
	self.ReportPath(fmt.Sprintf(format, arg...))
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
	self.ReportPathf("repository \"%s\" inaccessible", format.ColorValue(repositoryName))
}

func (self *Context) ReportFieldMissing() {
	self.ReportPath("field is required")
}

func (self *Context) ReportFieldUnsupported() {
	self.ReportPath("field is unsupported")
}

func (self *Context) ReportFieldUnsupportedValue() {
	self.ReportPathf("field has unsupported value: %s", self.FormatBadData())
}

func (self *Context) ReportFieldMalformedSequencedList() {
	self.ReportPathf("field must be a \"%s\" of single-key \"%s\" elements", format.ColorTypeName("sequenced list"), format.ColorTypeName("map"))
}

func (self *Context) ReportPrimitiveType() {
	self.ReportPath("primitive type cannot have properties")
}

func (self *Context) ReportMapKeyReused(key string) {
	self.ReportPathf("map key reused: %s", format.ColorValue(key))
}

//
// Namespaces
//

func (self *Context) ReportNameAmbiguous(type_ reflect.Type, name string, entityPtrs ...interface{}) {
	url := make([]string, len(entityPtrs))
	for i, entityPtr := range entityPtrs {
		url[i] = GetContext(entityPtr).URL.String()
	}
	self.Reportf("%s name \"%s\" is ambiguous, can be in %s", GetEntityTypeName(type_), format.ColorName(name), format.ColoredOptions(url, format.ColorValue))
}

func (self *Context) ReportFieldReferenceNotFound(types ...reflect.Type) {
	var entityTypeNames []string
	for _, type_ := range types {
		entityTypeNames = append(entityTypeNames, GetEntityTypeName(type_))
	}
	self.ReportPathf("field refers to unknown %s: %s", format.Options(entityTypeNames), self.FormatBadData())
}

//
// Inheritance
//

func (self *Context) ReportInheritanceLoop(parent interface{}) {
	self.ReportPathf("inheritance loop by deriving from \"%s\"", format.ColorTypeName(GetContext(parent).Name))
}

func (self *Context) ReportTypeIncomplete(parent interface{}) {
	self.ReportPathf("derives from incomplete type \"%s\"", format.ColorTypeName(GetContext(parent).Name))
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

//
// Normalize
//

func (self *Context) ReportUnsatisfiedRequirement() {
	self.ReportPathf("cannot satisfy requirement \"%s\"", format.ColorValue(self.Name))
}
