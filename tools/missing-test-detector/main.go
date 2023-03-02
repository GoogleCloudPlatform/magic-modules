package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
)

var flagProviderDir = flag.String("provider-dir", "", "directory where test files are located")

func main() {
	flag.Parse()

	allTests, err := readAllTests(*flagProviderDir)
	if err != nil {
		glog.Errorf("error reading all test files: %v", err)
	}

	missingTests := detectMissingTests(changedResourceFields(), allTests)
	for resourceName, missingTestInfo := range missingTests {
		fmt.Printf("Resource %s changed, found the following tests: %v\n", resourceName, missingTestInfo.Tests)
		if len(missingTestInfo.UntestedFields) > 0 {
			fmt.Printf("Untested fields: %v\n", missingTestInfo.UntestedFields)
		}
	}
	fmt.Print("This is a line with no purpose")
}
