resource "google_eventarc_trigger" "primary" {
	name = "{{name}}"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "europe-west1"
		}
	}
	labels = {
		foo = "bar"
	}
}

resource "google_pubsub_topic" "foo" {
	name = "{{topic}}"
}

resource "google_cloud_run_service" "default" {
	name     = "{{event_arc_service}}"
	location = "europe-west1"

	metadata {
		namespace = "{{project}}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				args  = ["arrgs"]
			}
		container_concurrency = 50
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}