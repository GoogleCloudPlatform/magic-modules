resource "google_sourcerepo_repository" "my-repo" {
  name = "my-repository-${local.name_suffix}"
}
