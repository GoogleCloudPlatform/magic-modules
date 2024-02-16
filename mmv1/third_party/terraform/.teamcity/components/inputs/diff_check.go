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
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	repo := os.Args
	services := fmt.Sprintf("../../../%v/services/...", repo[1])
	fmt.Println(services)
	cmd := exec.Command("go", "list", services)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
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

	service_file := os.Args[2]
	f, err := os.Open(service_file)
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
		fmt.Fprintf(os.Stderr, "error: diff in services_ga.kt\n")
		os.Exit(1)
	}
}
