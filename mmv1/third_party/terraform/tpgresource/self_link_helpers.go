package tpgresource

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResourceNameFromSelfLink(link string) string {
	parts := strings.Split(link, "/")
	return parts[len(parts)-1]
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

// Hash the relative path of a self link.
func SelfLinkRelativePathHash(selfLink interface{}) int {
	path, _ := GetRelativePath(selfLink.(string))
	return Hashcode(path)
}

func GetRelativePath(selfLink string) (string, error) {
	stringParts := strings.SplitAfterN(selfLink, "projects/", 2)
	if len(stringParts) != 2 {
		return "", fmt.Errorf("String was not a self link: %s", selfLink)
	}

	return "projects/" + stringParts[1], nil
}

func ConvertSelfLinkToV1(link string) string {
	reg := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/")
	return reg.ReplaceAllString(link, "/compute/v1/projects/")
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
