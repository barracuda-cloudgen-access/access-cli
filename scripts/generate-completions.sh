#!/usr/bin/env bash

mkdir dist/completions
go run . completion bash >dist/completions/access-cli.bash
go run . completion zsh >dist/completions/access-cli.zsh
