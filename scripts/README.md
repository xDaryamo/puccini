Build Guide
===========

Puccini tries to follow the conventions of the Go programming community. However, newcomers to Go
might not know where to start. Here are some scripts to get you up in running.

We recommend developing on Linux. The Go compiler is a cross-compiler so you can build all binaries
for all platforms and architectures from there. A Linux virtual machine would be find.

If you're on Windows you have two additional options. The better one is
[WSL](https://learn.microsoft.com/en-us/windows/wsl/about), which *is* Linux, from where you can do
cross-compiling. The second option is to 1) [install Go for Windows](https://go.dev/doc/install)
and then 2) install [git for Windows](https://git-scm.com/download/win), which provides a Unix-like
environment with bash (use "Git Bash").

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
you to press TAB to complete commands starting with `puccini-tosca`, `puccini-clout`, etc. You will need
to reload your bash shell to see the effect. `exec bash` may work.

[embed](embed)
--------------

You don't need to run this normally. Run it only if you change any of the files in the
[assets](../assets/) directory. It reads those files and wraps them in Go code so that they can be
compiled into Puccini's executables. So, after running this you would likely also want to re-run
`build`.

[test](test)
------------

This runs `build` and then some tests using Go's built-in testing tool.

[release](release)
------------------

This script installs and runs the amazing [GoReleaser](https://goreleaser.com/) tool in order to
cross-compile Puccini, create installation packages for Linux, Windows, and MacOS, and publish them
to [GitHub](https://github.com/tliron/puccini/releases). It uses the
[.goreleaser.yml](../.goreleaser.yml) file for its configuration.

Run it with the "-t" switch to do a test release, which will create local distributions (under
`dist/release/`) without publishing them to GitHub.
