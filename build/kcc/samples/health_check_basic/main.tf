resource "google_compute_health_check" "internal-health-check" {
 name = "internal-service-health-check"

 timeout_sec        = 1
 check_interval_sec = 1

 tcp_health_check {
   port = "80"
 }
}
