package compute

import (
	"fmt"
	"strings"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func GetInterconnectAttachmentLink(config *transport_tpg.Config, project, region, ic, userAgent string) (string, error) {
	if !strings.Contains(ic, "/") {
		icData, err := NewClient(config, userAgent).InterconnectAttachments.Get(
			project, region, ic).Do()
		if err != nil {
			return "", fmt.Errorf("Error reading interconnect attachment: %s", err)
		}
		ic = icData.SelfLink
	}

	return ic, nil
}
