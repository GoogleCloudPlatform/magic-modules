package compute_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/test"
)

// func TestAccComputeInstanceConfidentialInstanceConfigMain(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstanceConfidentialInstanceConfigMain",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_GracefulShutdownWithResetUpdate(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_GracefulShutdownWithResetUpdate",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_IP(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_IP",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_IPv6(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_IPv6",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_advancedMachineFeatures(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_advancedMachineFeatures",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_aliasIpRangeCommonAddresses(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_aliasIpRangeCommonAddresses",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_attachedDisk(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_attachedDisk",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_attachedDiskUpdate(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_attachedDiskUpdate",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_attachedDisk_modeRo(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_attachedDisk_modeRo",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_attachedDisk_sourceUrl(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_attachedDisk_sourceUrl",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_autoDeleteUpdate(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_autoDeleteUpdate",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

func TestAccComputeInstance_basic1(t *testing.T) {
	t.Parallel()

	test.AssertTestFile(
		t,
		"TestAccComputeInstance_basic1",
		"google_compute_instance",
		"compute.googleapis.com/Instance",
		[]string{
			"desired_status",
			"metadata",
		},
	)
}

// func TestAccComputeInstance_basic2(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_basic2",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_basic3(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_basic3",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_basic4(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_basic4",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_basic5(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_basic5",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_bootDisk_mode(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_bootDisk_mode",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_bootDisk_type(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_bootDisk_type",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_confidentialHyperDiskBootDisk(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_confidentialHyperDiskBootDisk",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_creationOnlyAttributionLabel(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_creationOnlyAttributionLabel",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_deletionProtectionExplicitFalse(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_deletionProtectionExplicitFalse",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_descriptionUpdate(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_descriptionUpdate",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_desiredStatusTerminatedOnCreation(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_desiredStatusTerminatedOnCreation",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_desiredStatusTerminatedUpdateFields(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_desiredStatusTerminatedUpdateFields",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_desiredStatusUpdateBasic(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_desiredStatusUpdateBasic",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_desiredStatus_suspended(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_desiredStatus_suspended",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_diskEncryption(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_diskEncryption",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_diskEncryptionRestart(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_diskEncryptionRestart",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_diskResourcePolicies(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_diskResourcePolicies",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_diskResourcePolicies_attachmentDiff(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_diskResourcePolicies_attachmentDiff",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_enableDisplay(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_enableDisplay",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_enableUefiNetworking(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_enableUefiNetworking",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_forceNewAndChangeMetadata(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_forceNewAndChangeMetadata",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_guestAccelerator(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_guestAccelerator",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_hostErrorTimeoutSecconds(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_hostErrorTimeoutSecconds",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_hostname(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_hostname",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{
// 			"desired_status",
// 			"metadata",
// 		},
// 	)
// }

// func TestAccComputeInstance_bootDisk_source(t *testing.T) {
// 	t.Parallel()

// 	test.AssertTestFile(
// 		t,
// 		"TestAccComputeInstance_bootDisk_source",
// 		"google_compute_instance",
// 		"compute.googleapis.com/Instance",
// 		[]string{},
// 	)
// }
