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

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"text/template"

	"github.com/golang/glog"
)

// Merges beta and GA resources for doc generation for a particular resource.
func mergeResource(res *Resource, resources map[Version][]*Resource, version *Version) *Resource {
	resourceAcrossVersions := make(map[Version]*Resource)
	for v, resList := range resources {
		for _, r := range resList {
			// Name is not unique, TerraformName must be
			if r.TerraformName() == res.TerraformName() {
				resourceAcrossVersions[v] = r
			}
		}
	}
	ga, gaExists := resourceAcrossVersions[GA_VERSION]
	beta, betaExists := resourceAcrossVersions[BETA_VERSION]
	alpha, alphaExists := resourceAcrossVersions[ALPHA_VERSION]
	private, privateExists := resourceAcrossVersions[PRIVATE_VERSION]
	if privateExists {
		return private
	}
	if alphaExists {
		return alpha
	}
	if gaExists {
		if betaExists {
			return mergeResources(ga, beta)
		}
		return ga
	}
	beta.Description = fmt.Sprintf("Beta only: %s", beta.Description)
	return beta
}

func mergeResources(ga, beta *Resource) *Resource {
	beta.Properties = mergeProperties(ga.Properties, beta.Properties)

	return beta
}

// Marks any sub properties as beta only
func mergeProperties(ga, beta []Property) []Property {
	gaProps := make(map[string]Property)
	for _, p := range ga {
		gaProps[p.title] = p
	}
	betaProps := make(map[string]Property)
	for _, p := range beta {
		betaProps[p.title] = p
	}
	inOrder := make([]string, 0)
	for k, _ := range betaProps {
		inOrder = append(inOrder, k)
	}
	sort.Strings(inOrder)
	modifiedProps := make([]Property, 0)
	for _, name := range inOrder {
		v := betaProps[name]
		if gaProp, ok := gaProps[name]; !ok {
			v.Description = fmt.Sprintf("(Beta only) %s", v.Description)
		} else if len(v.Properties) != 0 {
			// Look for sub-properties that might be beta only.
			// If the top-level property is beta only, sub-properties don't need to be marked
			v.Properties = mergeProperties(gaProp.Properties, v.Properties)
		}
		modifiedProps = append(modifiedProps, v)
	}

	return modifiedProps
}

func generateResourceWebsiteFile(res *Resource, resources map[Version][]*Resource, version *Version) {
	res = mergeResource(res, resources, version)

	if len(res.DocSamples()) <= 0 {
		fmt.Printf(" %-40s no samples, skipping doc generation\n", res.TerraformName())
		return
	}

	// Generate resource website file
	tmplInput := ResourceInput{
		Resource: *res,
	}

	tmpl, err := template.New("resource.html.markdown.tmpl").Funcs(TemplateFunctions).ParseFiles(
		"templates/resource.html.markdown.tmpl",
	)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, "resource.html.markdown.tmpl", tmplInput); err != nil {
		glog.Exit(err)
	}

	source := contents.Bytes()

	if oPath == nil || *oPath == "" {
		fmt.Printf("%v\n", string(source))
	} else {
		outname := fmt.Sprintf("%s_%s.html.markdown", res.ProductName(), res.Name())
		err := ioutil.WriteFile(path.Join(*oPath, "website/docs/r", outname), source, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}
}
