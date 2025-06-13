// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageBucketObjects() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleStorageBucketObjectsRead,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"match_glob": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket_objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"media_link": {
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

func datasourceGoogleStorageBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	bucketObjects := make([]map[string]interface{}, 0)
	bucket := d.Get("bucket").(string)
	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{StorageBasePath}}b/%s/o", bucket))
	if err != nil {
		return err
	}
	bucketObjects, err = transport_tpg.PluralDataSourceGetListMap(d, config, nil, userAgent, url, flattenDatasourceGoogleBucketObjectsList, params, "items")

	if err := d.Set("bucket_objects", bucketObjects); err != nil {
		return fmt.Errorf("Error retrieving bucket_objects: %s", err)
	}

	d.SetId(d.Get("bucket").(string))

	return nil
}

func flattenDatasourceGoogleBucketObjectsList(config *transport_tpg.Config, v interface{}) ([]map[string]interface{}, error) {
	if v == nil {
		return make([]map[string]interface{}, 0), nil
	}

	ls := v.([]interface{})
	buckets := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		o := raw.(map[string]interface{})

		var mLabels, mLocation, mName, mSelfLink, mStorageClass interface{}
		if oLabels, ok := o["labels"]; ok {
			mLabels = oLabels
		}
		if oLocation, ok := o["location"]; ok {
			mLocation = oLocation
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
		buckets = append(buckets, map[string]interface{}{
			"labels":        mLabels,
			"location":      mLocation,
			"name":          mName,
			"self_link":     mSelfLink,
			"storage_class": mStorageClass,
		})
	}

	return buckets, nil
}
