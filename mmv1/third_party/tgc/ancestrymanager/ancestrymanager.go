// Package ancestrymanager provides an interface to query the ancestry information for a resource.
package ancestrymanager

import (
	"fmt"
	"strconv"
	"strings"

	crmv1 "google.golang.org/api/cloudresourcemanager/v1"
	crmv3 "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/storage/v1"

	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"

	"go.uber.org/zap"
)

// AncestryManager is the interface that fetch ancestors for a resource.
type AncestryManager interface {
	// Ancestors returns a list of ancestors.
	Ancestors(config *transport_tpg.Config, tfData tpgresource.TerraformResourceData, cai *resources.Asset) ([]string, string, error)
}

type manager struct {
	// The logger.
	errorLogger *zap.Logger
	// cloud resource manager V3 client. If this field is nil, online lookups will
	// be disabled.
	// cloud resource manager V1 client. If this field is nil, online lookups will
	// be disabled.
	resourceManagerV3 *crmv3.Service
	resourceManagerV1 *crmv1.Service
	storageClient     *storage.Service
	// Cache to prevent multiple network calls for looking up the same
	// resource's ancestry. The map key is the resource itself, in the format of
	// "<type>/<id>", ancestors are sorted from closest to furthest.
	ancestorCache map[string][]string
}

// New returns AncestryManager that can be used to fetch ancestry information.
// Entries takes `projects/<number>` or `folders/<id>` as key and ancestry path
// as value to the offline cache. If the key is not prefix with `projects/` or
// `folders/`, it will be considered as a project. If offline is true, resource
// manager API requests for ancestry will be disabled.
func New(cfg *transport_tpg.Config, offline bool, entries map[string]string, errorLogger *zap.Logger) (AncestryManager, error) {
	am := &manager{
		ancestorCache: map[string][]string{},
		errorLogger:   errorLogger,
	}
	if !offline {
		am.resourceManagerV1 = cfg.NewResourceManagerClient(cfg.UserAgent)
		am.resourceManagerV3 = cfg.NewResourceManagerV3Client(cfg.UserAgent)
		am.storageClient = cfg.NewStorageClient(cfg.UserAgent)
	}
	err := am.initAncestryCache(entries)
	if err != nil {
		return nil, err
	}
	return am, nil
}

func (m *manager) initAncestryCache(entries map[string]string) error {
	for item, ancestry := range entries {
		if item != "" && ancestry != "" {
			ancestors, err := parseAncestryPath(ancestry)
			if err != nil {
				continue
			}
			key, err := parseAncestryKey(item)
			if err != nil {
				return err
			}
			// ancestry path should include the item itself
			if ancestors[0] != key {
				ancestors = append([]string{key}, ancestors...)
			}
			m.store(key, ancestors)
		}
	}
	return nil
}

func parseAncestryKey(val string) (string, error) {
	key := normalizeAncestry(val)
	ix := strings.LastIndex(key, "/")
	if ix == -1 {
		// If not containing /, then treat it as a project.
		return fmt.Sprintf("projects/%s", key), nil
	} else {
		k := key[:ix]
		if k == "projects" || k == "folders" || k == "organizations" {
			return key, nil
		}
		return "", fmt.Errorf("key with can only start with projects/, folders/, or organizations/")
	}
}

// Ancestors uses the resource manager API to get ancestors for resource.
// It implements a cache because many resources share the same ancestors.
func (m *manager) Ancestors(config *transport_tpg.Config, tfData tpgresource.TerraformResourceData, cai *resources.Asset) ([]string, string, error) {
	results, err := m.fetchAncestors(config, tfData, cai)
	if err != nil {
		return nil, "", err
	}

	parent, err := assetParent(cai, results)
	if err != nil {
		return nil, "", err
	}
	return results, parent, nil
}

