package google

import (
	"fmt"
	"net/http"
	"os"
)

// adapted from https://stackoverflow.com/questions/51325704/adding-a-default-http-header-in-go
type headerTransportLayer struct {
	http.Header
	baseTransit http.RoundTripper
}

func newTransportWithHeaders(baseTransit http.RoundTripper) headerTransportLayer {
	if baseTransit == nil {
		baseTransit = http.DefaultTransport
	}

	headers := make(http.Header)
	if requestReason := os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"); requestReason != "" {
		headers.Set("X-Goog-Request-Reason", requestReason)
	}
	// opt in extension for adding to the User-Agent header
	if ext := os.Getenv("GOOGLE_TPG_USERAGENT_EXTENSION"); ext != "" {
		headers.Set("User-Agent", ext)
	}

	return headerTransportLayer{Header: headers, baseTransit: baseTransit}
}

func (h headerTransportLayer) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Header {
		// only set headers that are not previously defined
		if _, ok := req.Header[key]; !ok {
			req.Header[key] = value
		} else if key == "User-Agent" {
			ua := req.Header.Get(key)
			ext := h.Header.Get(key)
			newUa := fmt.Sprintf("%s %s", ua, ext)
			req.Header.Set(key, newUa)
		}
	}
	return h.baseTransit.RoundTrip(req)
}
