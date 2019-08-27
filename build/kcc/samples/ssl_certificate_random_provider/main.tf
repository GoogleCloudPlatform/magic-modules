# You may also want to control name generation explicitly:
resource "google_compute_ssl_certificate" "default" {
  # The name will contain 8 random hex digits,
  # e.g. "my-certificate-48ab27cd2a"
  name        = "${random_id.certificate.hex}"
  private_key = "${file("path/to/private.key")}"
  certificate = "${file("path/to/certificate.crt")}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "random_id" "certificate" {
  byte_length = 4
  prefix      = "my-certificate-"

  # For security, do not expose raw certificate values in the output
  keepers = {
    private_key = "${base64sha256(file("path/to/private.key"))}"
    certificate = "${base64sha256(file("path/to/certificate.crt"))}"
  }
}
