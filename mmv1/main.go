package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/creasty/defaults"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/provider"
)

func main() {
	// TODO Q2: parse flags
	var version = "beta"
	var outputPath = "."
	var generateCode = true
	var generateDocs = true

	log.Printf("Initiating go MM compiler")

	// TODO Q1: allow specifying one product (flag or hardcoded)
	// var productsToGenerate []string
	// var allProducts = true
	var productsToGenerate = []string{"products/datafusion"}
	var allProducts = false

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

	log.Printf("Generating MM output to '%s'", outputPath)
	log.Printf("Using %s version", version)

	// Building compute takes a long time and can't be parallelized within the product
	// so lets build it first
	sort.Slice(allProductFiles, func(i int, j int) bool {
		if allProductFiles[i] == "compute" {
			return true
		}
		return false
	})

	yamlValidator := google.YamlValidator{}

	for _, productName := range allProductFiles {
		productYamlPath := path.Join(productName, "go_product.yaml")

		// TODO Q2: uncomment the error check that if the product.yaml exists for each product
		// after Go-converted product.yaml files are complete for all products
		// if _, err := os.Stat(productYamlPath); errors.Is(err, os.ErrNotExist) {
		// 	log.Fatalf("%s does not contain a product.yaml file", productName)
		// }

		// TODO Q2: product overrides

		if _, err := os.Stat(productYamlPath); err == nil {
			// TODO Q1: remove these lines, which are for debugging
			// log.Printf("productYamlPath %#v", productYamlPath)

			var resources []*api.Resource = make([]*api.Resource, 0)

			productYaml, err := os.ReadFile(productYamlPath)
			if err != nil {
				log.Fatalf("Cannot open the file: %v", productYaml)
			}
			productApi := &api.Product{}
			// Set default value of fields
			defaults.Set(productApi)
			yamlValidator.Parse(productYaml, productApi)

			// TODO Q1: remove these lines, which are for debugging
			prod, _ := json.Marshal(productApi)
			log.Printf("prod %s", string(prod))

			if !productApi.ExistsAtVersionOrLower(version) {
				log.Printf("%s does not have a '%s' version, skipping", productName, version)
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

				// TODO Q1: remove these lines, which are for debugging
				// log.Printf(" resourceYamlPath %s", resourceYamlPath)
				resourceYaml, err := os.ReadFile(resourceYamlPath)
				if err != nil {
					log.Fatalf("Cannot open the file: %v", resourceYamlPath)
				}
				resource := &api.Resource{}
				yamlValidator.Parse(resourceYaml, resource)

				// TODO Q1: remove these lines, which are for debugging
				// res, _ := json.Marshal(resource)
				// log.Printf("resource %s", string(res))

				// TODO Q1: add labels related fields

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
			providerToGenerate := provider.NewTerraform(productApi)

			if !slices.Contains(productsToGenerate, productName) {
				log.Printf("%s not specified, skipping generation", productName)
				continue
			}

			log.Printf("%s: Generating files", productName)
			providerToGenerate.Generate(outputPath, productName, generateCode, generateDocs)
		}

		// TODO Q2: copy common files
	}
}
