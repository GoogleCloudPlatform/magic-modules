resource "google_certificate_manager_certificate" "self-managed" {
  location = "global"
  name     = "self-managed"
  project  = "my-project"

  self_managed {
    pem_certificate = "-----BEGIN CERTIFICATE-----\nMIIFpzCCA4+gAwIBAgIUGrkv7D1l+G3QAQUT9f2jhTaVZ/gwDQYJKoZIhvcNAQEL\nBQAwYzELMAkGA1UEBhMCUEwxEDAOBgNVBAgMB01hc292aWExDzANBgNVBAcMBldh\ncnNhdzENMAsGA1UECgwEQUNNRTEMMAoGA1UECwwDQ0NNMRQwEgYDVQQDDAtleGFt\ncGxlLm9yZzAeFw0yNTA4MTExMjM5NTVaFw0zNTA4MDkxMjM5NTVaMGMxCzAJBgNV\nBAYTAlBMMRAwDgYDVQQIDAdNYXNvdmlhMQ8wDQYDVQQHDAZXYXJzYXcxDTALBgNV\nBAoMBEFDTUUxDDAKBgNVBAsMA0NDTTEUMBIGA1UEAwwLZXhhbXBsZS5vcmcwggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC60FhXjbBe/VtQPCxl8Vg/N7HT\nU8okjecIysf5J9OzZ3qU3RWht3ChxzoyAtztAKBPqjDwQrVIh/CB/gPaSJVsRm7u\nKNuIgYfiyI+tOFX5P9NBPqDdMM6gPizm9/zrEerbt3SttoYbxIo+39u6MQ9jH3t+\nJ0hNLthuWjkxigiPoWrdLIUKDkuR2Sbuvr5FdG+PXbeTkthLurobAJGrlzMhrLCD\nzCkDFLBfTcziWIruQFVVUaP9ubNNmNhVDlG9ey3B3YcQP7/Fi+lQvyNoDhA5fwGg\nUwUNkZpqd5YEaGTBrq3CNL+m9hU/MDq7IzdaUeGCevPR3xODfxjb9YPAmUt4liWU\nblleQhmfZreLKgWYffMPKsoMqCKs2QZqxjcifdRj6/T9YDvDsgigY5NtpvfZFkZC\nPH5X95jEg2YoEqtn/d4MrXnhdVNIct4I/ggFKOEhAx+a3NSGoadZEdMnhDXx2Y8O\nk5RNbKxvV8Lgd+v562i4TRnSeso6Ak+iir9JOyLDJK0VO8qC2I8oyJ9cRjTSIOzh\nYRTtlGNFnUYj1v2T1wYyf2HMMxnDx9faehOCjwaianPV4cOihokavQ2Wwp0A5hwE\njzYGrr7XskI8K2IT5RPOX36TbBIw0VtwYkCLJ6Jiv0e2lM71VwCD2LjbRMteV1Qm\nLISH+XV51tNHfh9TowIDAQABo1MwUTAdBgNVHQ4EFgQU7OL7lhMg5FquKzKmksqS\nDGIPR5AwHwYDVR0jBBgwFoAU7OL7lhMg5FquKzKmksqSDGIPR5AwDwYDVR0TAQH/\nBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAie6dX149q5VHGuNlgLNeLoBpeTMg\nFRLr8P7R63yzywlWg9qPYc2Bai6WI6Qt47B/7ta7lzvtAe8J9XSmmuldL4rSoQyl\nUQ5nBuJRtYJaThIrwHPFhyGi3up87S58VdOxOO/KjewooYl3lkJ3i9Mkk0lGVwa6\nhziw218yBoQpuUS+r2HPPgw4E3f6C1UVL4rGrhG6/UYjQ0b0tVry7pQPVbQPWeAS\nnqIpJbZhM+sVz7JlCgPSN+fB9qHyN5Rce4uDO4Hrsw2dFLTkyilbrcVfxLHyTAbY\nbSiJtYuXJYPXT26YWSAeJQWRrACLvLn6OCVQ5nRwZy1dqN0VjZi/YRWXaQeu7+7h\n8dih90jgpWufvf74pNouuhlGxZ3xKUQeYdKHn9b87lugXZLsN+34DcPjX3+iPIXt\nBeUqkwdob1aCawUp+f3h1r/TjCr8IYcb1uL0OfCFkk9/fgxTPQaG5QyQIrgEdWgA\nd2dDyGGDJy/cDjFYThU0oOgFcdVKrdS2Lw5PVnxJSPeji3s5xqW7lGWmyX8/qujZ\nwFeF+TxKVh1kfVnSFg09+1sgbpxp1LjWtFupwAJDiLSqAY0+dFPEMPcH122zTOpj\ndVfu0SFnpSgX/sE/tIiDIKmmTOcAJBb0pwGUjCWil7DZYgukY81F18f21wlsJFIQ\npNyP6e+zpDTjPNo=\n-----END CERTIFICATE-----\n"
    pem_private_key = "<private_key>"
  }
}

