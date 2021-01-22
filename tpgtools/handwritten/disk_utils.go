package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Is the new disk size smaller than the old one?
func isDiskShrinkage(old, new, _ interface{}) bool {
	// It's okay to remove size entirely.
	if old == nil || new == nil {
		return false
	}
	return new.(int) < old.(int)
}

// We cannot suppress the diff for the case when family name is not part of the image name since we can't
// make a network call in a DiffSuppressFunc.
func diskImageDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	// Understand that this function solves a messy problem ("how do we tell if the diff between two images
	// is 'ForceNew-worthy', without making a network call?") in the best way we can: through a series of special
	// cases and regexes.  If you find yourself here because you are trying to add a new special case,
	// you are probably looking for the diskImageFamilyEquals function and its subfunctions.
	// In order to keep this maintainable, we need to ensure that the positive and negative examples
	// in resource_compute_disk_test.go are as complete as possible.

	// 'old' is read from the API.
	// It always has the format 'https://www.googleapis.com/compute/v1/projects/(%s)/global/images/(%s)'
	matches := resolveImageLink.FindStringSubmatch(old)
	if matches == nil {
		// Image read from the API doesn't have the expected format. In practice, it should never happen
		return false
	}
	oldProject := matches[1]
	oldName := matches[2]

	// Partial or full self link family
	if resolveImageProjectFamily.MatchString(new) {
		// Value matches pattern "projects/{project}/global/images/family/{family-name}$"
		matches := resolveImageProjectFamily.FindStringSubmatch(new)
		newProject := matches[1]
		newFamilyName := matches[2]

		return diskImageProjectNameEquals(oldProject, newProject) && diskImageFamilyEquals(oldName, newFamilyName)
	}

	// Partial or full self link image
	if resolveImageProjectImage.MatchString(new) {
		// Value matches pattern "projects/{project}/global/images/{image-name}$"
		matches := resolveImageProjectImage.FindStringSubmatch(new)
		newProject := matches[1]
		newImageName := matches[2]

		return diskImageProjectNameEquals(oldProject, newProject) && diskImageEquals(oldName, newImageName)
	}

	// Partial link without project family
	if resolveImageGlobalFamily.MatchString(new) {
		// Value is "global/images/family/{family-name}"
		matches := resolveImageGlobalFamily.FindStringSubmatch(new)
		familyName := matches[1]

		return diskImageFamilyEquals(oldName, familyName)
	}

	// Partial link without project image
	if resolveImageGlobalImage.MatchString(new) {
		// Value is "global/images/{image-name}"
		matches := resolveImageGlobalImage.FindStringSubmatch(new)
		imageName := matches[1]

		return diskImageEquals(oldName, imageName)
	}

	// Family shorthand
	if resolveImageFamilyFamily.MatchString(new) {
		// Value is "family/{family-name}"
		matches := resolveImageFamilyFamily.FindStringSubmatch(new)
		familyName := matches[1]

		return diskImageFamilyEquals(oldName, familyName)
	}

	// Shorthand for image or family
	if resolveImageProjectImageShorthand.MatchString(new) {
		// Value is "{project}/{image-name}" or "{project}/{family-name}"
		matches := resolveImageProjectImageShorthand.FindStringSubmatch(new)
		newProject := matches[1]
		newName := matches[2]

		return diskImageProjectNameEquals(oldProject, newProject) &&
			(diskImageEquals(oldName, newName) || diskImageFamilyEquals(oldName, newName))
	}

	// Image or family only
	if diskImageEquals(oldName, new) || diskImageFamilyEquals(oldName, new) {
		// Value is "{image-name}" or "{family-name}"
		return true
	}

	return false
}

func diskImageProjectNameEquals(project1, project2 string) bool {
	// Convert short project name to full name
	// For instance, centos => centos-cloud
	fullProjectName, ok := imageMap[project2]
	if ok {
		project2 = fullProjectName
	}

	return project1 == project2
}

func diskImageEquals(oldImageName, newImageName string) bool {
	return oldImageName == newImageName
}

func diskImageFamilyEquals(imageName, familyName string) bool {
	// Handles the case when the image name includes the family name
	// e.g. image name: debian-9-drawfork-v20180109, family name: debian-9
	if strings.Contains(imageName, familyName) {
		return true
	}

	if suppressCanonicalFamilyDiff(imageName, familyName) {
		return true
	}

	if suppressWindowsSqlFamilyDiff(imageName, familyName) {
		return true
	}

	if suppressWindowsFamilyDiff(imageName, familyName) {
		return true
	}

	return false
}

