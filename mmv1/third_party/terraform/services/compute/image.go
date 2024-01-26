package compute

import (
	"fmt"
	"regexp"

	tpg_services_compute "github.com/hashicorp/terraform-provider-google-beta/google-beta/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/googleapi"
)

const (
	resolveImageFamilyRegex = "[-_a-zA-Z0-9]*"
	resolveImageImageRegex  = "[-_a-zA-Z0-9]*"
)

var (
	resolveImageProjectImage           = regexp.MustCompile(fmt.Sprintf("projects/(%s)/global/images/(%s)$", verify.ProjectRegex, resolveImageImageRegex))
	resolveImageProjectFamily          = regexp.MustCompile(fmt.Sprintf("projects/(%s)/global/images/family/(%s)$", verify.ProjectRegex, resolveImageFamilyRegex))
	resolveImageGlobalImage            = regexp.MustCompile(fmt.Sprintf("^global/images/(%s)$", resolveImageImageRegex))
	resolveImageGlobalFamily           = regexp.MustCompile(fmt.Sprintf("^global/images/family/(%s)$", resolveImageFamilyRegex))
	resolveImageFamilyFamily           = regexp.MustCompile(fmt.Sprintf("^family/(%s)$", resolveImageFamilyRegex))
	resolveImageProjectImageShorthand  = regexp.MustCompile(fmt.Sprintf("^(%s)/(%s)$", verify.ProjectRegex, resolveImageImageRegex))
	resolveImageProjectFamilyShorthand = regexp.MustCompile(fmt.Sprintf("^(%s)/(%s)$", verify.ProjectRegex, resolveImageFamilyRegex))
	resolveImageFamily                 = regexp.MustCompile(fmt.Sprintf("^(%s)$", resolveImageFamilyRegex))
	resolveImageImage                  = regexp.MustCompile(fmt.Sprintf("^(%s)$", resolveImageImageRegex))
	resolveImageLink                   = regexp.MustCompile(fmt.Sprintf("^https://www.googleapis.com/compute/[a-z0-9]+/projects/(%s)/global/images/(%s)", verify.ProjectRegex, resolveImageImageRegex))

	windowsSqlImage         = regexp.MustCompile("^sql-(?:server-)?([0-9]{4})-([a-z]+)-windows-(?:server-)?([0-9]{4})(?:-r([0-9]+))?-dc-v[0-9]+$")
	canonicalUbuntuLtsImage = regexp.MustCompile("^ubuntu-(minimal-)?([0-9]+)(?:.*(arm64))?.*$")
	cosLtsImage             = regexp.MustCompile("^cos-([0-9]+)-")
)

// built-in projects to look for images/families containing the string
// on the left in
var ImageMap = map[string]string{
	"centos":      "centos-cloud",
	"coreos":      "coreos-cloud",
	"debian":      "debian-cloud",
	"opensuse":    "opensuse-cloud",
	"rhel":        "rhel-cloud",
	"rocky-linux": "rocky-linux-cloud",
	"sles":        "suse-cloud",
	"ubuntu":      "ubuntu-os-cloud",
	"windows":     "windows-cloud",
	"windows-sql": "windows-sql-cloud",
}

func resolveImageImageExists(c *transport_tpg.Config, project, name, userAgent string) (bool, error) {
	if _, err := c.NewComputeClient(userAgent).Images.Get(project, name).Do(); err == nil {
		return true, nil
	} else if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
		return false, nil
	} else {
		return false, fmt.Errorf("Error checking if image %s exists: %s", name, err)
	}
}

func resolveImageFamilyExists(c *transport_tpg.Config, project, name, userAgent string) (bool, error) {
	if _, err := c.NewComputeClient(userAgent).Images.GetFromFamily(project, name).Do(); err == nil {
		return true, nil
	} else if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
		return false, nil
	} else {
		return false, fmt.Errorf("Error checking if family %s exists: %s", name, err)
	}
}

func sanityTestRegexMatches(expected int, got []string, regexType, name string) error {
	if len(got)-1 != expected { // subtract one, index zero is the entire matched expression
		return fmt.Errorf("Expected %d %s regex matches, got %d for %s", expected, regexType, len(got)-1, name)
	}
	return nil
}

// ResolveImage is a wrapper around compute.ResolveImage and recover from panic
// if client is il.
func ResolveImage(c *transport_tpg.Config, project, name, userAgent string) (ret string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if c.Client == nil {
				ret = name
				err = nil
				return
			}
			panic(r)
		}
	}()

	return tpg_services_compute.ResolveImage(c, project, name, userAgent)
}

// resolveImageRefToRelativeURI takes the output of ResolveImage and coerces it
// into a relative URI. In the event that a global/images/IMAGE or
// global/images/family/FAMILY reference is returned from ResolveImage,
// providerProject will be used as the project for the self_link.
func resolveImageRefToRelativeURI(providerProject, name string) (string, error) {
	switch {
	case resolveImageLink.MatchString(name): // https://www.googleapis.com/compute/v1/projects/xyz/global/images/xyz
		namePath, err := tpgresource.GetRelativePath(name)
		if err != nil {
			return "", err
		}

		return namePath, nil
	case resolveImageProjectImage.MatchString(name): // projects/xyz/global/images/xyz
		return name, nil
	case resolveImageProjectFamily.MatchString(name): // projects/xyz/global/images/family/xyz
		res := resolveImageProjectFamily.FindStringSubmatch(name)
		if err := sanityTestRegexMatches(2, res, "project family", name); err != nil {
			return "", err
		}
		return fmt.Sprintf("projects/%s/global/images/family/%s", res[1], res[2]), nil
	case resolveImageGlobalImage.MatchString(name): // global/images/xyz
		res := resolveImageGlobalImage.FindStringSubmatch(name)
		if err := sanityTestRegexMatches(1, res, "global image", name); err != nil {
			return "", err
		}
		return fmt.Sprintf("projects/%s/global/images/%s", providerProject, res[1]), nil
	case resolveImageGlobalFamily.MatchString(name): // global/images/family/xyz
		res := resolveImageGlobalFamily.FindStringSubmatch(name)
		if err := sanityTestRegexMatches(1, res, "global family", name); err != nil {
			return "", err
		}
		return fmt.Sprintf("projects/%s/global/images/family/%s", providerProject, res[1]), nil
	}
	return "", fmt.Errorf("Could not expand image or family %q into a relative URI", name)

}
