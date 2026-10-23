package compute

import (
	"fmt"
	"strings"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func GetInterconnectAttachmentLink(config *transport_tpg.Config, project, region, ic, userAgent string) (string, error) {
	if !strings.Contains(ic, "/") {
		url := fmt.Sprintf("%sprojects/%s/regions/%s/interconnectAttachments/%s", transport_tpg.BaseUrl(Product, config), project, region, ic)
		icData, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return "", fmt.Errorf("Error reading interconnect attachment: %s", err)
		}
		selfLink, ok := icData["selfLink"].(string)
		if !ok {
			return "", fmt.Errorf("Error reading interconnect attachment: selfLink not found in response")
		}
		ic = selfLink
	}

	return ic, nil
}
