package google

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func resourceConverterProject() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Project",
		Convert:           GetProjectCaiObject,
		MergeCreateUpdate: MergeProject,
	}
}

func resourceConverterProjectBillingInfo() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType: "cloudbilling.googleapis.com/ProjectBillingInfo",
		Convert:   GetProjectBillingInfoCaiObject,
	}
}

func GetProjectCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	// use project number if it's available; otherwise, fill in project id so that we
	// keep the CAI assets apart for different uncreated projects.
	var linkTmpl string
	if _, ok := d.GetOk("number"); ok {
		linkTmpl = "//cloudresourcemanager.googleapis.com/projects/{{number}}"
	} else {
		linkTmpl = "//cloudresourcemanager.googleapis.com/projects/{{project_id_or_project}}"
	}
	name, err := tpgresource.AssetName(d, config, linkTmpl)
	if err != nil {
		return []tpgresource.Asset{}, err
	}
	if obj, err := GetProjectApiObject(d, config); err == nil {
		return []tpgresource.Asset{{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &tpgresource.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []tpgresource.Asset{}, err
	}
}

func GetProjectApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	pid := d.Get("project_id").(string)

	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
	}

	if err := getParentResourceId(d, project); err != nil {
		return nil, err
	}

	if _, ok := d.GetOk("labels"); ok {
		project.Labels = tpgresource.ExpandLabels(d)
	}

	return tpgresource.JsonMap(project)
}

func getParentResourceId(d tpgresource.TerraformResourceData, p *cloudresourcemanager.Project) error {
	orgId := d.Get("org_id").(string)
	folderId := d.Get("folder_id").(string)

	if orgId != "" && folderId != "" {
		return fmt.Errorf("'org_id' and 'folder_id' cannot be both set.")
	}

	if orgId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   orgId,
			Type: "organization",
		}
	}

	if folderId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   strings.TrimPrefix(folderId, "folders/"),
			Type: "folder",
		}
	}

	return nil
}

func GetProjectBillingInfoCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	// use project number if it's available; otherwise, fill in project id so that we
	// keep the CAI assets apart for different uncreated projects.
	var linkTmpl string
	if _, ok := d.GetOk("number"); ok {
		linkTmpl = "//cloudbilling.googleapis.com/projects/{{number}}/billingInfo"
	} else {
		linkTmpl = "//cloudbilling.googleapis.com/projects/{{project_id_or_project}}/billingInfo"
	}
	name, err := tpgresource.AssetName(d, config, linkTmpl)
	if err != nil {
		return []tpgresource.Asset{}, err
	}
	project := strings.Split(name, "/")[4]
	if obj, err := GetProjectBillingInfoApiObject(d, project); err == nil {
		return []tpgresource.Asset{{
			Name: name,
			Type: "cloudbilling.googleapis.com/ProjectBillingInfo",
			Resource: &tpgresource.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/cloudbilling/v1/rest",
				DiscoveryName:        "ProjectBillingInfo",
				Data:                 obj,
			}},
		}, nil
	} else {
		return []tpgresource.Asset{}, err
	}
}

func GetProjectBillingInfoApiObject(d tpgresource.TerraformResourceData, project string) (map[string]interface{}, error) {
	if _, ok := d.GetOk("billing_account"); !ok {
		// TODO: If the project already exists, we could ask the API about it's
		// billing info here.
		return nil, tpgresource.ErrNoConversion
	}

	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", d.Get("billing_account")),
		Name:               fmt.Sprintf("projects/%s/billingInfo", project),
		ProjectId:          d.Get("project_id").(string),
	}

	return tpgresource.JsonMap(ba)
}

func MergeProject(existing, incoming tpgresource.Asset) tpgresource.Asset {
	existing.Resource = incoming.Resource
	return existing
}
