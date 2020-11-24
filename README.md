# kphp-cli

`kphp` is a [Go-inspired](https://golang.org/cmd/go/) command-line interface to the [KPHP](https://github.com/VKCOM/kphp/) compiler.

```bash
Usage:

  kphp COMMAND

Possible commands are:

  run        compile and run KPHP script
  env        print KPHP environment information
  version    print KPHP version information

Run 'kphp COMMAND --help' to see more information about a command.
```

## Installation

### Installation: get the `kphp` binary

Download the latest precompiled binary:

* [linux-amd64](TODO)
* [darwin-amd64](TODO)
* [windows-amd64](TODO)

Or you can install it from sources if you have [Go](https://golang.org/) installed:

```bash
# To install ruleguard binary under your $(go env GOPATH)/bin:
go get -v -u github.com/quasilyte/kphp-cli/cmd/kphp
```

### Installation: setting the `$KPHP_ROOT`

```bash
# Suppose that we want to install KPHP to $HOME/kphp
export KPHP_ROOT=$HOME/kphp

git clone https://github.com/VKCOM/kphp.git $KPHP_ROOT
cd $KPHP_ROOT
mkdir build
cd build
cmake ..
make
```

If you have problems related to the KPHP compilation, check the [official documentation](https://vkcom.github.io/kphp/kphp-internals/developing-and-extending-kphp/compiling-kphp-from-sources.html).

## `kphp run` subcommand

Example:

```bash
$ kphp run script.php
```

`kphp run --help`:

```
Usage of kphp run:
  -composer-root string
    	A folder that contains the root composer.json file
  -server
    	Whether to compile a script with --mode=server
```
