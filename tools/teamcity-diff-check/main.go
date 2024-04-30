package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

var serviceFile = flag.String("service_file", "services_ga", "kotlin service file to be parsed")

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

	file, err := os.Open(*serviceFile + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	googleServices := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		googleServices = append(googleServices, scanner.Text())
	}

	////////////////////////////////////////////////////////////////////////////////

	filePath := fmt.Sprintf("mmv1/third_party/terraform/.teamcity/components/inputs/%s.kt", *serviceFile)
	f, err := os.Open(fmt.Sprintf("../../%s", filePath)) // Need to make path relative to where the script is called
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
	pattern := regexp.MustCompile(`(?m)"(?P<service>\w+)"\sto\s+mapOf`)

	template := []byte("$service")

	dst := []byte{}
	teamcityServices := []string{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllSubmatchIndex(bs, -1) {
		service := pattern.Expand(dst, template, bs, submatches)
		teamcityServices = append(teamcityServices, string(service))
	}
	if len(teamcityServices) == 0 {
		fmt.Fprintf(os.Stderr, "error: script could not find any services listed in the file %s.kt .\n", filePath)
		os.Exit(1)
	}

	if diff := serviceDifference(googleServices, teamcityServices); len(diff) != 0 {
		fmt.Fprintf(os.Stderr, "error: missing services detected in %s\n", filePath)
		fmt.Fprintf(os.Stderr, "Please update file to include these new services: %s\n", diff)
		os.Exit(1)
	}

}
