package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/provider"
)

var versionFlag = flag.String("version", "", "provider version to generate")
var productNameFlag = flag.String("product_name", "", "name of the product referenced by --product")
var productFlag = flag.String("product", "", "path to product.yaml input file")
var productOverrideFlag = flag.String("product_override", "", "path to a product override input file")
var resourceFlag = flag.String("resource", "", "path to resource.yaml input file")
var resourceOverrideFlag = flag.String("resource_override", "", "path to a resource override input file")
var outputPathFlag = flag.String("output", "", "output path for generated files")
var typeFlag = flag.String("type", "", "type of output to generate [product|resource|operation]")
var providerFlag = flag.String("provider", "", "target provider")

func main() {
	flag.Parse()

	if *versionFlag == "" {
		log.Fatal("--version is required")
	}
	if *providerFlag == "" {
		log.Fatal("--provider is required")
	}
	if *outputPathFlag == "" {
		log.Fatal("--output is required")
	}
	if *productNameFlag == "" {
		log.Fatal("--product_name is required")
	}
	if *productFlag == "" {
		log.Fatal("--product is required")
	}

	switch *typeFlag {
	case "":
		log.Fatal("--type is required")
	case "product":
	case "resource":
		if *resourceFlag == "" {
			log.Fatal("--resource is required with --type=resource")
		}
	case "metadata":
		if *resourceFlag == "" {
			log.Fatal("--resource is required with --type=metadata")
		}
	case "operation":
		if *resourceFlag == "" {
			log.Fatal("--resource is required with --type=operation")
		}
	case "sweeper":
		if *resourceFlag == "" {
			log.Fatal("--resource is required with --type=sweeper")
		}
	default:
		log.Fatalf("unrecognized --type %q", *typeFlag)
	}

	var product api.Product
	api.Compile(*productFlag, &product)
	if *productOverrideFlag != "" {
		var override api.Product
		api.Compile(*productOverrideFlag, &override)
		api.Merge(reflect.ValueOf(product), reflect.ValueOf(override), *versionFlag)
	}
	if !product.ExistsAtVersionOrLower(*versionFlag) {
		log.Fatalf("product %q does not support version %q", *productNameFlag, *versionFlag)
	}
	product.Version = product.VersionObjOrClosest(*versionFlag)

	if *resourceFlag != "" {
		var resource api.Resource
		api.Compile(*resourceFlag, &resource)
		if *resourceOverrideFlag != "" {
			var override api.Resource
			api.Compile(*resourceOverrideFlag, &override)
			api.Merge(reflect.ValueOf(resource), reflect.ValueOf(override), *versionFlag)
		}
		resource.TargetVersionName = *versionFlag
		resource.SetDefault(&product)
		resource.Properties = resource.AddExtraFields(resource.PropertiesWithExcluded(), nil)
		resource.SetDefault(&product)
		product.Objects = []*api.Resource{&resource}
	}

	product.Validate()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not find wd: %v", err)
	}
	fsys := os.DirFS(filepath.Join(wd, "mmv1"))

	switch *providerFlag {
	case "tgc", "tgc_cai2hcl", "tgc_next", "oics":
		log.Fatalf("--provider %q is not yet supported", *providerFlag)
	case "tpg":
	default:
		log.Fatalf("unrecognized --provider %q", *providerFlag)
	}

	generator := provider.NewTerraform(&product, *versionFlag, time.Now(), fsys)

	switch *typeFlag {
	case "product":
		generator.GenerateProductFile(*outputPathFlag)
	case "resource":
		generator.GenerateResourceFile(*product.Objects[0], *outputPathFlag)
	case "metadata":
		generator.GenerateResourceMetadataFile(*product.Objects[0], *outputPathFlag)
	case "operation":
		generator.GenerateOperationFile(*product.Objects[0], *outputPathFlag)
	case "sweeper":
		generator.GenerateResourceSweeperFile(*product.Objects[0], *outputPathFlag)
	}
}
