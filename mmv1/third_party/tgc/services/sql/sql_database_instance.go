// ----------------------------------------------------------------------------
//
//	This file is copied here by Magic Modules. The code for building up a
//	sql database instance object is copied from the manually implemented
//	provider file:
//	third_party/terraform/resources/resource_sql_database_instance.go.erb.go
//
// ----------------------------------------------------------------------------
package sql

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

const SQLDatabaseInstanceAssetType string = "sqladmin.googleapis.com/Instance"

func ResourceConverterSQLDatabaseInstance() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: SQLDatabaseInstanceAssetType,
		Convert:   GetSQLDatabaseInstanceCaiObject,
	}
}

func GetSQLDatabaseInstanceCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//cloudsql.googleapis.com/projects/{{project}}/instances/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetSQLDatabaseInstanceApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: SQLDatabaseInstanceAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1beta4",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/sqladmin/v1beta4/rest",
				DiscoveryName:        "DatabaseInstance",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetSQLDatabaseInstanceApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return nil, err
	}

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		name = id.UniqueId()
	}

	instance := &sqladmin.DatabaseInstance{
		Project:              project,
		Name:                 name,
		Region:               region,
		Settings:             expandSqlDatabaseInstanceSettings(d.Get("settings").([]interface{}), !isFirstGen(d)),
		DatabaseVersion:      d.Get("database_version").(string),
		MasterInstanceName:   d.Get("master_instance_name").(string),
		ReplicaConfiguration: expandReplicaConfiguration(d.Get("replica_configuration").([]interface{})),
	}

	return cai.JsonMap(instance)
}

// Detects whether a database is 1st Generation by inspecting the tier name
func isFirstGen(d tpgresource.TerraformResourceData) bool {
	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})
	tier := settings["tier"].(string)

	// 1st Generation databases have tiers like 'D0', as opposed to 2nd Generation which are
	// prefixed with 'db'
	return !regexp.MustCompile("db*").Match([]byte(tier))
}

func expandSqlDatabaseInstanceSettings(configured []interface{}, secondGen bool) *sqladmin.Settings {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_settings := configured[0].(map[string]interface{})
	settings := &sqladmin.Settings{
		// Version is unset in Create but is set during update
		SettingsVersion:     int64(_settings["version"].(int)),
		Tier:                _settings["tier"].(string),
		ForceSendFields:     []string{"StorageAutoResize"},
		ActivationPolicy:    _settings["activation_policy"].(string),
		AvailabilityType:    _settings["availability_type"].(string),
		DataDiskSizeGb:      int64(_settings["disk_size"].(int)),
		DataDiskType:        _settings["disk_type"].(string),
		PricingPlan:         _settings["pricing_plan"].(string),
		UserLabels:          tpgresource.ConvertStringMap(_settings["user_labels"].(map[string]interface{})),
		BackupConfiguration: expandBackupConfiguration(_settings["backup_configuration"].([]interface{})),
		DatabaseFlags:       expandDatabaseFlags(_settings["database_flags"].(*schema.Set).List()),
		IpConfiguration:     expandIpConfiguration(_settings["ip_configuration"].([]interface{})),
		LocationPreference:  expandLocationPreference(_settings["location_preference"].([]interface{})),
		MaintenanceWindow:   expandMaintenanceWindow(_settings["maintenance_window"].([]interface{})),
	}

	// 1st Generation instances don't support the disk_autoresize parameter
	// and it defaults to true - so we shouldn't set it if this is first gen
	if secondGen {
		diskAutoresize := _settings["disk_autoresize"].(bool)
		settings.StorageAutoResize = &diskAutoresize
		settings.StorageAutoResizeLimit = int64(_settings["disk_autoresize_limit"].(int))
	}

	return settings
}

