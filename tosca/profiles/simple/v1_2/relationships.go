// This file was auto-generated from a YAML file

package v1_2

func init() {
	Profile["/tosca/simple/1.2/relationships.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_2

imports:

- capabilities.yaml
- interfaces.yaml

relationship_types:

  tosca.relationships.Root:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.7.1
    description: >-
      This is the default (root) TOSCA Relationship Type definition that all other TOSCA
      Relationship Types derive from.
    attributes:
      tosca_id:
        description: >-
          A unique identifier of the realized instance of a Relationship Template that derives from
          any TOSCA normative type.
        type: string
      tosca_name:
        description: >-
          This attribute reflects the name of the Relationship Template as defined in the TOSCA
          service template. This name is not unique to the realized instance model of corresponding
          deployed application as each template in the model can result in one or more instances
          (e.g., scaled) when orchestrated to a provider environment.
        type: string
      state:
        description: >-
          The state of the relationship instance.
        type: string
        default: initial
    interfaces:
      Configure:
        type: tosca.interfaces.relationship.Configure

  tosca.relationships.DependsOn:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.7.2
    description: >-
      This type represents a general dependency relationship between two nodes.
    derived_from: tosca.relationships.Root
    valid_target_types: [ tosca.capabilities.Node ]

  tosca.relationships.HostedOn:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.7.3
      role: host
    description: >-
      This type represents a hosting relationship between two nodes.
    derived_from: tosca.relationships.Root
    valid_target_types: [ tosca.capabilities.Container ]

  tosca.relationships.ConnectsTo:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.7.4
    description: >-
      This type represents a network connection relationship between two nodes.
    derived_from: tosca.relationships.Root
    valid_target_types: [ tosca.capabilities.Endpoint ]
    properties:
      credential:
        type: tosca.datatypes.Credential
        required: false

  tosca.relationships.AttachesTo:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.7.5
    description: >-
      This type represents an attachment relationship between two nodes. For example, an AttachesTo
      relationship type would be used for attaching a storage node to a Compute node.
    derived_from: tosca.relationships.Root
    valid_target_types: [ tosca.capabilities.Attachment ]
    properties:
      location:
        description: >-
          The relative location (e.g., path on the file system), which provides the root location to
          address an attached node. e.g., a mount point / path such as '/usr/data'. Note: The user
          must provide it and it cannot be "root".
        type: string
        constraints:
        - min_length: 1
      device:
        description: >-
          The logical device name which for the attached device (which is represented by the target
          node in the model). e.g., '/dev/hda1'.
        type: string
        required: false
    attributes:
      device:
        description: >-
          The logical name of the device as exposed to the instance. Note: A runtime property that
          gets set when the model gets instantiated by the orchestrator.
        type: string

  tosca.relationships.RoutesTo:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.7.6
    description: >-
      This type represents an intentional network routing between two Endpoints in different
      networks.
    derived_from: tosca.relationships.ConnectsTo
    valid_target_types: [ tosca.capabilities.Endpoint ]

  #
  # Network
  #

  tosca.relationships.network.LinksTo:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 8.5.4
    description: >-
      This relationship type represents an association relationship between Port and Network node
      types.
    derived_from: tosca.relationships.DependsOn
    valid_target_types: [ tosca.capabilities.network.Linkable ]

  tosca.relationships.network.BindsTo:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 8.5.5
    description: >-
      This type represents a network association relationship between Port and Compute node types.
    derived_from: tosca.relationships.DependsOn
    valid_target_types: [ tosca.capabilities.network.Bindable ]
`
}
