Development Guide
=================

Puccini tries to follow the conventions of the Go programming community. However, newcomers to Go
might not know where to start. Here are some scripts to get you up in running, in order:

[prepare.sh](prepare.sh)
------------------------

Though you can specify source files to the Go compiler, it's easiest if your projects are
conventionally organized in folders under the `GOPATH`. The default `GOPATH` is `~/go`, and you will
save yourself trouble by sticking to that default.

This script takes the current Puccini repository folder and symlinks it under the `GOPATH`. This
allows the Go tools to run as expected while allowing you to keep the repository where it's most
convenient for you.

[ensure.sh](ensure.sh)
----------------------

Puccini is deliberately 100% open source Go code, guaranteeing that it will build on all
Go-supported platforms, and so are its dependencies. But we need to download and build these
dependencies.

This script will do it for you by installing and running the [dep](https://github.com/golang/dep)
tool (also written in Go), which will analyze the code and will transiently identify dependencies.
Their source code will then be downloaded into the `vendors` directory, and then compiled into your
`GOPATH`.

This may take a minute or two the first time you run the tool, so please be patient.

The versions of dependencies are locked and specified in the [Gopkg.toml](../Gopkg.toml) and
[Gopkg.lock](../Gopkg.lock) files. Delete them and re-run to force the latest versions (which may of
course cause breakage).

[build.sh](build.sh)
--------------------

This builds the Puccini executables. You will find the output in `bin` under your `GOPATH`. If you
stick to the defaults, it would be `~/go/bin/puccini-tosca`, etc.

You might find it convenient to have `~/go/bin` on your search path. For Bash, add this to your
`~/.bashrc` file:

    export PATH="~/go/bin:$PATH" 

The Go compiler will only compile changed files. Also, it's a very fast compiler. So, generally you
should not be concerned about this step in your toolchain.

[install-bash-completion.sh](install-bash-completion.sh)
--------------------------------------------------------

Installs bash completion scripts for the current user, for the current build of Puccini. This allows
you to press TAB to complete commands starting with `puccini-tosca`, `puccini-js`, etc.

[embed.sh](embed.sh)
--------------------

You don't need to run this normally. Run it only if you change any of the files in the
[assets](../assets/) directory. It reads those files and wraps them in Go code so that they can be
compiled into Puccini's executables. So, after running this you would also have to re-run
`build.sh`.

[test.sh](test.sh)
------------------

This runs `build.sh` and then runs some tests using Go's built-in testing tool.

[release.sh](release.sh)
------------------------

This script installs and runs the amazing [GoReleaser](https://goreleaser.com/) tool in order to
cross-compile and create installation packages for Linux, Windows, and MacOS. It relies on the
[.goreleaser.yml](../.goreleaser.yml) file for its configuration.
