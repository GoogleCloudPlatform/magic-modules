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
	t := make(map[string]struct{}, len(tS))
	for _, s := range tS {
		t[s] = struct{}{}
	}

	var diff []string
	for _, s := range gS {
		if _, found := t[s]; !found {
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
	if err != nil {
		fmt.Println(string(stdout))
		fmt.Println(err.Error())
		return
	}

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
	if len(googleServices) == 0 {
		fmt.Fprintf(os.Stderr, "googleServices error: regex produced no matches.\n")
		os.Exit(1)
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
	if len(teamcityServices) == 0 {
		fmt.Fprintf(os.Stderr, "teamcityServices error: regex produced no matches.\n")
		os.Exit(1)
	}

	if diff := serviceDifference(googleServices, teamcityServices); len(diff) != 0 {
		fmt.Fprintf(os.Stderr, "error: diff in %s\n", *serviceFile)
		fmt.Fprintf(os.Stderr, "Missing Services: %s\n", diff)
		os.Exit(1)
	}

	fmt.Printf("No Diff in %s provider\n", *provider)
}
