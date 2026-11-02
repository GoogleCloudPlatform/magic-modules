// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
// ----------------------------------------------------------------------------
//
//	***     HANDWRITTEN CODE     ***
//
// ----------------------------------------------------------------------------
package sql

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func TestEncodeCustomerManagedCertificateContent(t *testing.T) {
	rawPem := "-----BEGIN CERTIFICATE-----\nSERVER_CA_CONTENT\n-----END CERTIFICATE-----"
	expectedBase64 := base64.StdEncoding.EncodeToString([]byte(rawPem))

	if got := encodeCustomerManagedCertificateContent(rawPem); got != expectedBase64 {
		t.Errorf("encodeCustomerManagedCertificateContent(rawPem) = %q, want %q", got, expectedBase64)
	}

	if got := encodeCustomerManagedCertificateContent(expectedBase64); got != expectedBase64 {
		t.Errorf("encodeCustomerManagedCertificateContent(expectedBase64) = %q, want %q", got, expectedBase64)
	}
}

func TestBuildCustomerManagedCertificatesPayload(t *testing.T) {
	rawCa := "-----BEGIN CERTIFICATE-----\nCA_DATA\n-----END CERTIFICATE-----"
	rawCert := "-----BEGIN CERTIFICATE-----\nCERT_DATA\n-----END CERTIFICATE-----"
	rawKey := "-----BEGIN PRIVATE KEY-----\nKEY_DATA\n-----END PRIVATE KEY-----"

	got := buildCustomerManagedCertificatesPayload(rawCa, rawCert, rawKey)
	want := []map[string]interface{}{
		{
			"certificateType": "CERTIFICATE_TYPE_PEER_SERVER_CA",
			"content":         base64.StdEncoding.EncodeToString([]byte(rawCa)),
		},
		{
			"certificateType": "CERTIFICATE_TYPE_PEER_CLIENT_CERTIFICATE",
			"content":         base64.StdEncoding.EncodeToString([]byte(rawCert)),
		},
		{
			"certificateType": "CERTIFICATE_TYPE_PEER_CLIENT_KEY",
			"content":         base64.StdEncoding.EncodeToString([]byte(rawKey)),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildCustomerManagedCertificatesPayload = %#v, want %#v", got, want)
	}
}

func TestExtractCustomerManagedCertFilePaths(t *testing.T) {
	apiList := []interface{}{
		map[string]interface{}{
			"certificateType": "CERTIFICATE_TYPE_PEER_SERVER_CA",
			"filePath":        "certs/bundle1/server_ca.pem",
		},
		map[string]interface{}{
			"certificateType": "CERTIFICATE_TYPE_PEER_CLIENT_KEY",
			"filePath":        "certs/bundle1/client_key.pem",
		},
		map[string]interface{}{
			"certificateType": "CERTIFICATE_TYPE_PEER_CLIENT_CERTIFICATE",
			"filePath":        "certs/bundle1/client_cert.pem",
		},
	}

	filePaths, serverCaPath, clientCertPath, clientKeyPath := extractCustomerManagedCertFilePaths(apiList)
	wantPaths := []string{
		"certs/bundle1/client_cert.pem",
		"certs/bundle1/client_key.pem",
		"certs/bundle1/server_ca.pem",
	}

	if !reflect.DeepEqual(filePaths, wantPaths) {
		t.Errorf("filePaths = %#v, want %#v", filePaths, wantPaths)
	}
	if serverCaPath != "certs/bundle1/server_ca.pem" {
		t.Errorf("serverCaPath = %q, want %q", serverCaPath, "certs/bundle1/server_ca.pem")
	}
	if clientCertPath != "certs/bundle1/client_cert.pem" {
		t.Errorf("clientCertPath = %q, want %q", clientCertPath, "certs/bundle1/client_cert.pem")
	}
	if clientKeyPath != "certs/bundle1/client_key.pem" {
		t.Errorf("clientKeyPath = %q, want %q", clientKeyPath, "certs/bundle1/client_key.pem")
	}
}

func TestFindSqlCustomerManagedCertificateBundle(t *testing.T) {
	res := map[string]interface{}{
		"customerManagedCertificates": []interface{}{
			map[string]interface{}{
				"bundle": "bundle1",
				"intent": "PG_LOGICAL",
			},
			map[string]interface{}{
				"bundle": "bundle2",
				"intent": "NATIVE_REPLICATION",
			},
		},
	}

	got := findSqlCustomerManagedCertificateBundle(res, "bundle2")
	if got == nil || got["intent"] != "NATIVE_REPLICATION" {
		t.Errorf("findSqlCustomerManagedCertificateBundle(bundle2) = %#v, want bundle2", got)
	}

	if missing := findSqlCustomerManagedCertificateBundle(res, "bundle3"); missing != nil {
		t.Errorf("findSqlCustomerManagedCertificateBundle(bundle3) = %#v, want nil", missing)
	}
}

func TestFlattenDataSourceSqlCustomerManagedCertificates(t *testing.T) {
	res := map[string]interface{}{
		"customerManagedCertificates": []interface{}{
			map[string]interface{}{
				"bundle": "bundle2",
				"intent": "PG_LOGICAL",
				"primaryDnsRecord": map[string]interface{}{
					"dns":         "peer2.example.com",
					"ipv4Address": "10.0.0.2",
				},
				"customerManagedCertificates": []interface{}{
					map[string]interface{}{
						"certificateType": "CERTIFICATE_TYPE_PEER_SERVER_CA",
						"filePath":        "certs/bundle2/server_ca.pem",
					},
				},
			},
			map[string]interface{}{
				"bundle": "bundle1",
				"intent": "PG_LOGICAL",
				"primaryDnsRecord": map[string]interface{}{
					"dns":         "peer1.example.com",
					"ipv4Address": "10.0.0.1",
				},
				"customerManagedCertificates": []interface{}{
					map[string]interface{}{
						"certificateType": "CERTIFICATE_TYPE_PEER_CLIENT_CERTIFICATE",
						"filePath":        "certs/bundle1/client_cert.pem",
					},
					map[string]interface{}{
						"certificateType": "CERTIFICATE_TYPE_PEER_SERVER_CA",
						"filePath":        "certs/bundle1/server_ca.pem",
					},
				},
			},
		},
	}

	got := flattenDataSourceSqlCustomerManagedCertificates(res)
	if len(got) != 2 {
		t.Fatalf("len(flattenDataSourceSqlCustomerManagedCertificates) = %d, want 2", len(got))
	}
	if got[0]["bundle"] != "bundle1" || got[1]["bundle"] != "bundle2" {
		t.Errorf("expected bundles sorted by bundle name [bundle1, bundle2], got [%v, %v]", got[0]["bundle"], got[1]["bundle"])
	}
	if got[0]["primary_dns"] != "peer1.example.com" || got[0]["primary_ipv4_address"] != "10.0.0.1" {
		t.Errorf("unexpected primary_dns / primary_ipv4_address: %#v", got[0])
	}
	if got[0]["server_ca_cert_file_path"] != "certs/bundle1/server_ca.pem" {
		t.Errorf("unexpected server_ca_cert_file_path: %v", got[0]["server_ca_cert_file_path"])
	}
}
