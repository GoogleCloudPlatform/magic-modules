package google

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

// Asset is the CAI representation of a resource.
type Asset struct {
	// The name, in a peculiar format: `\\<api>.googleapis.com/<self_link>`
	Name string `json:"name"`
	// The type name in `google.<api>.<resourcename>` format.
	Type      string         `json:"asset_type"`
	Resource  *AssetResource `json:"resource,omitempty"`
	IAMPolicy *IAMPolicy     `json:"iam_policy,omitempty"`
}

// AssetResource is the Asset's Resource field.
type AssetResource struct {
	// Api version
	Version string `json:"version"`
	// URI including scheme for the discovery doc - assembled from
	// product name and version.
	DiscoveryDocumentURI string `json:"discovery_document_uri"`
	// Resource name.
	DiscoveryName string `json:"discovery_name"`
	// Actual resource state as per Terraform.  Note that this does
	// not necessarily correspond perfectly with the CAI representation
	// as there are occasional deviations between CAI and API responses.
	// This returns the API response values instead.
	Data map[string]interface{} `json:"data,omitempty"`
}

type IAMPolicy struct {
	Bindings []IAMBinding `json:"bindings"`
}

type IAMBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

// assetName templates an asset.name by looking up and replacing all instances
// of {{field}}. In the case where a field would resolve to an empty string, a
// generated unique string will be used: "<placeholder-" + randomString() + ">".
// This is done to preserve uniqueness of asset.name for a given asset.asset_type.
func assetName(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string
	var err error

	if strings.Contains(linkTmpl, "{{project}}") {
		project, err = getProject(d, config)
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region, err = getRegion(d, config)
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone, err = getZone(d, config)
		if err != nil {
			return "", err
		}
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]

		var res string
		switch m {
		case "project":
			res = project
		case "region":
			res = region
		case "zone":
			res = zone
		default:
			v, ok := d.GetOk(m)
			if ok {
				res = fmt.Sprintf("%v", v)
			}
		}

		if res == "" {
			res = fmt.Sprintf("<placeholder-%s>", randString(8))
		}

		return res
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}

func randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
