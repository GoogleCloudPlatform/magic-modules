package storage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageBuckets() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleStorageBucketsRead,
		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleStorageBucketsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)

	url := "https://storage.googleapis.com/storage/v1/b"

	params["project"], err = tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for bucket: %s", err)
	}

	if v, ok := d.GetOk("prefix"); ok {
		params["prefix"] = v.(string)
	}

	url, err = transport_tpg.AddQueryParams(url, params)
	if err != nil {
		return err
	}

	opts := transport_tpg.GetPaginatedItemsOptions{
		ResourceData:   d,
		Config:         config,
		BillingProject: nil,
		UserAgent:      userAgent,
		URL:            url,
		ResourceToList: "items",
	}
	buckets, err := transport_tpg.GetPaginatedItems(opts)
	if err != nil {
		return fmt.Errorf("Error retrieving buckets: %s", err)
	}

	if err := d.Set("buckets", buckets); err != nil {
		return fmt.Errorf("Error retrieving buckets: %s", err)
	}

	d.SetId(params["project"])

	return nil
}

func flattenDatasourceGoogleBucketsList(config *transport_tpg.Config, v []map[string]interface{}) ([]map[string]interface{}, error) {
	if v == nil {
		return make([]map[string]interface{}, 0), nil
	}

	bucketObjects := make([]map[string]interface{}, 0, len(v))
	for _, raw := range v {
		o := raw

		var mContentType, mMediaLink, mName, mSelfLink, mStorageClass interface{}
		if oContentType, ok := o["contentType"]; ok {
			mContentType = oContentType
		}
		if oMediaLink, ok := o["mediaLink"]; ok {
			mMediaLink = oMediaLink
		}
		if oName, ok := o["name"]; ok {
			mName = oName
		}
		if oSelfLink, ok := o["selfLink"]; ok {
			mSelfLink = oSelfLink
		}
		if oStorageClass, ok := o["storageClass"]; ok {
			mStorageClass = oStorageClass
		}
		bucketObjects = append(bucketObjects, map[string]interface{}{
			"content_type":  mContentType,
			"media_link":    mMediaLink,
			"name":          mName,
			"self_link":     mSelfLink,
			"storage_class": mStorageClass,
		})
	}

	return bucketObjects, nil
}
