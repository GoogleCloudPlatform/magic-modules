package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var IamBigqueryDatasetSchema = map[string]*schema.Schema{
	"dataset_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"project": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

type BigqueryDatasetIamUpdater struct {
	project   string
	datasetId string
	Config    *Config
}

func NewBigqueryDatasetIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	d.Set("project", project)

	return &BigqueryDatasetIamUpdater{
		project:  project,
		datasetId: d.Get("dataset_id").(string),
		Config:   config,
	}, nil
}

func BigqueryDatasetIdParseFunc(d *schema.ResourceData, config *Config) error {
	fv, err := parseProjectFieldValue("datasets", d.Id(), "project", d, config, false)
	if err != nil {
		return err
	}

	d.Set("project", fv.Project)
	d.Set("dataset_id", fv.Name)

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *BigqueryDatasetIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	url := fmt.Sprintf("%s%s", u.Config.BigQueryBasePath, u.GetResourceId)

	res, err := sendRequest(u.Config, "GET", u.Project, url, nil)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	v, ok = res["access"]
	policy := &cloudresourcemanager.Policy{}
	if !ok {
		return policy
	}

	return policy, nil
}

func accessToPolicy(access map[string]interface{}) (*cloudresourcemanager.Policy) {
	if access == nil {
		return nil
	}

	return nil
}

func policyToAccess(policy *cloudresourcemanager.Policy) map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, binding := policy.Bindings {
		for _, member := binding.Members {
			access := map[string]int{
    		"role": binding.Role,
			}
			memberType, member := accessMemberToIam(member)
			access[memberType] = member
			res = append(res, access)
		}
	}

	return nil
}

func accessMemberToIam(member string) (string, string) {
	pieces = strings.SplitN(member, ":", 2)
	if len(pieces) > 1 {
		pieces[1] = strings.ToLower(pieces[1])
	}
	if 
}

func (u *BigqueryDatasetIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	bigtablePolicy, err := resourceManagerToBigtablePolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	req := &bigtableadmin.SetIamPolicyRequest{Policy: bigtablePolicy}
	_, err = u.Config.clientBigtableProjectsInstances.SetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BigqueryDatasetIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/datasets/%s", u.project, u.datasetId)
}

// Matches the mutex of google_big_query_dataset_access
func (u *BigqueryDatasetIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("%s", u.datasetId)
}

func (u *BigqueryDatasetIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Bigquery Dataset %s/%s", u.project, u.datasetId)
}
