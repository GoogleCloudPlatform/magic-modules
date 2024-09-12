package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/provider"
)

var wg sync.WaitGroup

// TODO rewrite: additional flags

// Example usage: --output $GOPATH/src/github.com/terraform-providers/terraform-provider-google-beta
var outputPath = flag.String("output", "", "path to output generated files to")

// Example usage: --version beta
var version = flag.String("version", "", "optional version name. If specified, this version is preferred for resource generation when applicable")

var product = flag.String("product", "", "optional product name. If specified, the resources under the specific product will be generated. Otherwise, resources under all products will be generated.")

var resourceToGenerate = flag.String("resource", "", "optional resource name. Limits generation to the specified resource within a particular product.")

var doNotGenerateCode = flag.Bool("no-code", false, "do not generate code")

var doNotGenerateDocs = flag.Bool("no-docs", false, "do not generate docs")

var forceProvider = flag.String("provider", "", "optional provider name. If specified, a non-default provider will be used.")

// Example usage: --yaml
var yamlMode = flag.Bool("yaml", false, "copy text over from ruby yaml to go yaml")

// Example usage: --template
var templateMode = flag.Bool("template", false, "copy templates over from .erb to go .tmpl")

// Example usage: --handwritten
var handwrittenMode = flag.Bool("handwritten", false, "copy handwritten files over from .erb to go .tmpl")

func main() {

	flag.Parse()

	if *yamlMode {
		CopyAllDescriptions()
	}

	if *templateMode {
		convertTemplates()
	}

	if *handwrittenMode {
		convertAllHandwrittenFiles()
	}

	if outputPath == nil || *outputPath == "" {
		log.Printf("No output path specified, exiting")
		return
	}

	if version == nil || *version == "" {
		log.Printf("No version specified, assuming ga")
		*version = "ga"
	}

	var generateCode = !*doNotGenerateCode
	var generateDocs = !*doNotGenerateDocs
	var productsToGenerate []string
	var allProducts = false
	if product == nil || *product == "" {
		allProducts = true
	} else {
		var productToGenerate = fmt.Sprintf("products/%s", *product)
		productsToGenerate = []string{productToGenerate}
	}

	var allProductFiles []string = make([]string, 0)

	files, err := filepath.Glob("products/**/product.yaml")
	if err != nil {
		return
	}
	for _, filePath := range files {
		dir := filepath.Dir(filePath)
		allProductFiles = append(allProductFiles, fmt.Sprintf("products/%s", filepath.Base(dir)))
	}
	// TODO rewrite: override directory

	if allProducts {
		productsToGenerate = allProductFiles
	}

	if productsToGenerate == nil || len(productsToGenerate) == 0 {
		log.Fatalf("No product.yaml file found.")
	}

	startTime := time.Now()
	log.Printf("Generating MM output to '%s'", *outputPath)
	log.Printf("Using %s version", *version)
	log.Printf("Using %s provider", *forceProvider)

	// Building compute takes a long time and can't be parallelized within the product
	// so lets build it first
	sort.Slice(allProductFiles, func(i int, j int) bool {
		if allProductFiles[i] == "products/compute" {
			return true
		}
		return false
	})

	var providerToGenerate provider.Provider
	var productsForVersion []*api.Product

	ch := make(chan string, len(allProductFiles))
	for _, pf := range allProductFiles {
		ch <- pf
	}

	for i := 0; i < len(allProductFiles); i++ {
		wg.Add(1)
		go GenerateProduct(ch, providerToGenerate, &productsForVersion, startTime, productsToGenerate, *resourceToGenerate, generateCode, generateDocs)
	}
	wg.Wait()

	close(ch)

	slices.SortFunc(productsForVersion, func(p1, p2 *api.Product) int {
		return strings.Compare(strings.ToLower(p1.Name), strings.ToLower(p2.Name))
	})

	// In order to only copy/compile files once per provider this must be called outside
	// of the products loop. This will get called with the provider from the final iteration
	// of the loop
	providerToGenerate = setProvider(*forceProvider, *version, productsForVersion[0], startTime)
	providerToGenerate.CopyCommonFiles(*outputPath, generateCode, generateDocs)

	log.Printf("Compiling common files for terraform")
	if generateCode {
		providerToGenerate.CompileCommonFiles(*outputPath, productsForVersion, "")

		// TODO rewrite: product overrides
	}
}

func GenerateProduct(productChannel chan string, providerToGenerate provider.Provider, productsForVersion *[]*api.Product, startTime time.Time, productsToGenerate []string, resourceToGenerate string, generateCode, generateDocs bool) {

	defer wg.Done()
	productName := <-productChannel

	productYamlPath := path.Join(productName, "go_product.yaml")

	// TODO rewrite: uncomment the error check that if the product.yaml exists for each product
	// after Go-converted product.yaml files are complete for all products
	// if _, err := os.Stat(productYamlPath); errors.Is(err, os.ErrNotExist) {
	// 	log.Fatalf("%s does not contain a product.yaml file", productName)
	// }

	// TODO rewrite: product overrides

	if _, err := os.Stat(productYamlPath); err == nil {
		var resources []*api.Resource = make([]*api.Resource, 0)

		productApi := &api.Product{}
		api.Compile(productYamlPath, productApi)

		if !productApi.ExistsAtVersionOrLower(*version) {
			log.Printf("%s does not have a '%s' version, skipping", productName, *version)
			return
		}

		resourceFiles, err := filepath.Glob(fmt.Sprintf("%s/*", productName))
		if err != nil {
			log.Fatalf("Cannot get resources files: %v", err)
		}
		for _, resourceYamlPath := range resourceFiles {
			if filepath.Base(resourceYamlPath) == "product.yaml" || filepath.Ext(resourceYamlPath) != ".yaml" {
				continue
			}

			// Prepend "go_" to the Go yaml files' name to distinguish with the ruby yaml files
			if filepath.Base(resourceYamlPath) == "go_product.yaml" || !strings.HasPrefix(filepath.Base(resourceYamlPath), "go_") {
				continue
			}

			resource := &api.Resource{}
			api.Compile(resourceYamlPath, resource)

			resource.TargetVersionName = *version
			resource.Properties = resource.AddLabelsRelatedFields(resource.PropertiesWithExcluded(), nil)
			resource.SetDefault(productApi)
			resource.Validate()
			resources = append(resources, resource)
		}

		// TODO rewrite: override resources

		// Sort resources by name
		sort.Slice(resources, func(i, j int) bool {
			return resources[i].Name < resources[j].Name
		})

		productApi.Objects = resources
		productApi.Validate()

		providerToGenerate = setProvider(*forceProvider, *version, productApi, startTime)

		*productsForVersion = append(*productsForVersion, productApi)

		if !slices.Contains(productsToGenerate, productName) {
			log.Printf("%s not specified, skipping generation", productName)
			return
		}

		log.Printf("%s: Generating files", productName)
		providerToGenerate.Generate(*outputPath, productName, resourceToGenerate, generateCode, generateDocs)
	}
}

// Sets provider via flag
func setProvider(forceProvider, version string, productApi *api.Product, startTime time.Time) provider.Provider {
	switch forceProvider {
	case "tgc":
		return provider.NewTerraformGoogleConversion(productApi, version, startTime)
	default:
		return provider.NewTerraform(productApi, version, startTime)
	}
}
