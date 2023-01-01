Puccini Tutorial
================

[Download and install Puccini](https://github.com/tliron/puccini/releases).

The distribution comes with three executables:

* [`puccini-tosca`](puccini-tosca/): compiles TOSCA into a Clout
* [`puccini-clout`](puccini-clout/): processes a Clout, e.g. by running scriptlets on it
* [`puccini-csar`](puccini-csar/): packs TOSCA sources and artifacts into a CSAR


Basic Usage
-----------

Let's start by compiling a self-contained local file:

    puccini-tosca compile examples/tosca/descriptions.yaml

What if the file imports other files? TOSCA `imports` can refer to either absolute
URLs or relative URLs ([RFC 1808](https://tools.ietf.org/html/rfc1808)). Note that
if the URL scheme is not provided it defaults to "file:", so that it can be treated
as a platform-independent file system path. Also note that relative URLs support
Unix-like `.` and `..` components, allowing you to refer to resources upwards in the
directory tree.

Let's compile a local example for OpenStack that uses imports (relative URLs):

    puccini-tosca compile examples/openstack/hello-world.yaml

Note that you can also use relative URLs in TOSCA artifacts (their `file` keyword).

Puccini can also compile directly from a URL. Let's use the same OpenStack example as
above:

    puccini-tosca compile https://raw.githubusercontent.com/tliron/puccini/main/examples/openstack/hello-world.yaml

You'll see that the relative URLs continue to work as expected even though the base
URL is not on the local filesystem.

The URL system is quite powerful and even supports access to git repositories (GitOps!).
Any valid git repository URL can follow the "git:" prefix, then follow with a "!" and the
path within the repository. Also note that in bash you need to escape the "!" character
or wrap it in single quotes:

    puccini-tosca compile 'git:https://github.com/tliron/puccini.git!examples/openstack/hello-world.yaml'

Puccini can also compile YAML from stdin:

    cat examples/tosca/descriptions.yaml | puccini-tosca compile

Be aware that a stdin source does not have a path and thus cannot support relative
URLs.

For the above examples we referred to a single, root YAML file. However, Puccini can also
compile from a CSAR package and, again, the CSAR can be a local file or at a URL. Let's create
a local CSAR and then compile it:

    puccini-csar create openstack.tar.gz examples/openstack
    puccini-tosca compile openstack.tar.gz

The `puccini-csar` tool will archive the entire directory and automatically create a
"TOSCA-Metadata" section for us, resulting in a compliant CSAR file.

For TOSCA files within the CSAR the relative URLs refer to the internal structure of the
archive. Note that Puccini does *not* unpack the archive into individual files, but rather
treats the entire archive as a self-contained filesystem. Thus, because we used relative URLs
in our TOSCA imports, the same exact OpenStack example works whether it's accessed locally,
at a URL, or from within a CSAR, and even a CSAR at a URL.


Controlling the Output
----------------------

The default output format is YAML but other formats are supported: JSON (and
[ARD](https://github.com/tliron/kutil/tree/master/ard/)-compatible JSON), XML, CBOR,
and MessagePack. Here's ARD-compatible JSON:

    puccini-tosca compile examples/tosca/descriptions.yaml --format=cjson

By default the output is nicely indented and and colorized for human readability. You can
turn off prettification if you're interested in the most compact output:

    puccini-tosca compile examples/tosca/descriptions.yaml --pretty=false

Note that colorization will *always* be disabled in contexts that do not support it. In
other words it will likely only appear in stdout for terminal emulators that support ANSI
color codes. However, you can also specifically turn off colorization:

    puccini-tosca compile examples/tosca/descriptions.yaml --colorize=false

By default the output is sent to stdout but you can also send it to a file (without
colorization):

    puccini-tosca compile examples/tosca/descriptions.yaml --output=clout.yaml

Of course if running in a shell you can also redirect stdout to a file (again, without
colorization):

    puccini-tosca compile examples/tosca/descriptions.yaml > clout.yaml

You can increase the verbosity of logging using `-v` or even `-vv`:

    puccini-tosca compile examples/tosca/descriptions.yaml -vv

By default all the log messages go to stderr but we can send them to a file:

    puccini-tosca compile examples/tosca/descriptions.yaml -vv --log=puccini.log
    cat puccini.log

If you only want to see the logs and not the Clout output:

    puccini-tosca compile examples/tosca/descriptions.yaml -vv > /dev/null

To suppress all output (if you're only interested in the return error code):

    puccini-tosca compile examples/tosca/descriptions.yaml --quiet

Also note that there is a `puccini-tosca parse` command that provides a lot
of internal diagnostic information about the language parser. It's generally
useful for Puccini developers rather than Puccini users, so it is out of scope
for this quickstart guide. See [here](puccini-tosca/) for more information.


More on Compilation
-------------------

Let's try to compile a TOSCA service template that requires inputs:

    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml

You'll see that Puccini reported a "problem" regarding the unassigned input. Any and all
compilation errors, whether they are syntactical, grammatical, or topological, are
gathered and organized by file, row, and column. Indeed, Puccini's strict and detailed
problem reporting is one of its most powerful features.

By default problems are reported in a human-readable format. However, like the Clout
output, problems can be formatted for easier consumption by other tools:

    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml --problems-format=json

Let's set that missing input:

    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml --input=ram=1gib

In this case the input is a string (actually a TOSCA `scalar-unit.size`), but note that
the the input format is YAML, which is also JSON-compatible, so that complex input
values can be provided, e.g. `--input=myinput={key1:value1,key2:value2}`. Also Note that
you can use the `--input` flag more than once to provide multiple inputs.

Inputs can also be loaded from a file (locally or at a URL) as straightforward YAML:

    echo 'ram: 1 gib' > inputs.yaml
    puccini-tosca compile examples/tosca/inputs-and-outputs.yaml --inputs=inputs.yaml

By default the compiler will "resolve" the topology, meaning that it will atempt to satisfy
all node template requirements and create relationships, thus completing the graph. However,
sometimes it may be useful to disable the resolution phase in order to avoid excessive problem
reports:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml --resolve=false

When you turn off the resolution phase you will indeed see no relationships in the Clout
(you'll see that the `edgesOut` for all vertexes is an empty list).

Read more about how Puccini implements resolution [here](assets/tosca/profiles/common/1.0/js/RESOLUTION.md).


TOSCA Functions and Constraints
-------------------------------

An important feature of `puccini-tosca compile` is that by default it does not call
TOSCA functions. Instead, function call stubs are inserted. This allows you to call
the functions at will (and repeatedly) using up-to-date data, including data that might
not even available during compilation. For example, the `get_attribute` function relies
on attribute values that would be provided by an orchestrator or cloud platform directly
from runtime resources.

You can see the call stubs by compiling this example:

    puccini-tosca compile examples/tosca/functions.yaml

You'll notice that the call stubs all have the special `$functionCall` key.

How do we call the functions? In Puccini we refer to this as "value coercion". As a
convenience we can use the `--coerce` flag to coerce the values during compilation:

    puccini-tosca compile examples/tosca/functions.yaml --coerce

You'll see that all properties now have their actual values rather than call stubs.

(In a real-world orchestration scenario we would want to coerce a Clout later as its
attribute values change. We'll discuss that below.)

Puccini handles TOSCA constraints in exactly the same way, the reason being that,
like functions, constraints would have to be applied to data that might not be
available during compilation, e.g. constraints associated with an attribute or its
data type.

Let's try this example:

    puccini-tosca compile examples/tosca/data-types.yaml --coerce

Now, edit `examples/tosca/data-types.yaml` and break a constraint. For example, the
`constrained_string` property requires a minimum length of 2 and a maximum length of
5, so let's set its value to a string with length 6, `ABCDEF` (at line 267), and
compile and coerce again:

    puccini-tosca compile examples/tosca/data-types.yaml --coerce

You'll see a problem reported telling you exactly which constraint failed and where.
Now, let's compile this same file without coercion (the default behavior):

    puccini-tosca compile examples/tosca/data-types.yaml

The problem was not reported this time.

**IMPORTANT! The implication is that by default you will not see constraint-related
problems reported during compilation, even for values that are known! Thus it's common
to use the `--coerce` flag with `puccini-tosca compile` when your goal is to validate
the TOSCA.**


Scriptlets
----------

The Clout format is essentially a graph database in a file. However, one powerful
(and optional) feature is the ability to embed JavaScript code scriptlets, either
as individual functions or as complete programs. Indeed, this is how the TOSCA
function and constraint call stubs are implemented.

Let's use the `puccini-clout` tool to list these embedded scriptlets:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml --output=clout.yaml
    puccini-clout scriptlet list clout.yaml

Note that `puccini-clout` can also accept Clout input from stdin, allowing us to pipe
the two tools:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml | puccini-clout scriptlet list

Let's extract a scriptlet's source code:

    puccini-clout scriptlet get tosca.function.concat clout.yaml

Finally, let's run the coercion scriptlet:

    puccini-clout scriptlet exec tosca.coerce clout.yaml

The command above is exactly equivalent to the `--coerce` flag we used previously
with `puccini-tosca compile`. Indeed, `--coerce` is simply a shortcut for executing
the `tosca.coerce` scriptlet. The difference is that here we are working with the
un-coerced Clout rather than the TOSCA source.

This difference is crucial to Day 2 orchestration processing. The TOSCA source is
compiled once and only once, and from then on the Clout lives on its own and can
indeed be modified by an orchestrator. At the minimum the orchestrator should fill
in the attribute values (and execute `tosca.coerce` to ensure that they match the
constraints).

The `puccini-clout scriptlet exec` command can also execute scriptlets that are not
embedded in the Clout. Let's use a scriptlet that creates an HTML page that visualizes
the topology:

    puccini-clout scriptlet exec assets/tosca/profiles/common/1.0/js/visualize.js clout.yaml --output=tosca.html
    xdg-open tosca.html

Note another shortcut for `puccini-tosca compile`: you can use the `--exec` flag to
execute scriptlets right after compilation, thus skipping the Clout intermediary:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml --exec=assets/tosca/profiles/common/1.0/js/visualize.js

See [here](puccini-clout/) for more information about the `puccini-clout` tool.


Why Scriptlets?
---------------

Strictly speaking scriptlets are not a necessary feature. For example, we could
have handled the TOSCA functions and constraints via a separate tool, which would
have contained implementations for all those those function call stubs. Similarly,
the `visualize.js` scriptlet we used above could have been implemented as an
independent Clout processing tool.

However, such external solutions come with a disadvantage: the Clout would have
limited portability because it would need to be distributed with that tool in order
to be fully functional. In a heterogeneous cloud orchestration environment this could
be burdensome. It's a security risk, too, as every additional binary increases the
attack surface.

Embedding the required code within the Clout makes it more portable and secure.
Indeed, when formatted as YAML (or JSON or XML) the entire self-contained Clout is
nothing more than a string, which is eminently transmittable and storable.

Of course you still need a tool to execute those JavaScript scriptlets, but it is
the same tool, `puccini-clout`, for all Clouts and all scriptlets. Indeed, the Clout
can contain custom vertexes, edges, and scriptlets, including those that did not
originate in a TOSCA service template. They do not even have to adhere to the TOSCA
structure.

For examples of how to create your own custom functions, constraints, and other
scriptlets for TOSCA, see [here](examples/javascript/).


More on CSARs
-------------

The [`tosca-csar`](puccini-csar/) tool supports tarball CSARs (`.tar.gz` or `.tar`) as well
as zip (`.zip` or the `.csar` alias). Note that tarballs have the advantage that they can be
streamed (e.g. from a HTTP URL) whereas using the zip format would require `puccini-tosca`
to first download the entire archive to the system's temporary directory. Both will work,
but tarballs are far more efficient.

Try zip via the `.csar` alias:

    puccini-csar create openstack.csar examples/openstack
    puccini-tosca compile openstack.csar

The CSAR format supports "Other-Definitions" metadata to specify additional service
templates beyond than the root "Entry-Definitions". To build the example:

    puccini-csar create cloud.tar.gz examples/csar \
        --entry=definitions=main.yaml \
        --other-definitions='other 1.yaml' \
        --other-definitions='other 2.yaml'

When compiling you can use the `--template` flag to select a non-default service
template. The flag can accept a valid "Other-Definitions" path, like so:

    puccini-tosca compile cloud.tar.gz --template="other 1.yaml"

It can also accept a number, where "0" would be the default service template, "1"
would be the first "Other-Definitions", and so on:

    puccini-tosca compile cloud.tar.gz --template=2

Puccini also supports "tar:" and "zip:" prefix schemes for URLs, allowing you to refer
to an entry within any archive file, CSAR or otherwise. Any valid URL can follow the
prefix, whether it's a local file URL, HTTP, etc. Note that for files it does require
absolute file system paths. Then follow with a "!" and the path within the archive.
Example:

    puccini-tosca compile "tar:$PWD/cloud.tar.gz\!main.yaml"

(We are using a backslash to escape the "!" for bash.)

Also useful is the `meta` command to validate and extract the CSAR metadata:

    puccini-csar meta cloud.tar.gz


Next Steps
----------

Now that you know how to use Puccini and understand how it works, check out all the
various included [examples](examples/).

And also check out [Turandot](https://turandot.puccini.cloud/), an orchestrator
for Kubernetes based on Puccini and Clout.
