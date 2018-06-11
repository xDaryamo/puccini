Clout: Cloud Topology ("clou-" + "t") Representation Language
=============================================================

File Format
-----------

Clout is ["agnostic raw data"](../ard/README.md) that can be stored as YAML/JSON/XML/etc.

### `version` (string)

Must be "1.0" to conform with this document.

### `metadata` (map of string to anything)

General metadata for the whole topology. It may include information about which frontend or
processor generated the Clout file, a timestamp, etc.

### `properties` (map of string to anything)

General implementation-specific properties for the whole topology.

The difference between `metadata` and `properties` is a matter of convention. Generally,
`properties` should be used for data that is implementation-specific while `metadata` should be
used for tooling. It is understood that this distinction might not always be clear and thus you
should not treat the two areas differently in terms of state management.  

### `vertexes` (map of string to `Vertex`)

It is **very important** that you *do not treat the keys of this map as data*, for example as the
unique name of a vertex. If you need a "name" for the vertex, it should be a property within the
vertex. The vertex map keys are an internal implementation detail of Clout.

The reason for this is critical to Clout's intended use. The vertex key is used *only* as a way to
map the topology internally within an instance of Clout. More specifically, it is used for the
`targetID` field in an edge so that the topology can graphed.

But a Clout processor may very well transform a Clout file and modify the topology. This could
involve adding new vertexes and edges or moving them around, for example to optimize a topology,
to heal a broken implementation, to scale out an overloaded system, etc. In doing so it may
regenerate these IDs. These IDs need only be unique to one specific Clout file, not generally.

If you do need to lookup a vertex by, say, its `name` property, then the correct way to do so is
to iterate through all vertexes and look for the first vertex that has that particular name.
Indeed, it is reasonable for Clout parsers to entirely hide these IDs from the user and perhaps
represent the vertex map as a list.

`Vertex`
--------

### `metadata` (map of string to anything)

Often you'll find information here about what kind of vertex this is, e.g. a TOSCA node.

### `properties` (map of string to anything)

Implementation-specific properties for the vertex.

### `edgesOut` (list of `Edge`)

Clout edges are directional, though you may choose to semantically ignore the direction. The edges
are stored in the *source* vertex, which is why this field is named `edgesOut`.

As a convenience, Clout parsers may very well add an in-memory `edgesIn` field, which would also be
a list of edges, after mapping the `targetID` fields of all edges to vertexes, or otherwise provide
a tool for looking up edges for which a certain vertex is a target.

`Edge`
------

### `metadata` (map of string to anything)

Often you'll find information here about what kind of edge this is, e.g. a TOSCA relationship.

### `properties` (map of string to anything)

Implementation-specific properties for the vertex.

### `targetID` (string)

The key in the vertexes map to which this edge is the target.

Note that there is no need for a `sourceID` because the edge is already located in the `edgesOut`
field of its source vertex. Clout parsers may very well add such a field for convenience.

Better yet, Clout parsers may do the ID lookup internally, provide direct access to the source
and target vertexes, and hide the `targetID` field.
