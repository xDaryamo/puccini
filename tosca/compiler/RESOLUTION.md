Topology Resolution
===================

This is where we create the flat topology: relationships from templates to capabilities (the
"sockets", if you will) in other node templates. We call this "resolving" the topology.

Resolution is handled via the **tosca.resolve** JavaScript embedded in the Clout. This allows you
to re-resolve an existing compiled Clout according to varying factors.


Loops
-----

Puccini allows for relationship loops: for example, two node templates can have `DependsOn`
relationships with each other. Semantically, this is not a problem. However, practically it could
be a problem for some orchestrators.

A good orchestrator should know what to do. For example, mutually dependent resources could be
provisioned simultaneously. Whether or not orchestrators can deal with such loops is beyond the
scope of Puccini and TOSCA.


Capability Occurrences
----------------------

For capabilities we take into account the `occurrences` field, which limits the number of times a
capability may be be used for relationships.

There's no elaboration in the TOSCA specification on what `occurrences` means. Our interpretation is
that it does *not* relate to the capacity of our actual resources. While it may be possible for an
orchestrator to provision an extra node to allow for more capacity, that would also change the
topology by creating additional relationships, and generally it would be an overly simplistic
strategy for scaling. TOSCA's role, and thus Puccini's, should merely be to validate the design.
Thus requirements-and-capabilities should have nothing to do with resource provisioning.

That said, `occurrences` does introduce a subtle restriction on how requirements are satisfied. It
means that some capabilities might have a *minimum* number of relationships attached to them, or
else a problem is reported. Likewise, some capabilities will allow for a *maximum* number of
possible relationships. Allocating these restricted slots in such a way that all requirements can be
satisfied while keeping all minimums fulfilled would require a non-trivial algorithm.

Currently, Puccini's algorithm is not ideal. It does report problems, and it does try to prioritize
some capabilities with `occurrences` restrictions over others. However, it still iterates
requirements one at a time, meaning that it may very well miss on finding a problem-free topology.
It at least guarantees that the results will be consistent by ensuring a reproducible order of
iteration via alphanumeric sorting.

A better algorithm would require either 1) trying various sort orders until one succeeds, or 2)
finding a more sophisticated way to prioritize certain pairs of requirements-and-capabilities.
Both approaches are difficult. 
