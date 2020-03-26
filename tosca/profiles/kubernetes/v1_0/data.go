// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/data.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

data_types:

  Count:
    derived_from: integer
    constraints:
    - greater_or_equal: 0

  Factor:
    derived_from: float
    constraints:
    - in_range: [ 0.0, 1.0 ]

  Amount:
    properties:
      factor:
        type: Factor
        required: false
      count:
        type: Count
        required: false

  IP:
    derived_from: string

  Port:
    derived_from: integer
    constraints:
    - in_range: [ 1, 65535 ]

  # https://stackoverflow.com/questions/2063213/regular-expression-for-validating-dns-label-host-name
  Hostname:
    derived_from: string
    constraints:
    - pattern: ^(?!-)[a-zA-Z0-9-]{1,63}(?<!-)$

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#labelselector-v1-meta
  LabelSelector:
    description: >-
      A label selector is a label query over a set of resources. The result of matchLabels and
      matchExpressions are ANDed. An empty label selector matches all objects. A null label selector
      matches no objects.
    properties:
      matchExpressions:
        description: >-
          matchExpressions is a list of label selector requirements. The requirements are ANDed.
        type: list
        entry_schema: LabelSelectorRequirement
        required: false
      matchLabels:
        description: >-
          matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is
          equivalent to an element of matchExpressions, whose key field is "key", the operator is
          "In", and the values array contains only "value". The requirements are ANDed.
        type: map
        entry_schema: string
        required: false

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#labelselectorrequirement-v1-meta
  LabelSelectorRequirement:
    description: >-
      A label selector requirement is a selector that contains values, a key, and an operator that
      relates the key and values.
    properties:
      key:
        description: >-
          key is the label key that the selector applies to.
        type: string
      operator:
        description: >-
          operator represents a key's relationship to a set of values. Valid operators are In,
          NotIn, Exists and DoesNotExist.
        type: string
        constraints:
        - valid_values: [ In, NotIn, Exists, DoesNotExist ]
      values:
        description: >-
          values is an array of string values. If the operator is In or NotIn, the values array
          must be non-empty. If the operator is Exists or DoesNotExist, the values array must be
          empty. This array is replaced during a strategic merge patch.
        type: list
        entry_schema: string

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#serviceport-v1-core
  ServicePort:
    description: >-
      ServicePort contains information on service's port.
    properties:
      name:
        description: >-
          The name of this port within the service. This must be a DNS_LABEL. All ports within a
          ServiceSpec must have unique names. This maps to the 'Name' field in EndpointPort
          objects. Optional if only one ServicePort is defined on this service.
        type: Hostname
        required: false
      port:
        description: >-
          The port that will be exposed by this service.
        type: Port
      nodePort:
        description: >-
          The port on each node on which this service is exposed when type=NodePort or LoadBalancer.
          Usually assigned by the system. If specified, it will be allocated to the service if
          unused or else creation of the service will fail. Default is to auto-allocate a port if
          the ServiceType of this Service requires one.
        type: Port
        required: false
      targetPort:
        description: >-
          Number or name of the port to access on the pods targeted by the service. Number must be
          in the range 1 to 65535. Name must be an IANA_SVC_NAME. If this is not specified, the
          value of the 'port' field is used (an identity map). This field is ignored for services
          with clusterIP=None, and should be omitted or set equal to the 'port' field.
        type: Port
        required: false
      targetPortName:
        description: >-
          If this is a string, it will be looked up as a named port in the target Pod's container
          ports.
        type: string
        required: false
      protocol:
        description: >-
          The IP protocol for this port. Supports "TCP" and "UDP". Default is TCP.
        type: string
        default: TCP
        constraints:
        - valid_values: [ TCP, UDP ]

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#containerport-v1-core
  ContainerPort:
    description: >-
      ContainerPort represents a network port in a single container.
    properties:
      name:
        description: >-
          If specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in
          a pod must have a unique name. Name for the port that can be referred to by services.
        type: string
        required: false
      containerPort:
        description: >-
          Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x
          < 65536.
        type: Port
      hostPort:
        description: >-
          Number of port to expose on the host. If specified, this must be a valid port number, 0 <
          x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do
          not need this.
        type: Port
        required: false
      hostIP:
        description: >-
          What host IP to bind the external port to.
        type: IP
        required: false
      protocol:
        description: >-
          Protocol for port. Must be UDP or TCP. Defaults to "TCP".
        type: string
        default: TCP
        constraints:
        - valid_values: [ TCP, UDP ]

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#deploymentstrategy-v1-apps
  DeploymentStrategy:
    description: >-
      The deployment strategy to use to replace existing pods with new ones.
    properties:
      type:
        description: >-
          Type of deployment. Can be "Recreate" or "RollingUpdate". Default is RollingUpdate.
        type: string
        default: RollingUpdate
        constraints:
        - valid_values: [ Recreate, RollingUpdate ]
      maxSurge:
        description: >-
          The maximum number of pods that can be scheduled above the desired number of pods. Value
          can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). This can not
          be 0 if MaxUnavailable is 0. Absolute number is calculated from percentage by rounding up.
          Defaults to 25%. Example: when this is set to 30%, the new RC can be scaled up immediately
          when the rolling update starts, such that the total number of old and new pods do not
          exceed 130% of desired pods. Once old pods have been killed, new RC can be scaled up
          further, ensuring that total number of pods running at any time during the update is at
          most 130% of desired pods.
        type: Amount
        default:
          factor: .25
      maxUnavailable:
        description: >-
          The maximum number of pods that can be unavailable during the update. Value can be an
          absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is
          calculated from percentage by rounding down. This can not be 0 if MaxSurge is 0. Defaults
          to 25%. Example: when this is set to 30%, the old RC can be scaled down to 70% of desired
          pods immediately when the rolling update starts. Once new pods are ready, old RC can be
          scaled down further, followed by scaling up the new RC, ensuring that the total number of
          pods available at all times during the update is at least 70% of desired pods.
        type: Amount
        default:
          factor: .25

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#podspec-v1-core
  Pod:
    description: >-
      Pod is a collection of containers that can run on a host. This resource is created by clients
      and scheduled onto hosts.
    properties:
      # Resources
      containers:
        description: >-
          List of containers belonging to the pod. Containers cannot currently be added or removed.
          There must be at least one container in a Pod. Cannot be updated.
        type: list
        entry_schema: Container
        constraints:
        - min_length: 1
      initContainers:
        description: >-
          List of initialization containers belonging to the pod. Init containers are executed in
          order prior to containers being started. If any init container fails, the pod is
          considered to have failed and is handled according to its restartPolicy. The name for an
          init container or normal container must be unique among all containers. Init containers
          may not have Lifecycle actions, Readiness probes, or Liveness probes. The
          resourceRequirements of an init container are taken into account during scheduling by
          finding the highest request/limit for each resource type, and then using the max of of
          that value or the sum of the normal containers. Limits are applied to init containers in a
          similar fashion. Init containers cannot currently be added or removed. Cannot be updated.
        type: list
        entry_schema: Container
        required: false
      volumes:
        description: >-
          List of volumes that can be mounted by containers belonging to the pod.
        type: list
        entry_schema: Volume
        required: false

      # Scheduling
      schedulerName:
        description: >-
          If specified, the pod will be dispatched by specified scheduler. If not specified, the pod
          will be dispatched by default scheduler.
        type: string
        required: false
      affinity:
        description: >-
          If specified, the pod's scheduling constraints
        type: Affinity
        required: false
      nodeName:
        description: >-
          NodeName is a request to schedule this pod onto a specific node. If it is non-empty, the
          scheduler simply schedules this pod onto that node, assuming that it fits resource
          requirements.
        type: string
        required: false
      nodeSelector:
        description: >-
          NodeSelector is a selector which must be true for the pod to fit on a node. Selector which
          must match a node's labels for the pod to be scheduled on that node.
        type: string # TODO
        required: false

      # Priority
      priorityClassName:
        description: >-
          If specified, indicates the pod's priority. "system-node-critical" and
          "system-cluster-critical" are two special keywords which indicate the highest priorities
          with the former being the highest priority. Any other name must be defined by creating a
          PriorityClass object with that name. If not specified, the pod priority will be default or
          zero if there is no default.
        type: string
        required: false
      priority:
        description: >-
          The priority value. Various system components use this field to find the priority of the
          pod. When Priority Admission Controller is enabled, it prevents users from setting this
          field. The admission controller populates this field from PriorityClassName. The higher
          the value, the higher the priority.
        type: integer
        required: false

      # Lifecycle
      restartPolicy:
        description: >-
          Restart policy for all containers within the pod. One of Always, OnFailure, Never. Default
          to Always.
        type: string
        default: Always
        constraints:
        - valid_values: [ Always, OnFailure, Never ]
      tolerations:
        description: >-
          If specified, the pod's tolerations.
        type: string # TODO
        required: false
      activeDeadlineSeconds:
        description: >-
          Optional duration in seconds the pod may be active on the node relative to StartTime
          before the system will actively try to mark it failed and kill associated containers.
          Value must be a positive integer.
        type: scalar-unit.time
        required: false
      terminationGracePeriodSeconds:
        description: >-
          Optional duration in seconds the pod needs to terminate gracefully. May be decreased in
          delete request. Value must be non-negative integer. The value zero indicates delete
          immediately. If this value is nil, the default grace period will be used instead. The
          grace period is the duration in seconds after the processes running in the pod are sent a
          termination signal and the time when the processes are forcibly halted with a kill signal.
          Set this value longer than the expected cleanup time for your process. Defaults to 30
          seconds.
        type: scalar-unit.time
        required: false

      # DNS
      hostname:
        description: >-
          Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a
          system-defined value.
        type: Hostname
        required: false
      subdomain:
        description: >-
          If specified, the fully qualified Pod hostname will be
          "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>". If not specified, the pod
          will not have a domainname at all.
        type: string
        required: false
      dnsConfig:
        description: >-
          Specifies the DNS parameters of a pod. Parameters specified here will be merged to the
          generated DNS configuration based on DNSPolicy.
        type: string # TODO
        required: false
      dnsPolicy:
        description: >-
          Set DNS policy for the pod. Defaults to "ClusterFirst". Valid values are
          'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. DNS parameters given in
          DNSConfig will be merged with the policy selected with DNSPolicy. To have DNS options set
          along with hostNetwork, you have to specify DNS policy explicitly to
          'ClusterFirstWithHostNet'.
        type: string
        default: ClusterFirst
        constraints:
        - valid_values: [ ClusterFirstWithHostNet, ClusterFirst, Default, None ]
      hostAliases:
        description: >-
          HostAliases is an optional list of hosts and IPs that will be injected into the pod's
          hosts file if specified. This is only valid for non-hostNetwork pods.
        type: string # TODO
        required: false

      # Host
      hostIPC:
        description: >-
          Use the host's ipc namespace. Optional: Default to false.
        type: boolean
        default: false
      hostNetwork:
        description: >-
          Host networking requested for this pod. Use the host's network namespace. If this option
          is set, the ports that will be used must be specified. Default to false.
        type: boolean
        default: false
      hostPID:
        description: >-
          Use the host's pid namespace. Optional: Default to false.
        type: boolean
        default: false
      shareProcessNamespace:
        description: >-
          Share a single process namespace between all of the containers in a pod. When this is set
          containers will be able to view and signal processes from other containers in the same
          pod, and the first process in each container will not be assigned PID 1. HostPID and
          ShareProcessNamespace cannot both be set. Optional: Default to false. This field is
          alpha-level and is honored only by servers that enable the PodShareProcessNamespace
          feature.
        type: boolean
        default: false

      # Security
      securityContext:
        description: >-
          SecurityContext holds pod-level security attributes and common container settings.
          Optional: Defaults to empty. See type description for default values of each field.
        type: string # TODO
        required: false
      serviceAccountName:
        description: >-
          ServiceAccountName is the name of the ServiceAccount to use to run this pod.
        type: string
        required: false
      automountServiceAccountToken:
        description: >-
          AutomountServiceAccountToken indicates whether a service account token should be
          automatically mounted.
        type: boolean
        required: false
      imagePullSecrets:
        description: >-
          ImagePullSecrets is an optional list of references to secrets in the same namespace to use
          for pulling any of the images used by this PodSpec. If specified, these secrets will be
          passed to individual puller implementations for them to use. For example, in the case of
          docker, only DockerConfig type secrets are honored.
        type: string # TODO
        required: false

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#container-v1-core
  Container:
    description: >-
      A single application container that you want to run within a pod.
    properties:
      name:
        description: >-
          Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique
          name (DNS_LABEL). Cannot be updated.
        type: Hostname

      # Image
      image:
        description: >-
          Docker image name. This field is optional to allow higher level config management to
          default or override container images in workload controllers like Deployments and
          StatefulSets.
        type: string
        required: false
      imagePullPolicy:
        description: >-
          Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag
          is specified, or IfNotPresent otherwise. Cannot be updated.
        type: string
        default: Always
        constraints:
        - valid_values: [ Always, Never, IfNotPresent ]

      # Resources
      resources:
        description: >-
          Compute Resources required by this container. Cannot be updated.
        type: string # TODO
        required: false
      volumeMounts:
        description: >-
          Pod volumes to mount into the container's filesystem. Cannot be updated.
        type: string # TODO
        required: false
      volumeDevices:
        description: >-
          volumeDevices is the list of block devices to be used by the container. This is an alpha
          feature and may change in the future.
        type: string # TODO
        required: false

      ports:
        description: >-
          List of ports to expose from the container. Exposing a port here gives the system
          additional information about the network connections a container uses, but is primarily
          informational. Not specifying a port here DOES NOT prevent that port from being exposed.
          Any port which is listening on the default "0.0.0.0" address inside a container will be
          accessible from the network. Cannot be updated.
        type: list
        entry_schema: ContainerPort
        required: false

      # Terminal
      stdin:
        description: >-
          Whether this container should allocate a buffer for stdin in the container runtime. If
          this is not set, reads from stdin in the container will always result in EOF. Default is
          false.
        type: boolean
        default: false
      stdinOnce:
        description: >-
          Whether the container runtime should close the stdin channel after it has been opened by a
          single attach. When stdin is true the stdin stream will remain open across multiple attach
          sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until
          the first client attaches to stdin, and then remains open and accepts data until the
          client disconnects, at which time stdin is closed and remains closed until the container
          is restarted. If this flag is false, a container processes that reads from stdin will
          never receive an EOF. Default is false
        type: boolean
        default: false
      tty:
        description: >-
          Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.
          Default is false.
        type: boolean
        default: false

      # Execution
      workingDir:
        description: >-
          Container's working directory. If not specified, the container runtime's default will be
          used, which might be configured in the container image. Cannot be updated.
        type: string # TODO
        required: false
      command:
        description: >-
          Entrypoint array. Not executed within a shell. The docker image's ENTRYPOINT is used if
          this is not provided. Variable references $(VAR_NAME) are expanded using the container's
          environment. If a variable cannot be resolved, the reference in the input string will be
          unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME).
          Escaped references will never be expanded, regardless of whether the variable exists or
          not. Cannot be updated.
        type: list
        entry_schema: string
        required: false
      args:
        description: >-
          Arguments to the entrypoint. The docker image's CMD is used if this is not provided.
          Variable references $(VAR_NAME) are expanded using the container's environment. If a
          variable cannot be resolved, the reference in the input string will be unchanged. The
          $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references
          will never be expanded, regardless of whether the variable exists or not. Cannot be
          updated.
        type: list
        entry_schema: string
        required: false
      env:
        description: >-
          List of environment variables to set in the container. Cannot be updated.
        type: list
        entry_schema: string # TODO
        required: false
      envFrom:
        description: >-
          List of sources to populate environment variables in the container. The keys defined
          within a source must be a C_IDENTIFIER. All invalid keys will be reported as an event when
          the container is starting. When a key exists in multiple sources, the value associated
          with the last source will take precedence. Values defined by an Env with a duplicate key
          will take precedence. Cannot be updated.
        type: list
        entry_schema: string # TODO
        required: false

      # Lifecycle
      lifecycle:
        description: >-
          Actions that the management system should take in response to container lifecycle events.
          Cannot be updated.
        type: string # TODO
        required: false
      livenessProbe:
        description: >-
          Periodic probe of container liveness. Container will be restarted if the probe fails.
          Cannot be updated.
        type: string # TODO
        required: false
      readinessProbe:
        description: >-
          Periodic probe of container service readiness. Container will be removed from service
          endpoints if the probe fails. Cannot be updated.
        type: string # TODO
        required: false
      terminationMessagePolicy:
        description: >-
          Indicate how the termination message should be populated. File will use the contents of
          terminationMessagePath to populate the container status message on both success and
          failure. FallbackToLogsOnError will use the last chunk of container log output if the
          termination message file is empty and the container exited with an error. The log output
          is limited to 2048 bytes or 80 lines, whichever is smaller. Defaults to File. Cannot be
          updated.
        type: string # TODO
        required: false
      terminationMessagePath:
        description: >-
          Optional: Path at which the file to which the container's termination message will be
          written is mounted into the container's filesystem. Message written is intended to be
          brief final status, such as an assertion failure message. Will be truncated by the node
          if greater than 4096 bytes. The total message length across all containers will be limited
          to 12kb. Defaults to /dev/termination-log. Cannot be updated.
        type: string # TODO
        required: false

      # Security
      securityContext:
        description: >-
          Security options the pod should run with.
        type: string # TODO
        required: false

  Volume: {}

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#affinity-v1-core
  Affinity:
    description: >-
      Affinity is a group of affinity scheduling rules.
    properties:
      nodeAffinity:
        description: >-
          Describes node affinity scheduling rules for the pod.
        type: NodeAffinity
        required: false
      podAffinity:
        description: >-
          Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone,
          etc. as some other pod(s)).
        type: PodAffinity
        required: false
      podAntiAffinity:
        description: >-
          Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same
          node, zone, etc. as some other pod(s)).
        type: PodAntiAffinity
        required: false

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#affinity-v1-core
  NodeAffinity:
    description: >-
      Node affinity is a group of node affinity scheduling rules.
    properties:
      preferredDuringSchedulingIgnoredDuringExecution:
        description: >-
          The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions
          specified by this field, but it may choose a node that violates one or more of the
          expressions. The node that is most preferred is the one with the greatest sum of weights,
          i.e. for each node that meets all of the scheduling requirements (resource request,
          requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through
          the elements of this field and adding "weight" to the sum if the node matches the
          corresponding matchExpressions; the node(s) with the highest sum are the most preferred.
        type: string # TODO
        required: false
      requiredDuringSchedulingIgnoredDuringExecution:
        description: >-
          If the affinity requirements specified by this field are not met at scheduling time, the
          pod will not be scheduled onto the node. If the affinity requirements specified by this
          field cease to be met at some point during pod execution (e.g. due to an update), the
          system may or may not try to eventually evict the pod from its node.
        type: string # TODO
        required: false

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#podaffinity-v1-core
  PodAffinity:
    description: >-
      Pod affinity is a group of inter pod affinity scheduling rules.
    properties:
      preferredDuringSchedulingIgnoredDuringExecution:
        description: >-
          The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions
          specified by this field, but it may choose a node that violates one or more of the
          expressions. The node that is most preferred is the one with the greatest sum of weights,
          i.e. for each node that meets all of the scheduling requirements (resource request,
          requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through
          the elements of this field and adding "weight" to the sum if the node has pods which
          matches the corresponding podAffinityTerm; the node(s) with the highest sum are the most
          preferred.
        type: string # TODO
        required: false
      requiredDuringSchedulingIgnoredDuringExecution:
        description: >-
          If the affinity requirements specified by this field are not met at scheduling time, the
          pod will not be scheduled onto the node. If the affinity requirements specified by this
          field cease to be met at some point during pod execution (e.g. due to a pod label update),
          the system may or may not try to eventually evict the pod from its node. When there are
          multiple elements, the lists of nodes corresponding to each podAffinityTerm are
          intersected, i.e. all terms must be satisfied.
        type: string # TODO
        required: false

  # https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#podantiaffinity-v1-core
  PodAntiAffinity:
    description: >-
      Pod anti affinity is a group of inter pod anti affinity scheduling rules.
    properties:
      preferredDuringSchedulingIgnoredDuringExecution:
        description: >-
          The scheduler will prefer to schedule pods to nodes that satisfy the anti-affinity
          expressions specified by this field, but it may choose a node that violates one or more of
          the expressions. The node that is most preferred is the one with the greatest sum of
          weights, i.e. for each node that meets all of the scheduling requirements (resource
          request, requiredDuringScheduling anti-affinity expressions, etc.), compute a sum by
          iterating through the elements of this field and adding "weight" to the sum if the node
          has pods which matches the corresponding podAffinityTerm; the node(s) with the highest sum
          are the most preferred.
        type: string # TODO
        required: false
      requiredDuringSchedulingIgnoredDuringExecution:
        description: >-
          If the anti-affinity requirements specified by this field are not met at scheduling time,
          the pod will not be scheduled onto the node. If the anti-affinity requirements specified
          by this field cease to be met at some point during pod execution (e.g. due to a pod label
          update), the system may or may not try to eventually evict the pod from its node. When
          there are multiple elements, the lists of nodes corresponding to each podAffinityTerm are
          intersected, i.e. all terms must be satisfied.
        type: string # TODO
        required: false
`
}
