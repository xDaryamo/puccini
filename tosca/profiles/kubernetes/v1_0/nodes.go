// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/nodes.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

imports:
- capabilities.yaml
- relationships.yaml

node_types:

  kubernetes.Service:
    description: >-
      Represents a collection of workloads (pods and controllers) and resources that all use the
      same selector.
    derived_from: tosca.nodes.Root
    capabilities:
      metadata: kubernetes.Metadata
      service: kubernetes.Service
      deployment: kubernetes.Deployment
    requirements:
    - route:
        capability: kubernetes.Service
        relationship: kubernetes.Route
        occurrences: [ 0, UNBOUNDED ]
`
}
