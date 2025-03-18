package acctest

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

type IamMember struct {
	Member, Role string
}

// BootstrapIamMembers ensures that a given set of member/role pairs exist in the default
// test project. This should be used to avoid race conditions that can happen on the
// default project due to parallel tests managing the same member/role pairings. Members
// will have `{project_number}` replaced with the default test project's project number.
func BootstrapIamMembers(t *testing.T, members []IamMember) {
	config := BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap a config for BootstrapIamMembers.")
	}
	client := config.NewResourceManagerClient(config.UserAgent)

	// Get the project since we need its number, id, and policy.
	project, err := client.Projects.Get(envvar.GetTestProjectFromEnv()).Do()
	if err != nil {
		t.Fatalf("Error getting project with id %q: %s", project.ProjectId, err)
	}

	// Get the organization ID from environment if any
	orgIdFromEnv := envvar.GetTestOrgFromEnv(t)

	// Separate bindings into project-level vs. org-level
	var projectBindings []*cloudresourcemanager.Binding
	var orgBindings []*cloudresourcemanager.Binding

	for _, member := range members {
		// Replace {project_number} and {organization_id} if present
		replacedMember := strings.ReplaceAll(member.Member, "{project_number}", strconv.FormatInt(project.ProjectNumber, 10))
		replacedMember = strings.ReplaceAll(replacedMember, "{organization_id}", orgIdFromEnv)

		// If the original member string indicates organization usage, store it as org binding
		if strings.Contains(member.Member, "{organization_id}") {
			orgBindings = append(orgBindings, &cloudresourcemanager.Binding{
				Role:    member.Role,
				Members: []string{replacedMember},
			})
		} else {
			// Otherwise, treat it as a project binding
			projectBindings = append(projectBindings, &cloudresourcemanager.Binding{
				Role:    member.Role,
				Members: []string{replacedMember},
			})
		}
	}

	// Apply project-level bindings if any
	if len(projectBindings) > 0 {
		applyProjectIamBindings(t, client, project.ProjectId, projectBindings)
	}

	// Apply org-level bindings if any
	if len(orgBindings) > 0 {
		if orgIdFromEnv == "" {
			t.Fatal("Error: Org-level IAM was requested, but no organization ID was set in the environment.")
		}
		orgName := "organizations/" + orgIdFromEnv
		applyOrgIamBindings(t, client, orgName, orgBindings)
	}
}

func applyProjectIamBindings(t *testing.T,
	client *cloudresourcemanager.Service,
	projectId string,
	newBindings []*cloudresourcemanager.Binding) {

	// Retry bootstrapping with exponential backoff for concurrent writes
	backoff := time.Second
	for {
		getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
		policy, err := client.Projects.GetIamPolicy(projectId, getPolicyRequest).Do()
		if transport_tpg.IsGoogleApiErrorWithCode(err, 429) {
			t.Logf("[DEBUG] 429 while attempting to read policy for project %s, waiting %v before attempting again", projectId, backoff)
			time.Sleep(backoff)
			continue
		} else if err != nil {
			t.Fatalf("Error getting iam policy for project %s: %v\n", projectId, err)
		}

		mergedBindings := tpgiamresource.MergeBindings(append(policy.Bindings, newBindings...))

		if tpgiamresource.CompareBindings(policy.Bindings, mergedBindings) {
			t.Logf("[DEBUG] All bindings already present for project %s", projectId)
			break
		}
		// The policy must change.
		policy.Bindings = mergedBindings
		setPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
		policy, err = client.Projects.SetIamPolicy(projectId, setPolicyRequest).Do()
		if err == nil {
			t.Logf("[DEBUG] Waiting for IAM bootstrapping to propagate for project %s.", projectId)
			time.Sleep(3 * time.Minute)
			break
		}
		if tpgresource.IsConflictError(err) {
			t.Logf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				t.Fatalf("Error applying IAM policy to %s: Too many conflicts.  Latest error: %s", projectId, err)
			}
			continue
		}
		t.Fatalf("Error setting project iam policy: %v", err)
	}
}

func applyOrgIamBindings(
	t *testing.T,
	client *cloudresourcemanager.Service,
	orgName string,
	newBindings []*cloudresourcemanager.Binding) {

	// Retry bootstrapping with exponential backoff for concurrent writes
	backoff := time.Second
	for {
		getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
		policy, err := client.Organizations.GetIamPolicy(orgName, getPolicyRequest).Do()
		if transport_tpg.IsGoogleApiErrorWithCode(err, 429) {
			t.Logf("[DEBUG] 429 while attempting to read policy for org %s, waiting %v before attempting again", orgName, backoff)
			time.Sleep(backoff)
			continue
		} else if err != nil {
			t.Fatalf("Error getting iam policy for org %s: %v\n", orgName, err)
		}

		mergedBindings := tpgiamresource.MergeBindings(append(policy.Bindings, newBindings...))

		if tpgiamresource.CompareBindings(policy.Bindings, mergedBindings) {
			t.Logf("[DEBUG] All bindings already present for org %s", orgName)
			break
		}
		// The policy must change.
		policy.Bindings = mergedBindings
		setPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
		policy, err = client.Organizations.SetIamPolicy(orgName, setPolicyRequest).Do()
		if err == nil {
			t.Logf("[DEBUG] Waiting for IAM bootstrapping to propagate for org %s.", orgName)
			time.Sleep(3 * time.Minute)
			break
		}
		if tpgresource.IsConflictError(err) {
			t.Logf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				t.Fatalf("Error applying IAM policy to %s: Too many conflicts.  Latest error: %s", orgName, err)
			}
			continue
		}
		t.Fatalf("Error setting org iam policy: %v", err)
	}
}
