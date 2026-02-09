// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package biglakeiceberg

import (
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var icebergNamespaceIgnoredProperties = []string{
	"location",
}

func isIgnoredProperty(k string) bool {
	for _, p := range icebergNamespaceIgnoredProperties {
		if k == p {
			return true
		}
	}
	return false
}

func icebergNamespacePropertiesDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// properties.KEY
	parts := strings.Split(k, ".")
	if len(parts) == 2 && isIgnoredProperty(parts[1]) {
		return true
	}
	return false
}

