package loader

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"golang.org/x/exp/slices"
)

type Loader struct {
	OverrideDirectory string
	Version           string
}

func (l *Loader) LoadProducts() map[string]*api.Product {
	if l.Version == "" {
		log.Printf("No version specified, assuming ga")
		l.Version = "ga"
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

	if l.OverrideDirectory != "" {
		log.Printf("Using override directory %s", l.OverrideDirectory)

		// Normalize override dir to a path that is relative to the magic-modules directory
		// This is needed for templates that concatenate pwd + override dir + path
		if filepath.IsAbs(l.OverrideDirectory) {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			l.OverrideDirectory, err = filepath.Rel(wd, l.OverrideDirectory)
			log.Printf("Override directory normalized to relative path %s", l.OverrideDirectory)
		}

		overrideFiles, err := filepath.Glob(fmt.Sprintf("%s/products/**/product.yaml", l.OverrideDirectory))
		if err != nil {
			panic(err)
		}
		for _, filePath := range overrideFiles {
			product, err := filepath.Rel(l.OverrideDirectory, filePath)
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

	return l.batchLoadProducts(allProductFiles)
}

func (l *Loader) batchLoadProducts(productNames []string) map[string]*api.Product {
	products := make(map[string]*api.Product)

	// Create result type for clarity
	type loadResult struct {
		name    string
		product *api.Product
		err     error
	}

	// Buffered channel to prevent goroutine blocking
	productChan := make(chan loadResult, len(productNames))

	// Use WaitGroup for proper synchronization
	var wg sync.WaitGroup

	// Launch all goroutines
	for _, productName := range productNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			product, err := l.loadProductOnly(name)
			productChan <- loadResult{
				name:    name,
				product: product,
				err:     err,
			}
		}(productName)
	}

	wg.Wait()
	close(productChan)

	// Collect results as they complete
	loadFailureCount := 0
	for result := range productChan {
		if result.err != nil {
			// Check if the error is the specific "version not found" error
			var versionErr *api.ErrProductVersionNotFound
			if errors.As(result.err, &versionErr) {
				continue
			}

			loadFailureCount++
			log.Printf("Error loading %s: %v", result.name, result.err)
			continue
		}
		products[result.name] = result.product
	}
	if loadFailureCount > 0 {
		log.Fatalf("Failed to load %d products", loadFailureCount)
	}

	return products
}

// loadProductOnly is a standalone function to just load a product without generation
// This can be used when you only need to load and validate a product configuration
func (l *Loader) loadProductOnly(productName string) (*api.Product, error) {
	product := &api.Product{}
	err := product.Load(productName, l.Version, l.OverrideDirectory)
	if err != nil {
		return nil, err
	}
	return product, nil
}
