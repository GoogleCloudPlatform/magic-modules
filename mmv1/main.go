package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/loader"
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

	GenerateProducts(*productFlag, *resourceFlag, *providerFlag, *versionFlag, *outputPathFlag, *overrideDirectoryFlag, !*doNotGenerateCode, !*doNotGenerateDocs)
}

func GenerateProducts(product, resource, providerName, version, outputPath, overrideDirectory string, generateCode, generateDocs bool) {
	if version == "" {
		log.Printf("No version specified, assuming ga")
		version = "ga"
	}

	startTime := time.Now()
	if providerName == "" {
		providerName = "default (terraform)"
	}
	log.Printf("Generating MM output to '%s'", outputPath)
	log.Printf("Building %s version", version)
	log.Printf("Building %s provider", providerName)

	loader := loader.NewLoader(loader.Config{Version: version, OverrideDirectory: overrideDirectory})
	loadedProducts := loader.LoadProducts()

	var productsToGenerate []string
	if product == "" {
		for _, p := range loadedProducts {
			productsToGenerate = append(productsToGenerate, p.PackagePath)
		}
	} else {
		var productToGenerate = fmt.Sprintf("products/%s", product)
		productsToGenerate = []string{productToGenerate}
	}

	for _, productApi := range loadedProducts {
		wg.Add(1)
		go GenerateProduct(version, providerName, productApi, outputPath, startTime, productsToGenerate, resource, generateCode, generateDocs)
	}
	wg.Wait()

	var productsForVersion []*api.Product
	for _, p := range loadedProducts {
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
}

// GenerateProduct generates code and documentation for a product
// This now uses the CompileProduct method to separate compilation from generation
func GenerateProduct(version, providerName string, productApi *api.Product, outputPath string,
	startTime time.Time, productsToGenerate []string, resourceToGenerate string,
	generateCode, generateDocs bool) {
	defer wg.Done()

	if !slices.Contains(productsToGenerate, productApi.PackagePath) {
		log.Printf("%s not specified, skipping generation", productApi.PackagePath)
		return
	}

	log.Printf("%s: Generating files", productApi.PackagePath)
	providerToGenerate := newProvider(providerName, version, productApi, startTime)
	providerToGenerate.Generate(outputPath, resourceToGenerate, generateCode, generateDocs)
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
