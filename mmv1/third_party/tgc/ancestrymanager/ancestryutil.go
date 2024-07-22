package ancestrymanager

import (
	"fmt"
	"strings"

	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
	"google.golang.org/api/googleapi"
)

// assetParent derives a resource's parent from its ancestors.
func assetParent(cai *resources.Asset, ancestors []string) (string, error) {
	if cai == nil {
		return "", fmt.Errorf("asset not provided")
	}
	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Folder":
		if len(ancestors) > 1 {
			parent := ancestors[1]
			if strings.HasPrefix(parent, "folders/") || strings.HasPrefix(parent, "organizations/") {
				return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[1]), nil
			}
		}
		if len(ancestors) == 1 && strings.HasPrefix(ancestors[0], "organizations/") {
			// organizations/unknown
			return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[0]), nil
		}
	case "cloudresourcemanager.googleapis.com/Organization":
		return "", nil
	case "cloudresourcemanager.googleapis.com/Project":
		if len(ancestors) < 1 {
			return "", fmt.Errorf("unexpected value for ancestors: %s", ancestors)
		}
		if strings.HasPrefix(ancestors[0], "projects/") {
			if len(ancestors) > 1 {
				return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[1]), nil
			}
		}
		return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[0]), nil
	default:
		if len(ancestors) < 1 {
			return "", fmt.Errorf("unexpected value for ancestors: %s", ancestors)
		}
		return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[0]), nil
	}
	return "", fmt.Errorf("unexpected value for ancestors: %v", ancestors)
}

// ConvertToAncestryPath composes a path containing organization/folder/project
// (i.e. "organization/my-org/folder/my-folder/project/my-prj").
func ConvertToAncestryPath(as []string) string {
	var path []string
	for i := len(as) - 1; i >= 0; i-- {
		path = append(path, as[i])
	}
	str := strings.Join(path, "/")
	return sanitizeAncestryPath(str)
}

func sanitizeAncestryPath(s string) string {
	ret := s
	// convert back to match existing ancestry path style.
	for _, r := range []struct {
		old string
		new string
	}{
		{"organizations/", "organization/"},
		{"folders/", "folder/"},
		{"projects/", "project/"},
	} {
		ret = strings.ReplaceAll(ret, r.old, r.new)
	}
	return ret
}

func getProjectFromSchema(projectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	res, ok := d.GetOk(projectSchemaField)
	if ok && projectSchemaField != "" {
		return res.(string), nil
	}
	res, ok = d.GetOk("parent")
	if ok && strings.HasPrefix(res.(string), "projects/") {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("required field '%s' is not set, you may use --project=my-project to provide a default project to resolve the issue", projectSchemaField)
}

// getOrganizationFromResource reads org_id field from terraform data.
func getOrganizationFromResource(tfData tpgresource.TerraformResourceData) (string, bool) {
	orgID, ok := tfData.GetOk("org_id")
	if ok {
		return orgID.(string), ok
	}
	orgID, ok = tfData.GetOk("parent")
	if ok && strings.HasPrefix(orgID.(string), "organizations/") {
		return orgID.(string), ok
	}
	return "", false
}

// getFolderFromResource reads folder_id, folder, parent field from terraform data.
func getFolderFromResource(tfData tpgresource.TerraformResourceData) (string, bool) {
	folderID, ok := tfData.GetOk("folder_id")
	if ok {
		return folderID.(string), ok
	}
	folderID, ok = tfData.GetOk("folder")
	if ok {
		return folderID.(string), ok
	}

	folderID, ok = tfData.GetOk("parent")
	if ok && strings.HasPrefix(folderID.(string), "folders/") {
		return folderID.(string), ok
	}
	return "", false
}

// isGoogleApiErrorWithCode cheks if the error code is of given type or not.
func isGoogleApiErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == errCode
}
