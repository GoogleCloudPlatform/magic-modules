package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/golang/glog"
)

// Samples is the array of samples read from the DCL
type Samples []Sample

// Sample is the object containing sample data from DCL samples
type Sample struct {
	// Name is the name of a sample
	Name *string

	// Description is a short description of the sample
	Description *string

	// DependencyFileNames contains the filenames of every resource in the sample
	DependencyFileNames []string `yaml:"dependencies"`

	// PrimaryResource is the filename of the sample's primary resource
	PrimaryResource *string `yaml:"resource"`

	// Substitutions contains every substition in the sample
	Substitutions []Substitution `yaml:"extra_substitutions"`

	// DependencyList is a list of objects containing metadata for each sample resource
	DependencyList []Dependency

	TestSlug string

	Versions []string

	Updates []map[string]string
}

// Substitution contains metadata that varies for the sample context
type Substitution struct {
	// Substitution is the text to be substituted, e.g. topic_name
	Substitution *string

	// TestValue is the value of the substituted text for tests
	TestValue *string `yaml:"test_value"`

	// TestValue is the value of the substituted text for docs
	DocsValue *string `yaml:"docs_value"`
}

// Dependency contains data that describes a single resource in a sample
type Dependency struct {
	// FileName is the name of the file as it appears in testcases.yaml
	FileName string

	// HCLLocalName is the local name of the HCL block, e.g. "basic" or "default"
	HCLLocalName string

	// DCLResourceType is the type represented in the DCL, e.g. "ComputeInstance"
	DCLResourceType string

	// TerraformResourceType is the type represented in Terraform, e.g. "google_compute_instance"
	TerraformResourceType string

	// HCLBlock is the snippet of HCL config that declares this resource
	HCLBlock string
}

// BuildDependency produces a Dependency using a file and filename
func BuildDependency(fileName string, v Version, b []byte) (*Dependency, error) {
	fileParts := strings.Split(fileName, ".")
	if len(fileParts) != 4 {
		return nil, fmt.Errorf("Invalid sample dependency file name: %s", fileName)
	}

	localName := fileParts[0]
	dclResourceType := strings.Title(fileParts[1]) + strings.Title(fileParts[2])
	// terraformResourceType := fmt.Sprintf("google_%s_%s", fileParts[1], fileParts[2])
	terraformResourceType, err := DCLToTerraformReference(dclResourceType, v.V)
	if err != nil {
		return nil, fmt.Errorf("Error generating sample dependency %s: %s", fileName, err)
	}

	block, err := ConvertSampleJSONToHCL(dclResourceType, v.V, b)
	if err != nil {
		return nil, fmt.Errorf("Error generating sample dependency %s: %s", fileName, err)
	}

	re := regexp.MustCompile(`(resource "` + terraformResourceType + `" ")(\w*)`)
	block = re.ReplaceAllString(block, "${1}"+localName)

	d := &Dependency{
		FileName:              fileName,
		HCLLocalName:          localName,
		DCLResourceType:       dclResourceType,
		TerraformResourceType: terraformResourceType,
		HCLBlock:              block,
	}
	return d, nil
}

// ReplaceReferences substitutes any reference tags for their HCL address
// This should only be called after every dependency for a sample is built
func (s Sample) ReplaceReferences(d *Dependency) error {
	re := regexp.MustCompile(`"{{ref:([a-z.]*):(\w*)}}"`)
	matches := re.FindAllStringSubmatch(d.HCLBlock, -1)

	for _, match := range matches {
		referenceFileName := match[1]
		idField := match[2]
		var replacingText string
		for _, dep := range s.DependencyList {
			if dep.FileName == referenceFileName {
				replacingText = dep.TerraformResourceType + "." + dep.HCLLocalName + "." + idField
				break
			}
		}
		if replacingText == "" {
			return fmt.Errorf("Could not find reference file name: %s", referenceFileName)
		}
		d.HCLBlock = re.ReplaceAllString(d.HCLBlock, replacingText)
	}
	return nil
}

func (s Sample) generateHCLTemplate() (string, error) {
	if len(s.DependencyList) == 0 {
		return "", fmt.Errorf("Could not generate HCL template for %s: there are no dependencies", *s.Name)
	}
	var hcl string
	// var primaryIndex int
	for index := range s.DependencyList {
		err := s.ReplaceReferences(&s.DependencyList[index])
		if err != nil {
			return "", fmt.Errorf("Could not generate HCL template for %s: %s", *s.Name, err)
		}
		// Skip appending the primary resource, it should go last
		// if strings.Contains(*s.PrimaryResource, s.DependencyList[index].FileName) {
		// 	primaryIndex = index
		// 	continue
		// }
		hcl = fmt.Sprintf("%s%s\n", hcl, s.DependencyList[index].HCLBlock)
	}

	// hcl = fmt.Sprintf("%s%s", hcl, s.DependencyList[primaryIndex].HCLBlock)
	return hcl, nil
}

// GenerateDocsHCL generates sample HCL using docs substitution metadata
func (s Sample) GenerateDocsHCL() string {
	hcl, err := s.generateHCLTemplate()
	if err != nil {
		glog.Exit(err)
	}

	for _, sub := range s.Substitutions {
		re := regexp.MustCompile(fmt.Sprintf(`{{%s}}`, *sub.Substitution))
		hcl = re.ReplaceAllString(hcl, *sub.DocsValue)
	}
	return hcl
}

func (s Sample) PrimaryResourceName() string {
	fileParts := strings.Split(*s.PrimaryResource, ".")
	return fileParts[0]
}
