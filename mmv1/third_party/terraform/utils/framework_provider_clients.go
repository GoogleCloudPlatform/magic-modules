package google

import (
	"fmt"
	"strings"

	"google.golang.org/api/dns/v1"
	firebase "google.golang.org/api/firebase/v1beta1"
	"google.golang.org/api/option"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Methods to create new services from config
// Some base paths below need the version and possibly more of the path
// set on them. The client libraries are inconsistent about which values they need;
// while most only want the host URL, some older ones also want the version and some
// of those "projects" as well. You can find out if this is required by looking at
// the basePath value in the client library file.

func (p *frameworkProvider) NewDnsClient(userAgent string, diags *diag.Diagnostics) *dns.Service {
	dnsClientBasePath := RemoveBasePathVersion(p.DNSBasePath)
	dnsClientBasePath = strings.ReplaceAll(dnsClientBasePath, "/dns/", "")
	tflog.Info(p.context, fmt.Sprintf("Instantiating Google Cloud DNS client for path %s", dnsClientBasePath))
	clientDns, err := dns.NewService(p.context, option.WithHTTPClient(p.client))
	if err != nil {
		diags.AddWarning("error creating client dns", err.Error())
		return nil
	}
	clientDns.UserAgent = userAgent
	clientDns.BasePath = dnsClientBasePath

	return clientDns
}

func (p *frameworkProvider) NewFirebaseClient(userAgent string, diags *diag.Diagnostics) *firebase.Service {
	firebaseClientBasePath := RemoveBasePathVersion(p.FirebaseBasePath)
	firebaseClientBasePath = strings.ReplaceAll(firebaseClientBasePath, "/firebase/", "")
	tflog.Info(p.context, fmt.Sprintf("Instantiating Google Cloud firebase client for path %s", firebaseClientBasePath))
	clientFirebase, err := firebase.NewService(p.context, option.WithHTTPClient(p.client))
	if err != nil {
		diags.AddWarning("error creating client firebase", err.Error())
		return nil
	}
	clientFirebase.UserAgent = userAgent
	clientFirebase.BasePath = firebaseClientBasePath

	return clientFirebase
}
