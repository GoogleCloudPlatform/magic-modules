// Contains common diff suppress functions.

package google

import (
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func EmptyOrDefaultStringSuppress(defaultVal string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return (old == "" && new == defaultVal) || (new == "" && old == defaultVal)
	}
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
