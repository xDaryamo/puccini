Puccini TOSCA Quirks
====================

These are activated via the `--quirk/-x` switch for
[**puccini-tosca**](../../puccini-tosca/):

* **imports.implicit.disable**: In TOSCA 1.0-1.3 the Simple Profile is implicitly imported by
  default. This quirk will disable implicit imports.

* **imports.permissive**: By default Puccini will report an error if a unit imports another
  unit with an incompatible grammar. This quirk will disable the check.

* **data_types.string.permissive**: By default Puccini is strict about "string"-typed values
  and will consider integers, floats, and boolean values to be problems. This quirk will accept
  such values and convert them as sensibly as possible to strings. This includes accepting floats
  and integers for the TOSCA "version" primitive type. Note that string conversions may very well
  *not* be identical to the literal YAML. For example, `1.0000` in YAML (a float) would become
  the string `1` in TOSCA.

* **data_types.timestamp.permissive**: By default Puccini requires all "timestamp" values to be
  specified as strings in the ISO 8601 format. However, some YAML environments may support the
  optional !!timestamp type. This quirk will allow such values. Note that such values will not have
  the "$originalString" key, because the literal YAML is not preserved by the YAML parser.

* **namespace.normative.ignore**: This will ignore any type that is has the
  "puccini.normative: true" metadata.

* **namespace.normative.shortcuts.disable**: In TOSCA 1.0-1.3 all the normative types have long
  names, such as "tosca.nodes.Compute", prefixed names ("tosca:Compute"), and also short names
  ("Compute"). Those short names are annoying because it means you can't use those names for your
  own types. This quirk disables the short names (the prefixed names remain).

* **substitution_mappings.requirements.list**: According to the examples in the TOSCA 1.0-1.3 specs,
  the `requirements` key under `substitution_mappings` is syntactically a map. However, this syntax
  is inconsistent because it doesn't match the syntax in node templates, which is a sequenced list.
  (In node types, too, it is a sequenced list, although grammatically it works like a map.) This
  quirk changes the expected syntax to be a sequenced list.
