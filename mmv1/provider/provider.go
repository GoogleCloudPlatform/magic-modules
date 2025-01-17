package provider

import (
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
)

type Provider interface {
	Generate(string, string, string, bool, bool)
	CopyCommonFiles(outputFolder string, generateCode, generateDocs bool)
	CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string)
}

// Shared constants and functions among the providers

const TERRAFORM_PROVIDER_GA = "github.com/hashicorp/terraform-provider-google"
const TERRAFORM_PROVIDER_BETA = "github.com/hashicorp/terraform-provider-google-beta"
const TERRAFORM_PROVIDER_PRIVATE = "internal/terraform-next"
const RESOURCE_DIRECTORY_GA = "google"
const RESOURCE_DIRECTORY_BETA = "google-beta"
const RESOURCE_DIRECTORY_PRIVATE = "google-private"

// # TODO(nelsonjr): Review all object interfaces and move to private methods
// # that should not be exposed outside the object hierarchy.
func ProviderName(t Provider) string {
	return reflect.TypeOf(t).Name()
}

func ImportPathFromVersion(v string) string {
	var tpg, dir string
	switch v {
	case "ga":
		tpg = TERRAFORM_PROVIDER_GA
		dir = RESOURCE_DIRECTORY_GA
	case "beta":
		tpg = TERRAFORM_PROVIDER_BETA
		dir = RESOURCE_DIRECTORY_BETA
	default:
		tpg = TERRAFORM_PROVIDER_PRIVATE
		dir = RESOURCE_DIRECTORY_PRIVATE
	}
	return fmt.Sprintf("%s/%s", tpg, dir)
}
