package bigquery

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBigQueryDatasets() *schema.Resource {

	dsSchema := map[string]*schema.Schema{
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the project in which the dataset is located. If it is not provided, the provider project is used.",
		},
		"datasets": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"labels": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"dataset_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Read:   DataSourceGoogleBigQueryDatasetsRead,
		Schema: dsSchema,
	}
}

func DataSourceGoogleBigQueryDatasetsRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)

	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}

	params := make(map[string]string)
	datasets := make([]map[string]interface{}, 0)

	for {
		url := fmt.Sprintf("https://bigquery.googleapis.com/bigquery/v2/projects/%s/datasets", project)

		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error retrieving datasets: %s", err)
		}

		pageDatasets := flattenDatasourceGoogleBigQueryDatasetsList(res["datasets"])
		datasets = append(datasets, pageDatasets...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("datasets", datasets); err != nil {
		return fmt.Errorf("Error retrieving datasets: %s", err)
	}

	id := fmt.Sprintf("projects/%s/datasets", project)

	d.SetId(id)

	return nil
}

func flattenDatasourceGoogleBigQueryDatasetsList(v interface{}) []map[string]interface{} {

	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})

	datasets := make([]map[string]interface{}, 0, len(ls))

	for _, raw := range ls {
		o := raw.(map[string]interface{})

		var mLabels map[string]interface{}
		var mDatasetName string

		if oLabels, ok := o["labels"].(map[string]interface{}); ok {
			mLabels = oLabels
		} else {
			mLabels = make(map[string]interface{}) // Initialize as an empty map if labels are missing
		}

		if oDatasetReference, ok := o["datasetReference"].(map[string]interface{}); ok {
			if datasetID, ok := oDatasetReference["datasetId"].(string); ok {
				mDatasetName = datasetID
			}
		}
		datasets = append(datasets, map[string]interface{}{
			"labels":     mLabels,
			"dataset_id": mDatasetName,
		})
	}

	return datasets
}
