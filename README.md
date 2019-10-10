Fyde CLI client [![Go Report Card](https://goreportcard.com/badge/github.com/fyde/fyde-cli)](https://goreportcard.com/report/github.com/fyde/fyde-cli) [![license](https://img.shields.io/github/license/fyde/fyde-cli.svg)](https://github.com/fyde/fyde-cli/blob/master/LICENSE)
===============

Cross-platform command line client for [Fyde Enterprise Console](https://fyde.github.io/docs/fyde-enterprise-console) APIs.

## Overview

fyde-cli enables interaction with the Fyde Enterprise Console using the command line.
It is designed with both interactive usage and scripting in mind.

fyde-cli is written in [Go](https://golang.org) and supports different architectures.
We provide pre-built binaries for x86-64 Windows, macOS and Linux, as well as ARM and ARM64 Linux.
You can use it in other architectures by compiling from source.

## Features

The goal for this project is to implement most operations possible with the web version of the [Fyde Enterprise Console](https://fyde.github.io/docs/fyde-enterprise-console), enabling batch mode for certain operations (like adding multiple users from a CSV file or a JSON dump).

Currently, the fyde-cli is at an early development stage.
Only read-only operations are implemented:

 - List users, groups, devices, resources, proxies and policies
 - Get info about specific user, group, device, resource, proxy or policy
 - List activity records and get info about specific record
 - Watch activity records as they happen in real-time

## Installation

If you are using an operating system and architecture for which we provide pre-built binaries, we recommend using those.
Just download the appropriate archive from the [releases page](https://github.com/fyde/fyde-cli/releases).
The fyde-cli binaries are statically compiled and have no external dependencies.

Inside the archive, there's the executable, a copy of this document and the license. Simply extract the executable and place it where you find most convenient (for example, in most Linux distros, you could use `/usr/local/bin` for a system-wide install).
Don't forget to add the executable to `$PATH`, or equivalent, if desired.

### Installing from source

If we do not provide a pre-built binary for your platform, or if you want to make changes to fyde-cli, you can compile it yourself.

#### Prerequisites

 - [Git](https://git-scm.com/)
 - [Go](https://golang.org) (version 1.13 or higher)
 - [go-swagger](https://github.com/go-swagger/go-swagger) - if you run into problems with the latest release, you can compile from the master branch, commit [5499ab](https://github.com/go-swagger/go-swagger/commit/5499abf2a8c86a57f3a8112aca47a624f609689e).

#### Obtaining the code

fyde-cli uses the modules support introduced in Go 1.11, which means you are not forced to place the code in a specific path under GOPATH. You can clone the repository into any folder:

```
git clone https://github.com/fyde/fyde-cli.git
cd fyde-cli
```

You can also clone the repo into the usual `$GOPATH/src/github.com/fyde/fyde-cli` path, but keep in mind that the project will not compile until the next step is complete (i.e. `go get github.com/fyde/fyde-cli` will always fail).

#### Generating code from the Swagger specification

After installing [go-swagger](https://github.com/go-swagger/go-swagger), run the following command on the root of this repo:

`swagger generate client -f swagger.yml`

This will generate the `client` and `models` packages.
The code in the `cmd` package depends on these.

#### Compiling

Simply run `go build`.
Because we are using Go modules, Go will take care of downloading the correct versions of the dependencies.

You should now have a `fyde-cli` executable.
You can `go install` the package, if you wish, which will place the binary in `$GOPATH/bin`.

## Usage

When run without arguments, fyde-cli presents a list of available commands.
It will also show where it is going to save and look for the configuration files, unless overridden.

To use the client with an endpoint other than the default, you should start by setting the endpoint:

```
$ fyde-cli endpoint set fydeconsole.example.com
Endpoint changed. Credentials cleared, please login again using `fyde-cli login`
```

You can then proceed to log in with your console credentials:

```
$ fyde-cli login
Username: you@example.com
Password:
Logged in successfully, access token stored in (...)fyde/fyde-cli/auth.yaml
```

You can now use other commands. For example, to list users, you can use `fyde-cli users list`.

All commands provide a help text with the available subcommands and flags.
For example, running `fyde-cli resources` will let you know about the `get` and `list` subcommands, and `fyde-cli resources list --help` will list all available flags for the list resources command, including pagination, sorting and filtering flags.

### Output formats

fyde-cli supports different output formats for different use cases:

 - Table, for interactive usage (`--output=table`)
 - CSV (`--output=csv`)
 - JSON (`--output=json` or `--output-json-pretty`)

By default, when an interactive terminal is detected, `table` output is used.
Otherwise, `json` is used.
JSON output generally contains the most information, sometimes including nested objects; CSV output corresponds to a CSV version of the table output.

All output formats are subject to pagination parameters, when those are available.

## Reporting issues

You can see existing issues and report new ones [on GitHub](https://github.com/fyde/fyde-cli/issues).

## License

fyde-cli is Copyright © 2019 Fyde, Inc. and is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), a OSI-approved license.