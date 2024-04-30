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
	"time"

	"golang.org/x/exp/slices"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/provider"
)

// TODO Q2: additional flags

// Example usage: --output $GOPATH/src/github.com/terraform-providers/terraform-provider-google-beta
var outputPath = flag.String("output", "", "path to output generated files to")

// Example usage: --version beta
var version = flag.String("version", "", "optional version name. If specified, this version is preferred for resource generation when applicable")

var product = flag.String("product", "", "optional product name. If specified, the resources under the specific product will be generated. Otherwise, resources under all products will be generated.")

func main() {
	flag.Parse()
	var generateCode = true
	var generateDocs = true

	if outputPath == nil || *outputPath == "" {
		log.Fatalf("No output path specified")
	}

	if version == nil || *version == "" {
		log.Fatalf("No version specified")
	}

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
	// TODO Q2: override directory

	if allProducts {
		productsToGenerate = allProductFiles
	}

	if productsToGenerate == nil || len(productsToGenerate) == 0 {
		log.Fatalf("No product.yaml file found.")
	}

	startTime := time.Now()
	log.Printf("Generating MM output to '%s'", *outputPath)
	log.Printf("Using %s version", *version)

	// Building compute takes a long time and can't be parallelized within the product
	// so lets build it first
	sort.Slice(allProductFiles, func(i int, j int) bool {
		if allProductFiles[i] == "compute" {
			return true
		}
		return false
	})

	var productsForVersion []map[string]interface{}
	for _, productName := range allProductFiles {
		productYamlPath := path.Join(productName, "go_product.yaml")

		// TODO Q2: uncomment the error check that if the product.yaml exists for each product
		// after Go-converted product.yaml files are complete for all products
		// if _, err := os.Stat(productYamlPath); errors.Is(err, os.ErrNotExist) {
		// 	log.Fatalf("%s does not contain a product.yaml file", productName)
		// }

		// TODO Q2: product overrides

		if _, err := os.Stat(productYamlPath); err == nil {
			var resources []*api.Resource = make([]*api.Resource, 0)

			productApi := &api.Product{}
			api.Compile(productYamlPath, productApi)

			if !productApi.ExistsAtVersionOrLower(*version) {
				log.Printf("%s does not have a '%s' version, skipping", productName, *version)
				continue
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

			// TODO Q2: override resources

			// Sort resources by name
			sort.Slice(resources, func(i, j int) bool {
				return resources[i].Name < resources[j].Name
			})

			productApi.Objects = resources
			productApi.Validate()

			// TODO Q2: set other providers via flag
			providerToGenerate := provider.NewTerraform(productApi, *version, startTime)

			if !slices.Contains(productsToGenerate, productName) {
				log.Printf("%s not specified, skipping generation", productName)
				continue
			}

			log.Printf("%s: Generating files", productName)
			providerToGenerate.Generate(*outputPath, productName, generateCode, generateDocs)

			// we need to preserve a single provider instance to use outside of this loop.
			productsForVersion = append(productsForVersion, map[string]interface{}{
				"Definitions": productApi,
				"Provider":    providerToGenerate,
			})
		}

		// TODO Q2: copy common files
	}

	slices.SortFunc(productsForVersion, func(p1, p2 map[string]interface{}) int {
		return strings.Compare(strings.ToLower(p1["Definitions"].(*api.Product).Name), strings.ToLower(p2["Definitions"].(*api.Product).Name))
	})

	// In order to only copy/compile files once per provider this must be called outside
	// of the products loop. This will get called with the provider from the final iteration
	// of the loop
	finalProduct := productsForVersion[len(productsForVersion)-1]
	provider := finalProduct["Provider"].(*provider.Terraform)

	provider.CopyCommonFiles(*outputPath, generateCode, generateDocs)

	log.Printf("Compiling common files for terraform")
	if generateCode {
		provider.CompileCommonFiles(*outputPath, productsForVersion, "")

		// TODO Q2: product overrides
	}
}
