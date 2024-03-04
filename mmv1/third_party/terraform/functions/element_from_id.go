package functions

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ValidateElementFromIdArguments is reusable validation logic used in provider-defined functions that use the GetElementFromId function
func ValidateElementFromIdArguments(ctx context.Context, input string, regex *regexp.Regexp, pattern string) *function.FuncError {
	submatches := regex.FindAllStringSubmatchIndex(input, -1)

	// Zero matches means unusable input; error returned
	if len(submatches) == 0 {
		return function.NewArgumentFuncError(0, fmt.Sprintf("The input string \"%s\" doesn't contain the expected pattern \"%s\".", input, pattern))
	}

	// >1 matches means input usable but not ideal; debug log
	if len(submatches) > 1 {
		tflog.Debug(ctx, fmt.Sprintf("The input string \"%s\" contains more than one match for the pattern \"%s\" expected by the provider function. Terraform will use the first found match.", input, pattern))
	}

	return nil
}

// GetElementFromId is reusable logic that is used in multiple provider-defined functions for pulling elements out of self links and ids of resources and data sources
func GetElementFromId(input string, regex *regexp.Regexp, template string) string {
	submatches := regex.FindAllStringSubmatchIndex(input, -1)
	submatch := submatches[0] // Take the only / left-most submatch
	dst := []byte{}
	return string(regex.ExpandString(dst, template, input, submatch))
}
