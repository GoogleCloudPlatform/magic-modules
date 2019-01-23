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

var IamKmsCryptoKeySchema = map[string]*schema.Schema{
	"crypto_key_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type KmsCryptoKeyIamUpdater struct {
	resourceId string
	iamHandle  *iam.Handle
	Config     *Config
}

func NewKmsCryptoKeyIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	cryptoKey := d.Get("crypto_key_id").(string)
	cryptoKeyId, err := parseKmsCryptoKeyId(cryptoKey, config)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error parsing resource ID for for %s: {{err}}", cryptoKey), err)
	}

	return &KmsCryptoKeyIamUpdater{
		resourceId: cryptoKeyId.cryptoKeyId(),
		iamHandle: config.clientKms.CryptoKeyIAM(&kmspb.CryptoKey{
			Name: cryptoKeyId.cryptoKeyId(),
		}),
		Config: config,
	}, nil
}

func CryptoIdParseFunc(d *schema.ResourceData, config *Config) error {
	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}
	d.Set("crypto_key_id", cryptoKeyId.cryptoKeyId())
	d.SetId(cryptoKeyId.cryptoKeyId())
	return nil
}

func (u *KmsCryptoKeyIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
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

func (u *KmsCryptoKeyIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
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

func (u *KmsCryptoKeyIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *KmsCryptoKeyIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-kms-crypto-key-%s", u.resourceId)
}

func (u *KmsCryptoKeyIamUpdater) DescribeResource() string {
	return fmt.Sprintf("KMS CryptoKey %q", u.resourceId)
}
