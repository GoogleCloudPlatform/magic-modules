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
	DependencyFileNames []string `yaml:"dependency"`

	// PrimaryResource is the filename of the sample's primary resource
	PrimaryResource *string `yaml:"resource"`

	// Substitutions contains every substition in the sample
	Substitutions []Substitution `yaml:"extra_substitutions"`

	// Dependencies is a list of objects containing metadata for each sample resource
	Dependencies []Dependency
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
func BuildDependency(fileName string, b []byte) (*Dependency, error) {
	fileParts := strings.Split(fileName, ".")
	if len(fileParts) != 4 {
		return nil, fmt.Errorf("Invalid sample dependency file name: %s", fileName)
	}

	localName := fileParts[0]
	dclResourceType := strings.Title(fileParts[1]) + strings.Title(fileParts[2])
	terraformResourceType := fmt.Sprintf("google_%s_%s", fileParts[1], fileParts[2])

	block, err := ConvertSampleJSONToHCL(dclResourceType, b)
	if err != nil {
		return nil, fmt.Errorf("Error generating sample dependency: %s", err)
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
	re := regexp.MustCompile(`{{ref:([a-z.]*):(\w*)}}`)
	matches := re.FindAllStringSubmatch(d.HCLBlock, -1)

	for _, match := range matches {
		referenceFileName := match[1]
		idField := match[2]
		var replacingText string
		for _, dep := range s.Dependencies {
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
	if len(s.Dependencies) == 0 {
		return "", fmt.Errorf("Could not generate HCL template for %s: there are no dependencies", *s.Name)
	}
	var hcl string
	var primaryIndex int
	for index := range s.Dependencies {
		err := s.ReplaceReferences(&s.Dependencies[index])
		if err != nil {
			return "", fmt.Errorf("Could not generate HCL template for %s: %s", *s.Name, err)
		}
		// Skip appending the primary resource, it should go last
		if s.Dependencies[index].FileName == *s.PrimaryResource {
			primaryIndex = index
			continue
		}
		hcl = fmt.Sprintf("%s%s\n", hcl, s.Dependencies[index].HCLBlock)
	}

	hcl = fmt.Sprintf("%s%s", hcl, s.Dependencies[primaryIndex].HCLBlock)
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
