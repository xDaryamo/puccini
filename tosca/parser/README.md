TOSCA Parser
============

Optimized for speed via caching and concurrency. Parsing even very big and complex service templates
is practically instantaneous, delayed at worst only by network and filesystem transfer.

Attempts to be as strictly compliant as possible, which is often challenging due to contradictions
and unclarity in the TOSCA specification. Where ambiguous we adhere to the spirit of the spec,
especially in regards to object-oriented polymorphism. Our prime directive is to ensure that an
inherited node type does not break the contract of the base node type.


Coding Principles
-----------------

There are over 60 different entity types in TOSCA 1.3. Writing custom code for each, even with
reusable utility functions, would quickly become a maintenance nightmare, not only for fixing bugs
but also for supporting future versions of TOSCA.

We have opted to combine utility functions with plenty of reflection. Reflection is used to
solve generic, repeatable actions, such as reading data of various types, looking up names in
the namespace, and inheriting fields from parent types. Generic code of this sort is tricky to get
right, but once you do the entire domain is managed via simple annotations (field tags).

There is a cost to using such annotations, and indeed even the simple field tags are controversial
within the Go community. The problem is that hidden, magical things happen that are not visible
in the immediate vicinity of the code in front of you. You have to look elsewhere for the systems
that read these tags and do things with them. The only way to reduce this cost is good
documentation: make coders aware of these systems and what they do. We will consider this
documentation to be a crucial component of the codebase.

* [Field Tags](TAGS.md)
* [Interfaces and Signatures](INTERFACES.md)


Phase 1: Read
-------------

This first phase validates syntax. Higher level grammar is validated in subsequent phases.

1. Read textual data from files and URLs
    * Files/URLs not found
    * I/O errors
    * Textual decoding errors
2. Parse YAML to [ARD](https://github.com/tliron/kutil/tree/master/ard/)
    * YAML parsing errors
3. Parse ARD to TOSCA data structures, normalizing all short notations to their full form
    * Required fields not set
    * Fields set to wrong YAML type
    * Unsupported fields used
4. Handle TOSCA imports recursively and concurrently
    * Import causes an endless loop


Phase 2: Namespaces
-------------------

The goal of this phase is to validate names (ensure that they are not ambiguous), and to provide
a mechanism for looking up names. We take into account the import `namespace_prefix` and support
multiple names per entity, as is necessary for the normative types.

Recursively, starting at tips of the import hierarchy:

1. Gather all names in the unit
2. Apply the import's `namespace_prefix` if defined
3. Set names in unit's namespace (every entity type has its own section)
4. Merge namespace into parent unit's
    * Ambiguous names (per entity type section)

And then:

5. Lookup fields from namespace
    * Name not found


Phase 3: Hierarchies
--------------------

The entire goal of this phase is validation. The entities already have hierarchical information
(their `Parent` field). But here we provide a hierarchy that is guaranteed to not have loops.

Recursively, starting at tips of the import hierarchy:

1. Gather all TOSCA types in the unit
2. Place types in hierarchy (every type has its own hierarchy)
    * Type's parent causes an endless loop
    * Type's parent is incomplete
3. Merge hierarchy into parent unit's


Phase 4: Inheritance
--------------------

From this phase we onward only deal with the root unit, not the imported units, because it already
has all names and types merged in.

This phase is complex because the order of inheritance cannot be determined generally. Not only do
types inherit from each other, but also definitions within the types are themselves typed, and those
types have their own hierarchy. Thus, a type cannot be resolved before its definitions are, and they
cannot be resolved before their types are, and their types' parents, etc.

This includes validating TOSCA's complex inheritance contract, which extends to embedded definitions
beyond simple type inheritance. For example, a capability definition within a node type must have a
capability type that is compatible with that defined at the parent node type.

Our solution to this complexity is to create a task graph. We will keep resolving independent tasks
(those that do not have dependencies), which should in turn make more tasks independent, continuing
until all tasks are resolved.

Note that in this phase we handle not only recursive type inheritance, but non-recursive
definition inheritance. For example, an interface definition inherits data from the interface type.

1. Copy over inherited fields from parent or definition type
    * Type's parent is incomplete (due to other problems reported in this phase)
2. If we are overriding, make sure that we are not breaking the contract
    * Between value of field in parent and child

Incredibly, the TOSCA spec does not describe inheritance. It is non-trivial and not obvious. In
Puccini we had to make our own implementation decisions. Unfortunately, other TOSCA parsers may
very well handle inheritance differently, leading to incompatibilities. The fault lies entirely
with the TOSCA spec.

For example, consider a capability definition within a node type. We inherit first from our parent
node type, and only then our capability definition will inherit from its capability type (which may
be a subtype of what our parent node type has). The first entity we come across in our traversal is
considered to be the final value (for subequently found entities we will treat the field as "already
assigned"). The bottom line is that the parent node type takes precedence over the capability type
of the capability definition.

At the end of this phase the types are considered "complete" in that we should not need to access
their parents for any data. All fields have been inherited.


Phase 5: Rendering
------------------

We call the act of applying a type to a template "rendering". ("Instantiation" is what happens next,
when we turn a template into an instance, which is out of the scope of Puccini and indeed out of the
scope of TOSCA.)

For entities that are "assignments" ("property assignments", "capability assignments", etc.) we
do more than just validate: we also change the value, for example by applying defaults.  

There are furthermore various custom validations per entity. One worth mentioning here is
requirement validation. Specifically, we take into account the `occurrences` field, which limits the
number of times a requirement may be assigned. The TOSCA spec mentions that the implied default
`occurrences` is \[1,1\] (the TOSCA spec oddly says that in a `range` the upper bound must be
greater than the lower bound, so that \[1,1] is impossible: one of many contradictions we must
resolve). However, we further assume that if `occurrences` is not specified then it is also intended
for the requirement to be automatically assigned if not explicitly specified.

At the end of this phase the templates are considered "complete" in that we should not have to
access their types for any data. All fields have been rendered.
