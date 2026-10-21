package container

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleContainerEngineVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerEngineVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_master_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_node_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_master_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"valid_node_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"release_channel_default_version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"release_channel_latest_version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"release_channel_upgrade_target_version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleContainerEngineVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if len(location) == 0 {
		return fmt.Errorf("Cannot determine location: set location in this data source or at provider-level")
	}

	location = fmt.Sprintf("projects/%s/locations/%s", project, location)
	url := fmt.Sprintf("%s%s/serverConfig", transport_tpg.BaseUrl(Product, config), location)
	resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error retrieving available container cluster versions: %s", err.Error())
	}

	validMasterVersions := make([]string, 0)
	if rawVersions, ok := resp["validMasterVersions"].([]interface{}); ok {
		for _, rawV := range rawVersions {
			v, _ := rawV.(string)
			if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
				validMasterVersions = append(validMasterVersions, v)
			}
		}
	}

	validNodeVersions := make([]string, 0)
	if rawVersions, ok := resp["validNodeVersions"].([]interface{}); ok {
		for _, rawV := range rawVersions {
			v, _ := rawV.(string)
			if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
				validNodeVersions = append(validNodeVersions, v)
			}
		}
	}

	if err := d.Set("valid_master_versions", validMasterVersions); err != nil {
		return fmt.Errorf("Error setting valid_master_versions: %s", err)
	}
	if len(validMasterVersions) > 0 {
		if err := d.Set("latest_master_version", validMasterVersions[0]); err != nil {
			return fmt.Errorf("Error setting latest_master_version: %s", err)
		}
	}

	if err := d.Set("valid_node_versions", validNodeVersions); err != nil {
		return fmt.Errorf("Error setting valid_node_versions: %s", err)
	}
	if len(validNodeVersions) > 0 {
		if err := d.Set("latest_node_version", validNodeVersions[0]); err != nil {
			return fmt.Errorf("Error setting latest_node_version: %s", err)
		}
	}

	if err := d.Set("default_cluster_version", resp["defaultClusterVersion"]); err != nil {
		return fmt.Errorf("Error setting default_cluster_version: %s", err)
	}

	releaseChannelDefaultVersion := map[string]string{}
	releaseChannelLatestVersion := map[string]string{}
	releaseChannelUpgradeTargetVersion := map[string]string{}
	if channels, ok := resp["channels"].([]interface{}); ok {
		for _, rawChannel := range channels {
			channelResp, ok := rawChannel.(map[string]interface{})
			if !ok {
				continue
			}
			channelName, _ := channelResp["channel"].(string)
			defaultVersion, _ := channelResp["defaultVersion"].(string)
			upgradeTargetVersion, _ := channelResp["upgradeTargetVersion"].(string)
			releaseChannelDefaultVersion[channelName] = defaultVersion
			releaseChannelUpgradeTargetVersion[channelName] = upgradeTargetVersion
			if validVersions, ok := channelResp["validVersions"].([]interface{}); ok {
				for _, rawV := range validVersions {
					v, _ := rawV.(string)
					if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
						releaseChannelLatestVersion[channelName] = v
						break
					}
				}
			}
		}
	}

	if err := d.Set("release_channel_default_version", releaseChannelDefaultVersion); err != nil {
		return fmt.Errorf("Error setting release_channel_default_version: %s", err)
	}
	if err := d.Set("release_channel_latest_version", releaseChannelLatestVersion); err != nil {
		return fmt.Errorf("Error setting release_channel_latest_version: %s", err)
	}
	if err := d.Set("release_channel_upgrade_target_version", releaseChannelUpgradeTargetVersion); err != nil {
		return fmt.Errorf("Error setting release_channel_upgrade_target_version: %s", err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_container_engine_versions",
		ProductName: "container",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleContainerEngineVersions(),
	}.Register()
}
