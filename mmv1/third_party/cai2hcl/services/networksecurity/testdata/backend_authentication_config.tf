resource "google_network_security_backend_authentication_config" "laurenzk-test1" {
  name             = "laurenzk-test1"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test2" {
  name             = "laurenzk-test2"
  trust_config     = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test3" {
  client_certificate = "projects/ccm-breakit/locations/global/certificates/anatolisaukhin-27101"
  name               = "laurenzk-test3"
  well_known_roots   = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test4" {
  client_certificate = "projects/ccm-breakit/locations/global/certificates/anatolisaukhin-27101"
  name               = "laurenzk-test4"
  trust_config       = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots   = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test5" {
  name             = "laurenzk-test5"
  trust_config     = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots = "NONE"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test7" {
  client_certificate = "projects/ccm-breakit/locations/global/certificates/anatolisaukhin-27101"
  name               = "laurenzk-test7"
  trust_config       = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots   = "NONE"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test8" {
  description = "My test description"

  labels = {
    foo = "bar"
  }

  name             = "laurenzk-test8"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test9" {
  location         = "europe-west1"
  name             = "laurenzk-test9"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test10" {
  name             = "laurenzk-test10"
  project          = "ccm-breakit"
  well_known_roots = "PUBLIC_ROOTS"
}
