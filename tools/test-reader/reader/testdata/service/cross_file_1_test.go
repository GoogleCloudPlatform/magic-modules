package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCrossFile(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"test1": testAccCrossFile1,
		"test2": testAccCrossFile2,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccCrossFile1(t *testing.T) {
	VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccCrossFileConfig1(),
			},
		},
	})
}
