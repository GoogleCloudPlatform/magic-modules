package storage

import (
	"fmt"
	"sort"

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

	opt := transport_tpg.GetPaginatedItemsOptions{
		ResourceData:   d,
		Config:         config,
		UserAgent:      userAgent,
		URL:            url,
		Params:         params,
		ResourceToList: "items",
		ListFlattener:  flattenDatasourceGoogleBucketObjectsList,
	}
	bucketObjects, err = transport_tpg.GetPaginatedItems(opt)
	if err != nil {
		return err
	}
	if err := d.Set("bucket_objects", bucketObjects); err != nil {
		return fmt.Errorf("Error retrieving bucket_objects: %s", err)
	}

	d.SetId(d.Get("bucket").(string))

	return nil
}

func flattenDatasourceGoogleBucketObjectsList(config *transport_tpg.Config, v []map[string]interface{}) ([]map[string]interface{}, error) {
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

	sort.Slice(bucketObjects, func(i, j int) bool {
		return bucketObjects[i]["name"].(string) < bucketObjects[j]["name"].(string)
	})

	return bucketObjects, nil
}
