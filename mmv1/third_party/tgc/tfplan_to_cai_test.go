// Package tfplan2cai converts tfplan to CAI assets.
package tfplan2cai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert_noPlanJSON(t *testing.T) {
	ctx := context.Background()
	jsonPlan := []byte{}
	options := &Options{Offline: true}
	assets, err := Convert(ctx, jsonPlan, options)
	assert.Empty(t, assets)
	assert.Error(t, err)
}

func TestConvert_noResourceChanges(t *testing.T) {
	ctx := context.Background()
	f := "./testdata/empty-0.13.7.tfplan.json"
	jsonPlan, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("Error parsing %s: %s", f, err)
	}
	options := &Options{Offline: true}
	assets, err := Convert(ctx, jsonPlan, options)
	assert.Empty(t, assets)
	assert.Empty(t, err)
}
