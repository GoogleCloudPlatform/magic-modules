package cloudidentity

import (
	"strings"
	"testing"
)

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_Create_MPA_DoneFalse(t *testing.T) {
	res := map[string]interface{}{
		"done": false,
		"metadata": map[string]interface{}{
			"state": "awaiting-multi-party-approval",
		},
	}
	_, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "creating")
	if err == nil {
		t.Fatalf("expected error for MPA awaiting approval on create, got nil")
	}
	if !strings.Contains(err.Error(), "awaiting-multi-party-approval") {
		t.Fatalf("expected error to contain awaiting-multi-party-approval, got: %v", err)
	}
}

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_Update_MPA_DoneFalse(t *testing.T) {
	res := map[string]interface{}{
		"done": false,
		"metadata": map[string]interface{}{
			"state": "awaiting-multi-party-approval",
		},
	}
	_, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "updating")
	if err == nil {
		t.Fatalf("expected error for MPA awaiting approval on update, got nil")
	}
	if !strings.Contains(err.Error(), "awaiting-multi-party-approval") {
		t.Fatalf("expected error to contain awaiting-multi-party-approval, got: %v", err)
	}
}

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_Create_DoneTrue(t *testing.T) {
	res := map[string]interface{}{
		"done": true,
		"response": map[string]interface{}{
			"name": "inboundOidcSsoProfiles/profile-123",
		},
	}
	respBody, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "creating")
	if err != nil {
		t.Fatalf("expected no error on create, got %v", err)
	}
	if respBody["name"] != "inboundOidcSsoProfiles/profile-123" {
		t.Fatalf("expected name inboundOidcSsoProfiles/profile-123, got %v", respBody["name"])
	}
}

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_Update_DoneTrue(t *testing.T) {
	res := map[string]interface{}{
		"done": true,
		"response": map[string]interface{}{
			"name": "inboundOidcSsoProfiles/profile-123",
		},
	}
	respBody, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "updating")
	if err != nil {
		t.Fatalf("expected no error on update, got %v", err)
	}
	if respBody["name"] != "inboundOidcSsoProfiles/profile-123" {
		t.Fatalf("expected name inboundOidcSsoProfiles/profile-123, got %v", respBody["name"])
	}
}

func TestCheckCloudIdentityInboundOidcSsoProfileOperation_Error(t *testing.T) {
	res := map[string]interface{}{
		"error": map[string]interface{}{
			"message": "permission denied",
		},
	}
	_, err := checkCloudIdentityInboundOidcSsoProfileOperation(res, "updating")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
