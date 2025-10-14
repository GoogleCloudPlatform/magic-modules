package resolvers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConvert_iamBinding(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Error initializing logger %s", err)
	}
	f := "iamBinding.tfplan.json"
	jsonPlan, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("Error parsing %s: %s", f, err)
	}
	resourceDataMap := NewDefaultPreResolver(logger).Resolve(jsonPlan)
	resourceDataMap = NewAdvancedResolver(logger).Resolve(jsonPlan, resourceDataMap)

	assert.Equal(t, 2, len(resourceDataMap), "Expected map size is 2")
	assert.Equal(t, 2, len(resourceDataMap["google_compute_instance_iam_member.foo"]), "Expected iam list to be size 2")
	assert.Equal(t, 0, len(resourceDataMap["google_compute_instance_iam_member.foo1"]), "Expected this key to return null")
}
