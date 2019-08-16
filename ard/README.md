Agnostic Raw Data (ARD)
=======================

What is "agnostic raw data"?

### Agnostic

Comprising primitives (string, integer, float, boolean, null, etc.) and structures (map with string
keys, list) that can be transmitted to practically any language or platform. It can also work with
a wide variety of formats, though with some limitations:

* YAML supports a rich set of primitive types, so ARD will survive a round trip to YAML. Indeed, the
  ARD type system should adhere to YAML's. However, note that YAML maps are ordered while ARD maps
  have arbitrary order. A round trip from YAML to ARD would thus lose order. 
* JSON can be read into ARD. However, because JSON has fewer types (no integers, only floats), ARD
  can be translated to JSON but some type information would be lost unless it were to be encoded
  within the data. This would then become an extended JSON format that would also have to be parsed
  and generated in a particular way.
* XML is more complicated: with a proper schema, ARD can survive a round trip. However, XML would
  have to be created specifically for that schema. Arbitrary XML cannot be parsed into ARD.

### Raw

The data is untreated and not validated. There's no schema.
