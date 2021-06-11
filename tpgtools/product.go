package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/nasa9084/go-openapi"
)

type ProductMetadata struct {
	// PackagePath is the path to the package relative to the dcl
	PackagePath string
	// PackageName is the namespace of the package within the dcl
	// the PackageName is normally a lowercase variant of ProductName
	PackageName string
	// ProductName is the case accounted (snake case) name of the product
	// that the resource belongs to.
	ProductName string
}

var productOverrides map[string]Overrides = make(map[string]Overrides, 0)

func GetProductMetadataFromDocument(document *openapi.Document, packagePath string) *ProductMetadata {
	titleParts := strings.Split(document.Info.Title, "/")
	if len(titleParts) < 0 {
		glog.Exitf("could not find product information for %s", packagePath)
	}
	productMetadata := NewProductMetadata(packagePath, jsonToSnakeCase(titleParts[0]))
	return productMetadata
}

func NewProductMetadata(packagePath, productName string) *ProductMetadata {
	if regexp.MustCompile("[A-Z]+").Match([]byte(productName)) {
		log.Fatalln("error - expected product name to be snakecase")
	}
	if _, ok := productOverrides[packagePath]; !ok {
		productOverrides[packagePath] = loadOverrides(packagePath, "tpgtools_product.yaml")
	}
	packageName := strings.Split(packagePath, "/")[0]
	return &ProductMetadata{
		PackagePath: packagePath,
		PackageName: packageName,
		ProductName: productName,
	}
}

func (pm *ProductMetadata) WriteBasePath() bool {
	po, ok := productOverrides[pm.PackagePath]
	if !ok {
		return true
	}
	skipBasePath := po.ResourceOverride(ProductSkipBasePath, "")
	return !skipBasePath
}

// ProductType is the title-cased product name of a resource. For example,
// "NetworkServices".
func (pm *ProductMetadata) ProductType() string {
	return snakeToTitleCase(pm.ProductName)
}

// ProductNameUpper is the all caps snakecase product name of a resource.
// For example, "NETWORK_SERVICES".
func (pm *ProductMetadata) ProductNameUpper() string {
	return strings.ToUpper(pm.ProductName)
}
