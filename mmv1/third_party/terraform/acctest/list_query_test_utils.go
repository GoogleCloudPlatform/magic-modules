package acctest

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// CapturedDisplayName captures a resource's display-name attribute during a
// create step and exposes a knownvalue.Check that asserts list-query results
// match the captured value.
type CapturedDisplayName struct {
	value string
}

// CaptureCheck returns a TestCheckFunc that copies the first non-empty value
// among attrCandidates from resourceAddr's state into the captured value.
// attrCandidates are checked in order from first index to last.
func (c *CapturedDisplayName) CaptureCheck(resourceAddr string, attrCandidates []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddr]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceAddr)
		}
		for _, k := range attrCandidates {
			if v, ok := rs.Primary.Attributes[k]; ok && v != "" {
				c.value = v
				return nil
			}
		}
		return fmt.Errorf("no display name attribute found in state for resource %s; tried %v", resourceAddr, attrCandidates)
	}
}

// ExpectKnownValue returns a knownvalue.Check that compares against the
// captured value. Fails if CaptureCheck has not run yet.
func (c *CapturedDisplayName) ExpectKnownValue() knownvalue.Check {
	return knownvalue.StringFunc(func(v string) error {
		if c.value == "" {
			return fmt.Errorf("display name was not captured from create step")
		}
		if v != c.value {
			return fmt.Errorf("expected display name %q, got %q", c.value, v)
		}
		return nil
	})
}
