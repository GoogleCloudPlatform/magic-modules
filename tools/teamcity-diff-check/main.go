package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

var version = flag.String("version", "", "the provider version under test. Must be `ga` or `beta`")
var teamcityServiceFile = flag.String("teamcity_services", "", "path to a kotlin service file to be parsed")
var providerServiceFile = flag.String("provider_services", "", "path to a .txt file listing all service packages in the provider")

// listDifference checks that all the items in list B are present in list A
func listDifference(listA, listB []string) error {

	a := make(map[string]struct{}, len(listA))
	for _, s := range listA {
		a[s] = struct{}{}
	}
	var diff []string
	for _, s := range listB {
		if _, found := a[s]; !found {
			diff = append(diff, s)
		}
	}

	if len(diff) > 0 {
		return fmt.Errorf("%v", diff)
	}

	return nil
}

func main() {
	flag.Parse()

	ga := *version == "ga"
	beta := *version == "beta"
	if !ga && !beta {
		fmt.Fprint(os.Stderr, "the flag `version` must be set to either `ga` or `beta`, and is case sensitive\n")
		os.Exit(1)
	}

	err := compareServices(*teamcityServiceFile, *providerServiceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errors when inspecting the %s version of the Google provider\n", *version)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "All services present in the %s provider codebase are present in TeamCity config, and vice versa\n", *version)
}

// compareServices contains most of the logic of the main function, but is separated to make the code more testable
func compareServices(teamcityServiceFile, providerServiceFile string) error {

	// Get array of services from the provider service list file
	file, err := os.Open(providerServiceFile)
	if err != nil {
		return fmt.Errorf("error opening provider service list file: %w", err)
	}
	defer file.Close()

	googleServices := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		googleServices = append(googleServices, scanner.Text())
	}
	if len(googleServices) == 0 {
		return fmt.Errorf("could not find any services in the provider service list file %s", providerServiceFile)
	}

	// Get array of services from the TeamCity service list file
	f, err := os.Open(teamcityServiceFile)
	if err != nil {
		return fmt.Errorf("error opening TeamCity service list file: %w", err)
	}

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error stating TeamCity service list file: %w", err)
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(f).Read(bs)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error processing TeamCity service list file: %w", err)
	}

	// Regex pattern captures "services" from the Kotlin service list file.
	pattern := regexp.MustCompile(`(?m)"(?P<service>\w+)"\sto\s+mapOf`)

	template := []byte("$service")

	dst := []byte{}
	teamcityServices := []string{}

	for _, submatches := range pattern.FindAllSubmatchIndex(bs, -1) {
		service := pattern.Expand(dst, template, bs, submatches)
		teamcityServices = append(teamcityServices, string(service))
	}
	if len(teamcityServices) == 0 {
		return fmt.Errorf("could not find any services in the TeamCity service list file %s", teamcityServiceFile)
	}

	// Determine diffs
	errTeamCity := listDifference(teamcityServices, googleServices)
	errProvider := listDifference(googleServices, teamcityServices)

	switch {
	case errTeamCity != nil && errProvider != nil:
		return fmt.Errorf(`mismatches detected:
TeamCity service file is missing services present in the provider: %s
Provider codebase is missing services present in the TeamCity service file: %s`,
			errTeamCity, errProvider)
	case errTeamCity != nil:
		return fmt.Errorf(`mismatches detected:
TeamCity service file is missing services present in the provider: %s`,
			errTeamCity)
	case errProvider != nil:
		return fmt.Errorf(`mismatches detected:
Provider codebase is missing services present in the TeamCity service file: %s`,
			errProvider)
	}

	return nil
}
