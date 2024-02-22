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

	f, err := os.Open(fmt.Sprintf("../../mmv1/third_party/terraform/.teamcity/components/inputs/%s", *serviceFile+".kt"))
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
		fmt.Fprintf(os.Stderr, "teamcityServices error: regex produced no matches.\n")
		os.Exit(1)
	}

	if diff := serviceDifference(googleServices, teamcityServices); len(diff) != 0 {
		fmt.Fprintf(os.Stderr, "error: diff in %s\n", *serviceFile)
		fmt.Fprintf(os.Stderr, "Missing Services: %s\n", diff)
		os.Exit(1)
	}

}
