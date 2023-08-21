// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"io"
	"log"
	"os"

	"github.com/hashicorp/go-changelog"
)

func main() {
	var err error

	// default to reading from stdin
	r := os.Stdin

	// read from file arg instead if provided
	// TODO: add --help text for [file] arg handling
	filepath := ""
	if len(os.Args) > 1 {
		filepath = os.Args[1]
		r, err = os.Open(filepath)
		if err != nil {
			log.Fatalf("error opening %s", filepath)
			os.Exit(1)
		}
	}

	b, err := io.ReadAll(r)
	if err != nil {
		if filepath != "" {
			log.Fatalf("error reading from %s", filepath)
		} else {
			log.Fatalf("error reading from stdin")
		}
		os.Exit(1)
	}

	entry := changelog.Entry{
		Body: string(b),
	}

	if err := entry.Validate(); err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}
