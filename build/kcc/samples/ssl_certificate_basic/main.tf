resource "google_compute_ssl_certificate" "default" {
  name_prefix = "my-certificate-"
  description = "a description"
  private_key = "${file("../static/ssl_cert/test.key")}"
  certificate = "${file("../static/ssl_cert/test.crt")}"

  lifecycle {
    create_before_destroy = true
  }
}
