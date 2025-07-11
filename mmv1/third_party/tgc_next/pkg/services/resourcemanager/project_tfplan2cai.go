package resourcemanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters/cai"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/transport"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func ProjectTfplan2caiConverter() cai.Tfplan2caiConverter {
	return cai.Tfplan2caiConverter{
		Convert: GetProjectAndBillingInfoCaiObjects,
	}
}

func GetProjectAndBillingInfoCaiObjects(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]caiasset.Asset, error) {
	if projectAsset, err := GetProjectCaiObject(d, config); err == nil {
		assets := []caiasset.Asset{projectAsset}
		if _, ok := d.GetOk("billing_account"); !ok {
			return assets, nil
		} else {
			if billingAsset, err := GetProjectBillingInfoCaiObject(d, config); err == nil {
				assets = append(assets, billingAsset)
				return assets, nil
			} else {
				return []caiasset.Asset{}, err
			}
		}
	} else {
		return []caiasset.Asset{}, err
	}
}

func GetProjectCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (caiasset.Asset, error) {
	linkTmpl := "//cloudresourcemanager.googleapis.com/projects/{{number}}"
	name, err := cai.AssetName(d, config, linkTmpl)
	if err != nil {
		return caiasset.Asset{}, err
	}
	if data, err := GetProjectData(d, config); err == nil {
		return caiasset.Asset{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://cloudresourcemanager.googleapis.com/$discovery/rest?version=v1",
				DiscoveryName:        "Project",
				Data:                 data,
			},
		}, nil
	} else {
		return caiasset.Asset{}, err
	}
}

func GetProjectData(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
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
	linkTmpl := "//cloudbilling.googleapis.com/projects/{{project_id_or_project}}/billingInfo"
	name, err := cai.AssetName(d, config, linkTmpl)
	if err != nil {
		return caiasset.Asset{}, err
	}
	project := strings.Split(name, "/")[4]
	if data, err := GetProjectBillingInfoData(d, project); err == nil {
		return caiasset.Asset{
			Name: name,
			Type: "cloudbilling.googleapis.com/ProjectBillingInfo",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://cloudbilling.googleapis.com/$discovery/rest",
				DiscoveryName:        "ProjectBillingInfo",
				Data:                 data,
				Location:             "global",
			},
		}, nil
	} else {
		return caiasset.Asset{}, err
	}
}

func GetProjectBillingInfoData(d tpgresource.TerraformResourceData, project string) (map[string]interface{}, error) {
	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", d.Get("billing_account")),
		Name:               fmt.Sprintf("projects/%s/billingInfo", project),
		ProjectId:          d.Get("project_id").(string),
	}

	return cai.JsonMap(ba)
}
