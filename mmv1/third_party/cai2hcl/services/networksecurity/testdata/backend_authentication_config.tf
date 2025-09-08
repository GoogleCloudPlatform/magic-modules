resource "google_network_security_backend_authentication_config" "laurenzk-test1" {
  location         = "global"
  name             = "laurenzk-test1"
  project          = "ccm-breakit"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test2" {
  location         = "global"
  name             = "laurenzk-test2"
  project          = "ccm-breakit"
  trust_config     = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test3" {
  client_certificate = "projects/ccm-breakit/locations/global/certificates/anatolisaukhin-27101"
  location           = "global"
  name               = "laurenzk-test3"
  project            = "ccm-breakit"
  well_known_roots   = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test4" {
  client_certificate = "projects/ccm-breakit/locations/global/certificates/anatolisaukhin-27101"
  location           = "global"
  name               = "laurenzk-test4"
  project            = "ccm-breakit"
  trust_config       = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots   = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test5" {
  location         = "global"
  name             = "laurenzk-test5"
  project          = "ccm-breakit"
  trust_config     = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots = "NONE"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test7" {
  client_certificate = "projects/ccm-breakit/locations/global/certificates/anatolisaukhin-27101"
  location           = "global"
  name               = "laurenzk-test7"
  project            = "ccm-breakit"
  trust_config       = "projects/ccm-breakit/locations/global/trustConfigs/id-2de0d4b7-89cf-476f-893d-4567b3791ca9"
  well_known_roots   = "NONE"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test8" {
  description = "My test description"

  labels = {
    foo = "bar"
  }

  location         = "global"
  name             = "laurenzk-test8"
  project          = "ccm-breakit"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test9" {
  location         = "europe-west1"
  name             = "laurenzk-test9"
  project          = "ccm-breakit"
  well_known_roots = "PUBLIC_ROOTS"
}

resource "google_network_security_backend_authentication_config" "laurenzk-test10" {
  location         = "global"
  name             = "laurenzk-test10"
  project          = "ccm-breakit"
  well_known_roots = "PUBLIC_ROOTS"
}
