Build Guide
===========

Puccini tries to follow the conventions of the Go programming community. However, newcomers to Go
might not know where to start. Here are some scripts to get you up in running, in order:

[build](build)
--------------

This builds the Puccini executables, first making sure to download any necessary dependencies.
You will find the output in `bin` under your `GOPATH`. If you stick to the defaults, it would
be `~/go/bin/puccini-tosca`, etc.

You might find it convenient to have `~/go/bin` on your search path. For Bash, add this to your
`~/.bashrc` file:

    export PATH="~/go/bin:$PATH" 

The Go compiler will only compile changed files. Also, it's a very fast compiler. So, generally you
should not be concerned about this step in your toolchain.

Dependency management is handled by [Go modules](https://github.com/golang/go/wiki/Modules),
introduced in Go 1.11. See the files `go.mod` and `go.sum`.

[install-bash-completion](install-bash-completion)
--------------------------------------------------

Installs bash completion scripts for the current user, for the current build of Puccini. This allows
you to press TAB to complete commands starting with `puccini-tosca`, `puccini-js`, etc.

[embed](embed)
--------------

You don't need to run this normally. Run it only if you change any of the files in the
[assets](../assets/) directory. It reads those files and wraps them in Go code so that they can be
compiled into Puccini's executables. So, after running this you would likely also want to re-run
`build`.

[test](test)
------------

This runs `build` and then runs some tests using Go's built-in testing tool.

[release](release)
------------------

This script installs and runs the amazing [GoReleaser](https://goreleaser.com/) tool in order to
cross-compile and create installation packages for Linux, Windows, and MacOS. It relies on the
[.goreleaser.yml](../.goreleaser.yml) file for its configuration.
