package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	_ "github.com/hashicorp/terraform-provider-google/google/services/firebase"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
