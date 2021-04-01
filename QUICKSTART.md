Puccini Quickstart
==================

[Download and install Puccini](https://github.com/tliron/puccini/releases).

The distribution comes with three executables:

* `puccini-tosca`: compiles TOSCA into a Clout
* `puccini-clout`: processes a Clout, e.g. by running scriptlets on it
* `puccini-csar`: packs TOSCA sources and artifacts into a CSAR

(Note that the first two are self-contained executables, the last is a bash script.)


Basic Usage
-----------

Let's start by compiling a self-contained local file:

    puccini-tosca compile examples/tosca/descriptions.yaml

What if the file imports other files? TOSCA `imports` can refer to absolutely located
files or URLs, but they can also be relative paths. Relative paths are processed as relative
to the importing file's location, including Unix-style support for `..` to access the parent
directory.

Here's an example with relative imports:

    puccini-tosca compile examples/openstack/hello-world.yaml

Note that support for relative paths is not only for imports but also for TOSCA artifact
references.

Puccini can also compile directly from a URL:

    puccini-tosca compile https://raw.githubusercontent.com/tliron/puccini/main/examples/openstack/hello-world.yaml

In the case of a URL the relative imports are processed as if URL's "path" component were a
Unix-style path. Indeed the OpenStack example above works whether you access it locally or at
a URL.

You can even compile from stdin:

    cat examples/tosca/descriptions.yaml | puccini-tosca compile

Though note that a stdin source does not have a path and thus cannot support relative
imports.

For the above examples we referred to a single, root YAML file. However, Puccini can also
compile a CSAR package (and again, it can be a local file or at a URL). Let's create a
local CSAR and then compile it:

    puccini-csar openstack.csar examples/openstack
    puccini-tosca compile openstack.csar

The `puccini-csar` tool will zip the entire directory and automatically create a
"TOSCA-Metadata" section for us, resulting in a compliant CSAR file.

For a CSAR the relative imports refer to the internal structure of the zip archive (note
that Puccini does *not* unpack the archive into individual files, but rather treats the
entire archive as a self-contained filesystem). So, once again, the same exact OpenStack
example works whether it's accessed locally, at a URL, or from within a CSAR.


Controlling the Output
----------------------

The default output format for the Clout is YAML but other formats are supported, too
(JSON, ARD-compatible JSON, XML, and CBOR). Here's JSON:

    puccini-tosca compile examples/tosca/descriptions.yaml --format=json

By default the output is prettified and colorized for human readability when possible.
To disable that:

    puccini-tosca compile examples/tosca/descriptions.yaml --format=json --pretty=false

By default the Clout is output to stdout but you can also output to a file:

    puccini-tosca compile examples/tosca/descriptions.yaml --output=clout.yaml

Of course if running in a shell you can also redirect to a file:

    puccini-tosca compile examples/tosca/descriptions.yaml > clout.yaml

You can increase the verbosity of logging using `-v` or even `-vv`:

    puccini-tosca compile examples/tosca/descriptions.yaml -vv

By default all the log messages go to stderr, but we can send them to a file:

    puccini-tosca compile examples/tosca/descriptions.yaml -vv --log=puccini.log
    cat puccini.log

To suppress all output (if you're only interested in the return error code):

    puccini-tosca compile examples/tosca/descriptions.yaml --quiet


More on Compilation
-------------------

Let's try to compile a TOSCA service template that requires inputs (and in this case the
inputs do not have default values):

    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml

You'll see that Puccini reported a "problem" regarding the unassigned input. Any and all
compilation errors, whether they are syntactical or grammatical or topological, are
gathered and organized by file, row, and column. Indeed, Puccini's strict and detailed
problem reporting is one of its most powerful features.

By default problems are reported in a human-readable format. However, like the Clout
output, problems can be formatted for easier consumption by other tools:

    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml --problems-format=json

Let's set that missing input:

    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml --input=ram=1gib

In this case the input is a string (actually a TOSCA `scalar-unit.size`), but note that
the the input format is YAML, which is also JSON-compatible, so that complex input
values are supported, e.g.: `--input=myinput={key:value}`.

Note that you can use the `--input` flag more than once to set multiple inputs.

Inputs can also be provided from a file (locally or at a URL) as straightforward YAML:

    echo 'ram: 1 gib' > inputs.yaml
    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml --inputs=inputs.yaml

By default the topology is "resolved", meaning that all node template requirements must be
satisfied in order to create relationships and thus a complete graph. However, while working
on a work-in-progress TOSCA service template you may want to disable the resolution phase to
avoid problems:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml --resolve=false

When you turn off the resolution phase you will indeed see no relationships in the Clout
(the `edgesOut` for all vertexes will be an empty list).


TOSCA Functions and Constraints
-------------------------------

An important feature of Clout is that by default it does not call TOSCA functions.
Instead, function call stubs are inserted. This allows you to call the functions at
will using data that is not available during compilation. For example, the
`get_attribute` function relies on attribute values that would be provided by an
orchestrator or cloud platform directly from runtime resources. Indeed, these values
can change, in which case we would want to call the functions again, resulting in
new return values.

You can see the stubs in this example:

    puccini-tosca compile examples/tosca/functions.yaml

The stubs all have the special `$functionCall` key. So, how do you call these
functions? In Puccini this is called "value coercion". As a convenience, `puccini-tosca`
supports a `--coerce` flag to coerce all values right after compilation. In a real-world
orchestration scenario we would want to coerce a Clout later as its attribute values
change. We'll discuss that below. For now, let's just see how coercion looks:

    puccini-tosca compile examples/tosca/functions.yaml --coerce

You'll see that all properties now have their actual values rather than stubs.

Puccini handles TOSCA constraints in exactly the same way. The reason is that,
like functions, constraints would have to be applied to data that is not available
during compilation. For example, an attribute can itself have constraints, but
also other values that call `get_attribute` would not have an actual value until
that attribute is set.

Let's try this example:

    puccini-tosca compile examples/tosca/data-types.yaml --coerce

Now, edit `examples/tosca/data-types.yaml` and introduce a coercion problem. For
example, the `constrained_string` property requires a minimum length of 2 and a
maximum length of 5. Let's set its value to a string with length 6, `ABCDEF` (at
line 267), and compile and coerce again:

    puccini-tosca compile examples/tosca/data-types.yaml --coerce

