package google

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func resourceComputeFirewallRuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["protocol"].(string))))

	// We need to make sure to sort the strings below so that we always
	// generate the same hash code no matter what is in the set.
	if v, ok := m["ports"]; ok && v != nil {
		s := convertStringArr(v.([]interface{}))
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	return hashcode(buf.String())
}
func resourceDNSManagedZoneNetworkHash(v interface{}) int {
	if v == nil {
		return 0
	}
	raw := v.(map[string]interface{})
	if url, ok := raw["network_url"]; ok {
		return selfLinkNameHash(url)
	}
	var buf bytes.Buffer
	schema.SerializeResourceForHash(&buf, raw, DnsManagedZonePrivateVisibilityConfigNetworksSchema())
	return hashcode(buf.String())
}

func resourceSourceRepoRepositoryPubSubConfigsHash(v interface{}) int {
	if v == nil {
		return 0
	}

	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", tpgresource.GetResourceNameFromSelfLink(m["topic"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", m["message_format"].(string)))
	if v, ok := m["service_account_email"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode(buf.String())
}