// fetchAncestors uses the resource manager API to get ancestors for resource.
// It implements a cache because many resources share the same ancestors.
func (m *manager) fetchAncestors(config *transport_tpg.Config, tfData tpgresource.TerraformResourceData, cai *resources.Asset) ([]string, error) {
	if cai == nil {
		return nil, fmt.Errorf("CAI asset is nil")
	}
	m.errorLogger.Info(fmt.Sprintf("Retrieving ancestry from resource (type=%s)", cai.Type))
	key := ""
	orgKey := ""
	folderKey := ""
	projectKey := ""

	orgID, orgOK := getOrganizationFromResource(tfData)
	if orgOK {
		orgKey = orgID
		if !strings.HasPrefix(orgKey, "organizations/") {
			orgKey = fmt.Sprintf("organizations/%s", orgKey)
		}
	}
	folderID, folderOK := getFolderFromResource(tfData)
	if folderOK {
		folderKey = folderID
		if !strings.HasPrefix(folderKey, "folders/") {
			folderKey = fmt.Sprintf("folders/%s", folderKey)
		}
	}
	project, _ := m.getProjectFromResource(tfData, config, cai)
	if project != "" {
		projectKey = project
		if !strings.HasPrefix(projectKey, "projects/") {
			projectKey = fmt.Sprintf("projects/%s", project)
		}
	}

	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Folder":
		if folderOK {
			key = folderKey
		} else if orgOK {
			key = orgKey
		} else {
			return []string{"organizations/unknown"}, nil
		}
	case "cloudresourcemanager.googleapis.com/Organization":
		if !orgOK {
			return nil, fmt.Errorf("organization id not found in terraform data")
		}
		key = orgKey
	case "iam.googleapis.com/Role":
		// google_organization_iam_custom_role or google_project_iam_custom_role
		if orgOK {
			key = orgKey
		} else if projectKey != "" {
			key = projectKey
		} else {
			return []string{"organizations/unknown"}, nil
		}
	case "cloudresourcemanager.googleapis.com/Project", "cloudbilling.googleapis.com/ProjectBillingInfo":
		// for google_project and google_project_iam resources
		var ancestors []string
		if projectKey != "" {
			ancestors = []string{projectKey}
			// Get ancestry from project level with v1 API first.
			// This is to avoid requiring folder level permission if
			// there is no folder change.
			m.getAncestorsWithCache(projectKey)
		}
		// only folder_id or org_id is allowed for google_project
		if orgOK {
			// no need to use API to fetch ancestors
			ancestors = append(ancestors, fmt.Sprintf("organizations/%s", orgID))
			return ancestors, nil
		}
		if folderOK {
			// If folder is changed, then it goes with v3 API, else it will use cache.
			key = folderKey
			ret, err := m.getAncestorsWithCache(key)
			if err != nil {
				return nil, err
			}
			ancestors = append(ancestors, ret...)
			return ancestors, nil
		}

		// neither folder_id nor org_id is specified
		if projectKey == "" {
			return []string{"organizations/unknown"}, nil
		}
		key = projectKey

	default:
		if projectKey == "" {
			return []string{"organizations/unknown"}, nil
		}
		key = projectKey
	}
	return m.getAncestorsWithCache(key)
}

func (m *manager) getAncestorsWithCache(key string) ([]string, error) {
	var ancestors []string
	cur := key
	for cur != "" {
		if cachedAncestors, ok := m.ancestorCache[cur]; ok {
			ancestors = append(ancestors, cachedAncestors...)
			break
		}
		if strings.HasPrefix(cur, "organizations/") {
			ancestors = append(ancestors, cur)
			break
		}
		if m.resourceManagerV3 == nil || m.resourceManagerV1 == nil {
			return nil, fmt.Errorf("resourceManager required to fetch ancestry for %s from the API", cur)
		}
		if strings.HasPrefix(cur, "projects") {
			// fall back to use v1 API GetAncestry to avoid requiring extra folder permission
			projectID := strings.TrimPrefix(cur, "projects/")
			resp, err := m.resourceManagerV1.Projects.GetAncestry(projectID, &crmv1.GetAncestryRequest{}).Do()
			if err != nil {
				return nil, handleCRMError(cur, err)
			}
			for _, anc := range resp.Ancestor {
				ancestor := normalizeAncestry(fmt.Sprintf("%s/%s", anc.ResourceId.Type, anc.ResourceId.Id))
				ancestors = append(ancestors, ancestor)
			}
			// break out of the loop
			cur = ""
		} else {
			project, err := m.resourceManagerV3.Projects.Get(cur).Do()
			if err != nil {
				return nil, handleCRMError(cur, err)
			}
			ancestors = append(ancestors, project.Name)
			cur = project.Parent
		}
	}
	m.store(key, ancestors)
	return ancestors, nil
}