You'll see a problem telling you exactly which constraint failed and where. Now,
let's compile this same file without coercsion (the default):

    puccini-tosca compile examples/tosca/data-types.yaml

The problem was not reported this time because the constraint stubs were not
called.

**IMPORTANT! What this means is that by default you will not see constraint-related
problems, even for values that are known during compilation! Thus it's common to use
the `--coerce` flag with `puccini-tosca compile` when your goal is to validate the
TOSCA.**


Scriptlets
----------

The Clout format is essentially a graph database in a file. However, one powerful
(and optional) feature is the ability to embed JavaScript code (scriptlets), either
as individual functions or complete programs. Indeed, this is how the TOSCA functions
and constraints are implemented.

Let's use the `puccini-clout` tool to list these scriptlets:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml --output=clout.yaml
    puccini-clout scriptlet list clout.yaml

Note that `puccini-clout` can also accept its input from stdin, allowing us to
pipe the two tools:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml | puccini-clout scriptlet list

Let's extract a scriptlet's source code:

    puccini-clout scriptlet get tosca.function.concat clout.yaml

Let's run the coercion scriptlet:

    puccini-clout scriptlet exec tosca.coerce clout.yaml

The command above is exactly equivalent to the `--coerce` flag we used previously,
which was indeed simply a shortcut for executing the `tosca.coerce` scriptlet. The
difference here, calling the scriptlet explicitly, is that we are longer are referring
to the TOSCA source. We needed to compile the TOSCA once and once, and from now on
orchestration can proceeed using Clout. As values in the Clout change (e.g. the
orchestrator udpates the attribute values with new runtime data) we can call the
`tosca.coerce` scriptlet and get the up-to-date values.

The `scriptlet exec` command can also execute scriptlets that are not embedded in
the Clout. Let's generate some HTML that visualizes the topology:

    puccini-clout scriptlet exec assets/tosca/profiles/common/1.0/js/visualize.js clout.yaml --output=tosca.html
    xdg-open tosca.html

Also note another shortcut, an `--exec` flag that allows us to execute any arbitrary
scriptlet right after compilation (skipping the Clout output):

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml --exec=assets/tosca/profiles/common/1.0/js/visualize.js


Why Scriptlets?
---------------

Strictly speaking this is not a necessary feature. For example, we could have
handled the TOSCA functions and constraints via a separate tool, which would have
scanned the Clout for those function call stubs and implemented them as necessary.
Similarly, the `visualize.js` scriptlet we used above could have been implemented
as an independent Clout processing tool.

However, such solutions come with a disadvantage: the Clout would have limited
portability, as it would need to be distributed with that tool. In a heterogeneous
cloud orchestration environment this could be burdensome.

Embedding all necessary code within the Clout makes it much more portable. Indeed,
when formatted as YAML (or JSON or XML) it is nothing more than a string, which is
eminently transmittable and storable.

Of course you still need a tool to execute those JavaScript scriptlets, but it is
the same tool (`puccini-clout`) for all Clouts, whatever version of TOSCA they come
from.

This approach works well with the cloud-native paradigm, e.g. by embedding the
`visualize.js` scriptlet into the Clout we are essentially allowing the Clout to
visualize itself.

For examples of how to create your own custom functions, constraints, and other
scriptlets, see [here](examples/javascript/).


Next Steps
----------

Now that you know how to use Puccini and more about how it works, check out all
the various included [examples](examples/).

And also check out [Turandot](https://turandot.puccini.cloud/), an orchestrator
for Kubernetes based on Puccini and Clout.
