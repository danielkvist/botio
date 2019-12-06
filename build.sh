#!/usr/bin/env bash

main() {
	for GOOS in darwin linux windows; do
    	for GOARCH in amd64 arm arm64; do
        	go build -v -o botio-$GOOS-$GOARCH .
    	done
	done
}

main "$@"