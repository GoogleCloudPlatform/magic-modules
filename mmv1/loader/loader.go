package loader

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"golang.org/x/exp/slices"
)

type Loader struct {
	// BaseDirectory points to mmv1 root, if cwd can be empty as relative paths are used
	BaseDirectory     string
	OverrideDirectory string
	Version           string
}

type Config struct {
	BaseDirectory     string // optional, defaults to current working directory
	OverrideDirectory string // optional
	Version           string // required
}

// NewLoader creates a new Loader instance, applying any
// provided options.
func NewLoader(config Config) *Loader {
	// Validation
	if config.Version == "" {
		panic("version is required")
	}

	l := &Loader{
		BaseDirectory:     config.BaseDirectory,
		OverrideDirectory: config.OverrideDirectory,
		Version:           config.Version,
	}

	// Normalize override dir to a path that is relative to the magic-modules directory
	// This is needed for templates that concatenate pwd + override dir + path
	if filepath.IsAbs(l.OverrideDirectory) {
		mmv1Dir := l.BaseDirectory
		if mmv1Dir == "" {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			mmv1Dir = wd
		}
		l.OverrideDirectory, _ = filepath.Rel(mmv1Dir, l.OverrideDirectory)
		log.Printf("Override directory normalized to relative path %s", l.OverrideDirectory)
	}

	return l
}

