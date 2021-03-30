package google

import (
	"fmt"
	"strings"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func GetProjectCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	// use project number if it's available; otherwise, fill in project id so that we
	// keep the CAI assets apart for different uncreated projects.
	var linkTmpl string
	if _, ok := d.GetOk("number"); ok {
		linkTmpl = "//cloudresourcemanager.googleapis.com/projects/{{number}}"
	} else {
		linkTmpl = "//cloudresourcemanager.googleapis.com/projects/{{project_id_or_project}}"
	}
	name, err := assetName(d, config, linkTmpl)
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetProjectApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetProjectApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	pid := d.Get("project_id").(string)

	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
	}

	if err := getParentResourceId(d, project); err != nil {
		return nil, err
	}

	if _, ok := d.GetOk("labels"); ok {
		project.Labels = expandLabels(d)
	}

	return jsonMap(project)
}

func getParentResourceId(d TerraformResourceData, p *cloudresourcemanager.Project) error {
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

func GetProjectBillingInfoCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	// use project number if it's available; otherwise, fill in project id so that we
	// keep the CAI assets apart for different uncreated projects.
	var linkTmpl string
	if _, ok := d.GetOk("number"); ok {
		linkTmpl = "//cloudbilling.googleapis.com/projects/{{number}}/billingInfo"
	} else {
		linkTmpl = "//cloudbilling.googleapis.com/projects/{{project_id_or_project}}/billingInfo"
	}
	name, err := assetName(d, config, linkTmpl)
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetProjectBillingInfoApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: "cloudbilling.googleapis.com/ProjectBillingInfo",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/cloudbilling/v1/rest",
				DiscoveryName:        "ProjectBillingInfo",
				Data:                 obj,
			}},
		}, nil
	} else {
		return []Asset{}, err
	}
}

func GetProjectBillingInfoApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	if _, ok := d.GetOk("billing_account"); !ok {
		// TODO: If the project already exists, we could ask the API about it's
		// billing info here.
		return nil, ErrNoConversion
	}

	var name string
	if number, ok := d.GetOk("number"); ok {
		name = fmt.Sprintf("projects/%s/billingInfo", number)
	} else {
		name = fmt.Sprintf("projects/%s/billingInfo", d.Get("project_id"))
	}

	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", d.Get("billing_account")),
		Name:               name,
		ProjectId:          d.Get("project_id").(string),
	}

	return jsonMap(ba)
}
