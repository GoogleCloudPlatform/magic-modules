package tfplan2cai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestConvert_iamBinding(t *testing.T) {
	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Error initializing logger %s", err)
	}
	f := "iamBinding.tfplan.json"
	jsonPlan, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("Error parsing %s: %s", f, err)
	}
	options := &Options{
		Offline:     true,
		ErrorLogger: logger,
		AncestryCache: map[string]string{
			"projects/terraform-dev-zhenhuali": "organizations/529579013760",
		},
	}
	assets, err := Convert(ctx, jsonPlan, options)
	if err != nil {
		t.Fatalf("Error parsing %s", err)
	}

	jsonData, err := json.MarshalIndent(assets, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
