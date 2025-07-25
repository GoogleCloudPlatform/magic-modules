package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/openapi_generate"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/provider"
)

var wg sync.WaitGroup

// TODO rewrite: additional flags

// Example usage: --output $GOPATH/src/github.com/terraform-providers/terraform-provider-google-beta
var outputPathFlag = flag.String("output", "", "path to output generated files to")

// Example usage: --version beta
var versionFlag = flag.String("version", "", "optional version name. If specified, this version is preferred for resource generation when applicable")

var overrideDirectoryFlag = flag.String("overrides", "", "directory containing yaml overrides")

var productFlag = flag.String("product", "", "optional product name. If specified, the resources under the specific product will be generated. Otherwise, resources under all products will be generated.")

var resourceFlag = flag.String("resource", "", "optional resource name. Limits generation to the specified resource within a particular product.")

var doNotGenerateCode = flag.Bool("no-code", false, "do not generate code")

var doNotGenerateDocs = flag.Bool("no-docs", false, "do not generate docs")

var providerFlag = flag.String("provider", "", "optional provider name. If specified, a non-default provider will be used.")

var openapiGenerate = flag.Bool("openapi-generate", false, "Generate MMv1 YAML from openapi directory (Experimental)")

var showImportDiffsFlag = flag.Bool("show-import-diffs", false, "write go import diffs to stdout")

func main() {

	// Handle all flags in main. Other functions must not access flag values directly.
	flag.Parse()

	if *openapiGenerate {
		parser := openapi_generate.NewOpenapiParser("openapi_generate/openapi", "products")
		parser.Run()
		return
	}

	if *outputPathFlag == "" {
		log.Printf("No output path specified, exiting")
		return
	}

	GenerateProducts(*productFlag, *resourceFlag, *providerFlag, *versionFlag, *outputPathFlag, *overrideDirectoryFlag, !*doNotGenerateCode, !*doNotGenerateDocs, *showImportDiffsFlag)
}

func GenerateProducts(product, resource, providerName, version, outputPath, overrideDirectory string, generateCode, generateDocs, showImportDiffs bool) {
	if version == "" {
		log.Printf("No version specified, assuming ga")
		version = "ga"
	}
	var productsToGenerate []string
	var allProducts = false
	if product == "" {
		allProducts = true
	} else {
		var productToGenerate = fmt.Sprintf("products/%s", product)
		productsToGenerate = []string{productToGenerate}
	}

	var allProductFiles []string = make([]string, 0)

	files, err := filepath.Glob("products/**/product.yaml")
	if err != nil {
		panic(err)
	}
	for _, filePath := range files {
		dir := filepath.Dir(filePath)
		allProductFiles = append(allProductFiles, fmt.Sprintf("products/%s", filepath.Base(dir)))
	}

	if overrideDirectory != "" {
		log.Printf("Using override directory %s", overrideDirectory)

		// Normalize override dir to a path that is relative to the magic-modules directory
		// This is needed for templates that concatenate pwd + override dir + path
		if filepath.IsAbs(overrideDirectory) {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			overrideDirectory, err = filepath.Rel(wd, overrideDirectory)
			log.Printf("Override directory normalized to relative path %s", overrideDirectory)
		}

		overrideFiles, err := filepath.Glob(fmt.Sprintf("%s/products/**/product.yaml", overrideDirectory))
		if err != nil {
			panic(err)
		}
		for _, filePath := range overrideFiles {
			product, err := filepath.Rel(overrideDirectory, filePath)
			if err != nil {
				panic(err)
			}
			dir := filepath.Dir(product)
			productFile := fmt.Sprintf("products/%s", filepath.Base(dir))
			if !slices.Contains(allProductFiles, productFile) {
				allProductFiles = append(allProductFiles, productFile)
			}
		}
	}

	if allProducts {
		productsToGenerate = allProductFiles
	}

	if productsToGenerate == nil || len(productsToGenerate) == 0 {
		log.Fatalf("No product.yaml file found.")
	}

	startTime := time.Now()
	if providerName == "" {
		providerName = "default (terraform)"
	}
	log.Printf("Generating MM output to '%s'", outputPath)
	log.Printf("Building %s version", version)
	log.Printf("Building %s provider", providerName)

	productsForVersionChannel := make(chan *api.Product, len(allProductFiles))
	for _, productFile := range allProductFiles {
		wg.Add(1)
		go GenerateProduct(version, providerName, productFile, outputPath, productsForVersionChannel, startTime, productsToGenerate, resource, overrideDirectory, generateCode, generateDocs)
	}
	wg.Wait()

	close(productsForVersionChannel)

	var productsForVersion []*api.Product
	for p := range productsForVersionChannel {
		productsForVersion = append(productsForVersion, p)
	}
	slices.SortFunc(productsForVersion, func(p1, p2 *api.Product) int {
		return strings.Compare(strings.ToLower(p1.Name), strings.ToLower(p2.Name))
	})

	// In order to only copy/compile files once per provider this must be called outside
	// of the products loop. Create an MMv1 provider with an arbitrary product (the first loaded).
	providerToGenerate := newProvider(providerName, version, productsForVersion[0], startTime)
	providerToGenerate.CopyCommonFiles(outputPath, generateCode, generateDocs)

	if generateCode {
		providerToGenerate.CompileCommonFiles(outputPath, productsForVersion, "")
	}

	provider.FixImports(outputPath, showImportDiffs)
}

