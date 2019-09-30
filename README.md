Fyde CLI client
===============

This tool will allow access to all Enterprise Console APIs from the command line.
It will cover all object types accessible from the web version and also help to run certain APIs in batch mode where it makes sense (ie. add multiple users from a CSV file or a JSON dump).

## Generating client code from the Swagger specification

Install [go-swagger](https://github.com/go-swagger/go-swagger), then run the following command on the root of this repo:

`swagger generate client -f .\swagger.yml`

## Compiling

Simply run `go build`. This project uses modules, make sure the installed Golang version supports them (1.11 and up).