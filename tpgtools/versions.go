package main

import "github.com/golang/glog"

type Version struct {
	V     string
	Order int
}

func fromString(v string) *Version {
	for _, version := range allVersions() {
		if v == version.V {
			return &version
		}
	}
	glog.Infof("Failed finding version: %s", v)
	return nil
}

type VersionOrder int

const (
	GA = iota
	BETA
)

var GA_VERSION = Version{V: "ga", Order: GA}
var BETA_VERSION = Version{V: "beta", Order: BETA}

func allVersions() []Version {
	return []Version{GA_VERSION, BETA_VERSION}
}
