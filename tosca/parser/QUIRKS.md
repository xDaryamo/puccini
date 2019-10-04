Puccini TOSCA Quirks
====================

These are activated via the `--quirk/-x` switch for
[**puccini-tosca**](../../puccini-tosca/):

* **data_types.string.permissive**: By default Puccini is strict about "string"-typed values
  and will consider integers, floats, and boolean values to be errors. This quirk will accept
  such values and convert them sensibly to strings. This includes handling the TOSCA "version"
  primitive type. Note that string conversions are *not* guaranteed to be identical to the YAML
  source code. For example, `1.0000` in YAML (a float) would become the string `1` in TOSCA.

* **substitution_mappings.requirements.list**: According to the examples in the spec, the
  `requirements` key under `substitution_mappings` is syntactically a map. However, this syntax is
  inconsistent because it doesn't match the syntax in node templates, which is a sequenced list.
  (In node types, too, it is a sequenced list, although grammatically it works like a map.) This
  quirk changes the accepted syntax to a sequenced list.
