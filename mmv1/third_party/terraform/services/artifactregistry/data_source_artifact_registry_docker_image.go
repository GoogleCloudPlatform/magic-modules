package artifactregistry

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages#DockerImage
type DockerImage struct {
	name           string
	uri            string
	tags           []string
	imageSizeBytes string
	mediaType      string
	uploadTime     string
	buildTime      string
	updateTime     string
}

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

	var res DockerImage

	// check that only one of digest or tag is set
	_, hasDigest := d.GetOk("digest")
	tag, hasTag := d.GetOk("tag")

	if !hasDigest && !hasTag {
		return fmt.Errorf("either tag or digest must be provided")
	}

	if hasDigest && hasTag {
		return fmt.Errorf("only one of tag or digest can be set")
	}

	if hasDigest {
		// fetch image by digest
		// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/get
		imageUrlSafe := url.QueryEscape(d.Get("image").(string))
		urlRequest, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{ArtifactRegistryBasePath}}projects/{{project}}/locations/{{region}}/repositories/{{repository}}/dockerImages/%s@{{digest}}", imageUrlSafe))
		if err != nil {
			return fmt.Errorf("Error setting api endpoint")
		}

		resGet, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    urlRequest,
			UserAgent: userAgent,
		})
		if err != nil {
			return err
		}

		res = convertResponseToStruct(resGet)
	} else {
		// fetch the list of images, ordered by update time
		// https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/list
		urlRequest, err := tpgresource.ReplaceVars(d, config, "{{ArtifactRegistryBasePath}}projects/{{project}}/locations/{{region}}/repositories/{{repository}}/dockerImages")
		if err != nil {
			return fmt.Errorf("Error setting api endpoint")
		}

		urlRequest, err = transport_tpg.AddQueryParams(urlRequest, map[string]string{"orderBy": "update_time desc"})
		if err != nil {
			return err
		}

		resListImages, token, err := retreiveListOfDockerImages(config, project, urlRequest, userAgent)
		if err != nil {
			return err
		}
	out:
		for {
			var resFiltered []DockerImage

			// The response gives a list of all images, filter by name
			for _, image := range resListImages {

				imageName, ok := d.Get("image").(string)
				if !ok {
					return fmt.Errorf("Error: Image name is not a string")
				}

				if strings.Contains(image.name, "/"+url.QueryEscape(imageName)+"@") {
					resFiltered = append(resFiltered, image)
				}
			}

			// If tag is provided, iterate over response and find the image containing the tag
			if hasTag {
				for _, image := range resFiltered {
					for _, iterTag := range image.tags {
						if iterTag == tag {
							res = image
							break out
						}
					}
				}
			} else { // use the first image in the response
				if len(resFiltered) > 0 {
					res = resFiltered[0]
					break out
				}
			}

			// No matching image was found on this page and no more pages are available
			if token == "" {
				return fmt.Errorf("Requested image was not found.")
			}

			// Retrieve the next page and search again
			urlNext, err := transport_tpg.AddQueryParams(urlRequest, map[string]string{"pageToken": token})
			if err != nil {
				return err
			}

			resListImages, token, err = retreiveListOfDockerImages(config, project, urlNext, userAgent)
			if err != nil {
				return err
			}
		}
	}

	// set the schema data using the response
	if err := d.Set("name", res.name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("self_link", res.uri); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	if err := d.Set("tags", res.tags); err != nil {
		return fmt.Errorf("Error setting tags: %s", err)
	}

	if err := d.Set("image_size_bytes", res.imageSizeBytes); err != nil {
		return fmt.Errorf("Error setting image_size_bytes: %s", err)
	}

	if err := d.Set("media_type", res.mediaType); err != nil {
		return fmt.Errorf("Error setting media_type: %s", err)
	}

	if err := d.Set("upload_time", res.uploadTime); err != nil {
		return fmt.Errorf("Error setting upload_time: %s", err)
	}

	if err := d.Set("build_time", res.buildTime); err != nil {
		return fmt.Errorf("Error setting build_time: %s", err)
	}

	if err := d.Set("update_time", res.updateTime); err != nil {
		return fmt.Errorf("Error setting update_time: %s", err)
	}

	id := fmt.Sprintf("projects/%s/locations/%s/repositories/dockerImages/%s", project, d.Get("region"), d.Get("image"))
	d.SetId(id)

	return nil
}

func retreiveListOfDockerImages(config *transport_tpg.Config, project string, urlRequest string, userAgent string) ([]DockerImage, string, error) {
	resList, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    urlRequest,
		UserAgent: userAgent,
	})
	if err != nil {
		return make([]DockerImage, 0), "", err
	}

	if nextPageToken, ok := resList["nextPageToken"].(string); ok {
		return flattenDatasourceListResponse(resList), nextPageToken, nil
	} else {
		return flattenDatasourceListResponse(resList), "", nil
	}
}

func flattenDatasourceListResponse(res map[string]interface{}) []DockerImage {
	var dockerImages []DockerImage

	resDockerImages, _ := res["dockerImages"].([]interface{})

	for _, resImage := range resDockerImages {
		image, _ := resImage.(map[string]interface{})
		dockerImages = append(dockerImages, convertResponseToStruct(image))
	}

	return dockerImages
}

func convertResponseToStruct(res map[string]interface{}) DockerImage {
	var dockerImage DockerImage

	if name, ok := res["name"].(string); ok {
		dockerImage.name = name
	}

	if uri, ok := res["uri"].(string); ok {
		dockerImage.uri = uri
	}

	if tags, ok := res["tags"].([]interface{}); ok {
		var stringTags []string

		for _, tag := range tags {
			strTag := tag.(string)
			stringTags = append(stringTags, strTag)
		}
		dockerImage.tags = stringTags
	}

	if imageSizeBytes, ok := res["imageSizeBytes"].(string); ok {
		dockerImage.imageSizeBytes = imageSizeBytes
	}

	if mediaType, ok := res["mediaType"].(string); ok {
		dockerImage.mediaType = mediaType
	}

	if uploadTime, ok := res["uploadTime"].(string); ok {
		dockerImage.uploadTime = uploadTime
	}

	if buildTime, ok := res["buildTime"].(string); ok {
		dockerImage.buildTime = buildTime
	}

	if updateTime, ok := res["updateTime"].(string); ok {
		dockerImage.updateTime = updateTime
	}

	return dockerImage
}
