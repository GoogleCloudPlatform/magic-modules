// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is controlled by MMv1, any changes made here will be overwritten

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

var serviceFile = flag.String("service_file", "services_ga.kt", "kotlin service file to be parsed")
var provider = flag.String("provider", "google", "Specify which provider to run diff_check on")

func main() {
	flag.Parse()
	var providerPath string
	if *provider == "google" {
		providerPath = "tgp"
	} else {
		providerPath = "tgbp"
	}
	services := fmt.Sprintf("../../%s/%s/services/...", providerPath, *provider)
	cmd := exec.Command("go", "list", services)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	pattern := regexp.MustCompile(`github\.com\/hashicorp\/terraform-provider-(google|google-beta)\/(google|google-beta)\/services\/(?P<service>\w+)`)

	// Template to convert "key: value" to "key=value" by
	// referencing the values captured by the regex pattern.
	template := []byte("$service\n")

	googleServices := []byte{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllSubmatchIndex(stdout, -1) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		googleServices = pattern.Expand(googleServices, template, stdout, submatches)
	}

	////////////////////////////////////////////////////////////////////////////////

	f, err := os.Open(fmt.Sprintf("../../%s/.teamcity/components/inputs/%s", providerPath, *serviceFile))
	if err != nil {
		panic(err)
	}

	// Get the file size
	stat, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(f).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	// Regex pattern captures "key: value" pair from the content.
	pattern = regexp.MustCompile(`(?m)"(?P<service>\w+)"\sto\s+mapOf`)

	// Template to convert "key: value" to "key=value" by
	// referencing the values captured by the regex pattern.
	template = []byte("$service\n")

	teamcityServices := []byte{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllSubmatchIndex(bs, -1) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		teamcityServices = pattern.Expand(teamcityServices, template, bs, submatches)
	}

	if !bytes.Equal(googleServices, teamcityServices) {
		fmt.Fprintf(os.Stderr, "error: diff in %s\n", *serviceFile)
		os.Exit(1)
	}
}
