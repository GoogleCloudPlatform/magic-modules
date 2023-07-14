package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func getRoleEntityPair(role_entity string) (*storage.RoleEntity, error) {
	return storage.GetRoleEntityPair(role_entity)
}
