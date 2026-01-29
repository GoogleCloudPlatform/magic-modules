package transcoder_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccTranscoderJobTemplate_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "transcoder-job-template-tagkey", map[string]interface{}{})
	jobTemplateID := "tf-test-" + acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org":             envvar.GetTestOrgFromEnv(t),
		"tagKey":          tagKey,
		"tagValue":        acctest.BootstrapSharedTestOrganizationTagValue(t, "transcoder-job-template-tagvalue", tagKey),
		"job_template_id": jobTemplateID,
		"project":         envvar.GetTestProjectFromEnv(),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTranscoderJobTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTranscoderJobTemplateTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_transcoder_job_template.test", "tags.%"),
					checkTranscoderJobTemplateTags(t),
				),
			},
			{
				ResourceName:            "google_transcoder_job_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func checkTranscoderJobTemplateTags(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_transcoder_job_template" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			// 1. Get the configured tag key and value from the state.
			var configuredTagValueNamespacedName string
			var tagKeyNamespacedName, tagValueShortName string
			for key, val := range rs.Primary.Attributes {
				if strings.HasPrefix(key, "tags.") && key != "tags.%" {
					tagKeyNamespacedName = strings.TrimPrefix(key, "tags.")
					tagValueShortName = val
					if tagValueShortName != "" {
						configuredTagValueNamespacedName = fmt.Sprintf("%s/%s", tagKeyNamespacedName, tagValueShortName)
						break
					}
				}
			}

			if configuredTagValueNamespacedName == "" {
				return fmt.Errorf("could not find a configured tag value in the state for resource %s", rs.Primary.ID)
			}

			if strings.Contains(configuredTagValueNamespacedName, "%{") {
				return fmt.Errorf("tag namespaced name contains unsubstituted variables: %q. Ensure the context map in the test step is populated", configuredTagValueNamespacedName)
			}

			safeNamespacedName := url.QueryEscape(configuredTagValueNamespacedName)
			describeTagValueURL := fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagValues/namespaced?name=%s", safeNamespacedName)

			respDescribe, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    describeTagValueURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error describing tag value using namespaced name %q: %v", configuredTagValueNamespacedName, err)
			}

			fullTagValueName, ok := respDescribe["name"].(string)
			if !ok || fullTagValueName == "" {
				return fmt.Errorf("tag value details (name) not found in response for namespaced name: %q, response: %v", configuredTagValueNamespacedName, respDescribe)
			}

			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3]
			instance_id := parts[5]

			parentURL := fmt.Sprintf("//transcoder.googleapis.com/projects/%s/locations/%s/jobTemplates/%s", project, location, instance_id)
			listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, url.QueryEscape(parentURL))

			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    listBindingsURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error calling TagBindings API: %v", err)
			}

			tagBindingsVal, exists := resp["tagBindings"]
			if !exists {
				tagBindingsVal = []interface{}{}
			}

			tagBindings, ok := tagBindingsVal.([]interface{})
			if !ok {
				return fmt.Errorf("'tagBindings' is not a slice in response for resource %s. Response: %v", rs.Primary.ID, resp)
			}

			foundMatch := false
			for _, binding := range tagBindings {
				bindingMap, ok := binding.(map[string]interface{})
				if !ok {
					continue
				}
				if bindingMap["tagValue"] == fullTagValueName {
					foundMatch = true
					break
				}
			}

			if !foundMatch {
				return fmt.Errorf("expected tag value %s (from namespaced %q) not found in tag bindings for resource %s. Bindings: %v", fullTagValueName, configuredTagValueNamespacedName, rs.Primary.ID, tagBindings)
			}

			t.Logf("Successfully found matching tag binding for %s with tagValue %s", rs.Primary.ID, fullTagValueName)
		}

		return nil
	}
}

func testAccTranscoderJobTemplateTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_transcoder_job_template" "test" {
  job_template_id = "%{job_template_id}"
  location        = "us-central1"
  project         = "%{project}"

  config {
    inputs {
      key = "input0"
    }
    edit_list {
      key               = "atom0"
      inputs            = ["input0"]
      start_time_offset = "0s"
    }
    ad_breaks {
      start_time_offset = "3.500s"
    }
    elementary_streams {
      key = "video-stream0"
      video_stream {
        h264 {
          width_pixels      = 640
          height_pixels     = 360
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 550000
          vbv_fullness_bits = 495000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"
        }
      }
    }
    elementary_streams {
      key = "video-stream1"
      video_stream {
        h264 {
          width_pixels      = 1280
          height_pixels     = 720
          bitrate_bps       = 550000
          frame_rate        = 60
          pixel_format      = "yuv420p"
          rate_control_mode = "vbr"
          crf_level         = 21
          gop_duration      = "3s"
          vbv_size_bits     = 2500000
          vbv_fullness_bits = 2250000
          entropy_coder     = "cabac"
          profile           = "high"
          preset            = "veryfast"

        }
      }
    }
    elementary_streams {
      key = "audio-stream0"
      audio_stream {
        codec             = "aac"
        bitrate_bps       = 64000
        channel_count     = 2
        channel_layout    = ["fl", "fr"]
        sample_rate_hertz = 48000
      }
    }
    mux_streams {
      key                = "sd"
      file_name          = "sd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream0", "audio-stream0"]
    }
    mux_streams {
      key                = "hd"
      file_name          = "hd.mp4"
      container          = "mp4"
      elementary_streams = ["video-stream1", "audio-stream0"]
    }
  }
  labels = {
    "label" = "key"
  }
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}`, context)
}
