package sql

import (
	"testing"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func TestSqlUserOptionalFields_DirectTypedAccess(t *testing.T) {
	user := &sqladmin.User{
		IamEmail:      "user@example.com",
		DatabaseRoles: []string{"role_a", "role_b"},
	}

	iamEmail, databaseRoles := sqlUserOptionalFields(user)
	if iamEmail != "user@example.com" {
		t.Fatalf("expected iam email to be preserved, got %q", iamEmail)
	}
	if len(databaseRoles) != 2 || databaseRoles[0] != "role_a" || databaseRoles[1] != "role_b" {
		t.Fatalf("expected database roles to be preserved, got %#v", databaseRoles)
	}

	// Ensure returned roles are copied and cannot mutate source user roles.
	databaseRoles[0] = "mutated"
	if user.DatabaseRoles[0] != "role_a" {
		t.Fatalf("expected user database roles to remain unchanged, got %#v", user.DatabaseRoles)
	}
}

func TestResourceDataForListedSqlUser_IsolatesOptionalFields(t *testing.T) {
	project := "test-project"
	instance := "test-instance"

	userWithOptionals := &sqladmin.User{
		Name:          "user_with",
		Host:          "%",
		Instance:      instance,
		Type:          "BUILT_IN",
		IamEmail:      "with@example.com",
		DatabaseRoles: []string{"role_a"},
		PasswordPolicy: &sqladmin.UserPasswordValidationPolicy{
			AllowedFailedAttempts: 1,
		},
	}
	d1, err := resourceDataForListedSqlUser(userWithOptionals, project)
	if err != nil {
		t.Fatalf("unexpected error building resource data for userWithOptionals: %v", err)
	}
	if got := d1.Get("iam_email").(string); got != "with@example.com" {
		t.Fatalf("expected iam_email for userWithOptionals, got %q", got)
	}
	if got, ok := d1.GetOk("database_roles"); !ok || len(got.([]interface{})) != 1 {
		t.Fatalf("expected one database role for userWithOptionals, got %#v (ok=%v)", got, ok)
	}
	if got, ok := d1.GetOk("password_policy"); !ok || len(got.([]interface{})) != 1 {
		t.Fatalf("expected password_policy for userWithOptionals, got %#v (ok=%v)", got, ok)
	}

	userWithoutOptionals := &sqladmin.User{
		Name:     "user_without",
		Host:     "",
		Instance: instance,
		Type:     "BUILT_IN",
	}
	d2, err := resourceDataForListedSqlUser(userWithoutOptionals, project)
	if err != nil {
		t.Fatalf("unexpected error building resource data for userWithoutOptionals: %v", err)
	}
	if got := d2.Get("iam_email").(string); got != "" {
		t.Fatalf("expected empty iam_email for userWithoutOptionals, got %q", got)
	}
	if got, ok := d2.GetOk("database_roles"); ok && len(got.([]interface{})) > 0 {
		t.Fatalf("expected no database_roles for userWithoutOptionals, got %#v", got)
	}
	if got, ok := d2.GetOk("password_policy"); ok && len(got.([]interface{})) > 0 {
		t.Fatalf("expected no password_policy for userWithoutOptionals, got %#v", got)
	}
	if got := d2.Get("host").(string); got != "" {
		t.Fatalf("expected empty host for userWithoutOptionals, got %q", got)
	}
}
