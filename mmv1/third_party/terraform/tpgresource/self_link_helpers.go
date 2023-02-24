package tpgresource

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Compare only the resource name of two self links/paths.
func compareResourceNames(_, old, new string, _ *schema.ResourceData) bool {
	return GetResourceNameFromSelfLink(old) == GetResourceNameFromSelfLink(new)
}

// Compare only the relative path of two self links.
func CompareSelfLinkRelativePaths(_, old, new string, _ *schema.ResourceData) bool {
	oldStripped, err := GetRelativePath(old)
	if err != nil {
		return false
	}

	newStripped, err := GetRelativePath(new)
	if err != nil {
		return false
	}

	if oldStripped == newStripped {
		return true
	}

	return false
}

// CompareSelfLinkOrResourceName checks if two resources are the same resource
//
// Use this method when the field accepts either a name or a self_link referencing a resource.
// The value we store (i.e. `old` in this method), must be a self_link.
func CompareSelfLinkOrResourceName(_, old, new string, _ *schema.ResourceData) bool {
	newParts := strings.Split(new, "/")

	if len(newParts) == 1 {
		// `new` is a name
		// `old` is always a self_link
		if GetResourceNameFromSelfLink(old) == newParts[0] {
			return true
		}
	}

	// The `new` string is a self_link
	return CompareSelfLinkRelativePaths("", old, new, nil)
}

// Hash the relative path of a self link.
func selfLinkRelativePathHash(selfLink interface{}) int {
	path, _ := GetRelativePath(selfLink.(string))
	return hashcode(path)
}

func GetRelativePath(selfLink string) (string, error) {
	stringParts := strings.SplitAfterN(selfLink, "projects/", 2)
	if len(stringParts) != 2 {
		return "", fmt.Errorf("String was not a self link: %s", selfLink)
	}

	return "projects/" + stringParts[1], nil
}

// Hash the name path of a self link.
func selfLinkNameHash(selfLink interface{}) int {
	name := GetResourceNameFromSelfLink(selfLink.(string))
	return hashcode(name)
}

func ConvertSelfLinkToV1(link string) string {
	reg := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/")
	return reg.ReplaceAllString(link, "/compute/v1/projects/")
}

func GetResourceNameFromSelfLink(link string) string {
	parts := strings.Split(link, "/")
	return parts[len(parts)-1]
}

func NameFromSelfLinkStateFunc(v interface{}) string {
	return GetResourceNameFromSelfLink(v.(string))
}

func StoreResourceName(resourceLink interface{}) string {
	return GetResourceNameFromSelfLink(resourceLink.(string))
}

type LocationType int

const (
	Zonal LocationType = iota
	Regional
	Global
)

func GetZonalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Zonal)
}

func GetRegionalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Regional)
}

func getResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *Config, locationType LocationType) (string, string, string, error) {
	if selfLink, ok := d.GetOk("self_link"); ok {
		return GetLocationalResourcePropertiesFromSelfLinkString(selfLink.(string))
	} else {
		project, err := GetProject(d, config)
		if err != nil {
			return "", "", "", err
		}

		location := ""
		if locationType == Regional {
			location, err = GetRegion(d, config)
			if err != nil {
				return "", "", "", err
			}
		} else if locationType == Zonal {
			location, err = GetZone(d, config)
			if err != nil {
				return "", "", "", err
			}
		}

		n, ok := d.GetOk("name")
		name := n.(string)
		if !ok {
			return "", "", "", errors.New("must provide either `self_link` or `name`")
		}
		return project, location, name, nil
	}
}

// given a full locational (non-global) self link, returns the project + region/zone + name or an error
func GetLocationalResourcePropertiesFromSelfLinkString(selfLink string) (string, string, string, error) {
	parsed, err := url.Parse(selfLink)
	if err != nil {
		return "", "", "", err
	}

	s := strings.Split(parsed.Path, "/")

	// This is a pretty bad way to tell if this is a self link, but stops us
	// from accessing an index out of bounds and causing a panic. generally, we
	// expect bad values to be partial URIs and names, so this will catch them
	if len(s) < 9 {
		return "", "", "", fmt.Errorf("value %s was not a self link", selfLink)
	}

	return s[4], s[6], s[8], nil
}

// return the region a selfLink is referring to
func GetRegionFromRegionSelfLink(selfLink string) string {
	re := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/[a-zA-Z0-9-]*/regions/([a-zA-Z0-9-]*)")
	switch {
	case re.MatchString(selfLink):
		if res := re.FindStringSubmatch(selfLink); len(res) == 2 && res[1] != "" {
			return res[1]
		}
	}
	return selfLink
}

// This function supports selflinks that have regions and locations in their paths
func GetRegionFromRegionalSelfLink(selfLink string) string {
	re := regexp.MustCompile("projects/[a-zA-Z0-9-]*/(?:locations|regions)/([a-zA-Z0-9-]*)")
	switch {
	case re.MatchString(selfLink):
		if res := re.FindStringSubmatch(selfLink); len(res) == 2 && res[1] != "" {
			return res[1]
		}
	}
	return selfLink
}

