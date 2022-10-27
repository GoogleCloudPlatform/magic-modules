package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
)

var flagProviderDir = flag.String("provider-dir", "", "directory where test files are located")

func main() {
	flag.Parse()

	missingTests, err := detectMissingTests(changedResourceFields(), *flagProviderDir)
	if err != nil {
		glog.Errorf("error detecting missing tests: %v", err)
	}
	for resourceName, missingTestInfo := range missingTests {
		fmt.Printf("Resource %s changed, found %d tests with %d total steps\n", resourceName, missingTestInfo.TestCount, missingTestInfo.StepCount)
		if len(missingTestInfo.UntestedFields) > 0 {
			fmt.Printf("Untested fields: %v\n", missingTestInfo.UntestedFields)
		}
	}
}
