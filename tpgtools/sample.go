package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
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

	IgnoreRead []string `yaml:"ignore_read"`

	// DependencyList is a list of objects containing metadata for each sample resource
	DependencyList []Dependency

	// The name of the test
	TestSlug RenderedString

	// The raw versions stated in the yaml
	Versions []string

	// A list of updates that the resource can transition between
	Updates []Update

	// HasGAEquivalent tells us if we should have `provider = google-beta`
	// in the testcase. (if the test doesn't have a ga version of the test)
	HasGAEquivalent bool

	// SamplesPath is the path to the directory where the original sample data is stored
	SamplesPath Filepath

	// resourceReference is the resource the sample belongs to
	resourceReference *Resource

	// CustomCheck allows you to add a terraform check function to all tests
	CustomCheck []string `yaml:"check"`

	// CodeInject references reletive raw files that should be injected into the sample test
	CodeInject []string `yaml:"code_inject"`

	// DocHide specifies a list of samples to hide from docs
	DocHide []string `yaml:"doc_hide"`

	// DocHideConditional specifies a list of samples to hide from docs when resource location matches
	DocHideConditional []DocHideCondition `yaml:"doc_hide_conditional"`

	// Testhide specifies a list of samples to hide from tests
	Testhide []string `yaml:"test_hide"`

	// TesthideConditional specifies a list of samples to hide from tests when resource location matches
	TestHideConditional []TestHideCondition `yaml:"test_hide_conditional"`

	// ExtraDependencies are the additional golang dependencies the injected code may require
	ExtraDependencies []string `yaml:"extra_dependencies"`

	// Type is the resource type.
	Type string `yaml:"type"`

	// Variables are the various attributes of the set of resources that need to be filled in.
	Variables []Variable `yaml:"variables"`
}

// Variable contains metadata about the types of variables in a sample.
type Variable struct {
	// Name is the variable name in the JSON.
	Name string `yaml:"name"`
	// Type is the variable type.
	Type string `yaml:"type"`
	// DocsValue is an optional value that should be substituted directly into
	// the documentation for this variable.  If not provided, tpgtools makes
	// its best guess about a suitable value.  Generally, this is only provided
	// if the "best guess" is a poor one.
	DocsValue string `yaml:"docs_value"`
}

type DocHideCondition struct {
	// Location is the location attribute to match, if matched, append Name to list of DocHide
	Location string `yaml:"location"`
	// Name specifies sample file name to add to DocHide if location matches.
	Name string `yaml:"file_name"`
}

type TestHideCondition struct {
        // Location is the location attribute to match, if matched, append Name to list of Testhide
        Location string `yaml:"location"`
        // Name specifies sample file name to add to Testhide if location matches.
        Name string `yaml:"file_name"`
}
// Dependency contains data that describes a single resource in a sample
type Dependency struct {
	// FileName is the name of the file as it appears in testcases.yaml
	FileName string

	// HCLLocalName is the local name of the HCL block, e.g. "basic" or "default"
	HCLLocalName string

	// TerraformResourceType is the type represented in Terraform, e.g. "google_compute_instance"
	TerraformResourceType string

	// HCLBlock is the snippet of HCL config that declares this resource
	HCLBlock string // Path to the directory where the sample data is stored
}

type Update struct {
	// The list of dependency resources to update.
	Dependencies []string `yaml:"dependencies"`

	// The resource to update.
	Resource string `yaml:"resource"`
}

func packageNameFromFilepath(fp Filepath, product SnakeCaseProductName) (DCLPackageName, error) {
	pm := NewProductMetadata(fp, string(product))
	return pm.PackageName, nil
}

