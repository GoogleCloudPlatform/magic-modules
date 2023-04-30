package google

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Compare only the resource name of two self links/paths.
func compareResourceNames(_, old, new string, _ *schema.ResourceData) bool {
	return GetResourceNameFromSelfLink(old) == GetResourceNameFromSelfLink(new)
}

// Compare only the relative path of two self links.
func compareSelfLinkRelativePaths(_, old, new string, _ *schema.ResourceData) bool {
	return tpgresource.CompareSelfLinkRelativePaths("", old, new, nil)
}

// compareSelfLinkOrResourceName checks if two resources are the same resource
//
// Use this method when the field accepts either a name or a self_link referencing a resource.
// The value we store (i.e. `old` in this method), must be a self_link.
func compareSelfLinkOrResourceName(_, old, new string, _ *schema.ResourceData) bool {
	newParts := strings.Split(new, "/")

	if len(newParts) == 1 {
		// `new` is a name
		// `old` is always a self_link
		if GetResourceNameFromSelfLink(old) == newParts[0] {
			return true
		}
	}

	// The `new` string is a self_link
	return compareSelfLinkRelativePaths("", old, new, nil)
}

// Hash the relative path of a self link.
func selfLinkRelativePathHash(selfLink interface{}) int {
	return tpgresource.SelfLinkRelativePathHash(selfLink)
}

func getRelativePath(selfLink string) (string, error) {
	return tpgresource.GetRelativePath(selfLink)
}

// Hash the name path of a self link.
func selfLinkNameHash(selfLink interface{}) int {
	name := GetResourceNameFromSelfLink(selfLink.(string))
	return tpgresource.Hashcode(name)
}

func ConvertSelfLinkToV1(link string) string {
	return tpgresource.ConvertSelfLinkToV1(link)
}

func GetResourceNameFromSelfLink(link string) string {
	return tpgresource.GetResourceNameFromSelfLink(link)
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

func GetZonalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Zonal)
}

func GetRegionalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Regional)
}

func getResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config, locationType LocationType) (string, string, string, error) {
	if selfLink, ok := d.GetOk("self_link"); ok {
		return GetLocationalResourcePropertiesFromSelfLinkString(selfLink.(string))
	} else {
		project, err := getProject(d, config)
		if err != nil {
			return "", "", "", err
		}

		location := ""
		if locationType == Regional {
			location, err = getRegion(d, config)
			if err != nil {
				return "", "", "", err
			}
		} else if locationType == Zonal {
			location, err = getZone(d, config)
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
	return tpgresource.GetLocationalResourcePropertiesFromSelfLinkString(selfLink)
}

// This function supports selflinks that have regions and locations in their paths
func GetRegionFromRegionalSelfLink(selfLink string) string {
	return tpgresource.GetRegionFromRegionalSelfLink(selfLink)
}
