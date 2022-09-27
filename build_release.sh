#!/bin/bash

platform=$1

if [ -z "$platform" ]; then
  platform="linux"
fi

echo "Building for platform $platform"

if [ "$platform" = 'mac' ]; then
  GOOS=darwin go build -a -installsuffix cgo -ldflags "-s -w" -o oogway cmd/main.go || exit 1
elif [ "$platform" = 'windows' ]; then
  GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o oogway.exe cmd/main.go || exit 1
else
  GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o oogway cmd/main.go || exit 1
fi

echo "Done!"
