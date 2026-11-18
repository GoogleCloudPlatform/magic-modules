package provider

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
)

type Provider interface {
	Generate(string, string, bool, bool)
	CopyCommonFiles(outputFolder string, generateCode, generateDocs bool)
	CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string)
}

// Shared constants and functions among the providers

const TERRAFORM_PROVIDER_GA = "github.com/hashicorp/terraform-provider-google"
const TERRAFORM_PROVIDER_BETA = "github.com/hashicorp/terraform-provider-google-beta"
const TGC_PROVIDER = "github.com/GoogleCloudPlatform/terraform-google-conversion/v7"
const RESOURCE_DIRECTORY_GA = "google"
const RESOURCE_DIRECTORY_BETA = "google-beta"
const RESOURCE_DIRECTORY_TGC = "pkg"

// # TODO: Review all object interfaces and move to private methods
// # that should not be exposed outside the object hierarchy.
func ProviderName(t Provider) string {
	return reflect.TypeOf(t).Name()
}

func ImportPathFromVersion(v string) string {
	var tpg, dir string
	switch v {
	case "ga":
		tpg = TERRAFORM_PROVIDER_GA
		dir = RESOURCE_DIRECTORY_GA
	case "beta":
		tpg = TERRAFORM_PROVIDER_BETA
		dir = RESOURCE_DIRECTORY_BETA
	default:
		tpg = "github.com/hashicorp/terraform-provider-google-" + v
		dir = "google-" + v
	}
	return fmt.Sprintf("%s/%s", tpg, dir)
}

func copyDir(srcDir, dstDir string, skip func(src string) bool) error {
	return filepath.WalkDir(srcDir, func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if skip != nil && skip(src) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(srcDir, src)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstDir, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(dst, info.Mode().Perm())
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, info.Mode().Perm())
	})
}
