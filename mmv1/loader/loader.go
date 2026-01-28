package loader

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/golang/glog"
	"golang.org/x/exp/slices"
)

type Loader struct {
	// baseDirectory points to mmv1 root, if cwd can be empty as relative paths are used
	baseDirectory     string
	overrideDirectory string
	Products          map[string]*api.Product
	version           string
	sysfs             google.ReadDirReadFileFS
}

type Config struct {
	BaseDirectory     string                   // required
	OverrideDirectory string                   // optional
	Version           string                   // required
	Sysfs             google.ReadDirReadFileFS // required
}

// NewLoader creates a new Loader instance, applying any
// provided options.
func NewLoader(config Config) *Loader {
	// Validation
	if config.Version == "" {
		panic("version is required")
	}
	if config.BaseDirectory == "" {
		panic("a base directory is required")
	}
	if config.Sysfs == nil {
		panic("sysfs is required")
	}
	l := &Loader{
		baseDirectory:     config.BaseDirectory,
		overrideDirectory: config.OverrideDirectory,
		version:           config.Version,
		sysfs:             config.Sysfs,
	}

	return l
}

func (l *Loader) LoadProducts() {
	if l.version == "" {
		log.Printf("No version specified, assuming ga")
		l.version = "ga"
	}

	var allProductFiles []string = make([]string, 0)

	files, err := filepath.Glob(filepath.Join(l.baseDirectory, "products/**/product.yaml"))
	if err != nil {
		panic(err)
	}
	for _, filePath := range files {
		dir := filepath.Dir(filePath)
		allProductFiles = append(allProductFiles, fmt.Sprintf("products/%s", filepath.Base(dir)))
	}

	log.Printf("Using base directory %q", l.baseDirectory)
	if l.overrideDirectory != "" {
		log.Printf("Using override directory %q", l.overrideDirectory)
		overrideFiles, err := filepath.Glob(filepath.Join(l.overrideDirectory, "products/**/product.yaml"))
		if err != nil {
			panic(err)
		}
		for _, filePath := range overrideFiles {
			product, err := filepath.Rel(l.overrideDirectory, filePath)
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

	l.Products = l.batchLoadProducts(allProductFiles)
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
	if l.overrideDirectory != "" {
		productOverridePath = filepath.Join(l.overrideDirectory, productYamlPath)
	}

	baseProductPath := filepath.Join(l.baseDirectory, productYamlPath)

	baseProductExists := Exists(baseProductPath)
	overrideProductExists := Exists(productOverridePath)

	if !(baseProductExists || overrideProductExists) {
		return nil, fmt.Errorf("%s does not contain a product.yaml file", productName)
	}

	// Compile the product configuration
	if overrideProductExists {
		if baseProductExists {
			api.Compile(baseProductPath, p)
			overrideApiProduct := &api.Product{}
			api.Compile(productOverridePath, overrideApiProduct)
			api.Merge(reflect.ValueOf(p).Elem(), reflect.ValueOf(*overrideApiProduct), l.version)
		} else {
			api.Compile(productOverridePath, p)
		}
	} else {
		api.Compile(baseProductPath, p)
	}

	// Check if product exists at the requested l.Version
	if !p.ExistsAtVersionOrLower(l.version) {
		return nil, &ErrProductVersionNotFound{ProductName: productName, Version: l.version}
	}

	// Compile all resources
	p.PackagePath = productName
	resources, err := l.loadResources(p)
	if err != nil {
		return nil, err
	}

	p.Version = p.VersionObjOrClosest(l.version)

	p.Objects = resources
	p.Validate()

	return p, nil
}

// loadResources loads all resources for a product
func (l *Loader) loadResources(product *api.Product) ([]*api.Resource, error) {
	var resources []*api.Resource = make([]*api.Resource, 0)

	// Get base resource files
	resourceFiles, err := filepath.Glob(filepath.Join(l.baseDirectory, product.PackagePath, "*"))
	if err != nil {
		return nil, fmt.Errorf("cannot get resource files: %v", err)
	}

	// Compile base resources (skip those that will be merged with overrides)
	for _, resourceYamlPath := range resourceFiles {
		if filepath.Base(resourceYamlPath) == "product.yaml" || filepath.Ext(resourceYamlPath) != ".yaml" {
			continue
		}
		relPath, err := filepath.Rel(l.baseDirectory, resourceYamlPath)
		if err != nil {
			return nil, fmt.Errorf("returned %q is not relative to %q", resourceYamlPath, l.baseDirectory)
		}

		// Skip if resource will be merged in the override loop
		if l.overrideDirectory != "" {
			overrideResourceExists := Exists(l.overrideDirectory, relPath)
			if overrideResourceExists {
				continue
			}
		}

		resource := l.loadResource(product, resourceYamlPath, "")
		resources = append(resources, resource)
	}

	// Compile override resources
	if l.overrideDirectory != "" {
		resources, err = l.reconcileOverrideResources(product, resources)
		if err != nil {
			return nil, err
		}
	}
	// Sort resources by name for consistent output
	slices.SortFunc(resources, func(a, b *api.Resource) int {
		return strings.Compare(a.Name, b.Name)
	})

	return resources, nil
}

// reconcileOverrideResources handles resolution of override resources
func (l *Loader) reconcileOverrideResources(product *api.Product, resources []*api.Resource) ([]*api.Resource, error) {
	overrideFiles, err := filepath.Glob(filepath.Join(l.overrideDirectory, product.PackagePath, "*"))
	if err != nil {
		return nil, fmt.Errorf("cannot get override files: %v", err)
	}

	for _, overrideYamlPath := range overrideFiles {
		if filepath.Base(overrideYamlPath) == "product.yaml" || filepath.Ext(overrideYamlPath) != ".yaml" {
			continue
		}

		baseResourcePath := filepath.Join(l.baseDirectory, product.PackagePath, filepath.Base(overrideYamlPath))
		resource := l.loadResource(product, baseResourcePath, overrideYamlPath)
		resources = append(resources, resource)
	}

	return resources, nil
}

// loadResource loads a single resource with optional override
// baseResourcePath and overrideResourcePath are expected to be absolute paths.
func (l *Loader) loadResource(product *api.Product, baseResourcePath string, overrideResourcePath string) *api.Resource {
	resource := &api.Resource{}

	// Check if base resource exists
	baseResourceExists := Exists(baseResourcePath)
	baseRelPath, _ := filepath.Rel(l.baseDirectory, baseResourcePath)

	if baseResourceExists {
		resource.SourceYamlFile = baseRelPath
	} else {
		relPath, _ := filepath.Rel(l.overrideDirectory, overrideResourcePath)
		resource.SourceYamlFile = relPath
	}

	if overrideResourcePath != "" {
		if baseResourceExists {
			// Merge base and override
			api.Compile(baseResourcePath, resource)
			overrideResource := &api.Resource{}
			api.Compile(overrideResourcePath, overrideResource)
			api.Merge(reflect.ValueOf(resource).Elem(), reflect.ValueOf(*overrideResource), l.version)
		} else {
			// Override only
			api.Compile(overrideResourcePath, resource)
		}
	} else {
		// Base only
		api.Compile(baseResourcePath, resource)
		resource.SourceYamlFile = baseRelPath
	}

	// Set resource defaults and validate
	resource.TargetVersionName = l.version
	// SetDefault before AddExtraFields to ensure relevant metadata is available on existing fields
	resource.SetDefault(product)
	resource.TestSampleSetUp(l.sysfs)

	for _, e := range resource.Examples {
		if err := e.LoadHCLText(l.sysfs); err != nil {
			glog.Exit(err)
		}
	}

	return resource
}

func (l *Loader) AddExtraFields() error {
	if l.Products == nil {
		return errors.New("products have not been loaded into memory")
	}

	for _, product := range l.Products {
		for _, resource := range product.Objects {
			resource.Properties = resource.AddExtraFields(resource.PropertiesWithExcluded(), nil)
			// SetDefault after AddExtraFields to ensure relevant metadata is available for the newly generated fields
			resource.SetDefault(product)
		}
	}

	return nil
}

func (l *Loader) Validate() {
	if l.Products == nil {
		log.Fatalln("products have not been loaded into memory")
	}

	for _, product := range l.Products {
		for _, resource := range product.Objects {
			resource.Validate()
		}
	}
}
