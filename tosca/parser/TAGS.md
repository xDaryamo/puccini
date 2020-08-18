Field Tags
==========

### name:"name"

Gives the whole type a human-readable name for reporting problems. Go doesn't support tags for the
type as a whole, so just put it on any field, even an anonymous field could make sense.

### read:"key", read:"key,reader"

Read phase. The key is used to extract data from the ARD.

The type of the data can usually be discovered via reflection on the field. However, non-primitive
types require a registered "reader" to be specified via the long form.

Note that we use pointers to signify the semantic different between unassigned fields and fields
with zero values. That's why a field would be a "*string" rather than a "string".

Set the key to "?" to collect all keys not used by other fields in the type. (E.g., this is how
a TOSCA interface would read all the operations.)

By default "read" treats *both* maps and list elements as unique, according to their mappable key
(retrieved by the Mappable interface). Sometimes we prefer to store the values in an array and
sometimes in a map, but in both cases the elements must be unique.

However, it is possible to override this behavior in two ways:

Specify "[]reader" for a list of non-unique elements. (E.g. TOSCA constraint clauses.)

Specify "{}reader" for sequenced list of non-unique elements. (E.g. TOSCA node\_type.requirements
and node\_template.requirements.)

Specify "<>reader" for sequenced list of unique elements. (E.g. TOSCA topology\_template.policies.)

Also, you can specify "!reader" to mark the field as important, meaning that it will be read before
any other fields. You can combine "!" with "[]", "{}", and "<>".

(Note: node\_type.requirements, as opposed to node\_template.requirements, are really "fake"
sequenced lists, because actually you cannot repeat the same definition name more than once.
The reason the TOSCA spec has this inconsistency is likely to make the syntax more like the
syntax in node\_template.)

### require:"key"

Read phase. Reports a problem if the key is not in the ARD. (Works in combination with "read" tag.)

The "key" value is optional. If not provided it will be read from the "read" tag.

### namespace:""

Namespace phase. By using this tag you are specifying that this entity should be registered on the
namespace. Uses the value of this field (a string) for the name.

### lookup:"key,field"

Lookup phase. Can lookup many names (maps, slices). The named "field" should match the type (string,
map of string, slice of strings).

The type of the data is discovered via reflection on the field. The key is used just for reporting
problems. It is possible for more than one field, of different types, to lookup from the same named
field: if that's the case, all will be processed together, and problems will be reported only if
a name fails for all of them. (E.g. TOSCA requirement.node, which can be either a node template
name or a node type name, so the lookup tag will be in two fields.)

Normally will report errors when lookups fail. To disable this behavior (making the lookup
optional) use "?field". (E.g. TOSCA requirement.capability.)

### hierarchy:""

Hierarchy phase. Marks this field as a container for types.

After the hierarchy is built all entities of this type (determined by reflection) will be merged
here. This allows us to import types from other units.

### inherit:"key,field"

Inheritance phase. The key is used just for reporting problems. The named "field" is a reference
field (pointer) to another entity. It is expected that our field will be matched by name and type
at the referred entity.

We use this both for derived type inheritance (field=Parent) and for inheriting from definitions.

### traverse:"ignore"

Tells the "traverse" mechanism not to recurse into this field.

This is not merely an optimization: in some cases we use this to avoid traversal loops.

