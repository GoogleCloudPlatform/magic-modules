resource "google_compute_url_map" "urlmap" {
  name        = "urlmap"
  description = "a description"

  default_service = "${google_compute_backend_service.home.self_link}"

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = "${google_compute_backend_service.home.self_link}"

    path_rule {
      paths   = ["/home"]
      service = "${google_compute_backend_service.home.self_link}"
    }

    path_rule {
      paths   = ["/login"]
      service = "${google_compute_backend_service.login.self_link}"
    }

    path_rule {
      paths   = ["/static"]
      service = "${google_compute_backend_bucket.static.self_link}"
    }
  }

  test {
    service = "${google_compute_backend_service.home.self_link}"
    host    = "hi.com"
    path    = "/home"
  }
}

resource "google_compute_backend_service" "login" {
  name        = "login"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_backend_service" "home" {
  name        = "home"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_http_health_check" "default" {
  name               = "health-check"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_backend_bucket" "static" {
  name        = "static-asset-backend-bucket"
  bucket_name = "${google_storage_bucket.static.name}"
  enable_cdn  = true
}

resource "google_storage_bucket" "static" {
  name     = "static-asset-bucket"
  location = "US"
}