func GenerateProduct(version, providerName, productName, outputPath string, productsForVersionChannel chan *api.Product, startTime time.Time, productsToGenerate []string, resourceToGenerate, overrideDirectory string, generateCode, generateDocs bool) {
	defer wg.Done()

	productYamlPath := path.Join(productName, "product.yaml")

	var productOverridePath string
	if overrideDirectory != "" {
		productOverridePath = filepath.Join(overrideDirectory, productName, "product.yaml")
	}

	_, baseProductErr := os.Stat(productYamlPath)
	baseProductExists := !errors.Is(baseProductErr, os.ErrNotExist)

	_, overrideProductErr := os.Stat(productOverridePath)
	overrideProductExists := !errors.Is(overrideProductErr, os.ErrNotExist)

	if !(baseProductExists || overrideProductExists) {
		log.Fatalf("%s does not contain a product.yaml file", productName)
	}

	productApi := &api.Product{}

	if overrideProductExists {
		if baseProductExists {
			api.Compile(productYamlPath, productApi, overrideDirectory)
			overrideApiProduct := &api.Product{}
			api.Compile(productOverridePath, overrideApiProduct, overrideDirectory)

			api.Merge(reflect.ValueOf(productApi), reflect.ValueOf(*overrideApiProduct))
		} else {
			api.Compile(productOverridePath, productApi, overrideDirectory)
		}
	} else {
		api.Compile(productYamlPath, productApi, overrideDirectory)
	}

	var resources []*api.Resource = make([]*api.Resource, 0)

	if !productApi.ExistsAtVersionOrLower(version) {
		log.Printf("%s does not have a '%s' version, skipping", productName, version)
		return
	}

	resourceFiles, err := filepath.Glob(fmt.Sprintf("%s/*", productName))
	if err != nil {
		log.Fatalf("Cannot get resources files: %v", err)
	}
	// Base resource loop
	for _, resourceYamlPath := range resourceFiles {
		if filepath.Base(resourceYamlPath) == "product.yaml" || filepath.Ext(resourceYamlPath) != ".yaml" {
			continue
		}

		if overrideDirectory != "" {
			// skip if resource will be merged in the override loop
			resourceOverridePath := filepath.Join(overrideDirectory, resourceYamlPath)
			_, overrideResourceErr := os.Stat(resourceOverridePath)
			overrideResourceExists := !errors.Is(overrideResourceErr, os.ErrNotExist)
			if overrideResourceExists {
				continue
			}
		}

		resource := &api.Resource{}
		api.Compile(resourceYamlPath, resource, overrideDirectory)
		resource.SourceYamlFile = resourceYamlPath

		resource.TargetVersionName = version
		resource.Properties = resource.AddLabelsRelatedFields(resource.PropertiesWithExcluded(), nil)
		resource.SetDefault(productApi)
		resource.Validate()
		resources = append(resources, resource)
	}

	// Override Resource Loop
	if overrideDirectory != "" {
		productOverrideDir := filepath.Dir(productOverridePath)
		overrideFiles, err := filepath.Glob(fmt.Sprintf("%s/*", productOverrideDir))
		if err != nil {
			log.Fatalf("Cannot get override files: %v", err)
		}
		for _, overrideYamlPath := range overrideFiles {
			if filepath.Base(overrideYamlPath) == "product.yaml" || filepath.Ext(overrideYamlPath) != ".yaml" {
				continue
			}

			resource := &api.Resource{}

			baseResourcePath := filepath.Join(productName, filepath.Base(overrideYamlPath))
			_, baseResourceErr := os.Stat(baseResourcePath)
			baseResourceExists := !errors.Is(baseResourceErr, os.ErrNotExist)
			if baseResourceExists {
				api.Compile(baseResourcePath, resource, overrideDirectory)
				overrideResource := &api.Resource{}
				api.Compile(overrideYamlPath, overrideResource, overrideDirectory)
				api.Merge(reflect.ValueOf(resource), reflect.ValueOf(*overrideResource))
			} else {
				api.Compile(overrideYamlPath, resource, overrideDirectory)
			}

			resource.TargetVersionName = version
			resource.Properties = resource.AddLabelsRelatedFields(resource.PropertiesWithExcluded(), nil)
			resource.SetDefault(productApi)
			resource.Validate()
			resources = append(resources, resource)
		}

		// Sort resources by name
		sort.Slice(resources, func(i, j int) bool {
			return resources[i].Name < resources[j].Name
		})

	}

	productApi.Objects = resources
	productApi.Validate()

	providerToGenerate := newProvider(providerName, version, productApi, startTime)
	productsForVersionChannel <- productApi

	if !slices.Contains(productsToGenerate, productName) {
		log.Printf("%s not specified, skipping generation", productName)
		return
	}

	log.Printf("%s: Generating files", productName)

	providerToGenerate.Generate(outputPath, productName, resourceToGenerate, generateCode, generateDocs)
}

func newProvider(providerName, version string, productApi *api.Product, startTime time.Time) provider.Provider {
	switch providerName {
	case "tgc":
		return provider.NewTerraformGoogleConversion(productApi, version, startTime)
	case "tgc_cai2hcl":
		return provider.NewCaiToTerraformConversion(productApi, version, startTime)
	case "tgc_next":
		return provider.NewTerraformGoogleConversionNext(productApi, version, startTime)
	case "oics":
		return provider.NewTerraformOiCS(productApi, version, startTime)
	default:
		return provider.NewTerraform(productApi, version, startTime)
	}
}