func (l *Loader) LoadProducts() map[string]*api.Product {
	if l.Version == "" {
		log.Printf("No version specified, assuming ga")
		l.Version = "ga"
	}

	var allProductFiles []string = make([]string, 0)

	files, err := filepath.Glob(filepath.Join(l.BaseDirectory, "products/**/product.yaml"))
	if err != nil {
		panic(err)
	}
	for _, filePath := range files {
		dir := filepath.Dir(filePath)
		allProductFiles = append(allProductFiles, fmt.Sprintf("products/%s", filepath.Base(dir)))
	}

	if l.OverrideDirectory != "" {
		log.Printf("Using override directory %s", l.OverrideDirectory)
		overrideFiles, err := filepath.Glob(filepath.Join(l.OverrideDirectory, "products/**/product.yaml"))
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

			product, err := l.LoadProduct(name)
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
			var versionErr *ErrProductVersionNotFound
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

// Load compiles a product with all its resources from the given path and optional overrides
// This loads the product configuration and all its resources into memory without generating any files
func (l *Loader) LoadProduct(productName string) (*api.Product, error) {
	p := &api.Product{}
	productYamlPath := filepath.Join(productName, "product.yaml")

	var productOverridePath string
	if l.OverrideDirectory != "" {
		productOverridePath = filepath.Join(l.OverrideDirectory, productYamlPath)
	}

	baseProductPath := filepath.Join(l.BaseDirectory, productYamlPath)

	baseProductExists := Exists(baseProductPath)
	overrideProductExists := Exists(productOverridePath)

	if !(baseProductExists || overrideProductExists) {
		return nil, fmt.Errorf("%s does not contain a product.yaml file", productName)
	}

	// Compile the product configuration
	if overrideProductExists {
		if baseProductExists {
			api.Compile(baseProductPath, p, l.OverrideDirectory)
			overrideApiProduct := &api.Product{}
			api.Compile(productOverridePath, overrideApiProduct, l.OverrideDirectory)
			api.Merge(reflect.ValueOf(p).Elem(), reflect.ValueOf(*overrideApiProduct), l.Version)
		} else {
			api.Compile(productOverridePath, p, l.OverrideDirectory)
		}
	} else {
		api.Compile(baseProductPath, p, l.OverrideDirectory)
	}

	// Check if product exists at the requested l.Version
	if !p.ExistsAtVersionOrLower(l.Version) {
		return nil, &ErrProductVersionNotFound{ProductName: productName, Version: l.Version}
	}

	// Compile all resources
	p.PackagePath = productName
	resources, err := l.loadResources(p)
	if err != nil {
		return nil, err
	}

	p.Objects = resources
	p.Validate()

	return p, nil
}

// loadResources loads all resources for a product
func (l *Loader) loadResources(product *api.Product) ([]*api.Resource, error) {
	var resources []*api.Resource = make([]*api.Resource, 0)

	// Get base resource files
	resourceFiles, err := filepath.Glob(filepath.Join(l.BaseDirectory, product.PackagePath, "*"))
	if err != nil {
		return nil, fmt.Errorf("cannot get resource files: %v", err)
	}

	// Compile base resources (skip those that will be merged with overrides)
	for _, resourceYamlPath := range resourceFiles {
		if filepath.Base(resourceYamlPath) == "product.yaml" || filepath.Ext(resourceYamlPath) != ".yaml" {
			continue
		}

		// Skip if resource will be merged in the override loop
		if l.OverrideDirectory != "" {
			overrideResourceExists := Exists(l.OverrideDirectory, resourceYamlPath)
			if overrideResourceExists {
				continue
			}
		}

		resource := l.loadResource(product, resourceYamlPath, "")
		resources = append(resources, resource)
	}

	// Compile override resources
	if l.OverrideDirectory != "" {
		resources, err = l.reconcileOverrideResources(product, resources)
		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}

// reconcileOverrideResources handles resolution of override resources
func (l *Loader) reconcileOverrideResources(product *api.Product, resources []*api.Resource) ([]*api.Resource, error) {
	productOverridePath := filepath.Join(l.OverrideDirectory, product.PackagePath, "product.yaml")
	productOverrideDir := filepath.Dir(productOverridePath)

	overrideFiles, err := filepath.Glob(filepath.Join(productOverrideDir, "*"))
	if err != nil {
		return nil, fmt.Errorf("cannot get override files: %v", err)
	}

	for _, overrideYamlPath := range overrideFiles {
		if filepath.Base(overrideYamlPath) == "product.yaml" || filepath.Ext(overrideYamlPath) != ".yaml" {
			continue
		}

		baseResourcePath := filepath.Join(product.PackagePath, filepath.Base(overrideYamlPath))
		resource := l.loadResource(product, baseResourcePath, overrideYamlPath)
		resources = append(resources, resource)
	}

	// Sort resources by name for consistent output
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	return resources, nil
}

// loadResource loads a single resource with optional override
func (l *Loader) loadResource(product *api.Product, baseResourcePath string, overrideResourcePath string) *api.Resource {
	resource := &api.Resource{}

	// Check if base resource exists
	baseResourceExists := Exists(l.BaseDirectory, baseResourcePath)

	if overrideResourcePath != "" {
		if baseResourceExists {
			// Merge base and override
			api.Compile(baseResourcePath, resource, l.OverrideDirectory)
			overrideResource := &api.Resource{}
			api.Compile(overrideResourcePath, overrideResource, l.OverrideDirectory)
			api.Merge(reflect.ValueOf(resource).Elem(), reflect.ValueOf(*overrideResource), l.Version)
			resource.SourceYamlFile = baseResourcePath
		} else {
			// Override only
			api.Compile(overrideResourcePath, resource, l.OverrideDirectory)
		}
	} else {
		// Base only
		api.Compile(baseResourcePath, resource, l.OverrideDirectory)
		resource.SourceYamlFile = baseResourcePath
	}

	// Set resource defaults and validate
	resource.TargetVersionName = l.Version
	// SetDefault before AddExtraFields to ensure relevant metadata is available on existing fields
	resource.SetDefault(product)
	resource.Properties = resource.AddExtraFields(resource.PropertiesWithExcluded(), nil)
	// SetDefault after AddExtraFields to ensure relevant metadata is available for the newly generated fields
	resource.SetDefault(product)
	resource.Validate()

	for _, e := range resource.Examples {
		e.LoadHCLText(l.BaseDirectory)
	}

	return resource
}
