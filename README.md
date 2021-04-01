Puccini
=======

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/puccini.svg)](https://github.com/tliron/puccini/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/puccini)](https://goreportcard.com/report/github.com/tliron/puccini)

Deliberately stateless cloud topology management and deployment tools based on
[TOSCA](https://www.oasis-open.org/committees/tosca/).

Want to dive in?

Head to the [quickstart guide](QUICKSTART.md).

Also check out this [live demo of Puccini TOSCA running in a browser](https://web.puccini.cloud/).

Note that Puccini is intentionally *not* an orchestrator. This is a "BYOO" kind of establishment
("Bring Your Own Orchestrator").

If you are looking for a comprehensive TOSCA orchestrator for
[Kubernetes](https://kubernetes.io/) based on Puccini, check out
[Turandot](https://turandot.puccini.cloud/). Puccini also
[enables TOSCA for Ansible](examples/ansible/) using custom extensions.

Also included are examples of [generating Ansible playbooks for OpenStack](examples/openstack/)
as well as [generating BPMN processes](examples/bpmn/) for middleware integration.


Get It
------

[![Download](assets/media/download.png "Download")](https://github.com/tliron/puccini/releases)

Each tool is a self-contained executable file, allowing them to be easily distributed and embedded
in toolchains, orchestration, and development environments. They are coded in 100%
[Go](https://golang.org/) and are very portable, even available for
[WebAssembly](https://webassembly.org/) (which is how the in-browser demo linked above works).

You can also embed Puccini into your program as a library. Puccini is immediately usable from Go,
but can be used in many other programming languages via self-contained shared C libraries. See
included wrappers and examples for [Java](wrappers/java/), [Python](wrappers/python/), and
[Ruby](wrappers/ruby/).

To build Puccini yourself see the [build guide](scripts/).


puccini-tosca
-------------

⮕ [Documentation](puccini-tosca/)

A TOSCA processor. Parses a TOSCA service template and compiles it to Clout (see below).

Why TOSCA? It's a high-level language made for modeling and validating cloud topologies using
reusable and extensible objects. It allows architects to focus on application design and
requirements without being bogged down by the ever-changing specificities of the infrastructure.

Puccini can compile several popular TOSCA and TOSCA-like dialects:
[TOSCA 1.3](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html),
[TOSCA 1.2](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html),
[TOSCA 1.1](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html),
[TOSCA 1.0](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.0/TOSCA-Simple-Profile-YAML-v1.0.html),
as well as the more limited grammars of
[Cloudify DSL 1.3](https://docs.cloudify.co/5.0.5/developer/blueprints/),
and
[OpenStack Heat HOT 2018-08-31](https://docs.openstack.org/heat/stein/template_guide/hot_guide.html).

Puccini is also following the progress on
[TOSCA 2.0](http://docs.oasis-open.org/tosca/TOSCA/v2.0/TOSCA-v2.0.html).

The TOSCA source can be accessed by URL, on the local file systems or via HTTP/HTTPS, as
individual files as well as packaged in
[CSAR files](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html#_Toc302251718).
Puccini also comes with a simple CSAR creation tool, **puccini-csar**.

### Standalone Parser

⮕ [Documentation](tosca/parser/)

Puccini's TOSCA parser is also available as an independent Go library. Its 5 phases do
normalization, validation, inheritance, and assignment of TOSCA's many types and templates, resulting
in a [flat, serializable data structure](tosca/normal/) that can easily be consumed by your program.
Validation error messages are precise and useful. It's a very, very fast multi-threaded parser, fast
enough that it can be usefully embedded in editors and IDEs for validating TOSCA while typing.
For an example, see the
[TOSCA Visual Studio Code Extension](https://github.com/tliron/puccini-vscode/).

TOSCA is a complex object-oriented language. We put considerable effort into adhering to every
aspect of the grammar, especially in regards to value type checking and type inheritance contracts,
which are key to delivering the object-oriented promise of extensibility while maintaining reliable
base type compatibility. Unfortunately, the TOSCA specification is famously inconsistent and
imprecise. For this reason, the Puccini parser also supports [quirk modes](tosca/QUIRKS.md)
that enable alternative behaviors based on differing interpretations of the spec.

### Compiler

The TOSCA-to-Clout compiler's main role is to take the parsed data structure and dump it into
Clout. The next step in the toolchain (which could be **puccini-clout**) would then connect the
Clout to your orchestration systems: deploying to your platforms, on-boarding to a service catalog,
etc. Thusly Clout functions as an "intermediate representation" (IR) for TOSCA.

By default the compiler also performs [topology resolution](tosca/compiler/RESOLUTION.md), which
attempts to satisfy requirements with capabilities, thus creating the relationships (Clout edges)
between node templates. This feature can be turned off in order to add more processing phases
before final resolution. Resolution is handled via the embedded **tosca.resolve** scriptlet.

### Visualization

You can graphically visualize the compiled TOSCA in a dynamic web page. A one-line example:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml --exec=assets/tosca/profiles/common/1.0/js/visualize.js > /tmp/puccini.html && xdg-open /tmp/puccini.html


puccini-clout
-------------

⮕ [Documentation](puccini-clout/)

A Clout processor. Can execute JavaScript scriptlets, whether they are in a Clout file (in the
metadata section) or provided as external files. It can evaluate TOSCA functions, apply constraints,
execute Kubernetes specification generation, translate workflows to BPMN, etc.

The tool can also be used to add/remove scriptlets by manipulating the metadata section in the
Clout. 

Also supported are implementation-specific JavaScript "plugins" that allow you to extend existing
scriptlet functionality without having to modify it. For example, you can add a plugin for
Kubernetes to handle custom application needs, such as adding sidecars, routers, loadbalancers, 
an Istio service mesh, etc. You can also use **puccini-clout** to add plugins to the Clout file,
either storing them permanently or piping through to add and execute them on-the-fly.

Note that **puccini-clout** is *not* a requirement for your toolchain. You can process and consume
the Clout output with your own tools.


Clout
-----

⮕ [Documentation](clout/)

Introducing the **clou**d **t**opology ("clou" + "t") representation language, which is simply a
representation of a generic graph database in YAML/JSON/XML.

Clout functions as the intermediary format for your deployments. As an analogy, consider a program
written in the C language. First, you must *compile* the C source into machine code for your
hardware architecture. Then, you *link* the compiled object, together with various libraries, into a
deployable executable for a specific target platform. Clout is the compiled object in this analogy.
If you only care about the final result then you won't see the Clout at all. However, the decoupling
allows for a more powerful toolchain. For example, some tools might change your Clout after the
initial compilation (to scale out, to optimize, to add platform hooks, debugging features, etc.) and
then you just need to "re-link" in order to update your deployment. This can happen without
requiring you to update your original source design. It may also possible to "de-compile" some cloud
deployments so that you can generate a Clout without any TOSCA "source code".

Clout is essentially a big, unopinionated, implementation-specific dump of vertexes and the edges
between them with un-typed, non-validated properties. Rule #1 of Clout is that everything and the
kitchen sink should be in one Clout file. Really, anything goes: specifications, configurations,
metadata, annotations, source code, documentation, and even text-encoded binaries. (The only
exception might be that you might want to store security certificates and keys elsewhere.)

In itself Clout is an unremarkable format. Think of it as a way to gather various deployment
specifications for disparate technologies in one place while allowing for the *relationships*
(edges) between entities to be specified and annotated. That's the topology.

Clout is not supposed to be human-readable or human-manageable. The idea is to use tools (Clout
frontends and processors) to deal with its complexity. We have some great ones for you here. For
example, with Puccini you can use just a little bit of TOSCA to generate a single big Clout file
that describes a complex Kubernetes service mesh.

If Clout's file size is an issue, it's good to know that Clout is usually eminently compressible,
comprising just text with quite a lot of repetition.

### Storage

Orchestrators may choose to store Clout opaquely, as is, in a key-value database or filesystem.
This could work well because cloud deployments change infrequently: often all that's needed is to
retrieve a Clout, parse and lookup data, and possibly update a TOSCA attribute and store it again.
Iterating many Clouts in sequence this way could be done quickly enough even for large
environments. Simple solutions are often best.

That said, it could also make sense to store Clout data in a graph database. This would allow for
sophisticated queries, using languages such [GraphQL](https://graphql.org/) and
[Gremlin](https://tinkerpop.apache.org/gremlin.html), as well as localized transactional updates.
This approach could be especially useful for highly composable and dynamic environments in which
Clouts combine together to form larger topologies and even relate to data coming from other systems.

Graph databases are quite diverse in features and Clout is very flexible, so one schema will not
fit all. Puccini instead comes with examples: see [storing in Neo4j](examples/neo4j/) and
[storing in Dgraph](examples/dgraph/).


FAQ
---

### Can Puccini deploy the same TOSCA to either Kubernetes, OpenStack, AWS, Azure, etc.?

If you didn't plan it that way, then: no. Firstly, Puccini is *not* an orchestrator.
(It's "BYOO" = "Bring Your Own Orchestrator"). Secondly, assuming you have an orchestrator
that supports all those cloud platforms, creating portable TOSCA is less useful than may
initially appear.

Indeed, many orchestrators come with a variety of adapters or plugins to support various
cloud platforms, with a very specific approach to portability. For example, Ansible's
inventory plugin mechanism can generate a list of host addresses from various cloud sources.
You could thus run the same playbook on virtual machines deployed on OpenStack, AWS,
baremetal, etc.

But what would portability mean for TOSCA? If your goal is to use TOSCA to generate
that inventory, for example with Puccini's [Ansible inventory plugin](examples/ansible/hosts/),
then how would you create your TOSCA topology template? What would be the node types?

The best answer would be to use types that are appropriate for your cloud platform,
that directly model its compute, networking, storage, etc. resources. So, just like you
would use a different Ansible inventory plugin for the different cloud platform you are
targeting, you would have a different TOSCA service template for each cloud platform.

However, if your desire is to use a single TOSCA service template for *all* clouds,
how would you model it? Your node types would have to match the lowest common set of features
for all cloud platforms and nothing more. You could not use any platform-specific feature
without breaking portability.

Until TOSCA version 1.3 the specification came with an implicit
[Simple Profile](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html)
(as well as a [Simple Profile for NFV](https://docs.oasis-open.org/tosca/tosca-nfv/v1.0/tosca-nfv-v1.0.html))
that intended to be exactly that: a set of agnostic models for compute, networking,
and storage resources that all cloud platforms were presumed to have.

Unfortunately, this approach didn't work well in practice, for several reasons. Firstly,
because cloud platforms are just too different, even in terms of the basics of how they
model simple resources. Secondly, because any non-trivial workload would want to leverage
platform-specific features. Indeed that is the competitive advantage and indeed a reason to
choose one platform over another. Thirdly, because some platform-specific features are
non-optional. In such cases the "portable" TOSCA would have to be supplemented with extra
artifacts required by each cloud platform, such that though you would have some shared TOSCA
you would still need to maintain some platform-specific code, compromising on the goal of true
portability.

For all these reasons, since TOSCA 2.0 it no longer comes with a normative Simple Profile
and does not commit to an expectation for transparent portability. Still, nothing is stopping
you from trying anyway. Perhaps you can come up with a better set of portable models than those
of the Simple Profile. It's just not a design goal for Puccini.

The value proposition of TOSCA is not "write-once-run-everywhere" portability but "use
the same language and tooling everywhere". Puccini's goal is provide some of that toolchain.

### Why Go?

[Go](https://golang.org/) is fast becoming the language of choice for cloud-native solutions.
It has the advantage of producing very deployable executables that make it easy to containerize
and integrate. Go features garbage collection and easy multi-threading (via lightweight
goroutines), but unlike Python, Ruby, and Perl it is a strictly typed language, which
encourages good programming practices and reduces the chance for bugs.

### Why JavaScript?

JavaScript lets you manipulate the Clout data structures directly using a full-blown, conventional
language. It's probably
[not anyone's favorite language](https://archive.org/details/wat_destroyallsoftware), but it's
familiar, mature, standardized (as [ECMAScript](https://en.wikipedia.org/wiki/ECMAScript)), and does
the job. From a certain angle it's essentially the Scheme language (because it has powerful closures
and functions are first class citizens) but with a crusty C syntax.

And because JavaScript is self-contained text, it's trivial to store it in a Clout file, which can
then be interpreted and run almost anywhere.

Our chosen ECMAScript engine is [goja](https://github.com/dop251/goja), which is 100% Go and does
not require any external dependencies.

### Is there an alternative to JavaScript if I just need to extract data from the Clout?

If the built-in JavaScript support is insufficient or unwanted, you can write your own custom YAML
processor in Python, Ruby, etc., to do exactly what you need, e.g.:

    puccini-tosca compile my-app.yaml | python myprocessor.py

Also check out [yq](https://mikefarah.gitbook.io/yq/), a great little tool for extracting YAML and
even performing simple manipulations. Example:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml | yq r - 'vertexes.(properties.name==light12)'

### Can I use simple text templating instead of TOSCA functions and YAML processing?

Nothing is stopping you. You can pipe the TOSCA or Clout to and from the text translator of your
choice at any point in the toolchain. Here's an example using
[gomplate](https://github.com/hairyhenderson/gomplate):

    puccini-tosca compile my-app.yaml | gomplate

Your TOSCA `my-app.yaml` can then include template expressions, such as:

	username: "{{strings.ReplaceAll "\"" "\\\"" .Env.USER}}"

Note the proper escaping of quotation marks to avoid invalid YAML. Also, remember that indentation
in YAML is significant, so it can be tricky to insert blocks into arbitrary locations. Generally,
using text templating to manipulate YAML is not a great idea for these reasons.

Puccini's decision to use an embedded interpreted programming language (JavaScript) is intentional
and important. Unlike some tools (see [Helm](https://helm.sh/)), we prefer not treat YAML
files as plain text to be manipulated by an anemic text templating language.

If you insist on text templating, a useful convention for your toolchain could be to add a file
extension to mark a file for template processing. For example, `.yaml.j2` could be recognized as
requiring Jinja2 template processing, after which the `.j2` extension would be stripped.

### Can I compose a single service from several interrelated Clout files?

TOSCA has a feature called "substitution mapping", which is useful for modeling service composition.
It is, however, a *design* feature. The implementation is up to your orchestration toolchain. See
our examples
[here](examples/tosca/substitution-mapping.yaml) and
[here](examples/tosca/substitution-mapping-client.yaml).

Puccini intentionally does *not* support service composition. Each Clout file is its own universe.
If you need to create edges between vertexes in one Clout file and vertexes in other Clout files,
then it's up to you and your tools to design and implement that integration. The solution could be
very elaborate indeed: the two Clouts might represent services with very different lifecycles, that
run in different clouds, that are handled by different orchestrators. And the connection might
require complex networking to achieve. There's simply no one-size-fits-all way Puccini could do
it—namespaces? proxies? catalogs? repositories?—so it insists on not having an opinion.

### TOSCA is so complicated! Help!

I know, right? Now imagine writing a parser for it... Not only is it a complex language, but the
[specification itself](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html)
(as of version 1.3) has many contradictions, errors, and gaps.

Please join [OASIS's TOSCA community](https://www.oasis-open.org/committees/tc_home.php?wg_abbrev=tosca)
to help improve the language!

Meanwhile, we've included [examples](examples/tosca/) of TOSCA core grammatical features,
with some running commentary. Treat them as your playground. Also, if you have 4 hours to spare,
grab some snacks, get comfortable, and watch the author's free online course for TOSCA 1.0:
[part 1](https://www.youtube.com/watch?v=aMkqLI6o-58),
[part 2](https://www.youtube.com/watch?v=6xGmpi--7-A).

(Author's note: This is my second take at writing a TOSCA parser. The first was
[AriaTosca](https://github.com/apache/incubator-ariatosca), an
incubation project under the Apache Software Foundation. I am grateful to
[Cloudify](https://cloudify.co/) for funding much of the AriaTosca project. Note, however, that
Puccini is a fresh start initiated by myself with no commercial backing. It does not use
AriaTosca code and has a radically different architecture as well as very different goals.)

### Why is it called "Puccini"?

[Giacomo Puccini](https://en.wikipedia.org/wiki/Giacomo_Puccini) was the composer of the
[*Tosca*](https://en.wikipedia.org/wiki/Tosca) opera (based on Victorien Sardou's play,
[*La Tosca*](https://en.wikipedia.org/wiki/La_Tosca)), as well as *La bohème*, *Madama Butterfly*,
and other famous works. The theme here is orchestration, orchestras, composition, and thus operas.
Capiche?

### How to pronounce "Puccini"?

For a demonstration of its authentic 19th-century Italian pronunciation see
[this clip](https://www.youtube.com/watch?v=dQw4w9WgXcQ).
