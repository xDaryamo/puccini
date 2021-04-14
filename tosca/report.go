package tosca

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
)

//
// Context
//

func (self *Context) ReportURL(skip int, item string, message string, row int, column int) bool {
	if self.URL != nil {
		return self.Problems.ReportFull(skip+1, self.URL.String(), item, message, row, column)
	} else {
		return self.Problems.Report(skip+1, item, message)
	}
}

func (self *Context) Report(skip int, item string, message string) bool {
	row, column := self.GetLocation()
	return self.ReportURL(skip+1, item, message, row, column)
}

func (self *Context) Reportf(skip int, f string, arg ...interface{}) bool {
	return self.Report(skip+1, "", fmt.Sprintf(f, arg...))
}

func (self *Context) ReportPath(skip int, message string) bool {
	path := self.Path.String()
	if path != "" {
		path = self.Problems.Stylist.Path(path)
	}
	return self.Report(skip+1, path, message)
}

func (self *Context) ReportPathf(skip int, f string, arg ...interface{}) bool {
	return self.ReportPath(skip+1, fmt.Sprintf(f, arg...))
}

func (self *Context) ReportProblematic(skip int, problematic problems.Problematic) bool {
	// Note: we are ignoring the problem's section and using the URL instead
	_, item, message, row, column := problematic.Problem(self.Problems.Stylist)
	return self.ReportURL(skip+1, item, message, row, column)
}

func (self *Context) ReportError(err error) bool {
	if problematic, ok := err.(problems.Problematic); ok {
		return self.ReportProblematic(1, problematic)
	} else {
		return self.ReportPath(1, err.Error())
	}
}

//
// Values
//

func (self *Context) FormatBadData() string {
	return self.Problems.Stylist.Error(fmt.Sprintf("%+v", self.Data))
}

func (self *Context) ReportValueWrongType(allowedTypeNames ...ard.TypeName) bool {
	return self.ReportPathf(1, "%s instead of %s", self.Problems.Stylist.TypeName(quote(ardGetTypeName(self.Data))), terminal.StylizedOptions(ardTypeNamesToStrings(allowedTypeNames), self.Problems.Stylist.TypeName))
}

func (self *Context) ReportAspectWrongType(aspect string, value ard.Value, allowedTypeNames ...ard.TypeName) bool {
	return self.ReportPathf(1, "%s is %s instead of %s", aspect, self.Problems.Stylist.TypeName(quote(ardGetTypeName(value))), terminal.StylizedOptions(ardTypeNamesToStrings(allowedTypeNames), self.Problems.Stylist.TypeName))
}

func (self *Context) ReportValueWrongFormat(format string) bool {
	return self.ReportPathf(1, "wrong format, must be %s: %s", quote(format), self.FormatBadData())
}

func (self *Context) ReportValueWrongLength(kind string, length int) bool {
	return self.ReportPathf(1, "%s does not have %d elements", kind, length)
}

func (self *Context) ReportValueInvalid(kind string, reason string) bool {
	if reason == "" {
		return self.ReportPathf(1, "invalid %s: %s", kind, self.FormatBadData())
	} else {
		return self.ReportPathf(1, "invalid %s, %s: %s", kind, reason, self.FormatBadData())
	}
}

func (self *Context) ReportValueMalformed(kind string, reason string) bool {
	if reason == "" {
		return self.ReportPathf(1, "malformed %s: %s", kind, self.FormatBadData())
	} else {
		return self.ReportPathf(1, "malformed %s, %s: %s", kind, reason, self.FormatBadData())
	}
}

//
// Read
//

func (self *Context) ReportImportIncompatible(url urlpkg.URL) bool {
	return self.Reportf(1, "incompatible import %s", self.Problems.Stylist.Value(quote(url.String())))
}

func (self *Context) ReportImportLoop(url urlpkg.URL) bool {
	return self.Reportf(1, "endless loop caused by importing %s", self.Problems.Stylist.Value(quote(url.String())))
}

func (self *Context) ReportRepositoryInaccessible(repositoryName string) bool {
	return self.ReportPathf(1, "inaccessible repository %s", self.Problems.Stylist.Value(quote(repositoryName)))
}

func (self *Context) ReportFieldMissing() bool {
	return self.ReportPath(1, "missing required field")
}

func (self *Context) ReportFieldUnsupported() bool {
	return self.ReportPath(1, "unsupported field")
}

func (self *Context) ReportFieldUnsupportedValue() bool {
	return self.ReportPathf(1, "unsupported value for field: %s", self.FormatBadData())
}

func (self *Context) ReportFieldMalformedSequencedList() bool {
	return self.ReportPathf(1, "field must be a %s of single-key %s elements", self.Problems.Stylist.TypeName(quote("sequenced list")), self.Problems.Stylist.TypeName(quote("map")))
}

func (self *Context) ReportPrimitiveType() bool {
	return self.ReportPath(1, "primitive type cannot have properties")
}

func (self *Context) ReportDuplicateMapKey(key string) bool {
	return self.ReportPathf(1, "duplicate map key: %s", self.Problems.Stylist.Value(key))
}

//
// Namespaces
//

func (self *Context) ReportNameAmbiguous(type_ reflect.Type, name string, entityPtrs ...EntityPtr) bool {
	return self.Reportf(1, "ambiguous %s name %s, can be in %s", GetEntityTypeName(type_), self.Problems.Stylist.Name(quote(name)), terminal.StylizedOptions(urlsOfEntityPtrs(entityPtrs), self.Problems.Stylist.Value))
}

