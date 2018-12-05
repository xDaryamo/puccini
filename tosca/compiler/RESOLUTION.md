Topology Resolution
===================

This is where we create the flat topology: relationships from templates to capabilities (the
"sockets", if you will) in other node templates. We call this "resolving" the topology.

Resolution is handled via the **tosca.resolve** JavaScript embedded in the Clout. This allows you
to re-resolve an existing compiled Clout according to varying factors.

For capabilities we take into account the `occurrences` field, which limits the number of times a
capability may be be used for relationships.

There's no elaboration in the TOSCA specification on what `occurrences` means. Our interpretation is
that it does *not* relate to the capacity of our actual resources. While it may be possible for an
orchestrator to provision an extra node to allow for more capacity, that would also change the
topology by creating additional relationships, and generally it would be an overly simplistic
strategy for scaling. TOSCA's role, and thus Puccini's, should merely be to validate the design.
Thus requirements-and-capabilities should have nothing to do with resource provisioning.

Relatedly, we also allow for relationship loops: for example, two node templates can have
`DependsOn` relationships with each other. This doesn't necessarily imply a problem: they could, for
example, be provisioned simultaneously. Whether or not orchestrators can deal with such loops is
beyond the scope of Puccini and TOSCA.
