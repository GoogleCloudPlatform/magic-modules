package tpgiamresource

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

// fakeIamUpdater is an in-memory ResourceIamUpdater. Like the real IAM APIs, it rejects a write
// with a stale etag with a 409.
type fakeIamUpdater struct {
	mu     sync.Mutex
	key    string
	policy *cloudresourcemanager.Policy
	etag   int
	gets   int
	sets   int
	// Called after the nth read (1-based) to simulate another writer; returns whether it changed p.
	afterGet func(n int, p *cloudresourcemanager.Policy) bool
}

func newFakeIamUpdater(t *testing.T, bindings ...*cloudresourcemanager.Binding) *fakeIamUpdater {
	// The mutex key is shared across tests through the global mutex store, so keep it unique.
	return &fakeIamUpdater{key: t.Name(), policy: &cloudresourcemanager.Policy{Bindings: bindings}}
}

func (f *fakeIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.gets++
	p, err := copyIamPolicy(f.policy)
	if err != nil {
		return nil, err
	}
	p.Etag = strconv.Itoa(f.etag)
	if f.afterGet != nil && f.afterGet(f.gets, f.policy) {
		f.etag++
	}
	return p, nil
}

func (f *fakeIamUpdater) SetResourceIamPolicy(p *cloudresourcemanager.Policy) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sets++
	if p.Etag != strconv.Itoa(f.etag) {
		return &googleapi.Error{Code: 409, Message: "There were concurrent policy changes."}
	}
	stored, err := copyIamPolicy(p)
	if err != nil {
		return err
	}
	f.policy = stored
	f.etag++
	return nil
}

func (f *fakeIamUpdater) GetMutexKey() string      { return f.key }
func (f *fakeIamUpdater) GetResourceId() string    { return f.key }
func (f *fakeIamUpdater) DescribeResource() string { return "fake resource " + f.key }

func setBinding(role string, members ...string) iamPolicyModifyFunc {
	b := &cloudresourcemanager.Binding{Role: role, Members: members}
	return func(p *cloudresourcemanager.Policy) error {
		p.Bindings = append(filterBindingsWithRoleAndCondition(p.Bindings, role, nil), b)
		return nil
	}
}

func addMember(role, member string) iamPolicyModifyFunc {
	return func(p *cloudresourcemanager.Policy) error {
		p.Bindings = MergeBindings(append(p.Bindings, &cloudresourcemanager.Binding{Role: role, Members: []string{member}}))
		return nil
	}
}

func removeBinding(role string) iamPolicyModifyFunc {
	return func(p *cloudresourcemanager.Policy) error {
		p.Bindings = filterBindingsWithRoleAndCondition(p.Bindings, role, nil)
		return nil
	}
}

func assertMembers(t *testing.T, p *cloudresourcemanager.Policy, role string, want ...string) {
	t.Helper()
	key := iamBindingKey{Role: role}
	got := createIamBindingsMap(p.Bindings)[key]
	wantSet := createIamBindingsMap([]*cloudresourcemanager.Binding{{Role: role, Members: want}})[key]
	if !reflect.DeepEqual(got, wantSet) {
		t.Errorf("members of %s: got %v, want %v", role, got, wantSet)
	}
}

func assertErrorContains(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("expected an error containing %q, got %v", want, err)
	}
}

