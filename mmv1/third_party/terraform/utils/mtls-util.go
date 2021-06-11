package google

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/api/option/internaloption"
	"google.golang.org/api/transport"
)

func isMtls() bool {
	regularEndpoint := "https://mockservice.googleapis.com/v1/"
	mtlsEndpoint := getMtlsEndpoint(regularEndpoint)
	_, endpoint, err := transport.NewHTTPClient(context.Background(),
		internaloption.WithDefaultEndpoint(regularEndpoint),
		internaloption.WithDefaultMTLSEndpoint(mtlsEndpoint),
	)
	if err != nil {
		return false
	}
	isMtls := strings.Contains(endpoint, "mtls")
	return isMtls
}

func getMtlsEndpoint(baseEndpoint string) string {
	u, err := url.Parse(baseEndpoint)
	if err != nil {
		return baseEndpoint
	}
	portParts := strings.Split(u.Host, ":")
	if len(portParts) == 0 || portParts[0] == "" {
		return baseEndpoint
	}
	domainParts := strings.Split(portParts[0], ".")
	if len(domainParts) > 1 {
		u.Host = fmt.Sprintf("%s.mtls.%s", domainParts[0], strings.Join(domainParts[1:], "."))
	} else {
		u.Host = fmt.Sprintf("%s.mtls", domainParts[0])
	}
	if len(portParts) > 1 {
		u.Host = fmt.Sprintf("%s:%s", u.Host, portParts[1])
	}
	return u.String()
}
