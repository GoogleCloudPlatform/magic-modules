package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

const (
	// Copied from the official Google Cloud auto-generated client.
	ProjectRegex         = verify.ProjectRegex
	ProjectRegexWildCard = verify.ProjectRegexWildCard
	RegionRegex          = verify.RegionRegex
	SubnetworkRegex      = verify.SubnetworkRegex

	SubnetworkLinkRegex = verify.SubnetworkLinkRegex

	RFC1035NameTemplate = verify.RFC1035NameTemplate
	CloudIoTIdRegex     = verify.CloudIoTIdRegex

	// Format of default Compute service accounts created by Google
	// ${PROJECT_ID}-compute@developer.gserviceaccount.com where PROJECT_ID is an int64 (max 20 digits)
	ComputeServiceAccountNameRegex = verify.ComputeServiceAccountNameRegex

	// https://cloud.google.com/iam/docs/understanding-custom-roles#naming_the_role
	IAMCustomRoleIDRegex = verify.IAMCustomRoleIDRegex

	// https://cloud.google.com/managed-microsoft-ad/reference/rest/v1/projects.locations.global.domains/create#query-parameters
	ADDomainNameRegex = verify.ADDomainNameRegex
)

var (
	// Service account name must have a length between 6 and 30.
	// The first and last characters have different restrictions, than
	// the middle characters. The middle characters length must be between
	// 4 and 28 since the first and last character are excluded.
	ServiceAccountNameRegex = verify.ServiceAccountNameRegex

	ServiceAccountLinkRegexPrefix = verify.ServiceAccountLinkRegexPrefix
	PossibleServiceAccountNames   = verify.PossibleServiceAccountNames
	ServiceAccountLinkRegex       = verify.ServiceAccountLinkRegex

	ServiceAccountKeyNameRegex = verify.ServiceAccountKeyNameRegex

	// Format of service accounts created through the API
	CreatedServiceAccountNameRegex = verify.CreatedServiceAccountNameRegex

	// Format of service-created service account
	// examples are:
	// 		$PROJECTID@cloudbuild.gserviceaccount.com
	// 		$PROJECTID@cloudservices.gserviceaccount.com
	// 		$PROJECTID@appspot.gserviceaccount.com
	ServiceDefaultAccountNameRegex = verify.ServiceDefaultAccountNameRegex

	ProjectNameInDNSFormRegex = verify.ProjectNameInDNSFormRegex
	ProjectNameRegex          = verify.ProjectNameRegex

	// Valid range for Cloud Router ASN values as per RFC6996
	// https://tools.ietf.org/html/rfc6996
	// Must be explicitly int64 to avoid overflow when building Terraform for 32bit architectures
	Rfc6996Asn16BitMin  = verify.Rfc6996Asn16BitMin
	Rfc6996Asn16BitMax  = verify.Rfc6996Asn16BitMax
	Rfc6996Asn32BitMin  = verify.Rfc6996Asn32BitMin
	Rfc6996Asn32BitMax  = verify.Rfc6996Asn32BitMax
	GcpRouterPartnerAsn = verify.GcpRouterPartnerAsn
)

var rfc1918Networks = verify.Rfc1918Networks

// validateGCEName ensures that a field matches the requirements for Compute Engine resource names
// https://cloud.google.com/compute/docs/naming-resources#resource-name-format
func validateGCEName(v interface{}, k string) (ws []string, errors []error) {
	return verify.ValidateGCEName(v, k)
}

// Ensure that the BGP ASN value of Cloud Router is a valid value as per RFC6996 or a value of 16550
func validateRFC6996Asn(v interface{}, k string) (ws []string, errors []error) {
	return verify.ValidateRFC6996Asn(v, k)
}

func validateRegexp(re string) schema.SchemaValidateFunc {
	return verify.ValidateRegexp(re)
}

func validateEnum(values []string) schema.SchemaValidateFunc {
	return verify.ValidateEnum(values)
}

func validateRFC1918Network(min, max int) schema.SchemaValidateFunc {
	return verify.ValidateRFC1918Network(min, max)
}

func validateRFC3339Time(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateRFC3339Time(v, k)
}

func validateRFC1035Name(min, max int) schema.SchemaValidateFunc {
	return verify.ValidateRFC1035Name(min, max)
}

func validateIpCidrRange(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateIpCidrRange(v, k)
}

func validateIAMCustomRoleID(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateIAMCustomRoleID(v, k)
}

func orEmpty(f schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return verify.OrEmpty(f)
}

func validateProjectID() schema.SchemaValidateFunc {
	return verify.ValidateProjectID()
}

func validateDSProjectID() schema.SchemaValidateFunc {
	return verify.ValidateDSProjectID()
}

func validateProjectName() schema.SchemaValidateFunc {
	return verify.ValidateProjectName()
}

func validateDuration() schema.SchemaValidateFunc {
	return verify.ValidateDuration()
}

func validateNonNegativeDuration() schema.SchemaValidateFunc {
	return verify.ValidateNonNegativeDuration()
}

func validateIpAddress(i interface{}, val string) ([]string, []error) {
	return verify.ValidateIpAddress(i, val)
}

func validateBase64String(i interface{}, val string) ([]string, []error) {
	return verify.ValidateBase64String(i, val)
}

// StringNotInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and that it matches none of the element in the invalid slice.
// if ignorecase is true, case is ignored.
func StringNotInSlice(invalid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return verify.StringNotInSlice(invalid, ignoreCase)
}

// Ensure that hourly timestamp strings "HH:MM" have the minutes zeroed out for hourly only inputs
func validateHourlyOnly(val interface{}, key string) (warns []string, errs []error) {
	return verify.ValidateHourlyOnly(val, key)
}

func validateRFC3339Date(v interface{}, k string) (warnings []string, errors []error) {
	return verify.ValidateRFC3339Date(v, k)
}

func validateADDomainName() schema.SchemaValidateFunc {
	return verify.ValidateADDomainName()
}

func testStringValidationCases(cases []StringValidationTestCase, validationFunc schema.SchemaValidateFunc) []error {
	es := make([]error, 0)
	for _, c := range cases {
		es = append(es, testStringValidation(c, validationFunc)...)
	}

	return es
}

func testStringValidation(testCase StringValidationTestCase, validationFunc schema.SchemaValidateFunc) []error {
	_, es := validationFunc(testCase.Value, testCase.TestName)
	if testCase.ExpectError {
		if len(es) > 0 {
			return nil
		} else {
			return []error{fmt.Errorf("Didn't see expected error in case \"%s\" with string \"%s\"", testCase.TestName, testCase.Value)}
		}
	}

	return es
}

type StringValidationTestCase struct {
	TestName    string
	Value       string
	ExpectError bool
}
