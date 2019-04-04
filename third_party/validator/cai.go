package google

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// Asset is the CAI representation of a resource.
type Asset struct {
	// The name, in a peculiar format: `\\<api>.googleapis.com/<self_link>`
	Name string `json:"name"`
	// The type name in `google.<api>.<resourcename>` format.
	Type         string         `json:"asset_type"`
	AncestryPath string         `json:"ancestry_path"`
	Resource     *AssetResource `json:"resource,omitempty"`
	IAMPolicy    *IAMPolicy     `json:"iam_policy,omitempty"`
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
	// Resource parent, example: "//cloudresourcemanager.googleapis.com/projects/my-project-id"
	Parent string `json:"parent"`
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

// replaceWithPlaceholder templates fields like asset.name by looking up
// and replacing all instances of {{field}}.
// In the case where a field would resolve to an empty string, a
// generated unique string will be used: "placeholder-" + randomString().
// This is done to preserve uniqueness of asset.name for a given asset.asset_type.
func replaceWithPlaceholder(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
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

var projectAncestryCache = make(map[string]string)

// getProjectAncestry uses the resource manager API to get ancestry paths for
// projects. It implements a cache because many resources share the same
// project.
func getProjectAncestry(d TerraformResourceData, config *Config) (string, error) {
	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}

	if path, ok := projectAncestryCache[project]; ok {
		return path, nil
	}

	ancestry, err := config.clientResourceManager.Projects.GetAncestry(project, &cloudresourcemanager.GetAncestryRequest{}).Do()
	if err != nil {
		return "", err
	}

	path := ancestryPath(ancestry.Ancestor)
	projectAncestryCache[project] = path

	return path, nil
}

// ancestryPath composes a path containing organization/folder/project
// (i.e. "organization/my-org/project/my-prj").
func ancestryPath(as []*cloudresourcemanager.Ancestor) string {
	var path []string
	for i := len(as) - 1; i >= 0; i-- {
		path = append(path, as[i].ResourceId.Type, as[i].ResourceId.Id)
	}
	return strings.Join(path, "/")
}