func findDCLReferencePackage(product SnakeCaseProductName) (DCLPackageName, error) {
	// Most packages are just the product name with all the underscores removed.
	// Try that first.
	// We can check if a package exists by checking the "productOverrides" map from product.go, which
	// will be populated by this point.  That takes a "Filepath", because the reference is to the
	// actual name of the directory that contains the overrides - by mandate, that's the same as the
	// dcl package name, so this conversion happens to work out.
	possibleFilepath := Filepath(strings.ReplaceAll(string(product), "_", ""))
	if _, ok := productOverrides[possibleFilepath]; ok {
		return packageNameFromFilepath(possibleFilepath, product)
	}
	baseFilepath := possibleFilepath
	possibleFilepath = Filepath(string(baseFilepath) + "/beta")
	if _, ok := productOverrides[possibleFilepath]; ok {
		return packageNameFromFilepath(possibleFilepath, product)
	}
	possibleFilepath = Filepath(string(baseFilepath) + "/alpha")
	if _, ok := productOverrides[possibleFilepath]; ok {
		return packageNameFromFilepath(possibleFilepath, product)
	}

	// Otherwise, just return an error.
	var productOverrideKeys []Filepath
	for k, _ := range productOverrides {
		productOverrideKeys = append(productOverrideKeys, k)
	}
	return DCLPackageName(""), fmt.Errorf("can't find %q in the overrides map, which contains %v", product, productOverrideKeys)
}

