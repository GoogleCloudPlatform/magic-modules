// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/converters/google/resources"
	"github.com/spf13/cobra"
)

type listSupportedResourcesOptions struct{}

func newListSupportedResourcesCmd() *cobra.Command {
	o := listSupportedResourcesOptions{}

	cmd := &cobra.Command{
		Use:   "list-supported-resources",
		Short: "List supported terraform resources.",
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}

func (o *listSupportedResourcesOptions) run() error {
	converters := resources.ResourceConverters()

	// Get a sorted list of terraform resources
	list := make([]string, 0, len(converters))
	for k := range converters {
		list = append(list, k)
	}
	sort.Strings(list)

	// go through and print
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, resource := range list {
		seenAssetTypes := make(map[string]bool)
		assetTypes := []string{}
		for _, resources := range converters[resource] {
			if _, exists := seenAssetTypes[resources.AssetType]; !exists {
				seenAssetTypes[resources.AssetType] = true
				assetTypes = append(assetTypes, resources.AssetType)
			}
		}
		sort.Strings(assetTypes)
		fmt.Fprintf(w, "%s\t%s\n", resource, strings.Join(assetTypes, ","))
	}
	w.Flush()

	return nil
}
