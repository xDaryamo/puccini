// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/simple/1.0/interfaces.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_0

interface_types:

  tosca.interfaces.Root:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 5.7.3
    description: >-
      This is the default (root) TOSCA Interface Type definition that all other TOSCA Interface
      Types derive from.

  tosca.interfaces.node.lifecycle.Standard:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 5.7.4
    description: >-
      This lifecycle interface defines the essential, normative operations that TOSCA nodes may
      support.
    derived_from: tosca.interfaces.Root
    create:
      description: >-
        Standard lifecycle create operation.
    configure:
      description: >-
        Standard lifecycle configure operation.
    start:
      description: >-
        Standard lifecycle start operation.
    stop:
      description: >-
        Standard lifecycle stop operation.
    delete:
      description: >-
        Standard lifecycle delete operation.

  tosca.interfaces.relationship.Configure:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      specification.location: 5.7.5
    description: >-
      The lifecycle interfaces define the essential, normative operations that each TOSCA
      Relationship Types may support.
    derived_from: tosca.interfaces.Root
    pre_configure_source:
      description: >-
        Operation to pre-configure the source endpoint.
    pre_configure_target:
      description: >-
        Operation to pre-configure the target endpoint.
    post_configure_source:
      description: >-
        Operation to post-configure the source endpoint.
    post_configure_target:
      description: >-
        Operation to post-configure the target endpoint.
    add_target:
      description: >-
        Operation to notify the source node of a target node being added via a relationship.
    add_source:
      description: >-
        Operation to notify the target node of a source node which is now available via a
        relationship.
    target_changed:
      description: >-
        Operation to notify source some property or attribute of the target changed
    remove_target:
      description: >-
        Operation to remove a target node.
    remove_source: # ERRATUM: does not appear in spec, but is mentioned
      description: >-
        Operation to remove the source node.
`
}
