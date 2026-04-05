package apigee

import (
	"log"
	"strings"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/googleapi"
)

// transformApigeeOrganizationReadError converts a 403 "permission denied (or
// it may not exist)" error into a 404 so that Terraform's standard
// HandleNotFoundError logic can detect that the organization has been deleted
// out-of-band and propose re-creation rather than surfacing an opaque
// permission error.
//
// Background: when an Apigee organization is deleted via the management API
// (not through Terraform), Terraform still holds state for the resource.  On
// the next plan/apply Terraform calls the Read function which GETs
// /v1/organizations/<name>.  Because the organization no longer exists the
// Apigee API returns HTTP 403 with a message containing "(or it may not
// exist)" instead of the more intuitive 404 — this is intentional API
// behaviour to avoid leaking existence information to callers that lack
// read access.  From Terraform's perspective this ambiguous 403 must be
// treated as "resource is gone" so that the plan shows an add rather than
// failing with an access-denied error.
func transformApigeeOrganizationReadError(err error) error {
	if gErr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error); ok {
		if gErr.Code == 403 && strings.Contains(gErr.Message, "(or it may not exist)") {
			// Rewrite the status code so HandleNotFoundError treats this as a
			// deleted resource and schedules re-creation on the next apply.
			gErr.Code = 404
		}

		log.Printf("[DEBUG] Transformed ApigeeOrganization read error: %v", gErr)
		return gErr
	}

	return err
}
