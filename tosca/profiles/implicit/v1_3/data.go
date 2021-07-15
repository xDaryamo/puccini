// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/implicit/1.3/data.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

metadata:

  puccini.scriptlet.import:tosca.comparer.version: internal:/tosca/implicit/2.0/js/comparers/version.js
  puccini.scriptlet.import:tosca.constraint._format: internal:/tosca/implicit/2.0/js/constraints/_format.js

data_types:

  #
  # Primitive
  #

  boolean:
    metadata:
      puccini.type: ard.boolean

  integer:
    metadata:
      puccini.type: ard.integer

  float:
    metadata:
      puccini.type: ard.float

  string:
    metadata:
      puccini.type: ard.string

  'null':
    metadata:
      puccini.type: ard.null

  timestamp:
    metadata:
      puccini.type: timestamp

  #
  # Special
  #

  version:
    metadata:
      puccini.type: version
      puccini.comparer: tosca.comparer.version
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.2

  range:
    metadata:
      puccini.type: range
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.3

  tosca.datatypes.json:
    metadata:
      tosca.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 5.3.2
    # ERRATUM: typo
    description: >-
      The json type is a TOSCA data Type used to define a string that contains data in the
      JavaScript Object Notation (JSON) format.
    derived_from: string
    constraints:
    - _format: json

  tosca.datatypes.xml:
    metadata:
      tosca.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 5.3.4
    # ERRATUM: typo
    description: >-
      The xml type is a TOSCA data Type used to define a string that contains data in the
      Extensible Markup Language (XML) format.
    derived_from: string
    constraints:
    - _format: xml

  #
  # With entry schema
  #

  list:
    metadata:
      puccini.type: ard.list
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.4

  map:
    metadata:
      puccini.type: ard.map
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.5

  #
  # Scalar
  #

  scalar-unit.size:
    metadata:
      puccini.type: scalar-unit.size
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.6.4

  scalar-unit.time:
    metadata:
      puccini.type: scalar-unit.time
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.6.5

  scalar-unit.frequency:
    metadata:
      puccini.type: scalar-unit.frequency
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.6.6

  scalar-unit.bitrate:
    metadata:
      puccini.type: scalar-unit.bitrate
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 3.3.6.7
`
}
