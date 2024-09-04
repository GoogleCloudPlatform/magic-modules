<% autogen_exception -%>
package dataproc

import (
        "testing"
)


func TestCloudDataprocBatchRuntimeConfigVersionDiffSuppress(t *testing.T) {
   cases := map[string]struct {
      Old, New           string
      ExpectDiffSuppress bool
   }{
      "empty_default_version": {
         Old:                "",
         New:                "2.2.100",
         ExpectDiffSuppress: true,
      },
      // Additional cases excluded for brevity
   }

   for tn, tc := range cases {
      if CloudDataprocBatchRuntimeConfigVersionDiffSuppressFunc(tc.Old, tc.New) != tc.ExpectDiffSuppress {
         t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
      }
   }
}
