package compute_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/test"
)

// Total 124 tests
func TestAccComputeInstance_serviceAccountEmail_0scopes(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateTerminated_desiredStatusNotSet_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_keyRevocationActionType(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_scheduling(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_networkIPAuto(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_resourcePolicySpread(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_subnet_auto(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_with375GbScratchDisk(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_bootDisk_type(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_forceNewAndChangeMetadata(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_attachedDisk_modeRo(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_autoDeleteUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"boot_disk.auto_delete",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_desiredStatusUpdateBasic(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_metadataStartupScript_gracefulSwitch(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstanceConfidentialInstanceConfigMain(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_advancedMachineFeatures(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_localSsdRecoveryTimeout(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_queueCount(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_spotVM_maxRunDuration_deleteTerminationAction(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_subnetworkUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_basic4(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_descriptionUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_deletionProtectionExplicitFalse(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"deletion_protection",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_desiredStatus_suspended(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_partnerMetadata_deletePartnerMetadata(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_secondaryAliasIpRange(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_aliasIpRangeCommonAddresses(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_confidentialHyperDiskBootDisk(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_maxRunDuration_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_multiNic(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_nictype_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_regionBootDisk(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_resourcePolicyCollocate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_shieldedVmConfig(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"shielded_instance_config.enable_vtpm", // It has default value
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_desiredStatusTerminatedUpdateFields(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_diskResourcePolicies(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_soleTenantNodeAffinities(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_creationOnlyAttributionLabel(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_performanceMonitoringUnit(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_stopInstanceToUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_IPv6(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_bootDisk_source(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_enableUefiNetworking(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_machineTypeUrl(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_metadataStartupScript_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"metadata_startup_script",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_noServiceAccount(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"service_account",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_reservationAffinities(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_resourcePolicyUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
			"resource_policies",
		},
	)
}

func TestAccComputeInstance_attachedDiskUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
			"attached_disk.disk_encryption_key_raw",
		},
	)
}

func TestAccComputeInstance_diskEncryptionRestart(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
			"boot_disk.disk_encryption_key_raw",
			"attached_disk.disk_encryption_key_raw",
		},
	)
}

func TestAccComputeInstance_spotVM(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_desiredStatusTerminatedOnCreation(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_internalIPv6(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_networkPerformanceConfig(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_IP(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_attachedDisk(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_network_ip_custom(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_serviceAccount(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateRunning_desiredStatusTerminated_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_bootDisk_mode(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"boot_disk.mode", // It has default value
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_bootDisk_sourceUrl(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_diskResourcePolicies_attachmentDiff(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_subnet_custom(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_basic1(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_basic5(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_hostname(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_ipv6ExternalReservation(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_kmsDiskEncryption(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_partnerMetadata(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"partner_metadata", // not in cai asset
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_serviceAccount_updated0to1to0scopes(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"service_account", // empty
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_attachedDisk_sourceUrl(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_hostErrorTimeoutSecconds(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"scheduling.automatic_restart", // It has the default value true
			"scheduling.host_error_timeout_seconds",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_localSsdVM_maxRunDuration_stopTerminationAction(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_networkIpUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_basic2(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_diskEncryption(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
			"boot_disk.disk_encryption_key_raw",     // missing in asset
			"attached_disk.disk_encryption_key_raw", // missing in asset
		},
	)
}

func TestAccComputeInstance_guestAccelerator(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_minCpuPlatform(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// User may specify AUTOMATIC using any case; it's empty in asset
			"min_cpu_platform",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_networkTier(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_primaryAliasIpRange(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_private_image_family(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_serviceAccount_updated(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_GracefulShutdownWithResetUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// graceful_shutdown is not in asset
			"scheduling.graceful_shutdown",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_enableDisplay(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateTerminated_desiredStatusNotSet_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_standardVM_maxRunDuration_stopTerminationAction(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_basic3(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			"can_ip_forward",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_resourceManagerTags(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_GracefulShutdownWithoutResetUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			"scheduling.graceful_shutdown", // not in asset
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_NetworkAttachmentUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_NetworkAttachment(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_NicStackTypeUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_PTRRecord(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_NicStackType_IPV6(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_attachedDisk_RSAencryption(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
			"attached_disk.disk_encryption_service_account", // not in cai asset
			"attached_disk.disk_encryption_key_rsa",         // not in cai asset
		},
	)
}

func TestAccComputeInstance_bootAndAttachedDisk_interface(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_bootDisk_storagePoolSpecified(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_bootDisk_storagePoolSpecified_nameOnly(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_creationOnlyAttributionLabelConfiguredOnUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_deletionProtectionExplicitTrueAndUpdateFalse(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"deletion_protection",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_desiredStatusSuspendedOnCreation(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_forceChangeMachineTypeManually(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_guestAcceleratorSkip(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
			"guest_accelerator", // Its count is 0
		},
	)
}

func TestAccComputeInstance_guestOsFeatures(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_imageEncryption(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_instanceEncryption(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_internalIPv6PrefixLength(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_localSsdRecoveryTimeout_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_partnerMetadata_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			"partner_metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_proactiveAttributionLabel(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_rsaBootDiskEncryption(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			"boot_disk.disk_encryption_key_rsa",         // not in asset
			"boot_disk.disk_encryption_service_account", // not in asset
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_schedulingTerminationTime(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_snapshot(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_snapshotEncryption(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
			"boot_disk.disk_encryption_key_raw",
		},
	)
}

func TestAccComputeInstance_spotVM_maxRunDuration_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_standardVM_maxRunDuration_deleteTerminationAction(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_spotVM_update(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateRunning_desiredStatusRunning_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateRunning_desiredStatusTerminated_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateRunning_desiredStatusRunning_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk",
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateRunning_desiredStatusNotSet_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateTerminated_desiredStatusTerminated_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateTerminated_desiredStatusRunning_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}

func TestAccComputeInstance_updateTerminated_desiredStatusRunning_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		[]string{
			"desired_status",
			"allow_stopping_for_update",
			"metadata",
			// params.resource_manager_tags is not in instance asset
			"params.resource_manager_tags",
			"params",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
		},
	)
}
