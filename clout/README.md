Clout: Cloud Topology ("clou-" + "t") Representation Language
=============================================================

File Format
-----------

Clout is ["agnostic raw data"](https://github.com/tliron/kutil/tree/master/ard/) that can be stored
as YAML/JSON/XML/etc.

### `version` (string)

Must be "1.0" to conform with this document.

### `metadata` (map of string to anything)

General metadata for the whole topology. It may include information about which frontend or
processor generated the Clout file, a timestamp, etc.

### `properties` (map of string to anything)

General implementation-specific properties for the whole topology.

The difference between `metadata` and `properties` is a matter of convention. Generally, `properties`
should be used for data that is implementation-specific while `metadata` should be used for tooling.
It is understood that this distinction might not always be clear and thus you should not treat the
two areas differently in terms of state management.  

### `vertexes` (map of string to Vertex)

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

If you do need to lookup a vertex by, say, its `name` property, then the correct way to do so is to
iterate through all vertexes and look for the first vertex that has that particular name. Indeed, it
is reasonable for Clout parsers to entirely hide these IDs from the user and perhaps represent the
vertex map as a list.


Vertex
------

### `metadata` (map of string to anything)

Often you'll find information here about what kind of vertex this is, e.g. a TOSCA node.

### `properties` (map of string to anything)

Implementation-specific properties for the vertex.

### `edgesOut` (list of Edge)

Clout edges are directional, though you may choose to semantically ignore the direction. The edges
are stored in the *source* vertex, which is why this field is named `edgesOut`.

As a convenience, Clout parsers may very well add an in-memory `edgesIn` field, which would also be
a list of edges, after mapping the `targetID` fields of all edges to vertexes, or otherwise provide
a tool for looking up edges for which a certain vertex is a target.


Edge
----

### `metadata` (map of string to anything)

Often you'll find information here about what kind of edge this is, e.g. a TOSCA relationship.

### `properties` (map of string to anything)

Implementation-specific properties for the vertex.

### `targetID` (string)

The key in the vertexes map to which this edge is the target.

Note that there is no need for a `sourceID` because the edge is already located in the `edgesOut`
field of its source vertex. Clout parsers may very well add such a field for convenience.

Better yet, Clout parsers may do the ID lookup internally, provide direct access to the source and
target vertexes, and hide the `targetID` field.


Coercibles
----------

A common feature in many Clout use cases is the inclusion of values that are meant to be "coerced"
at runtime. Coercion could include evaluating an expression, calling a function, testing for
validity of the value by applying constraints, etc.

Clout does not enforce a notation for such coercible values, however we do suggest a convention.
Puccini comes with tools to help you parse according to this notation and to perform coercion using
JavaScript.

The convention is recursive and assumes that each value is a map with one *and only one* of
the following fields:

* `$value`: this is an ARD literal value (boolean, integer, float, string, list, map, etc.)
* `$list`: this is a list of coercibles
* `$map`: this is a *list* of coercibles (not a map!) whereby each entry also includes a `$key` field,
   which itself is also a coercible (a list is used for better compatibility with systems that do
   not support maps with arbitrary keys; note that because this is a list duplicate keys are
   syntactically possible here, but that should be considered a semantic error)
* `$functionCall`: this is a function call notation (see below)

All coercibles may also have the following optional fields:

* `$constraints`: a list of coercibles in the `$functionCall` format (see notation below)
* `$information`: a map with descriptive data (see below)

If the value of `$value` is itself a map, that map may have following optional fields:

* `$string`: textual representation of the value for human-readability, comparison, sorting, etc.
* `$number`: numeric representation of the value (float or integer) for comparison, sorting, etc.
* `$originalString`: if the value was parsed from a string then this would be that string
* `$comparer`: name of scriptlet to be used for value comparisons

`$list` coercibles may also have the following optional fields:

* `$entryConstraints`: a list of coercibles in the `$functionCall` format (see notation below) intended
  to be applied to each entry in the list

`$map` coercibles may also have the following optional fields:
 
* `$keyConstraints`: a list of coercibles in the `$functionCall` format (see notation below) intended
  to be applied to each key in the map
* `$valueConstraints`: a list of coercibles in the `$functionCall` format (see notation below) intended
  to be applied to each value in the map

The value of `$functionCall` is a map with the following required fields:

* `name`: a string representing the name of the function
* `arguments`: a list of coercibles (can be an empty list but not null)

Additionally, it may have the follow optional fields (for debugging information):

* `path`: a string representing a semantic path within the source document (implementation-specific)
* `url`: a string representing the URL of the source document
* `row`: an integer representing the row within the source document
* `column`: an integer representing the column within the source document

The `$information` is a map with the following optional fields:

* `description`: a human-readable description of the coercible
* `definition`: type information map (see below) for the coercible's definition
* `type`: type information map (see below) for the coercible's type
* `entry`: type information map (see below) for a `$list` coercible's entries
* `key`: type information map (see below) for a `$map` coercible's keys
* `value`: type information map (see below) for a `$map` coercible's values
* `properties`: a map of `$map` coercible keys (strings) to an `$information` structure (recursive)

The "type information" mentioned above is a map with the following optional fields:

* `name`: the canonical name of the type
* `description`: a human-readable description of the type
* `schemaDescription`: a human-readable description of the schema (for `entry`, `key`, and `value`)
* `metadata`: a map of metadata strings

Example (generated from [this TOSCA example](../examples/tosca/data-types.yaml)):

    lowercase_string_map:
      $map:
        - $key:
            $functionCall:
              name: tosca.function.concat
              arguments:
                - $value: recip
                - $value: ient
              path: topology_template.node_templates["data"].properties["lowercase_string_map"]["concat:¶  - recip¶  - ient"]
              url: file:examples/tosca/data-types.yaml
              row: 180
              column: 9
          $value: Puccini
        - $key:
            $value: greeting
          $value: Hello
      $keyConstraints:
        - $functionCall:
            name: tosca.constraint.pattern
            arguments:
              - $value: '[a-z]*'
            path: topology_template.node_templates["data"].properties["lowercase_string_map"]
            url: file:examples/tosca/data-types.yaml
            row: 180
            column: 9
      $information:
        type:
          name: map
        key:
          name: LowerCase
          description: Lowercase string
        value:
          name: string
