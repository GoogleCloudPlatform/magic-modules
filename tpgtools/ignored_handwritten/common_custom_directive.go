package google

import (
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	dns "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dns"
)

func rrefCreateDirective(recordSet *dns.ResourceRecordSet) []dcl.ApplyOption {
	if recordSet.DnsType != nil && *recordSet.DnsType == "NS" {
		// NS type records may exist by default. In this case, we want to acquire
		// and modify existing records
		return []dcl.ApplyOption{
			dcl.WithLifecycleParam(dcl.BlockDestruction),
		}
	}
	return CreateDirective
}