resource "google_certificate_manager_certificate" "google-managed-1" {
  location = "global"

  managed {
    domains = ["example.org"]
  }

  name    = "google-managed-1"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-2" {
  location = "global"

  managed {
    dns_authorizations = ["projects/307841421122/locations/global/dnsAuthorizations/dns-authz-example-org"]
    domains            = ["example.org"]
  }

  name    = "google-managed-2"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-3" {
  location = "global"

  managed {
    domains         = ["foo.com"]
    issuance_config = "projects/307841421122/locations/global/certificateIssuanceConfigs/deletion-test4"
  }

  name    = "google-managed-3"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-w-description" {
  description = "My description"
  location    = "global"

  managed {
    domains = ["example.org"]
  }

  name    = "google-managed-w-description"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-w-labels" {
  labels = {
    foo = "bar"
  }

  location = "global"

  managed {
    domains = ["example.org"]
  }

  name    = "google-managed-w-labels"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-w-project" {
  location = "global"

  managed {
    dns_authorizations = ["projects/307841421122/locations/global/dnsAuthorizations/dns-authz-example-org"]
    domains            = ["example.org"]
  }

  name    = "google-managed-w-project"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-regional" {
  location = "europe-west1"

  managed {
    dns_authorizations = ["projects/307841421122/locations/europe-west1/dnsAuthorizations/dns-authz-example-org"]
    domains            = ["example.org"]
  }

  name    = "google-managed-regional"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-global" {
  location = "global"

  managed {
    dns_authorizations = ["projects/307841421122/locations/global/dnsAuthorizations/dns-authz-example-org"]
    domains            = ["example.org"]
  }

  name    = "google-managed-global"
  project = "my-project"
}

resource "google_certificate_manager_certificate" "google-managed-edge-cache" {
  location = "global"

  managed {
    dns_authorizations = ["projects/307841421122/locations/global/dnsAuthorizations/dns-authz-example-org", "projects/307841421122/locations/global/dnsAuthorizations/www-example-com"]
    domains            = ["example.org", "www.example.com"]
  }

  name    = "google-managed-edge-cache"
  project = "my-project"
  scope   = "EDGE_CACHE"
}

