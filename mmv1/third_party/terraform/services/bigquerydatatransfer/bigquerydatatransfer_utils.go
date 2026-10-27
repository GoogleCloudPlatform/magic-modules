package bigquerydatatransfer

import (
	"log"
	"strings"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/googleapi"
)

func transformBigqueryDataTransferDataSourceEnrollmentReadError(err error) error {
	if gErr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error); ok {
		if gErr.Code == 400 && strings.Contains(gErr.Message, "BigQuery DataTransfer is not enabled for") {
			// Reading a data source that is not enrolled returns 400 FAILED_PRECONDITION rather than
			// 404. HandleNotFoundError only recognises 404, so remap the code to get the desired
			// behaviour when the enrollment has been removed outside of Terraform.
			gErr.Code = 404
		}

		log.Printf("[DEBUG] Transformed BigqueryDataTransfer data source enrollment error")
		return gErr
	}

	return err
}
