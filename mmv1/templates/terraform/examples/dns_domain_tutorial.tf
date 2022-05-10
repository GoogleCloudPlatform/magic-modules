[START dns_domain_tutorial]
# to setup a web-server
resource "google_compute_instance" "default" {
  name         = "dns-compute-instance"
  machine_type = "g1-small"
  zone         = "us-central1-b"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = "default"
       access_config {
           // Ephemeral public IP
   }
  }
  metadata_startup_script = "sudo apt-get update && sudo apt-get install apache2 -y && echo '<!doctype html><html><body><h1>Hello World!</h1>'`date`'</body></html>' > /var/www/html/index.html"
}

# to allow http traffic
resource "google_compute_firewall" "default" {
  name    = "allow-http-traffic"
  network = "default"
  allow {
    ports    = ["80"]
    protocol = "tcp"
  }
  source_ranges = ["0.0.0.0/0"]
}


# to create a DNS zone
resource "google_dns_managed_zone" "default" {
  name          = "example-zone-googlecloudexample"
  dns_name      = "googlecloudexample.com."
  description   = "Example DNS zone"
  force_destroy = "true"
}

# to register web-server's ip address in DNS
resource "google_dns_record_set" "default" {
  managed_zone = google_dns_managed_zone.default.name
  name         = "server-ip-record.${google_dns_managed_zone.default.dns_name}"
  type         = "A"
  ttl          = 300
  rrdatas      = [google_compute_instance.default.network_interface[0].access_config[0].nat_ip]
}

[END dns_domain_tutorial]