// BuildDependency produces a Dependency using a file and filename
func BuildDependency(fileName string, product SnakeCaseProductName, localname, version string, hasGAEquivalent bool, b []byte) (*Dependency, error) {
	// Miscellaneous name rather than "resource name" because this is the name in the sample json file - which might not match the TF name!
	// we have to account for that.
	var resourceName miscellaneousNameSnakeCase
	var packageName DCLPackageName
	fileParts := strings.Split(fileName, ".")
	if len(fileParts) == 4 {
		p, rn := fileParts[1], fileParts[2]
		packageName = DCLPackageName(p)
		resourceName = miscellaneousNameSnakeCase(rn)
	} else if len(fileParts) == 3 {
		resourceName = miscellaneousNameSnakeCase(fileParts[1])
		var err error
		packageName, err = findDCLReferencePackage(product)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Invalid sample dependency file name: %s", fileName)
	}

	if localname == "" {
		localname = fileParts[0]
	}
	terraformResourceType, err := DCLToTerraformReference(packageName, resourceName, version)
	if err != nil {
		return nil, fmt.Errorf("Error generating sample dependency reference %s: %s", fileName, err)
	}

	block, err := ConvertSampleJSONToHCL(packageName, resourceName, version, hasGAEquivalent, b)
	if err != nil {
		glog.Errorf("failed to convert %q", fileName)
		return nil, fmt.Errorf("Error generating sample dependency %s: %s", fileName, err)
	}

	// Find all instances of `resource "foo" "bar"` and replace `bar` with localname.
	re := regexp.MustCompile(`(resource "` + terraformResourceType + `" ")(\w*)`)
	block = re.ReplaceAllString(block, "${1}"+localname)

	d := &Dependency{
		FileName:              fileName,
		HCLLocalName:          localname,
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
	dependencyBytes, err := ioutil.ReadFile(path.Join(string(s.SamplesPath), fileName))
	version := s.resourceReference.versionMetadata.V
	product := s.resourceReference.productMetadata.ProductName
	d, err := BuildDependency(fileName, product, localname, version, s.HasGAEquivalent, dependencyBytes)
	if err != nil {
		glog.Exit(err)
	}
	return *d
}

func (s *Sample) GetCodeToInject() []string {
	sampleAccessoryFolder := s.resourceReference.getSampleAccessoryFolder()
	var out []string
	for _, fileName := range s.CodeInject {
		filePath := path.Join(string(sampleAccessoryFolder), fileName)
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
	re := regexp.MustCompile(`"?{{\s*ref:([a-z_]*\.[a-z_]*\.[a-z_]*(?:\.[a-z_]*)?):([a-zA-Z0-9_\.\[\]]*)\s*}}"?`)
	matches := re.FindAllStringSubmatch(d.HCLBlock, -1)

	for _, match := range matches {
		referenceFileName := match[1]
		idField := dcl.TitleToSnakeCase(match[2])
		var tfReference string
		for _, dep := range s.DependencyList {
			if dep.FileName == referenceFileName {
				tfReference = dep.TerraformResourceType + "." + dep.HCLLocalName + "." + idField
				break
			}
		}
		if tfReference == "" {
			return fmt.Errorf("Could not find reference file name: %s", referenceFileName)
		}
		startsWithQuote := strings.HasPrefix(match[0], `"`)
		endsWithQuote := strings.HasSuffix(match[0], `"`)
		if !(startsWithQuote && endsWithQuote) {
			tfReference = fmt.Sprintf("${%s}", tfReference)
			if startsWithQuote {
				tfReference = `"` + tfReference
			}
			if endsWithQuote {
				tfReference += `"`
			}
		}
		d.HCLBlock = strings.Replace(d.HCLBlock, match[0], tfReference, 1)
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
		tc, err := ioutil.ReadFile(path.Join(string(s.SamplesPath), *s.PrimaryResource))
		if err != nil {
			glog.Exit(err)
		}
		hcl = string(tc)
	}
	for _, sub := range s.Variables {
		re := regexp.MustCompile(fmt.Sprintf(`{{%s}}`, sub.Name))
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
		primaryResource := update.Resource
		// TODO(magic-modules-eng): Consume new dependency list.
		newSample.PrimaryResource = &primaryResource
		if !newSample.isNativeHCL() {
			var newDeps []Dependency
			newDeps = append(newDeps, newSample.generateSampleDependencyWithName(*newSample.PrimaryResource, "primary"))
			for _, newDepFilename := range update.Dependencies {
				newDepFilename = strings.TrimPrefix(newDepFilename, "samples/")
				newDeps = append(newDeps, newSample.generateSampleDependencyWithName(newDepFilename, basicResourceName(newDepFilename)))
			}
			newSample.DependencyList = newDeps
		}
		newSample.TestSlug = RenderedString(fmt.Sprintf("%sUpdate%v", newSample.TestSlug, i))
		newSample.Updates = nil
		newSample.Variables = s.Variables
		out = append(out, newSample)
	}
	return out
}

func basicResourceName(depFilename string) string {
	re := regexp.MustCompile("^update(_\\d)?\\.")
	// update_1.resource.json -> basic.resource.json
	basicReplaced := re.ReplaceAllString(depFilename, "basic.")
	re = regexp.MustCompile("^update(_\\d)?_")
	// update_1_name.resource.json -> name.resource.json
	prefixTrimmed := re.ReplaceAllString(basicReplaced, "")
	return dcl.SnakeToJSONCase(strings.Split(prefixTrimmed, ".")[0])
}

// ExpandContext expands the context model used in the generated tests
func (s Sample) ExpandContext() map[string]string {
	out := map[string]string{}
	for _, sub := range s.Variables {
		translation, hasTranslation := translationMap[sub.Type]
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
	"org_id": {
		docsValue:    "123456789",
		contextKey:   "org_id",
		contextValue: "getTestOrgFromEnv(t)",
	},
	"org_name": {
		docsValue:    "example.com",
		contextKey:   "org_domain",
		contextValue: "getTestOrgDomainFromEnv(t)",
	},
	"region": {
		docsValue:    "us-west1",
		contextKey:   "region",
		contextValue: "getTestRegionFromEnv()",
	},
	"zone": {
		docsValue:    "us-west1-a",
		contextKey:   "zone",
		contextValue: "getTestZoneFromEnv()",
	},
	"org_target": {
		docsValue:    "123456789",
		contextKey:   "org_target",
		contextValue: "getTestOrgTargetFromEnv(t)",
	},
	"billing_account": {
		docsValue:    "000000-0000000-0000000-000000",
		contextKey:   "billing_acct",
		contextValue: "getTestBillingAccountFromEnv(t)",
	},
	"test_service_account": {
		docsValue:    "emailAddress:my@service-account.com",
		contextKey:   "service_acct",
		contextValue: "getTestServiceAccountFromEnv(t)",
	},
	"project": {
		docsValue:    "my-project-name",
		contextKey:   "project_name",
		contextValue: "getTestProjectFromEnv()",
	},
	"project_number": {
		docsValue:    "my-project-number",
		contextKey:   "project_number",
		contextValue: "getTestProjectNumberFromEnv()",
	},
	"customer_id": {
		docsValue:    "A01b123xz",
		contextKey:   "cust_id",
		contextValue: "getTestCustIdFromEnv(t)",
	},
	// Begin a long list of multicloud-only values which are not going to see reuse.
	// We can hardcode fake values because we are
	// always going to use the no-provisioning mode for unit testing, of these resources
	// where we don't have to actually have a real AWS account.

	"aws_account_id": {
		docsValue:    "012345678910",
		contextKey:   "aws_acct_id",
		contextValue: `"111111111111"`,
	},
	"aws_database_encryption_key": {
		docsValue:    "12345678-1234-1234-1234-123456789111",
		contextKey:   "aws_db_key",
		contextValue: `"00000000-0000-0000-0000-17aad2f0f61f"`,
	},
	"aws_region": {
		docsValue:    "my-aws-region",
		contextKey:   "aws_region",
		contextValue: `"us-west-2"`,
	},
	"aws_security_group": {
		docsValue:    "sg-00000000000000000",
		contextKey:   "aws_sg",
		contextValue: `"sg-0b3f63cb91b247628"`,
	},
	"aws_volume_encryption_key": {
		docsValue:    "12345678-1234-1234-1234-123456789111",
		contextKey:   "aws_vol_key",
		contextValue: `"00000000-0000-0000-0000-17aad2f0f61f"`,
	},
	"aws_vpc": {
		docsValue:    "vpc-00000000000000000",
		contextKey:   "aws_vpc",
		contextValue: `"vpc-0b3f63cb91b247628"`,
	},
	"aws_subnet": {
		docsValue:    "subnet-00000000000000000",
		contextKey:   "aws_subnet",
		contextValue: `"subnet-0b3f63cb91b247628"`,
	},
	"azure_application": {
		docsValue:    "12345678-1234-1234-1234-123456789111",
		contextKey:   "azure_app",
		contextValue: `"00000000-0000-0000-0000-17aad2f0f61f"`,
	},
	"azure_subscription": {
		docsValue:    "12345678-1234-1234-1234-123456789111",
		contextKey:   "azure_sub",
		contextValue: `"00000000-0000-0000-0000-17aad2f0f61f"`,
	},
	"azure_ad_tenant": {
		docsValue:    "12345678-1234-1234-1234-123456789111",
		contextKey:   "azure_tenant",
		contextValue: `"00000000-0000-0000-0000-17aad2f0f61f"`,
	},
	"azure_proxy_config_secret_version": {
		docsValue:    "0000000000000000000000000000000000",
		contextKey:   "azure_config_secret",
		contextValue: `"07d4b1f1a7cb4b1b91f070c30ae761a1"`,
	},
	"byo_multicloud_prefix": {
		docsValue:    "my-",
		contextKey:   "byo_prefix",
		contextValue: `"mmv2"`,
	},
}

// translateValue returns the value to embed in the hcl
func (sub *Variable) translateValue(isDocs bool) string {
	value := sub.Name
	translation, hasTranslation := translationMap[sub.Type]

	if isDocs {
		if sub.DocsValue != "" {
			return sub.DocsValue
		}
		if hasTranslation {
			return translation.docsValue
		}
		if sub.Type != "resource_name" {
			glog.Exitf("Cannot generate docs for variable of type %q.", sub.Type)
		}
		return value
	}

	if hasTranslation {
		return fmt.Sprintf("%%{%s}", translation.contextKey)
	}

	if sub.Type != "resource_name" {
		glog.Exitf("Cannot generate sample test with variable of type %q - please add to sample.go's translationMap.", sub.Type)
	}

	// Use '_' if already present, or '-' otherwise (some APIs require '-').
	if strings.Contains(value, "_") {
		value = fmt.Sprintf("tf_test_%s", value)
	} else {
		value = fmt.Sprintf("tf-test-%s", value)
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
