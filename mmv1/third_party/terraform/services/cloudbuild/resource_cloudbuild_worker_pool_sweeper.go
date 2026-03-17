package cloudbuild

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
)

func init() {
	sweeper.AddTestSweepersLegacy("CloudbuildWorkerPool", testSweepCloudbuildWorkerPool)
}

func testSweepCloudbuildWorkerPool(region string) error {
	log.Print("[INFO][SWEEPER_LOG] Starting sweeper for CloudbuildWorkerPool")

	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	t := &testing.T{}
	billingId := envvar.GetTestBillingAccountFromEnv(t)

	// Setup variables to be used for Delete arguments.
	d := map[string]string{
		"project":         config.Project,
		"region":          region,
		"location":        region,
		"zone":            "-",
		"billing_account": billingId,
	}

	client := NewDCLCloudbuildClient(config, config.UserAgent, "", 0)
	err = client.DeleteAllWorkerPool(context.Background(), d["project"], d["location"], isDeletableCloudbuildWorkerPool)
	if err != nil {
		return err
	}
	return nil
}

func isDeletableCloudbuildWorkerPool(r *WorkerPool) bool {
	return sweeper.IsSweepableTestResource(*r.Name)
}
