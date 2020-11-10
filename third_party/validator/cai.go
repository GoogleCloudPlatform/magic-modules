package google

import (
	"fmt"
	"math/rand"
	"regexp"
)

// Asset is the CAI representation of a resource.
type Asset struct {
	// The name, in a peculiar format: `\\<api>.googleapis.com/<self_link>`
	Name string `json:"name"`
	// The type name in `google.<api>.<resourcename>` format.
	Type      string         `json:"asset_type"`
	Resource  *AssetResource `json:"resource,omitempty"`
	IAMPolicy *IAMPolicy     `json:"iam_policy,omitempty"`
	OrgPolicy []*OrgPolicy   `json:"org_policy,omitempty"`
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

type OrgPolicy struct {
	Constraint     string                 `json:"constraint,omitempty"`
	ListPolicy     *Policy_ListPolicy     `json:"listPolicy"`
	BooleanPolicy  *Policy_BooleanPolicy  `json:"booleanPolicy"`
	RestoreDefault *Policy_RestoreDefault `json:"restoreDefault"`
	UpdateTime     *Timestamp             `json:"update_time,omitempty"`
}

type Timestamp struct {
	Seconds int64 `json:"seconds,omitempty"`
	Nanos   int64 `json:"nanos,omitempty"`
}

type Policy_ListPolicy_AllValues int32

type Policy_ListPolicy struct {
	AllowedValues     []string                    `json:"allowed_values,omitempty"`
	DeniedValues      []string                    `json:"denied_values,omitempty"`
	AllValues         Policy_ListPolicy_AllValues `json:"all_values,omitempty"`
	SuggestedValue    string                      `json:"suggested_value,omitempty"`
	InheritFromParent bool                        `json:"inherit_from_parent,omitempty"`
}

type Policy_BooleanPolicy struct {
	Enforced bool `json:"enforced,omitempty"`
}

type Policy_RestoreDefault struct {
}

// assetName templates an asset.name by looking up and replacing all instances
// of {{field}}. In the case where a field would resolve to an empty string, a
// generated unique string will be used: "placeholder-" + randomString().
// This is done to preserve uniqueness of asset.name for a given asset.asset_type.
func assetName(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")

	f, err := buildReplacementFunc(re, d, config, linkTmpl)
	if err != nil {
		return "", err
	}

	fWithPlaceholder := func(key string) string {
		val := f(key)
		if val == "" {
			val = fmt.Sprintf("placeholder-%s", randString(8))
		}
		return val
	}

	return re.ReplaceAllStringFunc(linkTmpl, fWithPlaceholder), nil
}

func randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
