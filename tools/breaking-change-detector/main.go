package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/golang/glog"
)

var providerVersion = flag.String("providerVersion", "google-beta", "The version of provider used, needed for documentation.")

func main() {
	flag.Parse()
	validateParameters()
	breakages := compare()
	sort.Strings(breakages)
	for _, breakage := range breakages {
		fmt.Println(breakage)
	}
}

func validateParameters() {
	if *providerVersion != "google" && *providerVersion != "google-beta" {
		glog.Exitln("only google and google-beta are supported provider versions")
	}
}
