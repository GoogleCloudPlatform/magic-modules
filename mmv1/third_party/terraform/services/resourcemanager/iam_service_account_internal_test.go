package resourcemanager

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func TestValidateServiceAccountResourceName(t *testing.T) {
	runSAValidation(t, "service_account_id", verify.ValidateRegexp(serviceAccountResourceNameRegex), []saValidationCase{
		{value: "projects/my-project/serviceAccounts/svcacct@my-project.iam.gserviceaccount.com"},
		{value: "projects/my-project/serviceAccounts/123456789012@serverless-robot-prod.iam.gserviceaccount.com"},
		{value: "projects/my-project/serviceAccounts/123456789012345678901"},
		{value: "projects/-/serviceAccounts/123456789012345678901"},

		{value: "svcacct@my-project.iam.gserviceaccount.com", expectError: true},
		{value: "projects/my-project/serviceAccounts/", expectError: true},
		{value: "evil/projects/my-project/serviceAccounts/123456789012345678901", expectError: true},
		{value: "projects/my-project/serviceAccounts/123456789012345678901/keys/abc", expectError: true},
		{value: "projects/my-project/serviceAccounts/123456789012345678901 ", expectError: true},
		{value: "projects/my-project/serviceAccounts/svc acct@my-project.iam.gserviceaccount.com", expectError: true},
		{value: "projects/my-project/serviceAccounts/123456789012345678901\n", expectError: true},
	})
}
