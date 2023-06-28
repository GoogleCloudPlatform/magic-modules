package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func resourceStorageNotificationParseID(id string) (string, string) {
	return storage.ResourceStorageNotificationParseID(id)
}