func TestIamPolicyReadModifyWrite_writeAllowed(t *testing.T) {
	t.Parallel()

	t.Run("veto skips the write", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t, &cloudresourcemanager.Binding{Role: "role-1", Members: []string{"user:a@example.com"}})
		var sawBefore, sawAfter *cloudresourcemanager.Policy
		err := iamPolicyReadModifyWrite(f, setBinding("role-1", "user:b@example.com"), func(before, after *cloudresourcemanager.Policy) error {
			sawBefore, sawAfter = before, after
			return iamPolicyNoExistingBindingChanged(before, after)
		})
		assertErrorContains(t, err, `role "role-1" already exists`)
		if f.sets != 0 {
			t.Errorf("expected no write, got %d", f.sets)
		}
		assertMembers(t, f.policy, "role-1", "user:a@example.com")
		// before is an independent copy of the policy as read; after has the change applied.
		assertMembers(t, sawBefore, "role-1", "user:a@example.com")
		assertMembers(t, sawAfter, "role-1", "user:b@example.com")
	})

	t.Run("allowed write happens", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t, &cloudresourcemanager.Binding{Role: "role-1", Members: []string{"user:a@example.com"}})
		if err := iamPolicyReadModifyWrite(f, setBinding("role-2", "user:b@example.com"), iamPolicyNoExistingBindingChanged); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if f.sets != 1 {
			t.Errorf("expected one write, got %d", f.sets)
		}
		assertMembers(t, f.policy, "role-1", "user:a@example.com")
		assertMembers(t, f.policy, "role-2", "user:b@example.com")
	})

	t.Run("no hook overwrites", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t, &cloudresourcemanager.Binding{Role: "role-1", Members: []string{"user:a@example.com"}})
		if err := iamPolicyReadModifyWrite(f, setBinding("role-1", "user:b@example.com"), nil); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		assertMembers(t, f.policy, "role-1", "user:b@example.com")
	})
}

func TestIamPolicyReadModifyWrite_conflictRerunsHook(t *testing.T) {
	t.Parallel()

	t.Run("conflicting concurrent change is vetoed on retry", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t)
		f.afterGet = func(n int, p *cloudresourcemanager.Policy) bool {
			if n != 1 {
				return false
			}
			p.Bindings = append(p.Bindings, &cloudresourcemanager.Binding{Role: "role-1", Members: []string{"user:x@example.com"}})
			return true
		}
		hookCalls := 0
		err := iamPolicyReadModifyWrite(f, setBinding("role-1", "user:a@example.com"), func(before, after *cloudresourcemanager.Policy) error {
			hookCalls++
			return iamPolicyNoExistingBindingChanged(before, after)
		})
		assertErrorContains(t, err, `role "role-1" already exists`)
		if hookCalls != 2 {
			t.Errorf("expected the hook to run on both attempts, got %d calls", hookCalls)
		}
		if f.sets != 1 {
			t.Errorf("expected only the rejected write, got %d", f.sets)
		}
		assertMembers(t, f.policy, "role-1", "user:x@example.com")
	})

	t.Run("unrelated concurrent change is kept", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t)
		f.afterGet = func(n int, p *cloudresourcemanager.Policy) bool {
			if n != 1 {
				return false
			}
			p.Bindings = append(p.Bindings, &cloudresourcemanager.Binding{Role: "role-2", Members: []string{"user:x@example.com"}})
			return true
		}
		hookCalls := 0
		err := iamPolicyReadModifyWrite(f, setBinding("role-1", "user:a@example.com"), func(before, after *cloudresourcemanager.Policy) error {
			hookCalls++
			return iamPolicyNoExistingBindingChanged(before, after)
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if hookCalls != 2 {
			t.Errorf("expected the hook to run on both attempts, got %d calls", hookCalls)
		}
		assertMembers(t, f.policy, "role-1", "user:a@example.com")
		assertMembers(t, f.policy, "role-2", "user:x@example.com")
	})
}

