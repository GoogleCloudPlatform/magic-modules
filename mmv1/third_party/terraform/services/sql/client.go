package sql

import (
	"log"

	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *sqladmin.Service {
	sqlClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c)))
	log.Printf("[INFO] Instantiating Google SqlAdmin client for path %s", sqlClientBasePath)
	clientSqlAdmin, err := sqladmin.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage: %s", err)
		return nil
	}
	clientSqlAdmin.UserAgent = userAgent
	clientSqlAdmin.BasePath = sqlClientBasePath

	return clientSqlAdmin
}