func expandReplicaConfiguration(configured []interface{}) *sqladmin.ReplicaConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_replicaConfiguration := configured[0].(map[string]interface{})
	return &sqladmin.ReplicaConfiguration{
		FailoverTarget: _replicaConfiguration["failover_target"].(bool),

		// MysqlReplicaConfiguration has been flattened in the TF schema, so
		// we'll keep it flat here instead of another expand method.
		MysqlReplicaConfiguration: &sqladmin.MySqlReplicaConfiguration{
			CaCertificate:           _replicaConfiguration["ca_certificate"].(string),
			ClientCertificate:       _replicaConfiguration["client_certificate"].(string),
			ClientKey:               _replicaConfiguration["client_key"].(string),
			ConnectRetryInterval:    int64(_replicaConfiguration["connect_retry_interval"].(int)),
			DumpFilePath:            _replicaConfiguration["dump_file_path"].(string),
			MasterHeartbeatPeriod:   int64(_replicaConfiguration["master_heartbeat_period"].(int)),
			Password:                _replicaConfiguration["password"].(string),
			SslCipher:               _replicaConfiguration["ssl_cipher"].(string),
			Username:                _replicaConfiguration["username"].(string),
			VerifyServerCertificate: _replicaConfiguration["verify_server_certificate"].(bool),
		},
	}
}

func expandMaintenanceWindow(configured []interface{}) *sqladmin.MaintenanceWindow {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	window := configured[0].(map[string]interface{})
	return &sqladmin.MaintenanceWindow{
		Day:             int64(window["day"].(int)),
		Hour:            int64(window["hour"].(int)),
		UpdateTrack:     window["update_track"].(string),
		ForceSendFields: []string{"Hour"},
	}
}

func expandLocationPreference(configured []interface{}) *sqladmin.LocationPreference {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_locationPreference := configured[0].(map[string]interface{})
	return &sqladmin.LocationPreference{
		FollowGaeApplication: _locationPreference["follow_gae_application"].(string),
		Zone:                 _locationPreference["zone"].(string),
		SecondaryZone:        _locationPreference["secondary_zone"].(string),
	}
}

func expandIpConfiguration(configured []interface{}) *sqladmin.IpConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_ipConfiguration := configured[0].(map[string]interface{})

	return &sqladmin.IpConfiguration{
		Ipv4Enabled:        _ipConfiguration["ipv4_enabled"].(bool),
		PrivateNetwork:     _ipConfiguration["private_network"].(string),
		AuthorizedNetworks: expandAuthorizedNetworks(_ipConfiguration["authorized_networks"].(*schema.Set).List()),
		ForceSendFields:    []string{"Ipv4Enabled"},
		NullFields:         []string{"RequireSsl"},
		SslMode:            _ipConfiguration["ssl_mode"].(string),
	}
}
func expandAuthorizedNetworks(configured []interface{}) []*sqladmin.AclEntry {
	an := make([]*sqladmin.AclEntry, 0, len(configured))
	for _, _acl := range configured {
		_entry := _acl.(map[string]interface{})
		an = append(an, &sqladmin.AclEntry{
			ExpirationTime: _entry["expiration_time"].(string),
			Name:           _entry["name"].(string),
			Value:          _entry["value"].(string),
		})
	}

	return an
}

func expandDatabaseFlags(configured []interface{}) []*sqladmin.DatabaseFlags {
	databaseFlags := make([]*sqladmin.DatabaseFlags, 0, len(configured))
	for _, _flag := range configured {
		_entry := _flag.(map[string]interface{})

		databaseFlags = append(databaseFlags, &sqladmin.DatabaseFlags{
			Name:  _entry["name"].(string),
			Value: _entry["value"].(string),
		})
	}
	return databaseFlags
}

func expandBackupConfiguration(configured []interface{}) *sqladmin.BackupConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_backupConfiguration := configured[0].(map[string]interface{})
	return &sqladmin.BackupConfiguration{
		BinaryLogEnabled: _backupConfiguration["binary_log_enabled"].(bool),
		Enabled:          _backupConfiguration["enabled"].(bool),
		StartTime:        _backupConfiguration["start_time"].(string),
		Location:         _backupConfiguration["location"].(string),
	}
}
