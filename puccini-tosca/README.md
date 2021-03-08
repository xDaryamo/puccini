puccini-tosca
=============

### Format

The default format for output is YAML, but you can select JSON, XML, or CBOR instead with
`--format/-f`. Note that Clout in JSON may lose some type information (e.g. JSON doesn't distinguish
between an integer and a float). For this reason we also support a "compatible JSON" format ("cjson")
that adds that type information. You would need specialized code to be able to consume this format.
XML output uses a bespoke structure for maps and lists, which also must be specially consumed.
(The `puccini-clout` tool supports all these formats as input.)

For YAML you can add the additional `--strict/-y` switch to output a stricter YAML, which adds
scalar type tags (such as `!!str`, `!!int`, `!!timestamp`) and outputs all strings in double quotes
with no `|` or `>` notations. This is useful if you are consuming the YAML output with a
non-compliant or buggy parser.

Another YAML-specific switch is `--timestamps`. By default Puccini will allow the YAML "!!timestamp"
type in its output. This type was included in YAML 1.1 but made optional in YAML 1.2. Most YAML
parsers should support it, but in case your YAML 1.2 parser doesn't you can disable this feature by
setting this switch to false, in which case a canonical ISO-8601 (RFC-3339) string will be output
instead.

The `--pretty` switch (enabled by default) attempts a more human-readable output, with indentation
and color highlighting in terminals. Disable this switch for a more compact output.

### TOSCA Quirks

**pucini-tosca** supports "quirks", via the `--quirk/-x` switch, which are variations on the default
grammar rules. The reason this is required is unfortunate: the low quality of the TOSCA spec,
riddled as it is with gaps, inconsistencies, and errors, means that there's too much room for
varying interpretations of the spec as well as missing functionality. Puccini aims to adhere as
closely as possible to the spec, literally and in spirit, but also must be pragmatic. Quirks allow
Puccini to smooth incompatibilities with other tools and work around a few TOSCA pain points.
Example of use:

    puccini-tosca compile weird.yaml -x data_types.string.permissive

The list of supported quirks is maintained [here](../tosca/QUIRKS.md).

### Errors and Debugging

If TOSCA compilation or parsing fails it will emit a colorful problem report to stderr and exit with
code 1. You can use the `--quiet/-q` to avoid output if all you want to do is check for success.

Logs are written to stderr (with colors) by default. Use `--log/-l` to output to a file (without
colors). Use `--verbose/-v` to add log verbosity. This can be used twice for maximum verbosity:
`-vv`.

A simple trick for if you just want to see the logs on the console: just redirect stdout to
`/dev/null` (stderr will still be present):

    puccini-tosca compile service.yaml -vv > /dev/null


`compile`
---------

The most common command is `compile`. The optional input is a filesystem path or URL to a TOSCA
service template YAML file or a CSAR file. If no input is provided will attempt to read YAML
from stdin. By default the compiled Clout will be output to stdout, but you can use the
`--output/-o` switch to specify a file (or direct to a file in the shell via `>`):

    puccini-tosca compile service.yaml --output clout.yaml

or:

    puccini-tosca compile service.yaml > clout.yaml

If your TOSCA service template YAML imports other units, and these imported paths are relative,
then Puccini will assume that they are relative to the base URL of the main service template file.
This allows you to move your files around, and even host them at a URL, without changing the file
contents. This works as expected even within a CSAR file because Puccini uses `zip:` URLs to locate
files within it.

### Inputs

Also useful (or necessary) is the `--input/-i` switch, which lets you set the inputs for the TOSCA
topology template. You may specify this switch multiple times for as many inputs as necessary.
The format is `name=value`, where `value` is YAML-encoded. For example, to set an integer:

    --input cores=4

To set a string:

    --input ram=4gb

To set a complex data type:

    --input 'port={type:tcp,number:8080}'

You can also load all inputs from YAML context located at a path or URL using the `--inputs/-n`
switch. The `--input/-i` switch is processed after, so it can be used override values in the YAML
content.

### Resolution

By default, the compiler will resolve the topology, which attempts to satisfy requirements with
capabilities and create relationships (Clout edges) between node templates. This can be disabled
via the `--resolve/-r` switch:

    puccini-tosca compile --resolve=false service.yaml

Topology resolution can be applied after compilation on an existing Clout via the embedded
**tosca.resolve** JavaScript:

    cat clout.yaml | puccini-clout scriptlet exec tosca.resolve
    
Read more about resolution [here](../tosca/compiler/RESOLUTION.md).

### Coercion

TOSCA functions and constraints are embedded in the Clout (as stubs) and are intended to be executed
when necessary. Thus they not called during compilation, unless they are needed for topology
resolution. If you want to call all of them and see the evaluated results, pipe the Clout through
**puccini-clout** and execute the embedded **tosca.coerce** scriptlet: 

    cat clout.yaml | puccini-clout scriptlet exec tosca.coerce

Values can be re-coerced later according to changing attributes and other external factors.

As a convenience, `compile` can call **tosca.coerce** for you via the `--coerce/-c` switch:

    puccini-tosca compile --coerce service.yaml

Note that doing so means you do *not* get the original compiled output.
 

`parse`
-------

If you need more diagnostics for TOSCA parsing use the `parse` command. It works similarly to
`compile` but does not emit Clout. Instead, it provides you various switches for examining the
internal workings of Puccini's TOSCA parser.

Use `--stop/-s` to specify a [phase](../tosca/parser/) (1-5) at which you wish the parser to stop.
This could be useful if you're getting too many problems in your report and wish to minimize them to
a more manageable list. Note that `-s 0` will skip the TOSCA parser entirely and just check that the
YAML input is readable.

`--dump/-d` is used to dump the internal data of phases. You may specify multiple phases to dump
using ",", e.g. `-d 2,3,4`. Per phase you will see:

* Phase 1: Read. The hierarchy of imported units starting at the service template URL.
* Phase 2: Namespaces. All names per type. Each imported unit has its own namespace.
* Phase 3: Hierarchies. Tree of all types by fully qualified name. Each imported unit has its own
  type hierarchy.
* Phase 4: Inheritance. A tree of all inheritance tasks and their dependencies by path.  
* Phase 5: Rendering. Dumps the rendered entities.
  More useful, perhaps, would be the `--filter/-t` switch (see below).

The `--filter/-t` switch can be used to filter for specific parsed entities. Each entity is given a
path that more-or-less follows JSON. For example, a path can be:

    topology_template.node_templates["store"].properties["name"]

The switch will search for all paths that contains your string, e.g. `-t properties`. You can even
include one or more "*" wildcards, e.g. `-t 'node*properties*data'`.
