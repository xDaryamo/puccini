// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/implicit/1.0/data.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_0

metadata:

  puccini.scriptlet.import:tosca.comparer.version: internal:/tosca/implicit/2.0/js/comparers/version.js

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
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 3.2.2

  range:
    metadata:
      puccini.type: range
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 3.2.3

  #
  # With entry schema
  #

  list:
    metadata:
      puccini.type: ard.list
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 3.2.4

  map:
    metadata:
      puccini.type: ard.map
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 3.2.5

  #
  # Scalar
  #

  scalar-unit.size:
    metadata:
      puccini.type: scalar-unit.size
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 3.2.6.4

  scalar-unit.time:
    metadata:
      puccini.type: scalar-unit.time
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 3.2.6.5

  scalar-unit.frequency:
    metadata:
      puccini.type: scalar-unit.frequency
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 3.2.6.6
`
}
