Puccini TOSCA Quirks
====================

These are activated via the `--quirk/-x` switch for [**puccini-tosca**](README.md):

* **substitution_mappings.requirements.list**: According to the examples in the spec, the
  `requirements` key under `substitution_mappings` is syntactically a map. However, this syntax is
  inconsistent because it doesn't match the syntax in node templates, which is a sequenced list.
  (In node types, too, it is a sequenced list, although grammatically it works like a map.) This
  quirk changes the accepted syntax to a sequenced list.