func ReplaceVars(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	return replaceVarsRecursive(d, config, linkTmpl, false, 0)
}

// relaceVarsForId shortens variables by running them through GetResourceNameFromSelfLink
// this allows us to use long forms of variables from configs without needing
// custom id formats. For instance:
// accessPolicies/{{access_policy}}/accessLevels/{{access_level}}
// with values:
// access_policy: accessPolicies/foo
// access_level: accessPolicies/foo/accessLevels/bar
// becomes accessPolicies/foo/accessLevels/bar
func ReplaceVarsForId(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	return replaceVarsRecursive(d, config, linkTmpl, true, 0)
}

// ReplaceVars must be done recursively because there are baseUrls that can contain references to regions
// (eg cloudrun service) there aren't any cases known for 2+ recursion but we will track a run away
// substitution as 10+ calls to allow for future use cases.
func replaceVarsRecursive(d TerraformResourceData, config *Config, linkTmpl string, shorten bool, depth int) (string, error) {
	if depth > 10 {
		return "", errors.New("Recursive substitution detcted")
	}

	// https://github.com/google/re2/wiki/Syntax
	re := regexp.MustCompile("{{([%[:word:]]+)}}")
	f, err := buildReplacementFunc(re, d, config, linkTmpl, shorten)
	if err != nil {
		return "", err
	}
	final := re.ReplaceAllStringFunc(linkTmpl, f)

	if re.Match([]byte(final)) {
		return replaceVarsRecursive(d, config, final, shorten, depth+1)
	}

	return final, nil
}

// This function isn't a test of transport.go; instead, it is used as an alternative
// to ReplaceVars inside tests.
func ReplaceVarsForTest(config *Config, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string

	if strings.Contains(linkTmpl, "{{project}}") {
		project = rs.Primary.Attributes["project"]
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region = GetResourceNameFromSelfLink(rs.Primary.Attributes["region"])
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone = GetResourceNameFromSelfLink(rs.Primary.Attributes["zone"])
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}

		if v, ok := rs.Primary.Attributes[m]; ok {
			return v
		}

		// Attempt to draw values from the provider config
		if f := reflect.Indirect(reflect.ValueOf(config)).FieldByName(m); f.IsValid() {
			return f.String()
		}

		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}

// This function replaces references to Terraform properties (in the form of {{var}}) with their value in Terraform
// It also replaces {{project}}, {{project_id_or_project}}, {{region}}, and {{zone}} with their appropriate values
// This function supports URL-encoding the result by prepending '%' to the field name e.g. {{%var}}
func buildReplacementFunc(re *regexp.Regexp, d TerraformResourceData, config *Config, linkTmpl string, shorten bool) (func(string) string, error) {
	var project, projectID, region, zone string
	var err error

	if strings.Contains(linkTmpl, "{{project}}") {
		project, err = GetProject(d, config)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{project_id_or_project}}") {
		v, ok := d.GetOkExists("project_id")
		if ok {
			projectID, _ = v.(string)
		}
		if projectID == "" {
			project, err = GetProject(d, config)
		}
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region, err = GetRegion(d, config)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone, err = GetZone(d, config)
		if err != nil {
			return nil, err
		}
	}

	f := func(s string) string {

		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "project_id_or_project" {
			if projectID != "" {
				return projectID
			}
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}
		if string(m[0]) == "%" {
			v, ok := d.GetOkExists(m[1:])
			if ok {
				return url.PathEscape(fmt.Sprintf("%v", v))
			}
		} else {
			v, ok := d.GetOkExists(m)
			if ok {
				if shorten {
					return GetResourceNameFromSelfLink(fmt.Sprintf("%v", v))
				} else {
					return fmt.Sprintf("%v", v)
				}
			}
		}

		// terraform-google-conversion doesn't provide a provider config in tests.
		if config != nil {
			// Attempt to draw values from the provider config if it's present.
			if f := reflect.Indirect(reflect.ValueOf(config)).FieldByName(m); f.IsValid() {
				return f.String()
			}
		}
		return ""
	}

	return f, nil
}

// This function isn't a test of transport.go; instead, it is used as an alternative
// to ReplaceVars inside tests.
func ReplaceVarsForFrameworkTest(prov *FrameworkProvider, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string

	if strings.Contains(linkTmpl, "{{project}}") {
		project = rs.Primary.Attributes["project"]
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region = GetResourceNameFromSelfLink(rs.Primary.Attributes["region"])
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone = GetResourceNameFromSelfLink(rs.Primary.Attributes["zone"])
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}

		if v, ok := rs.Primary.Attributes[m]; ok {
			return v
		}

		// Attempt to draw values from the provider
		if f := reflect.Indirect(reflect.ValueOf(prov)).FieldByName(m); f.IsValid() {
			return f.String()
		}

		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}
