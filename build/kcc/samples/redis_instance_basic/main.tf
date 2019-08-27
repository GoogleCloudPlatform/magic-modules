resource "google_redis_instance" "cache" {
  name           = "memory-cache"
  memory_size_gb = 1
}
