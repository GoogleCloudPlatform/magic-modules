package backupdr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceBackupDRAutoProtectionPolicy_ResourceTypeForceNew(t *testing.T) {
	res := ResourceBackupDRAutoProtectionPolicy()
	bpdSchema, ok := res.Schema["backup_plan_details"]
	if !ok {
		t.Fatal("expected backup_plan_details in AutoProtectionPolicy schema")
	}
	elem, ok := bpdSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatalf("expected backup_plan_details Elem to be *schema.Resource, got %T", bpdSchema.Elem)
	}

	rtSchema, ok := elem.Schema["resource_type"]
	if !ok {
		t.Fatal("expected resource_type in backup_plan_details schema")
	}
	if !rtSchema.ForceNew {
		t.Error("expected backup_plan_details.resource_type to have ForceNew: true")
	}

	bpSchema, ok := elem.Schema["backup_plan"]
	if !ok {
		t.Fatal("expected backup_plan in backup_plan_details schema")
	}
	if bpSchema.ForceNew {
		t.Error("expected backup_plan_details.backup_plan to remain updatable (ForceNew: false)")
	}
}
