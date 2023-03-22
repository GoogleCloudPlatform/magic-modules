package main

import (
	"flag"
	"fmt"
	"strings"

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
	if len(missingTests) > 0 {
		fmt.Println("## Missing test report\nYour PR includes resource fields which are not covered by any test.")
		for resourceName, missingTestInfo := range missingTests {
			fmt.Printf("\nResource: `%s` (%d total tests)\n", resourceName, len(missingTestInfo.Tests))
			glog.Infof("%s tests parsed: %v", resourceName, missingTestInfo.Tests)
			if len(missingTestInfo.UntestedFields) > 0 {
				fmt.Printf("Untested fields: %s\n", strings.Join(missingTestInfo.UntestedFields, ", "))
			}
		}
		fmt.Println("\nPlease add acceptance tests which include these fields.")
	}
}