func (self *Context) ReportFieldReferenceNotFound(types ...reflect.Type) bool {
	return self.ReportPathf(1, "reference to unknown %s: %s", terminal.Options(entityTypeNamesOfTypes(types)), self.FormatBadData())
}

//
// Inheritance
//

func (self *Context) ReportInheritanceLoop(parentType EntityPtr) bool {
	return self.ReportPathf(1, "inheritance loop by deriving from %s", self.Problems.Stylist.TypeName(quote(GetCanonicalName(parentType))))
}

func (self *Context) ReportTypeIncomplete(parentType EntityPtr) bool {
	return self.ReportPathf(1, "deriving from incomplete type %s", self.Problems.Stylist.TypeName(quote(GetCanonicalName(parentType))))
}

//
// Render
//

func (self *Context) ReportUndeclared(kind string) bool {
	return self.ReportPathf(1, "undeclared %s", kind)
}

func (self *Context) ReportUnknown(kind string) bool {
	return self.ReportPathf(1, "unknown %s: %s", kind, self.FormatBadData())
}

func (self *Context) ReportReferenceNotFound(kind string, entityPtr EntityPtr) bool {
	typeName := GetEntityTypeName(reflect.TypeOf(entityPtr).Elem())
	name := GetContext(entityPtr).Name
	return self.ReportPathf(1, "unknown %s reference in %s %s: %s", kind, typeName, self.Problems.Stylist.Name(quote(name)), self.FormatBadData())
}

func (self *Context) ReportReferenceAmbiguous(kind string, entityPtr EntityPtr) bool {
	typeName := GetEntityTypeName(reflect.TypeOf(entityPtr).Elem())
	name := GetContext(entityPtr).Name
	return self.ReportPathf(1, "ambiguous %s in %s %s: %s", kind, typeName, self.Problems.Stylist.Name(quote(name)), self.FormatBadData())
}

func (self *Context) ReportPropertyRequired(kind string) bool {
	return self.ReportPathf(1, "unassigned required %s", kind)
}

func (self *Context) ReportReservedMetadata() bool {
	return self.ReportPath(1, "reserved for use by Puccini")
}

func (self *Context) ReportUnknownDataType(dataTypeName string) bool {
	return self.ReportPathf(1, "unknown data type %s", self.Problems.Stylist.Error(quote(dataTypeName)))
}

func (self *Context) ReportMissingEntrySchema(kind string) bool {
	return self.ReportPathf(1, "missing entry schema for %s definition", kind)
}

func (self *Context) ReportUnsupportedType() bool {
	return self.ReportPathf(1, "unsupported puccini.type %s", self.Problems.Stylist.Error(quote(self.Name)))
}

func (self *Context) ReportIncompatibleType(type_ EntityPtr, parentType EntityPtr) bool {
	return self.ReportPathf(1, "type %s must be derived from type %s", self.Problems.Stylist.TypeName(quote(GetCanonicalName(type_))), self.Problems.Stylist.TypeName(quote(GetCanonicalName(parentType))))
}

func (self *Context) ReportIncompatibleTypeInSet(type_ EntityPtr) bool {
	return self.ReportPathf(1, "type %s must be derived from one of the types in the parent set", self.Problems.Stylist.TypeName(quote(GetCanonicalName(type_))))
}

func (self *Context) ReportIncompatible(name string, target string, kind string) bool {
	return self.ReportPathf(1, "%s cannot be %s of %s", self.Problems.Stylist.Name(quote(name)), kind, target)
}

func (self *Context) ReportIncompatibleExtension(extension string, requiredExtensions []string) bool {
	return self.ReportPathf(1, "extension %s is not %s", self.Problems.Stylist.Value(quote(extension)), terminal.StylizedOptions(requiredExtensions, self.Problems.Stylist.Value))
}

func (self *Context) ReportNotInRange(name string, value uint64, lower uint64, upper uint64) bool {
	return self.ReportPathf(1, "%s is %d, must be >= %d and <= %d", name, value, lower, upper)
}

func (self *Context) ReportCopyLoop(name string) bool {
	return self.ReportPathf(1, "endless loop caused by copying %s", self.Problems.Stylist.Value(quote(name)))
}

// Utils

func quote(value interface{}) string {
	return fmt.Sprintf("%q", value)
}

func ardTypeNameToString(typeName ard.TypeName) string {
	typeName_ := string(typeName)
	if strings.HasPrefix(typeName_, "ard.") {
		typeName_ = typeName_[4:]
	}
	return typeName_
}

func ardTypeNamesToStrings(typeNames []ard.TypeName) []string {
	strings_ := make([]string, len(typeNames))
	for index, typeName := range typeNames {
		strings_[index] = ardTypeNameToString(typeName)
	}
	return strings_
}

func ardGetTypeName(value ard.Value) string {
	return ardTypeNameToString(ard.GetTypeName(value))
}

func urlsOfEntityPtrs(entityPtrs []EntityPtr) []string {
	urls := make([]string, len(entityPtrs))
	for index, entityPtr := range entityPtrs {
		urls[index] = GetContext(entityPtr).URL.String()
	}
	return urls
}

func entityTypeNamesOfTypes(types []reflect.Type) []string {
	entityTypeNames := make([]string, len(types))
	for index, type_ := range types {
		entityTypeNames[index] = GetEntityTypeName(type_)
	}
	return entityTypeNames
}
