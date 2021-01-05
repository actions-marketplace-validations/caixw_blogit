#!/bin/sh

# SPDX-License-Identifier: MIT

cd `dirname $0`
builddate=`date -u '+%Y%m%d'`
commithash=`git rev-parse HEAD`
go build -ldflags "-X main.buildDate=${builddate}  -X main.commitHash=${commithash}" -v -o ./cmd/blogit/blogit ./cmd/blogit
