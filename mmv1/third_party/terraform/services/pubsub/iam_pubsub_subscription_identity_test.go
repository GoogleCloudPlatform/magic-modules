package pubsub

import "testing"

func TestPubsubSubscriptionIamMemberResource_HasIdentity(t *testing.T) {
	resource := PubsubSubscriptionIamMemberResource()
	if resource.Identity == nil {
		t.Fatalf("expected google_pubsub_subscription_iam_member resource identity to be configured")
	}
}

func TestStorageBucketIamParentResourceIdentityParser(t *testing.T) {
	tests := []struct {
		name              string
		projectValue      string
		subscriptionValue string
		want              string
	}{
		{
			name:              "short subscription name",
			projectValue:      "my-project",
			subscriptionValue: "my-subscription",
			want:              "projects/my-project/subscriptions/my-subscription",
		},
		{
			name:              "canonical subscription name",
			projectValue:      "my-project",
			subscriptionValue: "projects/my-project/subscriptions/my-subscription",
			want:              "projects/my-project/subscriptions/my-subscription",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := PubsubSubscriptionIamMemberResource()
			rd := resource.TestResourceData()
			identity, err := rd.Identity()
			if err != nil {
				t.Fatalf("creating identity: %v", err)
			}
			if err := identity.Set("project", tt.projectValue); err != nil {
				t.Fatalf("setting identity project: %v", err)
			}
			if err := identity.Set("subscription", tt.subscriptionValue); err != nil {
				t.Fatalf("setting identity subscription: %v", err)
			}

			got, err := PubsubSubscriptionIamParentResourceIdentityParser(rd, identity, nil)
			if err != nil {
				t.Fatalf("parsing identity: %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected parsed id: got %q, want %q", got, tt.want)
			}
		})
	}
}
