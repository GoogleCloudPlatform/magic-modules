// Contains functions that don't really belong anywhere else.

package google

import (
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// This is a Printf sibling (Nprintf; Named Printf), which handles strings like
// Nprintf("Hello %{target}!", map[string]interface{}{"target":"world"}) == "Hello world!".
// This is particularly useful for generated tests, where we don't want to use Printf,
// since that would require us to generate a very particular ordering of arguments.
//
// Deprecated: For backward compatibility Nprintf is still working,
// but all new code should use Nprintf in the acctest package instead.
func Nprintf(format string, params map[string]interface{}) string {
	return acctest.Nprintf(format, params)
}
