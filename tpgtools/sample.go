package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/golang/glog"
)

// Sample is the object containing sample data from DCL samples
type Sample struct {
	// Name of the file the sample was loaded from
	FileName string

	// Name is the name of a sample
	Name *string

	// Description is a short description of the sample
	Description *string

	// DependencyFileNames contains the filenames of every resource in the sample
	DependencyFileNames []string `yaml:"dependencies"`

	// PrimaryResource is the filename of the sample's primary resource
	PrimaryResource *string `yaml:"resource"`

	// Substitutions contains every substition in the sample
	Substitutions []Substitution `yaml:"substitutions"`

	IgnoreRead []string `yaml:"ignore_read"`

	// DependencyList is a list of objects containing metadata for each sample resource
	DependencyList []Dependency

	// The name of the test
	TestSlug string

	// The raw versions stated in the yaml
	Versions []string

	// A list of updates that the resource can transition between
	Updates []map[string]string

	// HasGAEquivalent tells us if we should have `provider = google-beta`
	// in the testcase. (if the test doesn't have a ga version of the test)
	HasGAEquivalent bool

	// SamplesPath is the path to the directory where the original sample data is stored
	SamplesPath string

	// resourceReference is the resource the sample belongs to
	resourceReference *Resource

	// CustomCheck allows you to add a terraform check function to all tests
	CustomCheck []string `yaml:"check"`

	// CodeInject references reletive raw files that should be injected into the sample test
	CodeInject []string `yaml:"code_inject"`

	// DocHide specifies a list of samples to hide from docs
	DocHide []string `yaml:"doc_hide"`

	// Testhide specifies a list of samples to hide from tests
	Testhide []string `yaml:"test_hide"`

	// ExtraDependencies are the additional golang dependencies the injected code may require
	ExtraDependencies []string `yaml:"extra_dependencies"`

	// Variables are the various attributes of the set of resources that need to be filled in.
	Variables []Variable `yaml:"variables"`
}

// Variable contains metadata about the types of variables in a sample.
type Variable struct {
	// Name is the variable name in the JSON.
	Name string `yaml:"name"`
	// Type is the variable type.
	Type string `yaml:"type"`
}

// Substitution contains metadata that varies for the sample context
type Substitution struct {
	// Substitution is the text to be substituted, e.g. topic_name
	Substitution *string

	// Value is the value of the substituted text
	Value *string `yaml:"value"`
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
	HCLBlock string // Path to the directory where the sample data is stored
}

// BuildDependency produces a Dependency using a file and filename
func BuildDependency(fileName, product, localname, version string, b []byte) (*Dependency, error) {
	var resourceName string
	fileParts := strings.Split(fileName, ".")
	if len(fileParts) == 4 {
		product = strings.Title(fileParts[1])
		resourceName = strings.Title(fileParts[2])
	} else if len(fileParts) == 3 {
		resourceName = strings.Title(fileParts[1])
	} else {
		return nil, fmt.Errorf("Invalid sample dependency file name: %s", fileName)
	}

	if localname == "" {
		localname = fileParts[0]
	}
	dclResourceType := product + snakeToTitleCase(resourceName)
	terraformResourceType, err := DCLToTerraformReference(dclResourceType, version)
	if err != nil {
		return nil, fmt.Errorf("Error generating sample dependency %s: %s", fileName, err)
	}

	block, err := ConvertSampleJSONToHCL(dclResourceType, version, b)
	if err != nil {
		return nil, fmt.Errorf("Error generating sample dependency %s: %s", fileName, err)
	}

	re := regexp.MustCompile(`(resource "` + terraformResourceType + `" ")(\w*)`)
	block = re.ReplaceAllString(block, "${1}"+localname)

	d := &Dependency{
		FileName:              fileName,
		HCLLocalName:          localname,
		DCLResourceType:       dclResourceType,
		TerraformResourceType: terraformResourceType,
		HCLBlock:              block,
	}
	return d, nil
}

func (s *Sample) generateSampleDependency(fileName string) Dependency {
	return s.generateSampleDependencyWithName(fileName, "")
}

func (s *Sample) generateSampleDependencyWithName(fileName, localname string) Dependency {
	dFileNameParts := strings.Split(fileName, "samples/")
	fileName = dFileNameParts[len(dFileNameParts)-1]
	dependencyBytes, err := ioutil.ReadFile(path.Join(s.SamplesPath, fileName))
	version := s.resourceReference.versionMetadata.V
	product := s.resourceReference.productMetadata.ProductType()
	d, err := BuildDependency(fileName, product, localname, version, dependencyBytes)
	if err != nil {
		glog.Exit(err)
	}
	return *d
}

func (s *Sample) GetCodeToInject() []string {
	sampleAccessoryFolder := s.resourceReference.getSampleAccessoryFolder()
	var out []string
	for _, fileName := range s.CodeInject {
		filePath := path.Join(sampleAccessoryFolder, fileName)
		tc, err := ioutil.ReadFile(filePath)
		if err != nil {
			glog.Exit(err)
		}
		out = append(out, string(tc))
	}
	return out
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
	for index := range s.DependencyList {
		err := s.ReplaceReferences(&s.DependencyList[index])
		if err != nil {
			return "", fmt.Errorf("Could not generate HCL template for %s: %s", *s.Name, err)
		}
		hcl = fmt.Sprintf("%s%s\n", hcl, s.DependencyList[index].HCLBlock)
	}

	return hcl, nil
}

