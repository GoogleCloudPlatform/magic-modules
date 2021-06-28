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

func (pm *ProductMetadata) ShouldWriteProductBasePath() bool {
	bp := pm.ProductBasePathDetails()
	if bp == nil {
		return true
	}
	return !bp.Skip
}

func (pm *ProductMetadata) BasePathIdentifierSnakeUpper() string {
	return strings.ToUpper(pm.BasePathIdentifierSnake())
}

func (pm *ProductMetadata) BasePathIdentifierSnake() string {
	bp := pm.ProductBasePathDetails()
	if bp != nil && bp.BasePathIdentifier != ""{
		return bp.BasePathIdentifier
	}
	return pm.ProductName
}

func (pm *ProductMetadata) BasePathIdentifier() string {
	bp := pm.ProductBasePathDetails()
	if bp != nil && bp.BasePathIdentifier != ""{
		return snakeToTitleCase(bp.BasePathIdentifier)
	}
	return pm.ProductType()
}

func (pm *ProductMetadata) ProductBasePathDetails() *ProductBasePathDetails {
	overrides, ok := productOverrides[pm.PackagePath]
	if !ok {
		// TODO maybe crash here?
		return nil
	}
	bp := ProductBasePathDetails{}
	bpOk, err := overrides.ProductOverrideWithDetails(ProductBasePath, &bp)
	if err != nil {
		log.Fatalln("error - failed to decode base path details")
	}

	if !bpOk {
		return nil
	}

	return &bp
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

// DCLPackage is the package name of the DCL client library to use for this
// resource. For example, the Package "access_context_manager" would have a
// DCLPackage of "accesscontextmanager"
func (pm *ProductMetadata) DCLPackage() string {
	return strings.Replace(pm.PackagePath, "_", "", -1)
}