// e.g. image: ubuntu-1404-trusty-v20180122, family: ubuntu-1404-lts
func suppressCanonicalFamilyDiff(imageName, familyName string) bool {
	parts := canonicalUbuntuLtsImage.FindStringSubmatch(imageName)
	if len(parts) == 3 {
		f := fmt.Sprintf("ubuntu-%s%s-lts", parts[1], parts[2])
		if f == familyName {
			return true
		}
	}

	return false
}

// e.g. image: sql-2017-standard-windows-2016-dc-v20180109, family: sql-std-2017-win-2016
// e.g. image: sql-2017-express-windows-2012-r2-dc-v20180109, family: sql-exp-2017-win-2012-r2
func suppressWindowsSqlFamilyDiff(imageName, familyName string) bool {
	parts := windowsSqlImage.FindStringSubmatch(imageName)
	if len(parts) == 5 {
		edition := parts[2] // enterprise, standard or web.
		sqlVersion := parts[1]
		windowsVersion := parts[3]

		// Translate edition
		switch edition {
		case "enterprise":
			edition = "ent"
		case "standard":
			edition = "std"
		case "express":
			edition = "exp"
		}

		var f string
		if revision := parts[4]; revision != "" {
			// With revision
			f = fmt.Sprintf("sql-%s-%s-win-%s-r%s", edition, sqlVersion, windowsVersion, revision)
		} else {
			// No revision
			f = fmt.Sprintf("sql-%s-%s-win-%s", edition, sqlVersion, windowsVersion)
		}

		if f == familyName {
			return true
		}
	}

	return false
}

// e.g. image: windows-server-1709-dc-core-v20180109, family: windows-1709-core
// e.g. image: windows-server-1709-dc-core-for-containers-v20180109, family: "windows-1709-core-for-containers
func suppressWindowsFamilyDiff(imageName, familyName string) bool {
	updatedFamilyString := strings.Replace(familyName, "windows-", "windows-server-", 1)
	updatedImageName := strings.Replace(imageName, "-dc-", "-", 1)

	return strings.Contains(updatedImageName, updatedFamilyString)
}

func expandComputeDiskType(v interface{}, d TerraformResourceData, config *Config) *string {
	if v == "" {
		return nil 
	}

	f, err := parseZonalFieldValue("diskTypes", v.(string), "project", "zone", d, config, true)
	if err != nil {
		return nil
	}

	rl := f.RelativeLink()
	return &rl
}

func expandComputeDiskSourceImage(v interface{}, d TerraformResourceData, config *Config) *string {
	if v == "" {
		return nil 
	}

	if v == nil {
		return nil
	}

	project, err := getProject(d, config)
	if err != nil {
		return nil
	}

	f, err := resolveImage(config, project, v.(string))
	if err != nil {
		return nil
	}

	return &f
}

func expandComputeDiskSnapshot(v interface{}, d TerraformResourceData, config *Config) *string {
	if v == "" {
		return nil 
	}

	f, err := parseGlobalFieldValue("snapshots", v.(string), "project", d, config, true)
	if err != nil {
		return nil
	}

	rl := f.RelativeLink()
	return &rl
}

func ConvertSelfLinkToV1UnlessNil(v interface{}) *string {
	if v == nil {
		return nil
	}
	vptr := v.(*string)

	if vptr == nil {
		return nil
	}
	
	val := ConvertSelfLinkToV1(*vptr)
	return &val
}

func flattenComputeDiskSnapshot(v interface{}, d *schema.ResourceData, meta interface{}) *string {
	config := meta.(*Config)
	if v == nil {
		return nil
	}

	vptr := v.(*string)
	if vptr == nil {
		return nil
	}

	val, err := parseGlobalFieldValue("snapshots", *vptr, "project", d, config, true)
	if err != nil {
		return nil
	}

	rl := val.RelativeLink()
	return &rl
}

func flattenComputeDiskImage(v interface{}, d *schema.ResourceData, meta interface{}) *string {
	config := meta.(*Config)
	if v == nil {
		return nil
	}
	project, err := getProject(d, config)
	if err != nil {
		return nil
	}

	vptr := v.(*string)
	if vptr == nil {
		return nil
	}

	f, err := resolveImage(config, project, *vptr)
	if err != nil {
		return nil
	}

	return &f
}
