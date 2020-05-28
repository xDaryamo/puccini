// This file was auto-generated from a YAML file

package v1_2

func init() {
	Profile["/tosca/simple/1.2/capabilities.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_2

imports:
- data.yaml

capability_types:

  tosca.capabilities.Root:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.1
    description: >-
      This is the default (root) TOSCA Capability Type definition that all other TOSCA Capability
      Types derive from.

  tosca.capabilities.Node:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.2
    description: >-
      The Node capability indicates the base capabilities of a TOSCA Node Type.
    derived_from: tosca.capabilities.Root

  tosca.capabilities.Compute:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.3
      role: host
    description: >-
      The Compute capability, when included on a Node Type or Template definition, indicates that
      the node can provide hosting on a named compute resource.
    derived_from: tosca.capabilities.Container
    properties:
      name:
        description: >-
          The optional name (or identifier) of a specific compute resource for hosting.
        type: string
        required: false
      num_cpus:
        description: >-
          Number of (actual or virtual) CPUs associated with the Compute node.
        type: integer
        constraints:
        - greater_or_equal: 1
        required: false
      cpu_frequency:
        description: >-
          Specifies the operating frequency of CPU's core. This property expresses the expected
          frequency of one (1) CPU as provided by the property "num_cpus".
        type: scalar-unit.frequency
        constraints:
        - greater_or_equal: 0.1 GHz
        required: false
      disk_size:
        description: >-
          Size of the local disk available to applications running on the Compute node (default unit
          is MB).
        type: scalar-unit.size
        constraints:
        - greater_or_equal: 0 MB
        required: false
      mem_size:
        description: >-
          Size of memory available to applications running on the Compute node (default unit is MB).
        type: scalar-unit.size
        constraints:
        - greater_or_equal: 0 MB
        required: false

  tosca.capabilities.Network:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.4
    description: >-
      The Network capability, when included on a Node Type or Template definition, indicates that
      the node can provide addressiblity for the resource a named network with the specified ports.
    derived_from: tosca.capabilities.Root
    properties:
      name:
        description: >-
          The otional name (or identifier) of a specific network resource.
        type: string
        required: false

  tosca.capabilities.Storage:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.5
    description: >-
      The Storage capability, when included on a Node Type or Template definition, indicates that
      the node can provide a named storage location with specified size range.
    derived_from: tosca.capabilities.Root
    properties:
      name:
        description: >-
          The optional name (or identifier) of a specific storage resource.
        type: string
        required: false

  tosca.capabilities.Container:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.6
    description: >-
      The Container capability, when included on a Node Type or Template definition, indicates that
      the node can act as a container for (or a host for) one or more other declared Node Types.
    derived_from: tosca.capabilities.Root

  tosca.capabilities.Endpoint:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.7
    description: >-
      This is the default TOSCA type that should be used or extended to define a network endpoint
      capability. This includes the information to express a basic endpoint with a single port or a
      complex endpoint with multiple ports. By default the Endpoint is assumed to represent an
      address on a private network unless otherwise specified.
    derived_from: tosca.capabilities.Root
    properties:
      protocol:
        description: >-
          The name of the protocol (i.e., the protocol prefix) that the endpoint accepts (any OSI
          Layer 4-7 protocols). Examples: http, https, ftp, tcp, udp, etc.
        type: string
        default: tcp
      port:
        description: >-
          The optional port of the endpoint.
        type: tosca.datatypes.network.PortDef
        required: false
      secure:
        description: >-
          Requests for the endpoint to be secure and use credentials supplied on the ConnectsTo
          relationship.
        type: boolean
        default: false
        required: false
      url_path:
        description: >-
          The optional URL path of the endpoint's address if applicable for the protocol.
        type: string
        required: false
      port_name:
        description: >-
          The optional name (or ID) of the network port this endpoint should be bound to.
        type: string
        required: false
      network_name:
        description: >-
          The optional name (or ID) of the network this endpoint should be bound to. network_name:
          PRIVATE | PUBLIC | <network_name> | <network_id>.
        type: string
        default: PRIVATE
        required: false
      initiator:
        description: >-
          The optional indicator of the direction of the connection.
        type: string
        constraints:
        - valid_values: [ source, target, peer ]
        default: source
        required: false
      ports:
        description: >-
          The optional map of ports the Endpoint supports (if more than one).
        type: map
        entry_schema:
          type: tosca.datatypes.network.PortSpec
        constraints:
        - min_length: 1
        required: false
    attributes:
      ip_address:
        description: >-
          Note: This is the IP address as propagated up by the associated node's host (Compute)
          container.
        type: string

  tosca.capabilities.Endpoint.Public:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.8
    description: >-
      This capability represents a public endpoint which is accessible to the general internet (and
      its public IP address ranges).

      This public endpoint capability also can be used to create a floating (IP) address that the
      underlying network assigns from a pool allocated from the application's underlying public
      network. This floating address is managed by the underlying network such that can be routed an
      application's private address and remains reliable to internet clients.
    derived_from: tosca.capabilities.Endpoint
    properties:
      network_name:
        type: string
        constraints:
        - equal: PUBLIC
        default: PUBLIC
      floating:
        description: >-
          Indicates that the public address should be allocated from a pool of floating IPs that are
          associated with the network.
        type: boolean
        default: false
        status: experimental
      dns_name:
        description: >-
          The optional name to register with DNS.
        type: string
        required: false
        status: experimental

  tosca.capabilities.Endpoint.Admin:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.9
    description: >-
      This is the default TOSCA type that should be used or extended to define a specialized
      administrator endpoint capability.
    derived_from: tosca.capabilities.Endpoint
    properties:
      secure:
        description: >-
          Requests for the endpoint to be secure and use credentials supplied on the ConnectsTo
          relationship.
        type: boolean
        constraints:
        - equal: true
        default: true

  tosca.capabilities.Endpoint.Database:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.10
    description: >-
      This is the default TOSCA type that should be used or extended to define a specialized
      database endpoint capability.
    derived_from: tosca.capabilities.Endpoint

  tosca.capabilities.Attachment:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.11
    description: >-
      This is the default TOSCA type that should be used or extended to define an attachment
      capability of a (logical) infrastructure device node (e.g., BlockStorage node).
    derived_from: tosca.capabilities.Root

  tosca.capabilities.OperatingSystem:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.12
    description: >-
      This is the default TOSCA type that should be used to express an Operating System capability
      for a node.
    derived_from: tosca.capabilities.Root
    properties:
      architecture:
        description: >-
          The Operating System (OS) architecture. Examples of valid values include: x86_32, x86_64,
          etc.
        type: string
        required: false
      type:
        description: >-
          The Operating System (OS) type. Examples of valid values include: linux, aix, mac,
          windows, etc.
        type: string
        required: false
      distribution:
        description: >-
          The Operating System (OS) distribution. Examples of valid values for a "type" of "Linux"
          would include: debian, fedora, rhel and ubuntu.
        type: string
        required: false
      version:
        description: >-
          The Operating System version.
        type: version
        required: false

  tosca.capabilities.Scalable:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.13
    description: >-
      This is the default TOSCA type that should be used to express a scalability capability for a
      node.
    derived_from: tosca.capabilities.Root
    properties:
      min_instances:
        description: >-
          This property is used to indicate the minimum number of instances that should be created
          for the associated TOSCA Node Template by a TOSCA orchestrator.
        type: integer
        default: 1
      max_instances:
        description: >-
          This property is used to indicate the maximum number of instances that should be created
          for the associated TOSCA Node Template by a TOSCA orchestrator.
        type: integer
        default: 1
      default_instances:
        description: >-
          An optional property that indicates the requested default number of instances that should
          be the starting number of instances a TOSCA orchestrator should attempt to allocate. Note:
          The value for this property MUST be in the range between the values set for
          "min_instances" and "max_instances" properties.
        type: integer
        required: false

  tosca.capabilities.network.Bindable:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 5.5.14
    description: >-
      A node type that includes the Bindable capability indicates that it can be bound to a logical
      network association via a network port.
    derived_from: tosca.capabilities.Node

  #
  # Network
  #

  tosca.capabilities.network.Linkable:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.2]'
      specification.location: 8.5.3
    description: >-
      A node type that includes the Linkable capability indicates that it can be pointed by
      tosca.relationships.network.LinksTo relationship type.
    derived_from: tosca.capabilities.Node
`
}