func handleCRMError(resource string, err error) error {
	if isGoogleApiErrorWithCode(err, 403) {
		helperURL := "https://cloud.google.com/docs/terraform/policy-validation/troubleshooting#ProjectCallerForbidden"
		return fmt.Errorf("user does not have the correct permissions for %s. For more info: %s", resource, helperURL)
	}
	return err
}

func (m *manager) store(key string, ancestors []string) {
	if key == "" || len(ancestors) == 0 {
		return
	}
	if _, ok := m.ancestorCache[key]; !ok {
		m.ancestorCache[key] = ancestors
	}
	// cache ancestors along the ancestry path
	for i, ancestor := range ancestors {
		if _, ok := m.ancestorCache[ancestor]; !ok {
			m.ancestorCache[ancestor] = ancestors[i:]
		}
	}
}

func parseAncestryPath(path string) ([]string, error) {
	normStr := normalizeAncestry(path)
	splits := strings.Split(normStr, "/")
	if len(splits)%2 != 0 {
		return nil, fmt.Errorf("unexpected format of ancestry path %s", path)
	}
	var ancestors []string
	allowedPrefixes := map[string]bool{
		"projects":      true,
		"folders":       true,
		"organizations": true,
	}
	for i := 0; i < len(splits); i = i + 2 {
		if _, ok := allowedPrefixes[splits[i]]; !ok {
			return nil, fmt.Errorf("invalid ancestry path %s with %s", path, splits[i])
		}
		ancestors = append(ancestors, fmt.Sprintf("%s/%s", splits[i], splits[i+1]))
	}
	// reverse the sequence
	i, j := 0, len(ancestors)-1
	for i < j {
		ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
		i++
		j--
	}
	return ancestors, nil
}

func normalizeAncestry(val string) string {
	for _, r := range []struct {
		old string
		new string
	}{
		{"organization/", "organizations/"},
		{"folder/", "folders/"},
		{"project/", "projects/"},
	} {
		val = strings.ReplaceAll(val, r.old, r.new)
	}
	return val
}

// getProjectFromResource reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func (m *manager) getProjectFromResource(d tpgresource.TerraformResourceData, config *transport_tpg.Config, cai *resources.Asset) (string, error) {

	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Project",
		"cloudbilling.googleapis.com/ProjectBillingInfo":
		res, ok := d.GetOk("number")
		if ok {
			return res.(string), nil
		}
		// Fall back to project_id if number is not available.
		res, ok = d.GetOk("project_id")
		if ok {
			return res.(string), nil
		} else {
			m.errorLogger.Warn(fmt.Sprintf("Failed to retrieve project_id for %s from resource", cai.Name))
		}
	case "storage.googleapis.com/Bucket":
		if cai.Resource != nil {
			res, ok := cai.Resource.Data["project"]
			if ok {
				return res.(string), nil
			}
		}
		m.errorLogger.Warn(fmt.Sprintf("Failed to retrieve project_id for %s from cai resource", cai.Name))

		bucketField, ok := d.GetOk("bucket")
		if ok {
			bucket := bucketField.(string)
			resp, err := m.storageClient.Buckets.Get(bucket).Do()
			if err == nil {
				projectNum := resp.ProjectNumber
				return strconv.Itoa(int(projectNum)), nil
			}
			m.errorLogger.Warn(fmt.Sprintf("Failed to get bucket %s", bucket))
		}
		m.errorLogger.Warn("Failed to retrieve bucket field from tf data")
	}

	return getProjectFromSchema("project", d, config)
}

type NoOpAncestryManager struct{}

func (*NoOpAncestryManager) Ancestors(config *transport_tpg.Config, tfData tpgresource.TerraformResourceData, cai *resources.Asset) ([]string, string, error) {
	return nil, "", nil
}
