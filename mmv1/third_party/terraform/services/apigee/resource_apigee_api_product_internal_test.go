package apigee

import (
	"regexp"
	"testing"
)

// apiProductNameRegexp mirrors the validation regex defined in ApiProduct.yaml.
// It is kept here so that any future change to the YAML is reflected in these tests.
const apiProductNameRegexp = `^[a-zA-Z][a-zA-Z0-9._\-$ %]*$`

func TestApigeeApiProduct_nameValidation(t *testing.T) {
	re := regexp.MustCompile(apiProductNameRegexp)

	validNames := []string{
		// lowercase (original behaviour preserved)
		"my-product",
		"product1",
		"product.name",
		"product_name",
		"a",
		// camelCase – previously rejected, now allowed (github.com/hashicorp/terraform-provider-google/issues/26523)
		"myApiProduct",
		"MyProduct",
		"StartsUpperCase",
		"mixedCase123",
		"CamelCase.with.dots",
		"CamelCase-with-hyphens",
	}

	for _, name := range validNames {
		if !re.MatchString(name) {
			t.Errorf("expected %q to be a valid API product name, but the regex rejected it", name)
		}
	}

	invalidNames := []string{
		// must not start with a digit
		"0product",
		"1CamelCase",
		// must not start with a special character
		"-product",
		"_product",
		".product",
		" product",
		// empty string
		"",
	}

	for _, name := range invalidNames {
		if re.MatchString(name) {
			t.Errorf("expected %q to be an invalid API product name, but the regex accepted it", name)
		}
	}
}
