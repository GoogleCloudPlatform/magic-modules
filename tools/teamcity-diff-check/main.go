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
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

var serviceFile = flag.String("service_file", "services_ga.kt", "kotlin service file to be parsed")
var provider = flag.String("provider", "google", "Specify which provider to run diff_check on")

func serviceDifference(gS, tS []string) []string {
	g := make(map[string]struct{}, len(gS))
	for _, s := range gS {
		g[s] = struct{}{}
	}

	var diff []string
	for _, s := range tS {
		if _, found := g[s]; !found {
			diff = append(diff, s)
		}
	}

	return diff
}

func main() {
	flag.Parse()

	servicesPath := fmt.Sprintf("../../provider/%s/services/", *provider)
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = servicesPath
	stdout, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println(string(stdout))
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// Regex pattern captures "services" from `go list ../../provider/{{*provider}}/services/..`
	pattern := regexp.MustCompile(`github\.com\/hashicorp\/terraform-provider-(google|google-beta)\/(google|google-beta)\/services\/(?P<service>\w+)`)

	template := []byte("$service")
	dst := []byte{}

	googleServices := []string{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllSubmatchIndex(stdout, -1) {
		service := pattern.Expand(dst, template, stdout, submatches)
		googleServices = append(googleServices, string(service))
	}

	////////////////////////////////////////////////////////////////////////////////

	f, err := os.Open(fmt.Sprintf("../../mmv1/third_party/terraform/.teamcity/components/inputs/%s", *serviceFile))
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

	// Regex pattern captures "services" from *serviceFile.
	pattern = regexp.MustCompile(`(?m)"(?P<service>\w+)"\sto\s+mapOf`)

	template = []byte("$service")

	dst = []byte{}
	teamcityServices := []string{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllSubmatchIndex(bs, -1) {
		service := pattern.Expand(dst, template, bs, submatches)
		teamcityServices = append(teamcityServices, string(service))
	}

	if diff := serviceDifference(googleServices, teamcityServices); len(diff) != 0 {
		fmt.Fprintf(os.Stderr, "error: diff in %s\n", *serviceFile)
		fmt.Fprintf(os.Stderr, "Missing Services: %s\n", diff)
		os.Exit(1)
	}

	fmt.Printf("No Diff in %s provider\n", *provider)
}
