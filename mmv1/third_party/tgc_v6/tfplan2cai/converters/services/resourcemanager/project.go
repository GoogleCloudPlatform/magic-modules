package resourcemanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters/cai"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// ProjectAssetType is the CAI asset type name for project.
const ProjectAssetType string = "cloudresourcemanager.googleapis.com/Project"
const ProjectBillingInfoAssetType string = "cloudbilling.googleapis.com/ProjectBillingInfo"

func ResourceConverterProject() cai.ResourceConverter {
	return cai.ResourceConverter{
		Convert: GetProjectAndBillingInfoCaiObjects,
	}
}

func GetProjectAndBillingInfoCaiObjects(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]caiasset.Asset, error) {
	assets := []caiasset.Asset{}
	if projectAsset, err := GetProjectCaiObject(d, config); err == nil {
		if billingAsset, err := GetProjectBillingInfoCaiObject(d, config); err == nil {
			assets = append(assets, projectAsset)
			assets = append(assets, billingAsset)
		} else {
			return assets, err
		}
	} else {
		return assets, err
	}
	return assets, nil
}

func GetProjectCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (caiasset.Asset, error) {
	// use project number if it's available; otherwise, fill in project id so that we
	// keep the CAI assets apart for different uncreated projects.
	var linkTmpl string
	if _, ok := d.GetOk("number"); ok {
		linkTmpl = "//cloudresourcemanager.googleapis.com/projects/{{number}}"
	} else {
		linkTmpl = "//cloudresourcemanager.googleapis.com/projects/{{project_id_or_project}}"
	}

	name, err := cai.AssetName(d, config, linkTmpl)
	if err != nil {
		return caiasset.Asset{}, err
	}
	if obj, err := GetProjectApiObject(d, config); err == nil {
		return caiasset.Asset{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Data:                 obj,
			},
		}, nil
	} else {
		return caiasset.Asset{}, err
	}
}

func GetProjectApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	pid := d.Get("project_id").(string)

	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
	}

	if res, ok := d.GetOk("number"); ok {
		num, err := strconv.ParseInt(res.(string), 10, 64)
		if err != nil {
			return nil, err
		}

		project.ProjectNumber = num
	}

	if err := getParentResourceId(d, project); err != nil {
		return nil, err
	}

	if _, ok := d.GetOk("effective_labels"); ok {
		project.Labels = tpgresource.ExpandEffectiveLabels(d)
	}

	return cai.JsonMap(project)
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

func GetProjectBillingInfoCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (caiasset.Asset, error) {
	// use project number if it's available; otherwise, fill in project id so that we
	// keep the CAI assets apart for different uncreated projects.
	var linkTmpl string
	if _, ok := d.GetOk("number"); ok {
		linkTmpl = "//cloudbilling.googleapis.com/projects/{{number}}/billingInfo"
	} else {
		linkTmpl = "//cloudbilling.googleapis.com/projects/{{project_id_or_project}}/billingInfo"
	}

	name, err := cai.AssetName(d, config, linkTmpl)
	if err != nil {
		return caiasset.Asset{}, err
	}
	project := strings.Split(name, "/")[4]
	if obj, err := GetProjectBillingInfoApiObject(d, project); err == nil {
		return caiasset.Asset{
			Name: name,
			Type: ProjectBillingInfoAssetType,
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/cloudbilling/v1/rest",
				DiscoveryName:        "ProjectBillingInfo",
				Data:                 obj,
			},
		}, nil
	} else {
		return caiasset.Asset{}, err
	}
}

func GetProjectBillingInfoApiObject(d tpgresource.TerraformResourceData, project string) (map[string]interface{}, error) {
	if _, ok := d.GetOk("billing_account"); !ok {
		// TODO: If the project already exists, we could ask the API about it's
		// billing info here.
		return nil, cai.ErrNoConversion
	}

	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", d.Get("billing_account")),
		Name:               fmt.Sprintf("projects/%s/billingInfo", project),
		ProjectId:          d.Get("project_id").(string),
	}

	return cai.JsonMap(ba)
}
