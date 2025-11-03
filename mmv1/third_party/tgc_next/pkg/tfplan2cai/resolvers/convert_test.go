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

	idToResourceChangeMap := NewIamAdvancedResolver(logger).Resolve(jsonPlan)

	assert.Equal(t, 1, len(idToResourceChangeMap), "Expected map size is 1")
	assert.Equal(t, 2, len(idToResourceChangeMap["instance_name/google_compute_instance.tgc-iam.name/project/terraform-dev-zhenhuali/zone/us-central1-a/"]), "Expected iam list to be size 2")
	assert.Equal(t, 0, len(idToResourceChangeMap["google_compute_instance_iam_member.foo1"]), "Expected this key to return null")
}