resource "google_certificate_manager_certificate" "self-managed-all-regions" {
  location = "global"
  name     = "self-managed-all-regions"
  project  = "my-project"
  scope    = "ALL_REGIONS"

  self_managed {
    pem_certificate = "-----BEGIN CERTIFICATE-----\nMIIFpzCCA4+gAwIBAgIUGrkv7D1l+G3QAQUT9f2jhTaVZ/gwDQYJKoZIhvcNAQEL\nBQAwYzELMAkGA1UEBhMCUEwxEDAOBgNVBAgMB01hc292aWExDzANBgNVBAcMBldh\ncnNhdzENMAsGA1UECgwEQUNNRTEMMAoGA1UECwwDQ0NNMRQwEgYDVQQDDAtleGFt\ncGxlLm9yZzAeFw0yNTA4MTExMjM5NTVaFw0zNTA4MDkxMjM5NTVaMGMxCzAJBgNV\nBAYTAlBMMRAwDgYDVQQIDAdNYXNvdmlhMQ8wDQYDVQQHDAZXYXJzYXcxDTALBgNV\nBAoMBEFDTUUxDDAKBgNVBAsMA0NDTTEUMBIGA1UEAwwLZXhhbXBsZS5vcmcwggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC60FhXjbBe/VtQPCxl8Vg/N7HT\nU8okjecIysf5J9OzZ3qU3RWht3ChxzoyAtztAKBPqjDwQrVIh/CB/gPaSJVsRm7u\nKNuIgYfiyI+tOFX5P9NBPqDdMM6gPizm9/zrEerbt3SttoYbxIo+39u6MQ9jH3t+\nJ0hNLthuWjkxigiPoWrdLIUKDkuR2Sbuvr5FdG+PXbeTkthLurobAJGrlzMhrLCD\nzCkDFLBfTcziWIruQFVVUaP9ubNNmNhVDlG9ey3B3YcQP7/Fi+lQvyNoDhA5fwGg\nUwUNkZpqd5YEaGTBrq3CNL+m9hU/MDq7IzdaUeGCevPR3xODfxjb9YPAmUt4liWU\nblleQhmfZreLKgWYffMPKsoMqCKs2QZqxjcifdRj6/T9YDvDsgigY5NtpvfZFkZC\nPH5X95jEg2YoEqtn/d4MrXnhdVNIct4I/ggFKOEhAx+a3NSGoadZEdMnhDXx2Y8O\nk5RNbKxvV8Lgd+v562i4TRnSeso6Ak+iir9JOyLDJK0VO8qC2I8oyJ9cRjTSIOzh\nYRTtlGNFnUYj1v2T1wYyf2HMMxnDx9faehOCjwaianPV4cOihokavQ2Wwp0A5hwE\njzYGrr7XskI8K2IT5RPOX36TbBIw0VtwYkCLJ6Jiv0e2lM71VwCD2LjbRMteV1Qm\nLISH+XV51tNHfh9TowIDAQABo1MwUTAdBgNVHQ4EFgQU7OL7lhMg5FquKzKmksqS\nDGIPR5AwHwYDVR0jBBgwFoAU7OL7lhMg5FquKzKmksqSDGIPR5AwDwYDVR0TAQH/\nBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAie6dX149q5VHGuNlgLNeLoBpeTMg\nFRLr8P7R63yzywlWg9qPYc2Bai6WI6Qt47B/7ta7lzvtAe8J9XSmmuldL4rSoQyl\nUQ5nBuJRtYJaThIrwHPFhyGi3up87S58VdOxOO/KjewooYl3lkJ3i9Mkk0lGVwa6\nhziw218yBoQpuUS+r2HPPgw4E3f6C1UVL4rGrhG6/UYjQ0b0tVry7pQPVbQPWeAS\nnqIpJbZhM+sVz7JlCgPSN+fB9qHyN5Rce4uDO4Hrsw2dFLTkyilbrcVfxLHyTAbY\nbSiJtYuXJYPXT26YWSAeJQWRrACLvLn6OCVQ5nRwZy1dqN0VjZi/YRWXaQeu7+7h\n8dih90jgpWufvf74pNouuhlGxZ3xKUQeYdKHn9b87lugXZLsN+34DcPjX3+iPIXt\nBeUqkwdob1aCawUp+f3h1r/TjCr8IYcb1uL0OfCFkk9/fgxTPQaG5QyQIrgEdWgA\nd2dDyGGDJy/cDjFYThU0oOgFcdVKrdS2Lw5PVnxJSPeji3s5xqW7lGWmyX8/qujZ\nwFeF+TxKVh1kfVnSFg09+1sgbpxp1LjWtFupwAJDiLSqAY0+dFPEMPcH122zTOpj\ndVfu0SFnpSgX/sE/tIiDIKmmTOcAJBb0pwGUjCWil7DZYgukY81F18f21wlsJFIQ\npNyP6e+zpDTjPNo=\n-----END CERTIFICATE-----\n"
    pem_private_key = "<private_key>"
  }
}

