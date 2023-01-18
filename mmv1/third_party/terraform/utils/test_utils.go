package google

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

var ProjectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
}

var CredsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var orgEnvVars = []string{
	"GOOGLE_ORG",
}

var RegionEnvVars = []string{
	"GOOGLE_REGION",
	"GCLOUD_REGION",
	"CLOUDSDK_COMPUTE_REGION",
}

var ZoneEnvVars = []string{
	"GOOGLE_ZONE",
	"GCLOUD_ZONE",
	"CLOUDSDK_COMPUTE_ZONE",
}

var Sources map[string]VcrSource

func init() {
	Sources = make(map[string]VcrSource)
}

// TestAccPreCheck ensures at least one of the project env variables is set.
func GetTestProjectFromEnv() string {
	return MultiEnvSearch(ProjectEnvVars)
}

func GetTestOrgFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, orgEnvVars...)
	return MultiEnvSearch(orgEnvVars)
}

// TestAccPreCheck ensures at least one of the credentials env variables is set.
func GetTestCredsFromEnv() string {
	// Return empty string if GOOGLE_USE_DEFAULT_CREDENTIALS is set to true.
	if MultiEnvSearch(CredsEnvVars) == "true" {
		return ""
	}
	return MultiEnvSearch(CredsEnvVars)
}

// TestAccPreCheck ensures at least one of the region env variables is set.
func GetTestRegionFromEnv() string {
	return MultiEnvSearch(RegionEnvVars)
}

func GetTestZoneFromEnv() string {
	return MultiEnvSearch(ZoneEnvVars)
}

var (
	Pname          = "Terraform Acceptance Tests"
	originalPolicy *cloudresourcemanager.Policy
	testPrefix     = "tf-test"
)

type ResourceDataMock struct {
	FieldsInSchema      map[string]interface{}
	FieldsWithHasChange []string
	id                  string
}

func (d *ResourceDataMock) HasChange(key string) bool {
	exists := false
	for _, val := range d.FieldsWithHasChange {
		if key == val {
			exists = true
		}
	}

	return exists
}

func (d *ResourceDataMock) Get(key string) interface{} {
	v, _ := d.GetOk(key)
	return v
}

func (d *ResourceDataMock) GetOk(key string) (interface{}, bool) {
	v, ok := d.GetOkExists(key)
	if ok && !IsEmptyValue(reflect.ValueOf(v)) {
		return v, true
	} else {
		return v, false
	}
}

func (d *ResourceDataMock) GetOkExists(key string) (interface{}, bool) {
	for k, v := range d.FieldsInSchema {
		if key == k {
			return v, true
		}
	}

	return nil, false
}

func (d *ResourceDataMock) Set(key string, value interface{}) error {
	d.FieldsInSchema[key] = value
	return nil
}

func (d *ResourceDataMock) SetId(v string) {
	d.id = v
}

func (d *ResourceDataMock) Id() string {
	return d.id
}

func (d *ResourceDataMock) GetProviderMeta(dst interface{}) error {
	return nil
}

func (d *ResourceDataMock) Timeout(key string) time.Duration {
	return time.Duration(1)
}

type ResourceDiffMock struct {
	Before     map[string]interface{}
	After      map[string]interface{}
	Cleared    map[string]interface{}
	IsForceNew bool
}

func (d *ResourceDiffMock) GetChange(key string) (interface{}, interface{}) {
	return d.Before[key], d.After[key]
}

func (d *ResourceDiffMock) HasChange(key string) bool {
	old, new := d.GetChange(key)
	return old != new
}

func (d *ResourceDiffMock) Get(key string) interface{} {
	return d.After[key]
}

func (d *ResourceDiffMock) GetOk(key string) (interface{}, bool) {
	v, ok := d.After[key]
	return v, ok
}

func (d *ResourceDiffMock) Clear(key string) error {
	if d.Cleared == nil {
		d.Cleared = map[string]interface{}{}
	}
	d.Cleared[key] = true
	return nil
}

func (d *ResourceDiffMock) ForceNew(key string) error {
	d.IsForceNew = true
	return nil
}

func CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName string) func(*terraform.State) error {
	return checkDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName, map[string]struct{}{})
}

func checkDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		errMsg := ""
		// Data Sources are often derived from resources, so iterate over the resource fields to
		// make sure all fields are accounted for in the data source.
		// If a field exists in the data source but not in the resource, its expected value should
		// be checked separately.
		for k := range rsAttr {
			if _, ok := ignoreFields[k]; ok {
				continue
			}
			if k == "%" {
				continue
			}
			if dsAttr[k] != rsAttr[k] {
				// ignore data Sources where an empty list is being compared against a null list.
				if k[len(k)-1:] == "#" && (dsAttr[k] == "" || dsAttr[k] == "0") && (rsAttr[k] == "" || rsAttr[k] == "0") {
					continue
				}
				errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr[k], rsAttr[k])
			}
		}

		if errMsg != "" {
			return errors.New(errMsg)
		}

		return nil
	}
}

func SkipIfEnvNotSet(t *testing.T, envs ...string) {
	if t == nil {
		log.Printf("[DEBUG] Not running inside of test - skip skipping")
		return
	}

	for _, k := range envs {
		if os.Getenv(k) == "" {
			log.Printf("[DEBUG] Warning - environment variable %s is not set - skipping test %s", k, t.Name())
			t.Skipf("Environment variable %s is not set", k)
		}
	}
}

func RandString(t *testing.T, length int) string {
	if !IsVcrEnabled() {
		return acctest.RandString(length)
	}
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	s, err := CreateVcrSource(t, envPath, vcrMode)
	if err != nil {
		// At this point we haven't created any resources, so fail fast
		t.Fatal(err)
	}

	r := rand.New(s.Source)
	result := make([]byte, length)
	set := "abcdefghijklmnopqrstuvwxyz012346789"
	for i := 0; i < length; i++ {
		result[i] = set[r.Intn(len(set))]
	}
	return string(result)
}

func IsVcrEnabled() bool {
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	return envPath != "" && vcrMode != ""
}

// A source for a given VCR test with the value that seeded it
type VcrSource struct {
	Seed   int64
	Source rand.Source
}

var SourcesLock = sync.RWMutex{}

func readSeedFromFile(fileName string) (int64, error) {
	// Max number of digits for int64 is 19
	data := make([]byte, 19)
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	_, err = f.Read(data)
	if err != nil {
		return 0, err
	}
	// Remove NULL characters from seed
	data = bytes.Trim(data, "\x00")
	seed := string(data)
	return StringToFixed64(seed)
}

// Produces a rand.Source for VCR testing based on the given mode.
// In RECORDING mode, generates a new seed and saves it to a file, using the seed for the source
// In REPLAYING mode, reads a seed from a file and creates a source from it
func CreateVcrSource(t *testing.T, path, mode string) (*VcrSource, error) {
	SourcesLock.RLock()
	s, ok := Sources[t.Name()]
	SourcesLock.RUnlock()
	if ok {
		return &s, nil
	}
	switch mode {
	case "RECORDING":
		seed := rand.Int63()
		s := rand.NewSource(seed)
		vcrSource := VcrSource{Seed: seed, Source: s}
		SourcesLock.Lock()
		Sources[t.Name()] = vcrSource
		SourcesLock.Unlock()
		return &vcrSource, nil
	case "REPLAYING":
		seed, err := readSeedFromFile(VcrSeedFile(path, t.Name()))
		if err != nil {
			return nil, fmt.Errorf("no cassette found on disk for %s, please replay this testcase in recording mode - %w", t.Name(), err)
		}
		s := rand.NewSource(seed)
		vcrSource := VcrSource{Seed: seed, Source: s}
		SourcesLock.Lock()
		Sources[t.Name()] = vcrSource
		SourcesLock.Unlock()
		return &vcrSource, nil
	default:
		log.Printf("[DEBUG] No valid environment var set for VCR_MODE, expected RECORDING or REPLAYING, skipping VCR. VCR_MODE: %s", mode)
		return nil, errors.New("No valid VCR_MODE set")
	}
}

// Retrieves a unique test name used for writing files
// replaces all `/` characters that would cause filepath issues
// This matters during tests that dispatch multiple tests, for example TestAccLoggingFolderExclusion
func VcrSeedFile(path, name string) string {
	return filepath.Join(path, fmt.Sprintf("%s.Seed", VcrFileName(name)))
}

func VcrFileName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}
