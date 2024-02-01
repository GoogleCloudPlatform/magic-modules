// Contains common diff suppress functions.

package google

import (
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func OptionalPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return prefix+old == new || prefix+new == old
	}
}

func OptionalSurroundingSpacesSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func EmptyOrDefaultStringSuppress(defaultVal string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return (old == "" && new == defaultVal) || (new == "" && old == defaultVal)
	}
}

func IpCidrRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// The range may be a:
	// A) single IP address (e.g. 10.2.3.4)
	// B) CIDR format string (e.g. 10.1.2.0/24)
	// C) netmask (e.g. /24)
	//
	// For A) and B), no diff to suppress, they have to match completely.
	// For C), The API picks a network IP address and this creates a diff of the form:
	// network_interface.0.alias_ip_range.0.ip_cidr_range: "10.128.1.0/24" => "/24"
	// We should only compare the mask portion for this case.
	if len(new) > 0 && new[0] == '/' {
		oldNetmaskStartPos := strings.LastIndex(old, "/")

		if oldNetmaskStartPos != -1 {
			oldNetmask := old[strings.LastIndex(old, "/"):]
			if oldNetmask == new {
				return true
			}
		}
	}

	return false
}

func CaseDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToUpper(old) == strings.ToUpper(new)
}

func TimestampDiffSuppress(format string) schema.SchemaDiffSuppressFunc {
	return func(_, old, new string, _ *schema.ResourceData) bool {
		oldT, err := time.Parse(format, old)
		if err != nil {
			return false
		}

		newT, err := time.Parse(format, new)
		if err != nil {
			return false
		}

		return oldT == newT
	}
}

func comparePubsubSubscriptionExpirationPolicy(_, old, new string, _ *schema.ResourceData) bool {
	trimmedNew := strings.TrimLeft(new, "0")
	trimmedOld := strings.TrimLeft(old, "0")
	if strings.Contains(trimmedNew, ".") {
		trimmedNew = strings.TrimRight(strings.TrimSuffix(trimmedNew, "s"), "0") + "s"
	}
	if strings.Contains(trimmedOld, ".") {
		trimmedOld = strings.TrimRight(strings.TrimSuffix(trimmedOld, "s"), "0") + "s"
	}
	return trimmedNew == trimmedOld
}

func rrefDiffSuppressfunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("type") == "AAAA" {
		return ipv6AddressDiffSuppress(k, old, new, d)
	}
	return false
}

// This is separate from rrefDiffSuppressfunc for unit testing
func ipv6AddressDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldIp := net.ParseIP(old)
	newIp := net.ParseIP(new)

	return oldIp.Equal(newIp)
}

// This is separate from CaseDiffSuppress as it strips quotation marks
func dnsRecordSetRrefsDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToLower(strings.Trim(old, `"`)) == strings.ToLower(strings.Trim(new, `"`))
}