func TestAllowBatchedWrite(t *testing.T) {
	strict := iamPolicyWriteAllowedFunc(iamPolicyNoExistingBindingChanged)
	existing := []*cloudresourcemanager.Binding{{Role: "role-1", Members: []string{"user:x@example.com"}}}
	testCases := []struct {
		name    string
		before  []*cloudresourcemanager.Binding
		changes []iamPolicyChange
		wantErr string
	}{
		{
			name: "member added after a strict binding",
			changes: []iamPolicyChange{
				{modify: setBinding("role-1", "user:a@example.com"), allowWrite: strict},
				{modify: addMember("role-1", "user:b@example.com")},
			},
		},
		{
			name: "strict binding after a member for the same role",
			changes: []iamPolicyChange{
				{modify: addMember("role-1", "user:b@example.com")},
				{modify: setBinding("role-1", "user:a@example.com"), allowWrite: strict},
			},
			wantErr: `role "role-1" already exists`,
		},
		{
			name:   "strict binding after the same role is deleted",
			before: existing,
			changes: []iamPolicyChange{
				{modify: removeBinding("role-1")},
				{modify: setBinding("role-1", "user:a@example.com"), allowWrite: strict},
			},
		},
		{
			name:    "strict binding on an existing role",
			before:  existing,
			changes: []iamPolicyChange{{modify: setBinding("role-1", "user:a@example.com"), allowWrite: strict}},
			wantErr: `role "role-1" already exists`,
		},
		{
			name: "strict bindings on different roles",
			changes: []iamPolicyChange{
				{modify: setBinding("role-1", "user:a@example.com"), allowWrite: strict},
				{modify: setBinding("role-2", "user:b@example.com"), allowWrite: strict},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			before := &cloudresourcemanager.Policy{Bindings: tc.before}
			original, err := copyIamPolicy(before)
			if err != nil {
				t.Fatal(err)
			}
			err = allowBatchedWrite(tc.changes)(before, nil)
			if tc.wantErr == "" && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if tc.wantErr != "" {
				assertErrorContains(t, err, tc.wantErr)
			}
			if !CompareBindings(before.Bindings, original.Bindings) {
				t.Errorf("before was modified: %s", DebugPrintBindings(before.Bindings))
			}
		})
	}
}

func TestSendBatchModifyIamPolicy(t *testing.T) {
	t.Parallel()

	t.Run("veto fails the batch without writing", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t, &cloudresourcemanager.Binding{Role: "role-1", Members: []string{"user:x@example.com"}})
		_, err := sendBatchModifyIamPolicy(f)(f.key, []iamPolicyChange{
			{modify: addMember("role-2", "user:b@example.com")},
			{modify: setBinding("role-1", "user:a@example.com"), allowWrite: iamPolicyNoExistingBindingChanged},
		})
		assertErrorContains(t, err, `role "role-1" already exists`)
		if f.sets != 0 {
			t.Errorf("expected no write, got %d", f.sets)
		}
	})

	t.Run("changes without hooks are written together", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t)
		_, err := sendBatchModifyIamPolicy(f)(f.key, []iamPolicyChange{
			{modify: setBinding("role-1", "user:a@example.com")},
			{modify: addMember("role-2", "user:b@example.com")},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if f.sets != 1 {
			t.Errorf("expected one write, got %d", f.sets)
		}
		assertMembers(t, f.policy, "role-1", "user:a@example.com")
		assertMembers(t, f.policy, "role-2", "user:b@example.com")
	})

	t.Run("unexpected body type", func(t *testing.T) {
		t.Parallel()
		f := newFakeIamUpdater(t)
		if _, err := sendBatchModifyIamPolicy(f)(f.key, "not a list of changes"); err == nil {
			t.Error("expected an error")
		}
	})
}

func TestIamBindingOverwriteOnCreate(t *testing.T) {
	config := func(v cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{"overwrite_on_create": v})
	}
	testCases := []struct {
		name string
		raw  cty.Value
		want bool
	}{
		{name: "no raw config", raw: cty.NilVal, want: true},
		{name: "unset", raw: config(cty.NullVal(cty.Bool)), want: true},
		{name: "unknown", raw: config(cty.UnknownVal(cty.Bool)), want: true},
		{name: "true", raw: config(cty.True), want: true},
		{name: "false", raw: config(cty.False), want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := (&schema.Resource{Schema: iamBindingSchema}).Data(&terraform.InstanceState{ID: "id", RawConfig: tc.raw})
			if got := iamBindingOverwriteOnCreate(d); got != tc.want {
				t.Errorf("got %t, want %t", got, tc.want)
			}
		})
	}
}
