#!/bin/bash

version=$1
echo "build windowns app for version : $version"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$version"
echo "build M1 app for version : $version"
go build -ldflags "-X main.version=$version" -o tssh-appleSilicon
