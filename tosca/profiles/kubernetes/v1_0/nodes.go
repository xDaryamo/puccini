// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/nodes.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

imports:

- capabilities.yaml
- relationships.yaml
- interfaces.yaml

node_types:

  Service:
    description: >-
      Represents a collection of workloads (pods and controllers) and resources that all use the
      same selector.
    interfaces:
      lifecycle:
        type: Lifecycle
    capabilities:
      metadata: Metadata
      service: Service
      deployment: Deployment
    requirements:
    - route:
        capability: Service
        relationship: Route
        occurrences: [ 0, UNBOUNDED ]
`
}
