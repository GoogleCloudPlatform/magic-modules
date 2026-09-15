package resourcemanager

import "testing"

func TestServiceAccountIamMemberResource_HasIdentity(t *testing.T) {
	resource := ServiceAccountIamMemberResource()
	if resource.Identity == nil {
		t.Fatalf("expected google_service_account_iam_member resource identity to be configured")
	}
}

func TestServiceAccountIamParentResourceIdentityParser(t *testing.T) {
	resource := ServiceAccountIamMemberResource()
	rd := resource.TestResourceData()
	identity, err := rd.Identity()
	if err != nil {
		t.Fatalf("creating identity: %v", err)
	}

	serviceAccountID := "projects/my-project/serviceAccounts/my-sa@my-project.iam.gserviceaccount.com"
	if err := identity.Set("service_account_id", serviceAccountID); err != nil {
		t.Fatalf("setting identity service_account_id: %v", err)
	}

	got, err := ServiceAccountIamParentResourceIdentityParser(rd, identity, nil)
	if err != nil {
		t.Fatalf("parsing identity: %v", err)
	}
	if got != serviceAccountID {
		t.Fatalf("unexpected parsed id: got %q, want %q", got, serviceAccountID)
	}
}
