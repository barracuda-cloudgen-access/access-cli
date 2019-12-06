#!/usr/bin/env bash

mkdir dist/completions
go run . completion bash >dist/completions/fyde-cli.bash
go run . completion zsh >dist/completions/fyde-cli.zsh
