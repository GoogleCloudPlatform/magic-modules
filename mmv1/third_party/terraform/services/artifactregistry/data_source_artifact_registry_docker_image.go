package artifactregistry

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceArtifactRegistryDockerImage() *schema.Resource {

	return &schema.Resource{
		Read: DataSourceArtifactRegistryDockerImageRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The repository name.`,
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The location of the artifact registry repository. For example, "us-west1".`,
			},
			"image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The image name to fetch.`,
			},
			"digest": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The image digest to fetch.  This cannot be used if tag is provided.`,
			},
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The tag of the version of the image to fetch. This cannot be used if digest is provided.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fully qualified name of the fetched image.`,
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI to access the image.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `All tags associated with the image.`,
			},
			"image_size_bytes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Calculated size of the image in bytes.`,
			},
			"media_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Media type of this image.`,
			},
			"upload_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time, as a RFC 3339 string, the image was uploaded.`,
			},
			"build_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time, as a RFC 3339 string, this image was built.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time, as a RFC 3339 string, this image was updated.`,
			},
		},
	}
}

// ArtifactRegistryBasePath
func DataSourceArtifactRegistryDockerImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	var res map[string]interface{}

	// check that only one of digest or tag is set
	if _, ok := d.GetOk("digest"); ok {
		if _, ok := d.GetOk("tag"); ok {
			return fmt.Errorf("only one of tag or digest can be set")
		} else {
			// fetch image by digest
			// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/get
			url, err := tpgresource.ReplaceVars(d, config, "{{ArtifactRegistryBasePath}}projects/{{project}}/locations/{{region}}/repositories/{{repository}}/dockerImages/{{image}}@{{digest}}")
			if err != nil {
				return fmt.Errorf("Error setting api endpoint")
			}

			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   project,
				RawURL:    url,
				UserAgent: userAgent,
			})
			if err != nil {
				return err
			}
		}
	} else {
		// fetch the list of images, ordered by update time
		// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/list
		url, err := tpgresource.ReplaceVars(d, config, "{{ArtifactRegistryBasePath}}projects/{{project}}/locations/{{region}}/repositories/{repository}}/dockerImages")
		if err != nil {
			return fmt.Errorf("Error setting api endpoint")
		}

		u, err := transport_tpg.AddQueryParams(url, map[string]string{"orderBy": "update_time desc"})
		if err != nil {
			return err
		}

		reslist, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "LIST",
			Project:   project,
			RawURL:    u,
			UserAgent: userAgent,
		})
		if err != nil {
			return err
		}

		// If tag is provided, iterate over response and find the image containing the tag
		if _, ok := d.GetOk("tag"); ok {
			tag := d.Get("tag")
		out:
			for _, image := range reslist["dockerImages"] {
				for _, tags := range image["tags"] {
					if tags.(string) == tag {
						res := image
						break out
					}
				}
			}
		} else {
			// use the first image in the response
			res := reslist["dockerImages"][0]
		}
	}

	// set the schema data using the response
	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("self_link", res["uri"]); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	if err := d.Set("tags", res["tags"]); err != nil {
		return fmt.Errorf("Error setting tags: %s", err)
	}

	if err := d.Set("image_size_bytes", res["imageSizeBytes"]); err != nil {
		return fmt.Errorf("Error setting image_size_bytes: %s", err)
	}

	if err := d.Set("media_type", res["mediaType"]); err != nil {
		return fmt.Errorf("Error setting media_type: %s", err)
	}

	if err := d.Set("upload_time", res["uploadTime"]); err != nil {
		return fmt.Errorf("Error setting upload_time: %s", err)
	}

	if err := d.Set("build_time", res["buildTime"]); err != nil {
		return fmt.Errorf("Error setting build_time: %s", err)
	}

	if err := d.Set("update_time", res["updateTime"]); err != nil {
		return fmt.Errorf("Error setting update_time: %s", err)
	}

	return nil
}