resource "google_certificate_manager_certificate" "self-managed-client-auth" {
  location = "global"
  name     = "self-managed-client-auth"
  project  = "my-project"
  scope    = "CLIENT_AUTH"

  self_managed {
    pem_certificate = "-----BEGIN CERTIFICATE-----\nMIIFpzCCA4+gAwIBAgIUGrkv7D1l+G3QAQUT9f2jhTaVZ/gwDQYJKoZIhvcNAQEL\nBQAwYzELMAkGA1UEBhMCUEwxEDAOBgNVBAgMB01hc292aWExDzANBgNVBAcMBldh\ncnNhdzENMAsGA1UECgwEQUNNRTEMMAoGA1UECwwDQ0NNMRQwEgYDVQQDDAtleGFt\ncGxlLm9yZzAeFw0yNTA4MTExMjM5NTVaFw0zNTA4MDkxMjM5NTVaMGMxCzAJBgNV\nBAYTAlBMMRAwDgYDVQQIDAdNYXNvdmlhMQ8wDQYDVQQHDAZXYXJzYXcxDTALBgNV\nBAoMBEFDTUUxDDAKBgNVBAsMA0NDTTEUMBIGA1UEAwwLZXhhbXBsZS5vcmcwggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC60FhXjbBe/VtQPCxl8Vg/N7HT\nU8okjecIysf5J9OzZ3qU3RWht3ChxzoyAtztAKBPqjDwQrVIh/CB/gPaSJVsRm7u\nKNuIgYfiyI+tOFX5P9NBPqDdMM6gPizm9/zrEerbt3SttoYbxIo+39u6MQ9jH3t+\nJ0hNLthuWjkxigiPoWrdLIUKDkuR2Sbuvr5FdG+PXbeTkthLurobAJGrlzMhrLCD\nzCkDFLBfTcziWIruQFVVUaP9ubNNmNhVDlG9ey3B3YcQP7/Fi+lQvyNoDhA5fwGg\nUwUNkZpqd5YEaGTBrq3CNL+m9hU/MDq7IzdaUeGCevPR3xODfxjb9YPAmUt4liWU\nblleQhmfZreLKgWYffMPKsoMqCKs2QZqxjcifdRj6/T9YDvDsgigY5NtpvfZFkZC\nPH5X95jEg2YoEqtn/d4MrXnhdVNIct4I/ggFKOEhAx+a3NSGoadZEdMnhDXx2Y8O\nk5RNbKxvV8Lgd+v562i4TRnSeso6Ak+iir9JOyLDJK0VO8qC2I8oyJ9cRjTSIOzh\nYRTtlGNFnUYj1v2T1wYyf2HMMxnDx9faehOCjwaianPV4cOihokavQ2Wwp0A5hwE\njzYGrr7XskI8K2IT5RPOX36TbBIw0VtwYkCLJ6Jiv0e2lM71VwCD2LjbRMteV1Qm\nLISH+XV51tNHfh9TowIDAQABo1MwUTAdBgNVHQ4EFgQU7OL7lhMg5FquKzKmksqS\nDGIPR5AwHwYDVR0jBBgwFoAU7OL7lhMg5FquKzKmksqSDGIPR5AwDwYDVR0TAQH/\nBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAie6dX149q5VHGuNlgLNeLoBpeTMg\nFRLr8P7R63yzywlWg9qPYc2Bai6WI6Qt47B/7ta7lzvtAe8J9XSmmuldL4rSoQyl\nUQ5nBuJRtYJaThIrwHPFhyGi3up87S58VdOxOO/KjewooYl3lkJ3i9Mkk0lGVwa6\nhziw218yBoQpuUS+r2HPPgw4E3f6C1UVL4rGrhG6/UYjQ0b0tVry7pQPVbQPWeAS\nnqIpJbZhM+sVz7JlCgPSN+fB9qHyN5Rce4uDO4Hrsw2dFLTkyilbrcVfxLHyTAbY\nbSiJtYuXJYPXT26YWSAeJQWRrACLvLn6OCVQ5nRwZy1dqN0VjZi/YRWXaQeu7+7h\n8dih90jgpWufvf74pNouuhlGxZ3xKUQeYdKHn9b87lugXZLsN+34DcPjX3+iPIXt\nBeUqkwdob1aCawUp+f3h1r/TjCr8IYcb1uL0OfCFkk9/fgxTPQaG5QyQIrgEdWgA\nd2dDyGGDJy/cDjFYThU0oOgFcdVKrdS2Lw5PVnxJSPeji3s5xqW7lGWmyX8/qujZ\nwFeF+TxKVh1kfVnSFg09+1sgbpxp1LjWtFupwAJDiLSqAY0+dFPEMPcH122zTOpj\ndVfu0SFnpSgX/sE/tIiDIKmmTOcAJBb0pwGUjCWil7DZYgukY81F18f21wlsJFIQ\npNyP6e+zpDTjPNo=\n-----END CERTIFICATE-----\n"
    pem_private_key = "<private_key>"
  }
}

resource "google_certificate_manager_certificate" "google-managed-default" {
  location = "global"

  managed {
    domains = ["example.org"]
  }

  name    = "google-managed-default"
  project = "my-project"
}
