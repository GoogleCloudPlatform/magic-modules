package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Samples []Sample

type Sample struct {
	// name is the name of a sample.
	Name *string

	Description *string

	DependencyFileNames []string `yaml:"dependency"`

	PrimaryResource *string `yaml:"resource"`

	Substitutions []Substitution `yaml:"extra_substitutions"`

	Dependencies []Dependency
}

type Substitution struct {
	Substitution *string
	TestValue    *string `yaml:"test_value"`
	DocsValue    *string `yaml:"docs_value"`
}

type Dependency struct {
	FileName              string
	HCLLocalName          string
	DCLResourceType       string
	TerraformResourceType string
	HCLBlock              string
}

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

func (s Sample) GenerateDocsHCL() string {
	hcl, _ := s.generateHCLTemplate()
	return hcl
}

// func generateTestHCL() string{

// }
