package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
)

//
// DataDefinition
//

type DataDefinition interface {
	ToValueMeta() *normal.ValueMeta
	GetDescription() string
	GetTypeMetadata() Metadata
	GetValidationClause() *ValidationClause
	GetKeySchema() *Schema
	GetEntrySchema() *Schema
}