// GenerateHCL generates sample HCL using docs substitution metadata
func (s Sample) GenerateHCL(isDocs bool) string {
	var hcl string
	var err error
	if !s.isNativeHCL() {
		hcl, err = s.generateHCLTemplate()
		if err != nil {
			glog.Exit(err)
		}
	} else {
		tc, err := ioutil.ReadFile(path.Join(s.SamplesPath, *s.PrimaryResource))
		if err != nil {
			glog.Exit(err)
		}
		hcl = string(tc)
	}
	for _, sub := range s.Substitutions {
		re := regexp.MustCompile(fmt.Sprintf(`{{%s}}`, *sub.Substitution))
		hcl = re.ReplaceAllString(hcl, sub.translateValue(isDocs))
	}
	return hcl
}

// isNativeHCL returns whether the resource file is terraform synatax
func (s Sample) isNativeHCL() bool {
	return strings.HasSuffix(*s.PrimaryResource, ".tf.tmpl")
}

// EnumerateWithUpdateSamples returns an array of new samples expanded with
// any subsequent samples
func (s *Sample) EnumerateWithUpdateSamples() []Sample {
	out := []Sample{*s}
	for i, update := range s.Updates {
		newSample := *s
		primaryResource := update["resource"]
		newSample.PrimaryResource = &primaryResource
		if !newSample.isNativeHCL() {
			var newDeps []Dependency
			newDeps = append(newDeps, newSample.DependencyList...)
			newDeps[0] = newSample.generateSampleDependencyWithName(*newSample.PrimaryResource, "primary")
			newSample.DependencyList = newDeps
		}
		newSample.TestSlug = fmt.Sprintf("%sUpdate%v", newSample.TestSlug, i)
		newSample.Updates = nil
		out = append(out, newSample)
	}
	return out
}

// ExpandContext expands the context model used in the generated tests
func (s Sample) ExpandContext() map[string]string {
	out := map[string]string{}
	for _, sub := range s.Substitutions {
		translation, hasTranslation := translationMap[*sub.Value]
		if hasTranslation {
			out[translation.contextKey] = translation.contextValue
		}
	}
	return out
}

type translationIndex struct {
	docsValue    string
	contextKey   string
	contextValue string
}

var translationMap = map[string]translationIndex{
	":ORG_ID": {
		docsValue:    "123456789",
		contextKey:   "org_id",
		contextValue: "getTestOrgFromEnv(t)",
	},
	":ORG_DOMAIN": {
		docsValue:    "example.com",
		contextKey:   "org_domain",
		contextValue: "getTestOrgDomainFromEnv(t)",
	},
	":CREDENTIALS": {
		docsValue:    "my/credentials/filename.json",
		contextKey:   "credentials",
		contextValue: "getTestCredsFromEnv(t)",
	},
	":REGION": {
		docsValue:    "us-west1",
		contextKey:   "region",
		contextValue: "getTestRegionFromEnv()",
	},
	":ORG_TARGET": {
		docsValue:    "123456789",
		contextKey:   "org_target",
		contextValue: "getTestOrgTargetFromEnv(t)",
	},
	":BILLING_ACCT": {
		docsValue:    "000000-0000000-0000000-000000",
		contextKey:   "billing_acct",
		contextValue: "getTestBillingAccountFromEnv(t)",
	},
	":SERVICE_ACCT": {
		docsValue:    "emailAddress:my@service-account.com",
		contextKey:   "service_acct",
		contextValue: "getTestServiceAccountFromEnv(t)",
	},
	":PROJECT": {
		docsValue:    "my-project-name",
		contextKey:   "project_name",
		contextValue: "getTestProjectFromEnv()",
	},
	":PROJECT_NAME": {
		docsValue:    "my-project-name",
		contextKey:   "project_name",
		contextValue: "getTestProjectFromEnv()",
	},
	":FIRESTORE_PROJECT_NAME": {
		docsValue:    "my-project-name",
		contextKey:   "firestore_project_name",
		contextValue: "getTestFirestoreProjectFromEnv(t)",
	},
	":CUST_ID": {
		docsValue:    "A01b123xz",
		contextKey:   "cust_id",
		contextValue: "getTestCustIdFromEnv(t)",
	},
	":IDENTITY_USER": {
		docsValue:    "cloud_identity_user",
		contextKey:   "identity_user",
		contextValue: "getTestIdentityUserFromEnv(t)",
	},
}

// translateValue returns the value to embed in the hcl
func (sub *Substitution) translateValue(isDocs bool) string {
	value := *sub.Value
	translation, hasTranslation := translationMap[value]

	if isDocs {
		if hasTranslation {
			return translation.docsValue
		}
		return value
	}

	if hasTranslation {
		return fmt.Sprintf("%%{%s}", translation.contextKey)
	}

	if strings.Contains(value, "-") {
		value = fmt.Sprintf("tf-test-%s", value)
	} else if strings.Contains(value, "_") {
		value = fmt.Sprintf("tf_test_%s", value)
	}

	// Random suffix is 10 characters and standard name length <= 64
	if len(value) > 54 {
		value = value[:54]
	}
	return fmt.Sprintf("%s%%{random_suffix}", value)
}

func (s Sample) PrimaryResourceName() string {
	fileParts := strings.Split(*s.PrimaryResource, ".")
	return fileParts[0]
}
