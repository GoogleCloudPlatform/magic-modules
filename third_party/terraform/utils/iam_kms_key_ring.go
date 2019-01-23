package google

import (
	"context"
	"fmt"

	"cloud.google.com/go/iam"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/grpc/metadata"
)

var IamKmsKeyRingSchema = map[string]*schema.Schema{
	"key_ring_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type KmsKeyRingIamUpdater struct {
	resourceId string
	iamHandle  *iam.Handle
	Config     *Config
}

func NewKmsKeyRingIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	keyRing := d.Get("key_ring_id").(string)
	keyRingId, err := parseKmsKeyRingId(keyRing, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", keyRing), err)
	}

	return &KmsKeyRingIamUpdater{
		resourceId: keyRingId.keyRingId(),
		iamHandle: config.clientKms.KeyRingIAM(&kmspb.KeyRing{
			Name: keyRingId.keyRingId(),
		}),
		Config: config,
	}, nil
}

func KeyRingIdParseFunc(d *schema.ResourceData, config *Config) error {
	keyRingId, err := parseKmsKeyRingId(d.Id(), config)
	if err != nil {
		return err
	}

	d.Set("key_ring_id", keyRingId.keyRingId())
	d.SetId(keyRingId.keyRingId())
	return nil
}

func (u *KmsKeyRingIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"x-goog-request-params", fmt.Sprintf("name=%v", u.resourceId))
	p, err := u.iamHandle.Policy(ctx)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := kmsToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *KmsKeyRingIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	kmsPolicy, err := resourceManagerToKmsPolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"x-goog-request-params", fmt.Sprintf("name=%v", u.resourceId))
	if err := u.iamHandle.SetPolicy(ctx, kmsPolicy); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *KmsKeyRingIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *KmsKeyRingIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-kms-key-ring-%s", u.resourceId)
}

func (u *KmsKeyRingIamUpdater) DescribeResource() string {
	return fmt.Sprintf("KMS KeyRing %q", u.resourceId)
}

func resourceManagerToKmsPolicy(p *cloudresourcemanager.Policy) (*iam.Policy, error) {
	var out iam.Policy
	if err := Convert(p, out); err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a kms policy: {{err}}", err)
	}
	return &out, nil
}

func kmsToResourceManagerPolicy(p *iam.Policy) (*cloudresourcemanager.Policy, error) {
	var out cloudresourcemanager.Policy
	if err := Convert(p, out); err != nil {
		return nil, errwrap.Wrapf("Cannot convert a kms policy to a v1 policy: {{err}}", err)
	}
	return &out, nil
}
