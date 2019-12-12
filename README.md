Fyde CLI client [![Go Report Card](https://goreportcard.com/badge/github.com/fyde/fyde-cli)](https://goreportcard.com/report/github.com/fyde/fyde-cli) [![license](https://img.shields.io/github/license/fyde/fyde-cli.svg)](https://github.com/fyde/fyde-cli/blob/master/LICENSE)
===============

Cross-platform command line client for [Fyde Enterprise Console](https://fyde.github.io/docs/fyde-enterprise-console) APIs.

## Overview

fyde-cli enables interaction with the Fyde Enterprise Console using the command line.
It is designed with both interactive usage and scripting in mind.

fyde-cli is written in [Go](https://golang.org) and supports different architectures.
We provide pre-built i386 and x86-64 Windows and Linux binaries and x86-64 macOS binaries.
deb and rpm packages are also provided.
You can use fyde-cli in other architectures by compiling from source.

## Features

fyde-cli supports most operations possible with the web version of the [Fyde Enterprise Console](https://fyde.github.io/docs/fyde-enterprise-console), enabling batch mode for certain operations (like adding multiple users from a CSV file or a JSON dump).
The following operations are implemented:

 - List users, groups, devices, resources, proxies and policies
 - Get info about specific user, group, device, resource, proxy or policy
 - Create users, groups, resources, policies, proxies and domains, using command line flags or in batch mode, from files
 - Edit users, groups, resources, policies and proxies, using command line flags or in batch mode, from files
 - Delete users, groups, devices, resources, policies, proxies and domains
 - Generate, view, and revoke user enrollment links
 - Revoke device authentication
 - List activity records and get info about specific record
 - Watch activity records as they happen in real-time

fyde-cli will be continually updated to support new management console features.

## Installation

If you use [Homebrew](https://brew.sh/) on Linux or Mac, we recommend installing fyde-cli through our Homebrew tap - instructions [in the tap repo](https://github.com/fyde/homebrew-tap).

If you are using an operating system and architecture for which we provide pre-built binaries, we recommend using those.
Just download the appropriate archive from the [releases page](https://github.com/fyde/fyde-cli/releases).
We also provide deb and rpm packages.
The fyde-cli binaries are statically compiled and have no external dependencies.

Inside each archive, you'll find the executable for the corresponding platform, a copy of this document and the license. Simply extract the executable and place it where you find most convenient (for example, in most Linux distros, you could use `/usr/local/bin` for a system-wide install).
Don't forget to add the executable to `$PATH`, or equivalent, if desired.

### Installing from source

If we do not provide a pre-built binary for your platform, or if you want to make changes to fyde-cli, you can compile it yourself, following the [instructions in the wiki](https://github.com/fyde/fyde-cli/wiki/Compiling-from-source).

## Usage

When run without arguments, fyde-cli presents a list of available commands.
It will also show where it is going to save and look for the configuration files, unless overridden.

To use the client with an endpoint other than the default, you should start by setting the endpoint:

```
$ fyde-cli endpoint set fydeconsole.example.com
Endpoint changed to fydeconsole.example.com.
Credentials cleared, please login again using `fyde-cli login`
```

See [this page](https://github.com/fyde/fyde-cli/wiki/Working-with-different-MC-endpoints) for more information about using different endpoints.

You can then proceed to log in with your console credentials:

```
$ fyde-cli login
Email address: you@example.com
Password:
Logged in successfully, access token stored in (...)fyde/fyde-cli/auth.yaml
```

You can now use other commands. For example, to list users, you can use `fyde-cli users list`.

All commands provide a help text with the available subcommands and flags.
For example, running `fyde-cli resources` will let you know about the `get`, `list`, `add`, `edit` and `delete` subcommands, and `fyde-cli resources list --help` will list all available flags for the list resources command, including pagination, sorting and filtering flags.

### Output formats

fyde-cli supports different output formats for different use cases:

 - Table, for interactive usage (`--output=table`)
 - CSV (`--output=csv`)
 - JSON (`--output=json` or `--output=json-pretty`)

By default, when an interactive terminal is detected, `table` output is used.
Otherwise, `json` is used.
JSON output generally contains the most information, sometimes including nested objects; CSV output corresponds to a CSV version of the table output.

All output formats are subject to pagination parameters, when those are available.

Additional output options are available for record creation and editing commands:
 - `--errors-only` - output will be restricted to records whose creation/editing failed

### Input formats

When adding or editing records, fyde-cli can receive input in three different ways:

 - Interactively, through command line flags
   - Users will be prompted to interactively provide mandatory fields that were not included in the passed flags
 - From JSON files, using `--from-file=filename.json --file-format=json` (`--file-format=json` is the default and can be omitted)
   - In this case, fyde-cli will expect a JSON array containing the different records
 - From CSV files, using `--from-file=filename.csv --file-format=csv`
   - In this case, fyde-cli will expect a file containing comma-separated values, with one record per line. The first record must be a header mapping each column to the correct record field

The expected formats when using JSON and CSV files are documented in [the fyde-cli wiki](https://github.com/fyde/fyde-cli/wiki#batch-mode-operations).

### Behavior on error

When creating, editing or deleting multiple records in one go, by default fyde-cli will stop on the first error.
However, one may want to perform the operation in a "best effort" basis, where fyde-cli will continue processing the remaining records/arguments regardless of previous server-issued errors.
This can be enabled using the `--continue-on-error` flag.
When this flag is passed, fyde-cli never exits with a non-zero code, as long as the input is correctly formatted and all errors come from server-side operations.

## Reporting issues

You can see existing issues and report new ones [on GitHub](https://github.com/fyde/fyde-cli/issues).

## License

fyde-cli is Copyright © 2019 Fyde, Inc. and is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), a OSI-approved license.