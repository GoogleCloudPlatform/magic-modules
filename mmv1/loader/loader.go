package loader

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/golang/glog"
)

type Loader struct {
	baseFS     fs.FS
	overrideFS fs.FS
	version    string
	OverlayFS  google.ReadDirReadFileFS
}

type Config struct {
	BaseDirectory     string // required
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
	if config.BaseDirectory == "" {
		panic("a base directory is required")
	}
	l := &Loader{
		baseFS:  os.DirFS(config.BaseDirectory),
		version: config.Version,
	}
	log.Printf("Using base directory %q", config.BaseDirectory)
	if config.OverrideDirectory != "" {
		log.Printf("Using override directory %q", config.OverrideDirectory)
		l.overrideFS = os.DirFS(config.OverrideDirectory)
	}
	var err error
	l.OverlayFS, err = google.NewOverlayFS(l.overrideFS, l.baseFS)
	if err != nil {
		panic(err)
	}
	return l
}

func (l *Loader) LoadProducts() map[string]*api.Product {
	if l.version == "" {
		log.Printf("No version specified, assuming ga")
		l.version = "ga"
	}

	var allProductFiles []string = make([]string, 0)

	files, err := fs.Glob(l.baseFS, "products/**/product.yaml")
	if err != nil {
		panic(err)
	}
	for _, filePath := range files {
		product := filepath.Base(filepath.Dir(filePath))
		allProductFiles = append(allProductFiles, fmt.Sprintf("products/%s", product))
	}

	if l.overrideFS != nil {
		overrideFiles, err := fs.Glob(l.overrideFS, "products/**/product.yaml")
		if err != nil {
			panic(err)
		}
		for _, filePath := range overrideFiles {
			product := filepath.Base(filepath.Dir(filePath))
			productFile := fmt.Sprintf("products/%s", product)
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
// productName looks like `products/foo`
func (l *Loader) LoadProduct(productName string) (*api.Product, error) {
	productYamlPath := filepath.Join(productName, "product.yaml")

	var baseContents, overrideContents []byte
	baseContents, _ = fs.ReadFile(l.baseFS, productYamlPath)
	if l.overrideFS != nil {
		overrideContents, _ = fs.ReadFile(l.overrideFS, productYamlPath)
	}
	if baseContents == nil && overrideContents == nil {
		return nil, fmt.Errorf("%s does not contain a product.yaml file", productName)
	}

	var p, overrideProduct *api.Product
	if overrideContents != nil {
		overrideProduct = &api.Product{}
		api.CompileContents(overrideContents, overrideProduct, fmt.Sprintf("${OVERRIDE}/%s", productYamlPath))
	}
	if baseContents != nil {
		p = &api.Product{}
		api.CompileContents(baseContents, p, fmt.Sprintf("${BASE}/%s", productYamlPath))
		if overrideProduct != nil {
			api.Merge(reflect.ValueOf(p).Elem(), reflect.ValueOf(*overrideProduct), l.version)
		}
	} else {
		p = overrideProduct
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

	p.Objects = resources
	p.Validate()

	return p, nil
}

type contents struct {
	base, override []byte
}

// loadResources loads all resources for a product
func (l *Loader) loadResources(product *api.Product) ([]*api.Resource, error) {
	var resources []*api.Resource = make([]*api.Resource, 0)

	resourceContents := make(map[string]*contents)
	if l.overrideFS != nil {
		overrideFiles, err := fs.Glob(l.overrideFS, filepath.Join(product.PackagePath, "*"))
		if err != nil {
			return nil, fmt.Errorf("cannot get override files: %v", err)
		}

		for _, overrideYamlPath := range overrideFiles {
			base := filepath.Base(overrideYamlPath)
			if base == "product.yaml" || filepath.Ext(overrideYamlPath) != ".yaml" {
				continue
			}
			full := filepath.Join(product.PackagePath, base)
			c, err := fs.ReadFile(l.overrideFS, full)
			if err != nil {
				return nil, fmt.Errorf("cannot read override file %q: %v", full, err)
			}
			resourceContents[base] = &contents{override: c}
		}
	}
	// Get base resource files
	resourceFiles, err := fs.Glob(l.baseFS, filepath.Join(product.PackagePath, "*"))
	if err != nil {
		return nil, fmt.Errorf("cannot get base resource files: %v", err)
	}
	for _, resourceYamlPath := range resourceFiles {
		base := filepath.Base(resourceYamlPath)
		if base == "product.yaml" || filepath.Ext(resourceYamlPath) != ".yaml" {
			continue
		}
		full := filepath.Join(product.PackagePath, base)
		c, err := fs.ReadFile(l.baseFS, full)
		if err != nil {
			return nil, fmt.Errorf("cannot read base file %q: %v", full, err)
		}
		p, ok := resourceContents[base]
		if !ok {
			resourceContents[base] = &contents{base: c}
		} else {
			p.base = c
		}
	}
	for name, c := range resourceContents {
		resources = append(resources, l.loadResource(product, name, c))
	}
	// Sort resources by name for consistent output
	slices.SortFunc(resources, func(a, b *api.Resource) int {
		return strings.Compare(a.Name, b.Name)
	})

	return resources, nil
}

// loadResource loads a single resource with optional override
func (l *Loader) loadResource(product *api.Product, name string, c *contents) *api.Resource {
	resource := &api.Resource{}
	resource.SourceYamlFile = filepath.Join(product.PackagePath, name)

	if c.override != nil {
		if c.base != nil {
			// Merge base and override
			api.CompileContents(c.base, resource, fmt.Sprintf("${BASE}/%s", resource.SourceYamlFile))
			overrideResource := &api.Resource{}
			api.CompileContents(c.override, overrideResource, fmt.Sprintf("${OVERRIDE}/%s", resource.SourceYamlFile))
			api.Merge(reflect.ValueOf(resource).Elem(), reflect.ValueOf(*overrideResource), l.version)
		} else {
			// Override only
			api.CompileContents(c.override, resource, fmt.Sprintf("${OVERRIDE}/%s", resource.SourceYamlFile))
		}
	} else {
		// Base only
		api.CompileContents(c.base, resource, fmt.Sprintf("${BASE}/%s", resource.SourceYamlFile))
	}

	// Set resource defaults and validate
	resource.TargetVersionName = l.version
	// SetDefault before AddExtraFields to ensure relevant metadata is available on existing fields
	resource.SetDefault(product)
	resource.Properties = resource.AddExtraFields(resource.PropertiesWithExcluded(), nil)
	// SetDefault after AddExtraFields to ensure relevant metadata is available for the newly generated fields
	resource.SetDefault(product)
	resource.Validate()
	resource.TestSampleSetUp(l.OverlayFS)

	for _, e := range resource.Examples {
		if err := e.LoadHCLText(l.OverlayFS); err != nil {
			glog.Exit(err)
		}
	}

	return resource
}
