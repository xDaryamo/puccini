Puccini
=======

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/puccini.svg)](https://github.com/tliron/puccini/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/puccini)](https://goreportcard.com/report/github.com/tliron/puccini)

Deliberately stateless cloud topology management and deployment tools for Kubernetes and TOSCA.

Check out the [quickstart guide](QUICKSTART.md).

puccini-tosca
-------------

Clout frontend for TOSCA. Parses a TOSCA service template and compiles it to Clout (see below).

Why [TOSCA](https://www.oasis-open.org/committees/tosca/)? It's a high-level language designed
exactly for modeling and validating cloud topologies with reusable and inheritable objects. It
allows architects to focus on application logic and requirements without being bogged down by the
ever-changing specificities of the infrastructure.

**puccini-tosca** comes with
[a TOSCA profile for Kubernetes/Istio](assets/tosca/profiles/kubernetes/1.10/), which includes node,
capability, relationship, policy, and other types, as well as straightforward JavaScript code to
generate [Kubernetes](https://kubernetes.io/)/[Istio](https://istio.io/) specs.

We support
[CSAR files](http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/os/TOSCA-Simple-Profile-YAML-v1.1-os.html#_Toc489606742)
(TOSCA packages) in addition to YAML files. We're even including a simple CSAR creation tool,
**puccini-csar**.

How do TOSCA, Clout, JavaScript, and Kubernetes all fit together in Puccini? Consider this: with a
single command line you can take a TOSCA service template, compile it with **puccini-tosca**, pipe
the Clout through the **puccini-js** processor, which will run JavaScript to generate Kubernetes
specs, then pipe those to [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/),
which will finally upload the specs to a running Kubernetes cluster. Like so:

     puccini-tosca compile my-app.yaml | puccini-js exec kubernetes.generate | kubectl apply -f -

Et voilà, your abstract design became a running deployment.

### Standalone Parser

Puccini's [TOSCA parser](tosca/parser/README.md) is available as an independent Go library. Its 6
phases do normalization, validation, inheritance, and assignment of TOSCA's many types and
templates, finally satisyfing requirements with capabilities, resulting in a
[flat, serializable data structure](tosca/normal/README.md) that can easily be consumed by your
program. Validation error messages are precise and useful. It's a very, very fast parser, enough
that it can be usefully embedded in editors and IDEs for validating TOSCA while typing.

TOSCA is a complex object-oriented language. Considerable effort has been put into adhering to every
aspect the grammar, especially in regards to value type checking and type inheritance contracts.

### Compiler

The TOSCA-to-Clout compiler just takes the parsed data structure and dumps it into Clout. It also
includes any JavaScript required to process the Clout. Thusly Clout functions as an "intermediate
representation" (IR) for TOSCA.

* [**puccini-tosca** documentation](puccini-tosca/README.md)
* [TOSCA parser documentation](tosca/parser/README.md)

puccini-js
----------

Clout processor for JavaScript. Executes existing JavaScript in a Clout file. For example, it can
execute the Kubernetes spec generation code inserted by **puccini-tosca**. It also supports
executing intrinsic functions and value constraints (for example, TOSCA's).

Also supported are implementation-specific JavaScript "plugins" that allow you to extend existing
functionality. For example, you can add a plugin for Kubernetes to handle custom application needs,
such as adding sidecars, routers, loadbalancers, etc. Indeed, Istio support is implemented as a
plugin. You can also use **puccini-js** to add plugins to the Clout file, either storing them
permanently or piping through to add and execute them on-the-fly.

### TOSCA Intrinsic Functions and Constraints

These are implemented in JavaScript so that they can be put into the Clout and then be executed
by **puccini-js**, allowing a compiled-from-TOSCA Clout file to be entirely independent from TOSCA.
The Clout lives on its own.

To call these functions we provide the **tosca.coerce** JavaScript, which calls all functions and
replaces the call stubs with the returned values:

    puccini-js exec tosca.coerce my-clout.yaml --output=coerced-clout.yaml

A useful side benefit of this implementation is we allow you to easily extend TOSCA by
[adding your own functions/constraints](examples/grammar/functions-custom.yaml). Obviously, such
custom functions are not part of the TOSCA spec and will not be compatible with other TOSCA
implementations.

### TOSCA Attributes

TOSCA attributes (as opposed to properties) represent live data in a running deployment. And the
intrinsic function, `get_attribute`, allows other values to make use of this live data. The
implication is that some values in the Clout should change as these attributes change. But also,
attribute definitions in TOSCA allow you to define constraints on the value, so we must also make
sure that the new data complies with them.

Our solution has two steps. First, we have JavaScript (**kubernetes.update**) that extracts these
attributes from a Kubernetes cluster (by calling **kubectl**) and updates the Clout. Second, we
run **tosca.coerce**, which not only calls instrinsic functions but also applies the constraints.

Putting it all together, let's refresh a Clout:

    puccini-js exec kubernetes.update my-clout.yaml | puccini-js exec tosca.coerce -o coerced-clout.yaml

### Workflows

*WORK IN PROGRESS*

Also supported are TOSCA workflows on Kubernetes. The actual operations are executed on the
containers using `kubectl exec` after deploying the artifacts with `kubectl cp`. Note that this
straightforward feature is not intended to replace robust orchestration solutions, but it can
form the basis of one. Of course, you can write your own Clout processor to run workflows on your
favorite workflow engine, such as [Ansible](https://www.ansible.com/) or a
[BPEL](https://en.wikipedia.org/wiki/Business_Process_Execution_Language)/[BPMN](https://en.wikipedia.org/wiki/Business_Process_Model_and_Notation)
implementation.

* [**puccini-js** documentation](puccini-js/README.md)

Clout
-----

Introducing the cloud topology ("clou-" + "t") representation language, which can be formatted as
YAML/JSON/XML.

Clout is an intermediary format for your deployments. As an analogy, consider a program written in
the C language. First, you must *compile* the C source into machine code for your hardware
architecture. Then, you *link* the compiled object, together with various libraries, into a
deployable executable for a specific target platform. Clout here is the compiled object. If you only
care about the final result then you won't see the Clout at all. However, this decoupling allows for
a more powerful tool chain. For example, some tools might change your Clout after the initial
compilation (to scale out, to optimize, to add platform hooks, debugging features, etc.) and then
you just need to "re-link" in order to update your deployment. This can happen without requiring
you to update your original source design. It may also possible to "de-compile" some cloud
deployments so that you can generate a Clout without "source code".

Clout is essentially a big, unopinionated, implementation-specific dump of vertexes and the edges
between them with un-typed, non-validated properties. Rule #1 of Clout is that everything and the
kitchen sink should be in one Clout file. Really, anything goes: specifications, configurations,
metadata, annotations, source code, documentation, and even text-encoded binaries. (The only
possible exception might be that you would want to store security certificates and keys
elsewhere.)

In itself Clout is an unremarkable format. Think of it as a way to gather various deployment specs
for disparate technologies in one place while allowing for the *relationships* (edges) between
entities to be specified and annotated. That's the topology.

Clout is not supposed to be human-readable or human-manageable. The idea is to use tools (Clout
frontends and processors) to deal with its complexity. We have some great ones for you here. For
example, with Puccini you can use just a little bit of TOSCA to generate a single big Clout file
that describes a complex Kubernetes service mesh.

If Clout's file size is an issue, it's good to know that Clout is usually eminently compressible,
comprising just text with quite a lot of repetition.

* [Clout documentation](clout/README.md)

FAQ
---

### Can Puccini deploy my existing TOSCA to Kubernetes?

If you didn't plan it that way, then: no. Our TOSCA Kubernetes/Istio profiles do *not* make use of
TOSCA's
[Simple Profile](http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html)
or [Simple Profile for NFV](http://docs.oasis-open.org/tosca/tosca-nfv/v1.0/tosca-nfv-v1.0.html)
types (Compute, BlockStorage, VDU, etc.). Still, if you find these so-called "normative" types
useful, they are included in Puccini and will be compiled into Clout. You may write your own
JavaScript to deploy them to your cloud orchestration environment. For example, it could make sense
to generate an OpenStack [Heat](https://wiki.openstack.org/wiki/Heat) template, or even call the
OpenStack APIs directly, allowing you to go from TOSCA through Clout to OpenStack, just as we go
from TOSCA to Kubernetes.

Generally speaking, the notion that a single set of normative types could be used for all the
various cloud and container platforms out there is a pipe dream. The devil is in the details, and
the amount of detail needed for cloud deployment keeps growing and diversifying. Thus every Clout
file is vehemently platform-specific. However, by bringing the tiniest implementation details all
into one place we can at least have a common tool chain for all platforms. That's the gist of
Puccini.

### JavaScript? Really?

The decision to use an interpreted programming language is intentional and important. Unlike some
Kubernetes tools ([Helm](https://helm.sh/)), we do not treat YAML files as plain text to be
manipulated by an anemic text templating language, where working around YAML's strict
indentation is a nightmare.

JavaScript lets you manipulate data structures directly using a full-blown, conventional language.
It's probably
[not anyone's favorite language](https://archive.org/details/wat_destroyallsoftware), but it's
familiar, mature, standardized (as [ECMAScript](https://en.wikipedia.org/wiki/ECMAScript)), and does
the job. From a certain angle it's essentially Scheme (because it has powerful closures and
functions are first class citizens), just with a crusty C syntax.

And because JavaScript is self-contained text, it's trivial to store it in a Clout file, which can
then be interpreted and run almost anywhere.

Our chosen ECMAScript engine is [goja](https://github.com/dop251/goja), which is 100% Go and does
not require any external dependencies.

### Can't I use simple text templating instead of JavaScript?

Nothing is stopping you. You can pipe the input and output to and from the text translator of your
choice at any point in the tool chain. Here's an example using
[gomplate](https://github.com/hairyhenderson/gomplate):

    puccini-tosca compile my-app.yaml | gomplate | puccini-js exec kubernetes.generate

Your TOSCA can then inject expressions into values, such as `username: '{{.Env.USER}}'`.

### Instead of one big Clout file can I compose a service from several interrelated Clout files?

Clout intentionally does *not* support service composition. Each Clout file is its own universe. If you
need to relate vertexes in one Clout file to others, then it's up to you and your tools to describe
that composition, possibly by combining multiple Clout files into one. There's simply no
one-size-fits-all way Clout could do it—namespaces? proxies? catalogs? repositories?—so it
insists on not having an opinion.

TOSCA has a feature called "substitution mapping", which is useful for modeling service composition.
However, it's a design feature. The implementation, which would likely be very complex, is up to
your orchestration tool chain. See our examples
[here](examples/grammar/substitution-mapping.yaml) and
[here](examples/grammar/substitution-mapping-client.yaml).

### TOSCA is so complicated! Help?

I know, right? Now imagine writing a parser for it... Not only is it a complex language, but the
[spec itself](http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html)
(as of version 1.1) has many contradictions, errors, and gaps.

To help you out we've included [examples](examples/grammar/) of TOSCA core grammatical features,
with some running commentary. Treat them as your playground. Also, if you have 4 hours to spare,
grab some snacks, get comfortable, and watch this free online course for TOSCA 1.0:
[part 1](https://www.youtube.com/watch?v=aMkqLI6o-58),
[part 2](https://www.youtube.com/watch?v=6xGmpi--7-A).

(Author's note: This is actually my second take at writing a TOSCA parser, after learning a great
deal from my previous efforts in [AriaTosca](https://github.com/apache/incubator-ariatosca), an
incubation project under the Apache Software Foundation. I am grateful to
[Cloudify](https://cloudify.co/) for funding much of the AriaTosca effort. Note, however, that
Puccini is a fresh start initiated by myself with no commercial backing. It does not use any of the
AriaTosca code and has a radically different architecture as well as very different goals.)

### Why doesn't the TOSCA parser tell me which line number in the relevant file a problem occurred?

Unfortunately, [our YAML parser](https://gopkg.in/yaml.v2) doesn't expose this information. There
is an [open issue](https://github.com/go-yaml/yaml/issues/108) for it, and if it's resolved we will
add this important feature in the future.

### Why is it called "Puccini"?

[Giacomo Puccini](https://en.wikipedia.org/wiki/Giacomo_Puccini) was the composer of the
[*Tosca*](https://en.wikipedia.org/wiki/Tosca) opera (based on Victorien Sardou's play,
[*La Tosca*](https://en.wikipedia.org/wiki/La_Tosca)), as well as *La bohème*, *Madama Butterfly*,
and other famous works. The theme here is orchestration, orchestras, composition, and thus operas.
Capiche?

### How to pronounce "Puccini"?

For a demonstration of its authentic 19th-century Italian pronunciation see
[this clip](https://www.youtube.com/watch?v=dQw4w9WgXcQ).
