package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

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
		funcs := template.FuncMap{
			"join": strings.Join,
			"backTickAll": func(ss []string) []string {
				rs := make([]string, len(ss))
				for i, s := range ss {
					rs[i] = fmt.Sprintf("`%s`", s)
				}
				return rs
			},
		}
		outputTemplate, err := template.New("output.tmpl").Funcs(funcs).ParseFiles("output.tmpl")
		if err != nil {
			glog.Exitf("Error parsing missing test template file: %s", err)
		}
		if err := outputTemplate.Execute(os.Stdout, missingTests); err != nil {
			glog.Exitf("Error executing missing test output template: %s", err)
		}
		for resourceName, missingTestInfo := range missingTests {
			glog.Infof("%s tests parsed: %v", resourceName, missingTestInfo.Tests)
		}
	}
}
