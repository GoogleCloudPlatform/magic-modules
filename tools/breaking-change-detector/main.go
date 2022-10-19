package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/GoogleCloudPlatform/magic-modules/.ci/breaking-change-detector/docs"
	"github.com/golang/glog"
)

var docMode = flag.Bool("docs", false, "Switches the mode from running the comparison to creating a markdown file detailing the breaking change rules")
var providerFolder = flag.String("providerFolder", "", "The location of the provider folder to output documentation into.. if not provided the documentation will be output to console")
var providerVersion = flag.String("providerVersion", "google-beta", "The version of provider used, needed for documentation.")

func main() {
	flag.Parse()
	validateParameters()

	if *docMode {
		docs.Generate(*providerFolder)
	} else {
		breakages := compare()
		sort.Strings(breakages)
		for _, breakage := range breakages {
			fmt.Println(breakage)
		}
	}

}

func validateParameters() {
	if *providerFolder != "" && !*docMode {
		glog.Exitln("parameter -docs must be set when specifying -providerFolder")
	}
	if *providerVersion != "google" && *providerVersion != "google-beta" {
		glog.Exitln("only google and google-beta are supported provider versions")
	}
}
