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

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

func main() {
	var products_to_generate []string
	var all_products = true

	var all_product_files []string = make([]string, 0)

	files, err := filepath.Glob("products/**/product.yaml")
	if err != nil {
		return
	}
	for _, file_path := range files {
		dir := filepath.Dir(file_path)
		all_product_files = append(all_product_files, fmt.Sprintf("products/%s", filepath.Base(dir)))
	}

	if all_products {
		products_to_generate = all_product_files
	}

	if products_to_generate == nil || len(products_to_generate) == 0 {
		log.Fatalf("No product.yaml file found.")
	}

	// Building compute takes a long time and can't be parallelized within the product
	// so lets build it first
	sort.Slice(all_product_files, func(i int, j int) bool {
		if all_product_files[i] == "compute" {
			return true
		}
		return false
	})

	yamlValidator := google.YamlValidator{}

	for _, product_name := range all_product_files {
		product_yaml_path := path.Join(product_name, "go_product.yaml")

		// TODO: uncomment the error check that if the product.yaml exists for each product
		// after Go-converted product.yaml files are complete for all products

		// if _, err := os.Stat(product_yaml_path); errors.Is(err, os.ErrNotExist) {
		// 	log.Fatalf("%s does not contain a product.yaml file", product_name)
		// }

		if _, err := os.Stat(product_yaml_path); err == nil {
			log.Printf("product_yaml_path %#v", product_yaml_path)

			productYaml, err := os.ReadFile(product_yaml_path)
			if err != nil {
				log.Fatalf("Cannot open the file: %v", productYaml)
			}
			productApi := api.Product{}
			yamlValidator.Parse(productYaml, &productApi)

			// TODO: remove these lines, which are for debugging
			prod, _ := json.Marshal(&productApi)
			log.Printf("prod %s", string(prod))

			resourceFiles, err := filepath.Glob(fmt.Sprintf("%s/*", product_name))
			if err != nil {
				log.Fatalf("Cannot get resources files: %v", err)
			}
			for _, resourceYamlPath := range resourceFiles {
				if filepath.Base(resourceYamlPath) == "product.yaml" || filepath.Ext(resourceYamlPath) != ".yaml" {
					continue
				}

				// TODO REMOVE: limiting test block
				// if !strings.Contains(resourceYamlPath, "datafusion") {
				// 	continue
				// }

				// Prepend "go_" to the Go yaml files' name to distinguish with the ruby yaml files
				if filepath.Base(resourceYamlPath) == "go_product.yaml" || !strings.HasPrefix(filepath.Base(resourceYamlPath), "go_") {
					continue
				}

				log.Printf(" resourceYamlPath %s", resourceYamlPath)
				resourceYaml, err := os.ReadFile(resourceYamlPath)
				if err != nil {
					log.Fatalf("Cannot open the file: %v", resourceYamlPath)
				}
				resource := api.Resource{}
				yamlValidator.Parse(resourceYaml, &resource)

				// TODO: remove these lines, which are for debugging
				res, _ := json.Marshal(&resource)
				log.Printf("resource %s", string(res))
			}
		}
	}
}
