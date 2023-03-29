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

	changedFields := changedResourceFields()

	missingTests, err := detectMissingTests(changedFields, allTests)
	if err != nil {
		glog.Errorf("error detecting missing tests: %v", err)
	}
	for resourceName, missingTestInfo := range missingTests {
		fmt.Printf("Resource %s changed\n", resourceName)
		glog.Infof("Tests parsed: %v", missingTestInfo.Tests)
		if len(missingTestInfo.UntestedFields) > 0 {
			fmt.Printf("Untested fields: %v\n", missingTestInfo.UntestedFields)
		}
	}
}
