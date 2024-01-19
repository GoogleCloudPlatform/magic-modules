package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"

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
		product_yaml_path := path.Join(product_name, "product_go.yaml")

		// TODO: uncomment the error check that if the product.yaml exists for each product
		// after Go-converted product.yaml files are complete for all products

		// if _, err := os.Stat(product_yaml_path); errors.Is(err, os.ErrNotExist) {
		// 	log.Fatalf("%s does not contain a product.yaml file", product_name)
		// }

		if _, err := os.Stat(product_yaml_path); err == nil {
			log.Printf(" product_yaml_path %#v", product_yaml_path)

			productYaml, err := os.ReadFile(product_yaml_path)
			if err != nil {
				return
			}

			productApi := api.Product{}
			yamlValidator.Parse(productYaml, &productApi)
			log.Printf(" productApi %#v", productApi)

			resourceFiles, err := filepath.Glob(fmt.Sprintf("%s/*", product_name))
			if err != nil {
				return
			}
			for _, file_path := range resourceFiles {
				if filepath.Base(file_path) == "product.yaml" || filepath.Base(file_path) == "product_go.yaml" || filepath.Ext(file_path) != ".yaml" {
					continue
				}
			}
		}
	}
}
