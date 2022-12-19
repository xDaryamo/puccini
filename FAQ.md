Puccini FAQ
===========

### Why TOSCA?

It's a high-level language made for modeling and validating cloud topologies using reusable and
extensible objects. It allows architects to focus on application design and requirements without
being bogged down by the ever-changing specificities of cloud platforms, while also allowing
engineers to provide rich, up-to-date integrations with those platforms.

### TOSCA is so complicated! Help!

I know, right? Now imagine writing a parser for it... Not only is it a complex language, but the
[specification itself](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html)
(as of version 1.3) has contradictions, errors, and gaps.

Please join [OASIS's TOSCA community](https://www.oasis-open.org/committees/tc_home.php?wg_abbrev=tosca)
to help improve the language!

Meanwhile, Puccini includes [examples](examples/tosca/) of TOSCA's grammatical features with some
running commentary. Treat them as your playground. Also, if you have 4 hours to spare, grab some
snacks, get comfortable, and watch the author's free online course for TOSCA 1.0:
[part 1](https://www.youtube.com/watch?v=aMkqLI6o-58),
[part 2](https://www.youtube.com/watch?v=6xGmpi--7-A).

(Author's note: This is my second take at writing a TOSCA parser. The first was
[AriaTosca](https://github.com/apache/incubator-ariatosca), an
incubation project under the Apache Software Foundation. I am grateful to
[Cloudify](https://cloudify.co/) for funding much of the AriaTosca project. Note, however, that
Puccini is a fresh start initiated by myself with no commercial backing. It does not use
AriaTosca code and has a radically different architecture as well as very different goals.)

### Can Puccini deploy the same TOSCA to either Kubernetes, OpenStack, AWS, Azure, etc.?

If you didn't plan it that way, then: no. Firstly, Puccini is *not* an orchestrator.
(It's "BYOO" = "Bring Your Own Orchestrator"). Secondly, creating portable TOSCA may be
less useful than expected.

Many orchestrators come with a variety of adapters or plugins to support various cloud
platforms coupled with a very specific approach to portability. For example, Ansible's
inventory plugin mechanism can generate a list of host addresses from various cloud
sources. You could thus run the same playbook on virtual machines deployed on OpenStack,
AWS, baremetal, etc.

But what would portability mean for TOSCA? If your goal is to use TOSCA to generate
that inventory, for example with Puccini's [Ansible inventory plugin](examples/ansible/hosts/),
then how would you create your TOSCA topology template for it to be portable? What would be the
node types?

The best answer would be to use types that are appropriate for your cloud platform,
that directly model its compute, networking, storage, etc. resources. So, just like you
would use a different Ansible inventory plugin for each cloud platform you are targeting,
you would have a different TOSCA service template for each cloud platform.

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

And so TOSCA 2.0 no longer comes with a normative Simple Profile and does not commit to an
expectation for transparent portability. That said, nothing is stopping you from trying anyway.
Perhaps you can come up with a better set of portable models than those of the Simple Profile.
It's just not a goal for TOSCA 2.0 and not a goal for Puccini.

The value proposition of TOSCA is not "write-once-run-everywhere" portability but "use
the same language and tooling everywhere", and the value proposition of Puccini is to provide
some of that toolchain.

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

If the built-in JavaScript support is insufficient or undesirable, you can write your own custom YAML
processor in Python, Ruby, etc., to do exactly what you need, e.g.:

    puccini-tosca compile my-app.yaml | python myprocessor.py

Also check out [yq](https://mikefarah.gitbook.io/yq/), a great little tool for extracting YAML and
even performing simple manipulations. Example:

    puccini-tosca compile examples/tosca/requirements-and-capabilities.yaml | yq '.vertexes.[]|select(.properties.name=="light6")'

### Can I use text templating instead of TOSCA functions like `get_input`?

TOSCA's avoidance of templating is deliberate. Unlike some tools (see
[Helm](https://helm.sh/)), YAML is not treated as plain text to be manipulated by an anemic text
templating language, but rather as a representation of a semantic structure.

That said, nothing is stopping you from using templates. You can pipe the TOSCA or Clout to and
from the tool of your choice at any point. Just note the proper escaping of quotation marks to avoid
invalid YAML. Also remember that indentation in YAML is significant, so it can be tricky to insert
blocks into arbitrary locations. Generally, using text templating to manipulate YAML is not a great
idea for these reasons.

Here's an example using
[gomplate](https://github.com/hairyhenderson/gomplate):

    puccini-tosca compile my-app.yaml | gomplate

Your TOSCA `my-app.yaml` can then include template expressions, such as:

	username: "{{strings.ReplaceAll "\"" "\\\"" .Env.USER}}"

If you insist on text templating, a useful convention for your toolchain could be to add a file
extension to mark a file for template processing. For example, `.yaml.j2` could be recognized as
requiring Jinja2 template processing, after which the `.j2` extension would be stripped.

### Can I compose a single service from several interrelated Clout files?

TOSCA has a feature called "substitution mapping", which is useful for modeling service composition.
It is, however, a *design* feature. The implementation is up to your orchestration toolchain. See
our examples
[here](examples/tosca/substitution-mapping.yaml) and
[here](examples/tosca/substitution-mapping-client.yaml).

Puccini intentionally does *not* support service composition. Each Clout file is its own graph
universe. If you need to create graph edges between vertexes in one Clout file and vertexes in other
Clout files, then it's up to you and your tools to design and implement that integration. The solution
could be very elaborate indeed: the two Clouts might represent services with very different lifecycles,
that run in different clouds, that are handled by different orchestrators. And the connection might
require complex networking to achieve. There's simply no one-size-fits-all way Puccini could do
it—namespaces? proxies? catalogs? repositories?—so it insists on not having an opinion.

### Why is it called "Puccini"?

[Giacomo Puccini](https://en.wikipedia.org/wiki/Giacomo_Puccini) was the composer of the
[*Tosca*](https://en.wikipedia.org/wiki/Tosca) opera (based on Victorien Sardou's play,
[*La Tosca*](https://en.wikipedia.org/wiki/La_Tosca)), as well as *La bohème*, *Madama Butterfly*,
and other famous works. The theme here is orchestration, orchestras, composition, and thus operas.
Capiche?

### How to pronounce "Puccini"?

For a demonstration of its authentic 19th-century Italian pronunciation see
[this clip](https://www.youtube.com/watch?v=dQw4w9WgXcQ).
