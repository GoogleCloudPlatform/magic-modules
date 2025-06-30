package cai

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

// AssetName templates an asset.name by looking up and replacing all instances
// of {{field}}. In the case where a field would resolve to an empty string, "null" will be used.
func AssetName(d tpgresource.TerraformResourceData, config *transport_tpg.Config, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([%[:word:]]+)}}")
	f, err := tpgresource.BuildReplacementFunc(re, d, config, linkTmpl, false)
	if err != nil {
		return "", err
	}

	fWithPlaceholder := func(key string) string {
		val := f(key)
		if val == "" {
			val = "null"
		}
		return val
	}

	return re.ReplaceAllStringFunc(linkTmpl, fWithPlaceholder), nil
}
