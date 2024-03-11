package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"missing-test-detector/reader"

	"github.com/golang/glog"
)

var flagServicesDir = flag.String("services-dir", "", "directory where service directories are located")

func main() {
	flag.Parse()

	allTests, errs := reader.ReadAllTests(*flagServicesDir)
	for path, err := range errs {
		glog.Infof("error reading path: %s, err: %v", path, err)
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
