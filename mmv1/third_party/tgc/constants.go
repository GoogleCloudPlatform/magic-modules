package google

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

// ErrNoConversion can be returned if a conversion is unable to be returned.

// because of the current state of the system.
// Example: The conversion requires that the resource has already been created
// and is now being updated).
var ErrNoConversion = tpgresource.ErrNoConversion

// ErrEmptyIdentityField can be returned when fetching a resource is not possible
// due to the identity field of that resource returning empty.
var ErrEmptyIdentityField = tpgresource.ErrEmptyIdentityField

// ErrResourceInaccessible can be returned when fetching an IAM resource
// on a project that has not yet been created or if the service account
// lacks sufficient permissions
var ErrResourceInaccessible = tpgresource.ErrResourceInaccessible

// Global MutexKV
//
// Deprecated: For backward compatibility mutexKV is still working,
// but all new code should use MutexStore in the transport_tpg package instead.
var mutexKV = transport_tpg.MutexStore
