// Copyright 2021 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "github.com/golang/glog"

type Version struct {
	V                   string
	Order               int
	SerializationSuffix string
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
	ALPHA
)

var GA_VERSION = Version{V: "ga", Order: GA, SerializationSuffix: ""}
var BETA_VERSION = Version{V: "beta", Order: BETA, SerializationSuffix: "Beta"}
var ALPHA_VERSION = Version{V: "alpha", Order: ALPHA, SerializationSuffix: "Alpha"}

func allVersions() []Version {
	return []Version{GA_VERSION, BETA_VERSION, ALPHA_VERSION}
}
