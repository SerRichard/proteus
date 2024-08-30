#!/usr/bin/env bash
FILES=$(go list ./cli)
exec go build -o proteus $FILES