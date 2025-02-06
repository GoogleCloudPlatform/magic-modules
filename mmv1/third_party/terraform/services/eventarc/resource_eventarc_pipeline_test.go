package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventarcPipeline_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcPipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcPipeline_update0(context),
			},
			{
				ResourceName:            "google_eventarc_pipeline.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccEventarcPipeline_update1(context),
			},
			{
				ResourceName:            "google_eventarc_pipeline.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcPipeline_update0(context map[string]interface{}) string {
	return acctest.Nprintf(`
`, context)
}

func testAccEventarcPipeline_update1(context map[string]interface{}) string {
	return acctest.Nprintf(`
`, context)
}
