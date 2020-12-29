// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/data.yaml"] = `
tosca_definitions_version: tosca_2_0

metadata:

  puccini.scriptlet.import:tosca.comparer.version: js/comparers/version.js
  puccini.scriptlet.import:tosca.constraint._format: js/constraints/_format.js

data_types:

  #
  # Primitive
  #

  boolean:
    metadata:
      puccini.type: ard.boolean
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  integer:
    metadata:
      puccini.type: ard.integer
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  float:
    metadata:
      puccini.type: ard.float
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  string:
    metadata:
      puccini.type: ard.string
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  bytes:
    metadata:
      puccini.type: bytes
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  nil:
    metadata:
      puccini.type: ard.null
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  #
  # Special
  #

  version?:
    metadata:
      puccini.type: version
      puccini.comparer: tosca.comparer.version
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  timestamp:
    metadata:
      puccini.type: timestamp
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  range:
    metadata:
      puccini.type: range
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  tosca.datatypes.json?:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'
    # ERRATUM: typo
    description: >-
      The json type is a TOSCA data Type used to define a string that contains data in the
      JavaScript Object Notation (JSON) format.
    derived_from: string
    constraints:
    - _format: json

  tosca.datatypes.xml?:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'
    # ERRATUM: typo
    description: >-
      The xml type is a TOSCA data Type used to define a string that contains data in the
      Extensible Markup Language (XML) format.
    derived_from: string
    constraints:
    - _format: xml

  #
  # Scalar
  #

  scalar-unit.size:
    metadata:
      puccini.type: scalar-unit.size
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  scalar-unit.time:
    metadata:
      puccini.type: scalar-unit.time
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  scalar-unit.frequency:
    metadata:
      puccini.type: scalar-unit.frequency
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  scalar-unit.bitrate:
    metadata:
      puccini.type: scalar-unit.bitrate
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  #
  # Collections
  #

  list:
    metadata:
      puccini.type: ard.list
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'

  map:
    metadata:
      puccini.type: ard.map
      specification.citation: '[TOSCA-v2.0]'
      specification.location: '?'
`
}
